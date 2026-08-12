//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

package cloudstack

import (
	"encoding/base64"
	"fmt"
	"log"
	"slices"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gopkg.in/yaml.v3"
)

func dataSourceCloudstackKubernetesClusterConfig() *schema.Resource {
	return &schema.Resource{
		Read: datasourceCloudStackKubernetesClusterConfigRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the Kubernetes cluster to retrieve the config for.",
			},

			//Computed values
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the Kubernetes cluster.",
			},

			"config_data": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The raw kubeconfig of the Kubernetes cluster.",
			},

			"endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the Kubernetes API server, taken from the kubeconfig.",
			},

			"cluster_ca_certificate": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The PEM encoded certificate authority of the Kubernetes API server.",
			},

			"client_certificate": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The PEM encoded client certificate used to authenticate against the Kubernetes API server.",
			},

			"client_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The PEM encoded client key used to authenticate against the Kubernetes API server.",
			},
		},
	}
}

// kubeConfig models only the subset of a kubeconfig document that this data
// source exposes as separate attributes.
type kubeConfig struct {
	CurrentContext string              `yaml:"current-context"`
	Clusters       []kubeConfigCluster `yaml:"clusters"`
	Contexts       []kubeConfigContext `yaml:"contexts"`
	Users          []kubeConfigUser    `yaml:"users"`
}

type kubeConfigCluster struct {
	Name    string                `yaml:"name"`
	Cluster kubeConfigClusterInfo `yaml:"cluster"`
}

type kubeConfigClusterInfo struct {
	Server                   string `yaml:"server"`
	CertificateAuthorityData string `yaml:"certificate-authority-data"`
}

type kubeConfigContext struct {
	Name    string                `yaml:"name"`
	Context kubeConfigContextInfo `yaml:"context"`
}

type kubeConfigContextInfo struct {
	Cluster string `yaml:"cluster"`
	User    string `yaml:"user"`
}

type kubeConfigUser struct {
	Name string             `yaml:"name"`
	User kubeConfigUserInfo `yaml:"user"`
}

type kubeConfigUserInfo struct {
	ClientCertificateData string `yaml:"client-certificate-data"`
	ClientKeyData         string `yaml:"client-key-data"`
}

// kubernetesClusterCredentials holds the values extracted from a kubeconfig that
// are needed to connect to the Kubernetes API server.
type kubernetesClusterCredentials struct {
	Endpoint             string
	ClusterCACertificate string
	ClientCertificate    string
	ClientKey            string
}

func datasourceCloudStackKubernetesClusterConfigRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	clusterID := d.Get("cluster_id").(string)

	log.Printf("[DEBUG] Retrieving config of Kubernetes Cluster %s", clusterID)

	p := cs.Kubernetes.NewGetKubernetesClusterConfigParams()
	p.SetId(clusterID)

	config, err := cs.Kubernetes.GetKubernetesClusterConfig(p)
	if err != nil {
		// CloudStack refuses to hand out a config while the cluster is still
		// starting, and when the Kubernetes service plugin is disabled.
		return fmt.Errorf("Failed to get the config of Kubernetes Cluster %s: %s", clusterID, err)
	}

	if config.Configdata == "" {
		return fmt.Errorf("Kubernetes Cluster %s returned an empty config; the cluster ID may not "+
			"exist, the cluster may still be starting, or the Kubernetes service plugin may be disabled", clusterID)
	}

	credentials, err := parseKubernetesClusterConfig(config.Configdata)
	if err != nil {
		return fmt.Errorf("Failed to parse the config of Kubernetes Cluster %s: %s", clusterID, err)
	}

	if *credentials == (kubernetesClusterCredentials{}) {
		return fmt.Errorf("The config of Kubernetes Cluster %s does not contain a cluster endpoint, "+
			"CA certificate, client certificate or client key; use config_data directly instead", clusterID)
	}

	d.SetId(config.Id)
	d.Set("name", config.Name)
	d.Set("config_data", config.Configdata)
	d.Set("endpoint", credentials.Endpoint)
	d.Set("cluster_ca_certificate", credentials.ClusterCACertificate)
	d.Set("client_certificate", credentials.ClientCertificate)
	d.Set("client_key", credentials.ClientKey)

	return nil
}

