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
	"errors"

	"github.com/go-ini/ini"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_url": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("CLOUDSTACK_API_URL", nil),
				ConflictsWith: []string{"config", "profile"},
			},

			"api_key": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("CLOUDSTACK_API_KEY", nil),
				ConflictsWith: []string{"config", "profile"},
				Sensitive:     true,
			},

			"secret_key": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("CLOUDSTACK_SECRET_KEY", nil),
				ConflictsWith: []string{"config", "profile"},
				Sensitive:     true,
			},

			"config": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"api_url", "api_key", "secret_key"},
			},

			"profile": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"api_url", "api_key", "secret_key"},
			},

			"http_get_only": {
				Type:        schema.TypeBool,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CLOUDSTACK_HTTP_GET_ONLY", false),
			},

			"timeout": {
				Type:        schema.TypeInt,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CLOUDSTACK_TIMEOUT", 900),
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"cloudstack_autoscale_policy":     dataSourceCloudstackAutoscalePolicy(),
			"cloudstack_autoscale_vm_group":   dataSourceCloudstackAutoscaleVMGroup(),
			"cloudstack_autoscale_vm_profile": dataSourceCloudstackAutoscaleVMProfile(),
			"cloudstack_condition":            dataSourceCloudstackCondition(),
			"cloudstack_counter":              dataSourceCloudstackCounter(),
			"cloudstack_user_data":            dataSourceCloudstackUserData(),
			"cloudstack_template":             dataSourceCloudstackTemplate(),
			"cloudstack_ssh_keypair":          dataSourceCloudstackSSHKeyPair(),
			"cloudstack_instance":             dataSourceCloudstackInstance(),
			"cloudstack_network_offering":     dataSourceCloudstackNetworkOffering(),
			"cloudstack_zone":                 dataSourceCloudStackZone(),
			"cloudstack_service_offering":     dataSourceCloudstackServiceOffering(),
			"cloudstack_volume":               dataSourceCloudstackVolume(),
			"cloudstack_vpc":                  dataSourceCloudstackVPC(),
			"cloudstack_ipaddress":            dataSourceCloudstackIPAddress(),
			"cloudstack_user":                 dataSourceCloudstackUser(),
			"cloudstack_vpn_connection":       dataSourceCloudstackVPNConnection(),
			"cloudstack_pod":                  dataSourceCloudstackPod(),
			"cloudstack_domain":               dataSourceCloudstackDomain(),
			"cloudstack_project":              dataSourceCloudstackProject(),
			"cloudstack_physical_network":     dataSourceCloudStackPhysicalNetwork(),
			"cloudstack_role":                 dataSourceCloudstackRole(),
			"cloudstack_cluster":              dataSourceCloudstackCluster(),
			"cloudstack_limits":               dataSourceCloudStackLimits(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"cloudstack_affinity_group":                 resourceCloudStackAffinityGroup(),
			"cloudstack_attach_volume":                  resourceCloudStackAttachVolume(),
			"cloudstack_autoscale_policy":               resourceCloudStackAutoScalePolicy(),
			"cloudstack_autoscale_vm_group":             resourceCloudStackAutoScaleVMGroup(),
			"cloudstack_autoscale_vm_profile":           resourceCloudStackAutoScaleVMProfile(),
			"cloudstack_cni_configuration":              resourceCloudStackCniConfiguration(),
			"cloudstack_condition":                      resourceCloudStackCondition(),
			"cloudstack_configuration":                  resourceCloudStackConfiguration(),
			"cloudstack_counter":                        resourceCloudStackCounter(),
			"cloudstack_cluster":                        resourceCloudStackCluster(),
			"cloudstack_disk":                           resourceCloudStackDisk(),
			"cloudstack_egress_firewall":                resourceCloudStackEgressFirewall(),
			"cloudstack_firewall":                       resourceCloudStackFirewall(),
			"cloudstack_host":                           resourceCloudStackHost(),
			"cloudstack_instance":                       resourceCloudStackInstance(),
			"cloudstack_ipaddress":                      resourceCloudStackIPAddress(),
			"cloudstack_kubernetes_cluster":             resourceCloudStackKubernetesCluster(),
			"cloudstack_kubernetes_version":             resourceCloudStackKubernetesVersion(),
			"cloudstack_loadbalancer_rule":              resourceCloudStackLoadBalancerRule(),
			"cloudstack_network":                        resourceCloudStackNetwork(),
			"cloudstack_network_acl":                    resourceCloudStackNetworkACL(),
			"cloudstack_network_acl_rule":               resourceCloudStackNetworkACLRule(),
			"cloudstack_nic":                            resourceCloudStackNIC(),
			"cloudstack_physical_network":               resourceCloudStackPhysicalNetwork(),
			"cloudstack_pod":                            resourceCloudStackPod(),
			"cloudstack_port_forward":                   resourceCloudStackPortForward(),
			"cloudstack_network_service_provider_state": resourceCloudStackNetworkServiceProviderState(),
			"cloudstack_private_gateway":                resourceCloudStackPrivateGateway(),
			"cloudstack_secondary_ipaddress":            resourceCloudStackSecondaryIPAddress(),
			"cloudstack_secondary_storage":              resourceCloudStackSecondaryStorage(),
			"cloudstack_security_group":                 resourceCloudStackSecurityGroup(),
			"cloudstack_security_group_rule":            resourceCloudStackSecurityGroupRule(),
			"cloudstack_ssh_keypair":                    resourceCloudStackSSHKeyPair(),
			"cloudstack_static_nat":                     resourceCloudStackStaticNAT(),
			"cloudstack_static_route":                   resourceCloudStackStaticRoute(),
			"cloudstack_storage_pool":                   resourceCloudStackStoragePool(),
			"cloudstack_template":                       resourceCloudStackTemplate(),
			"cloudstack_traffic_type":                   resourceCloudStackTrafficType(),
			"cloudstack_vpc":                            resourceCloudStackVPC(),
			"cloudstack_vpn_connection":                 resourceCloudStackVPNConnection(),
			"cloudstack_vpn_customer_gateway":           resourceCloudStackVPNCustomerGateway(),
			"cloudstack_vpn_gateway":                    resourceCloudStackVPNGateway(),
			"cloudstack_network_offering":               resourceCloudStackNetworkOffering(),
			"cloudstack_disk_offering":                  resourceCloudStackDiskOffering(),
			"cloudstack_vlan_ip_range":                  resourceCloudstackVlanIpRange(),
			"cloudstack_volume":                         resourceCloudStackVolume(),
			"cloudstack_zone":                           resourceCloudStackZone(),
			"cloudstack_service_offering":               resourceCloudStackServiceOffering(),
			"cloudstack_account":                        resourceCloudStackAccount(),
			"cloudstack_project":                        resourceCloudStackProject(),
			"cloudstack_user":                           resourceCloudStackUser(),
			"cloudstack_user_data":                      resourceCloudStackUserData(),
			"cloudstack_user_data_template_link":        resourceCloudStackUserDataTemplateLink(),
			"cloudstack_domain":                         resourceCloudStackDomain(),
			"cloudstack_network_service_provider":       resourceCloudStackNetworkServiceProvider(),
			"cloudstack_role":                           resourceCloudStackRole(),
			"cloudstack_limits":                         resourceCloudStackLimits(),
			"cloudstack_snapshot_policy":                resourceCloudStackSnapshotPolicy(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (any, error) {
	apiURL, apiURLOK := d.GetOk("api_url")
	apiKey, apiKeyOK := d.GetOk("api_key")
	secretKey, secretKeyOK := d.GetOk("secret_key")
	config, configOK := d.GetOk("config")
	profile, profileOK := d.GetOk("profile")

	switch {
	case apiURLOK, apiKeyOK, secretKeyOK:
		if !(apiURLOK && apiKeyOK && secretKeyOK) {
			return nil, errors.New("'api_url', 'api_key' and 'secret_key' should all have values")
		}
	case configOK, profileOK:
		if !(configOK && profileOK) {
			return nil, errors.New("'config' and 'profile' should both have a value")
		}
	default:
		return nil, errors.New(
			"either 'api_url', 'api_key' and 'secret_key' or 'config' and 'profile' should have values")
	}

	if configOK && profileOK {
		cfg, err := ini.Load(config.(string))
		if err != nil {
			return nil, err
		}

		section, err := cfg.GetSection(profile.(string))
		if err != nil {
			return nil, err
		}

		apiURL = section.Key("url").String()
		apiKey = section.Key("apikey").String()
		secretKey = section.Key("secretkey").String()
	}

	cfg := Config{
		APIURL:      apiURL.(string),
		APIKey:      apiKey.(string),
		SecretKey:   secretKey.(string),
		HTTPGETOnly: d.Get("http_get_only").(bool),
		Timeout:     int64(d.Get("timeout").(int)),
	}

	return cfg.NewClient()
}