// parseKubernetesClusterConfig extracts the endpoint and the client credentials
// from a kubeconfig document. An unparsable document, an undecodable
// certificate, or a current context that names a cluster/user absent from a
// multi-entry list is an error; a kubeconfig that simply does not carry a
// cluster or a user entry at all only yields empty values, so that config_data
// stays usable.
func parseKubernetesClusterConfig(configData string) (*kubernetesClusterCredentials, error) {
	var config kubeConfig
	if err := yaml.Unmarshal([]byte(configData), &config); err != nil {
		return nil, fmt.Errorf("Invalid kubeconfig: %s", err)
	}

	// Resolve which cluster and user the current context points at.
	clusterName, userName := "", ""
	for _, context := range config.Contexts {
		if context.Name == config.CurrentContext {
			clusterName = context.Context.Cluster
			userName = context.Context.User
			break
		}
	}

	credentials := &kubernetesClusterCredentials{}

	cluster, err := findKubeConfigEntry(config.Clusters, clusterName, "cluster",
		func(c kubeConfigCluster) string { return c.Name })
	if err != nil {
		return nil, err
	}
	if cluster >= 0 {
		caCertificate, err := decodeKubernetesClusterConfigValue(
			config.Clusters[cluster].Cluster.CertificateAuthorityData, "certificate-authority-data")
		if err != nil {
			return nil, err
		}

		credentials.Endpoint = config.Clusters[cluster].Cluster.Server
		credentials.ClusterCACertificate = caCertificate
	}

	user, err := findKubeConfigEntry(config.Users, userName, "user",
		func(u kubeConfigUser) string { return u.Name })
	if err != nil {
		return nil, err
	}
	if user >= 0 {
		clientCertificate, err := decodeKubernetesClusterConfigValue(
			config.Users[user].User.ClientCertificateData, "client-certificate-data")
		if err != nil {
			return nil, err
		}

		clientKey, err := decodeKubernetesClusterConfigValue(
			config.Users[user].User.ClientKeyData, "client-key-data")
		if err != nil {
			return nil, err
		}

		credentials.ClientCertificate = clientCertificate
		credentials.ClientKey = clientKey
	}

	return credentials, nil
}

// findKubeConfigEntry returns the index of the named entry in a kubeconfig
// list. It returns -1 with no error when the list is empty, so that a
// kubeconfig which does not carry the entry at all still leaves config_data
// usable. When the current context does not name an entry at all, or there is
// only one entry to begin with, it falls back to the first one, since there is
// only one reasonable candidate either way. But when the current context names
// a specific entry that is absent from a list of more than one, guessing is
// refused: picking some other entry could silently pair one cluster's endpoint
// with a different cluster's client credentials.
func findKubeConfigEntry[T any](entries []T, name string, kind string, nameOf func(T) string) (int, error) {
	if len(entries) == 0 {
		log.Printf("[WARN] Kubeconfig does not contain any %s", kind)
		return -1, nil
	}

	if name != "" {
		if i := slices.IndexFunc(entries, func(e T) bool { return nameOf(e) == name }); i >= 0 {
			return i, nil
		}

		if len(entries) > 1 {
			return -1, fmt.Errorf("kubeconfig context references %s %q, which is not defined among %d %s entries",
				kind, name, len(entries), kind)
		}

		log.Printf("[WARN] Kubeconfig does not contain the %s %q referenced by its current context, using the only %s instead", kind, name, kind)
		return 0, nil
	}

	if len(entries) > 1 {
		log.Printf("[WARN] Kubeconfig has %d %s entries and no current context to select one, using the first %s", len(entries), kind, kind)
	}

	return 0, nil
}

// decodeKubernetesClusterConfigValue decodes a base64 encoded kubeconfig field,
// passing an absent field through as an empty string.
func decodeKubernetesClusterConfigValue(value string, field string) (string, error) {
	if value == "" {
		log.Printf("[WARN] Kubeconfig does not contain %s", field)
		return "", nil
	}

	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", fmt.Errorf("Failed to decode %s: %s", field, err)
	}

	return string(decoded), nil
}
