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
	"fmt"
	"log"
	"strings"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudStackKubernetesCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackKubernetesClusterCreate,
		Read:   resourceCloudStackKubernetesClusterRead,
		Update: resourceCloudStackKubernetesClusterUpdate,
		Delete: resourceCloudStackKubernetesClusterDelete,
		Importer: &schema.ResourceImporter{
			State: importStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"kubernetes_version": {
				Type:     schema.TypeString,
				Required: true,
			},

			"service_offering": {
				Type:     schema.TypeString,
				Required: true,
			},

			// Begin optional params
			"size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},

			"autoscaling_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"min_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"max_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"control_nodes_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true, // For now
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"keypair": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"network_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				// Default:  "Running",
			},

			"account": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "An optional account for the Kubernetes cluster. Must be used with domain_id.",
			},

			"domain_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "An optional domain ID for the Kubernetes cluster. If the account parameter is used, domain_id must also be used. Hosts dedicated to the specified domain will be used for deploying the cluster",
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"noderootdisksize": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  8,
			},

			"docker_registry_url": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"docker_registry_username": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"docker_registry_password": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
			},

			"as_number": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "An optional as number for the Kubernetes cluster",
			},

			"cni_config_details": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "An optional map of CNI configuration details. It is used to specify the parameters values for the variables in userdata",
			},

			"cni_configuration_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "An optional CNI configuration ID for the Kubernetes cluster. If not specified, the default CNI configuration will be used",
			},

			"etcd_nodes_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true, // For now
				Description: "Number of etcd nodes in the Kubernetes cluster. Default is 0",
			},

			"hypervisor": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The hypervisor on which to deploy the cluster.",
			},

			"node_offerings": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "An optional map of node roles to service offerings. If not specified, the service_offering parameter will be used for all node roles. Valid roles are: worker, control, etcd",
			},

			"node_templates": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "An optional map of node roles to instance templates. If not specified, system VM template will be used. Valid roles are: worker, control, etcd",
			},
		},
	}
}

func resourceCloudStackKubernetesClusterCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// State is always Running when created
	if state, ok := d.GetOk("state"); ok {
		if state.(string) != "Running" {
			return fmt.Errorf("State must be 'Running' when first creating a cluster")
		}
	}

	name := d.Get("name").(string)
	size := int64(d.Get("size").(int))
	serviceOfferingID, e := retrieveID(cs, "service_offering", d.Get("service_offering").(string))
	if e != nil {
		return e.Error()
	}
	zoneID, e := retrieveID(cs, "zone", d.Get("zone").(string))
	if e != nil {
		return e.Error()
	}
	kubernetesVersionID, e := retrieveID(cs, "kubernetes_version", d.Get("kubernetes_version").(string))
	if e != nil {
		return e.Error()
	}

	// Create a new parameter struct
	p := cs.Kubernetes.NewCreateKubernetesClusterParams(name, kubernetesVersionID, name, serviceOfferingID, size, zoneID)

	// Set optional params
	if description, ok := d.GetOk("description"); ok {
		p.SetDescription(description.(string))
	}
	if keypair, ok := d.GetOk("keypair"); ok {
		p.SetKeypair(keypair.(string))
	}
	if networkID, ok := d.GetOk("network_id"); ok {
		p.SetNetworkid(networkID.(string))
	}
	if controlNodesSize, ok := d.GetOk("control_nodes_size"); ok {
		p.SetControlnodes(int64(controlNodesSize.(int)))
	}
	if noderootdisksize, ok := d.GetOk("noderootdisksize"); ok {
		p.SetNoderootdisksize(int64(noderootdisksize.(int)))
	}
	if dockerurl, ok := d.GetOk("docker_registry_url"); ok {
		p.SetDockerregistryurl(dockerurl.(string))
	}
	if dockerusername, ok := d.GetOk("docker_registry_username"); ok {
		p.SetDockerregistryusername(dockerusername.(string))
	}
	if dockerpassword, ok := d.GetOk("docker_registry_password"); ok {
		p.SetDockerregistrypassword(dockerpassword.(string))
	}

	// If there is a project supplied, we retrieve and set the project id
	if err := setProjectid(p, cs, d); err != nil {
		return err
	}

	if account, ok := d.GetOk("account"); ok {
		p.SetAccount(account.(string))
	}
	if domainID, ok := d.GetOk("domain_id"); ok {
		p.SetDomainid(domainID.(string))
	}

	if asNumber, ok := d.GetOk("as_number"); ok {
		p.SetAsnumber(int64(asNumber.(int)))
	}

	if etcdNodesSize, ok := d.GetOk("etcd_nodes_size"); ok {
		p.SetEtcdnodes(int64(etcdNodesSize.(int)))
	}

	if hypervisor, ok := d.GetOk("hypervisor"); ok {
		p.SetHypervisor(hypervisor.(string))
	}

	if cniConfigID, ok := d.GetOk("cni_configuration_id"); ok {
		p.SetCniconfigurationid(cniConfigID.(string))
	}

	if nodeOfferings, ok := d.GetOk("node_offerings"); ok {
		nodeOfferingsMap := nodeOfferings.(map[string]interface{})
		nodeOfferingsFormatted := make(map[string]string)
		for nodeType, offeringName := range nodeOfferingsMap {
			// Retrieve the offering ID
			offeringID, e := retrieveID(cs, "service_offering", offeringName.(string))
			if e != nil {
				return e.Error()
			}
			nodeOfferingsFormatted[nodeType] = offeringID
		}
		p.SetNodeofferings(nodeOfferingsFormatted)
	}

	if nodeTemplates, ok := d.GetOk("node_templates"); ok {
		nodeTemplatesMap := nodeTemplates.(map[string]interface{})
		nodeTemplatesFormatted := make(map[string]string)
		for nodeType, templateName := range nodeTemplatesMap {
			zoneID, err := retrieveID(cs, "zone", d.Get("zone").(string))
			if err != nil {
				return err.Error()
			}
			templateID, e := retrieveTemplateID(cs, zoneID, templateName.(string))
			if e != nil {
				return e.Error()
			}
			nodeTemplatesFormatted[nodeType] = templateID
		}
		p.SetNodetemplates(nodeTemplatesFormatted)
	}

	if cniConfigDetails, ok := d.GetOk("cni_config_details"); ok {
		cniConfigDetailsMap := cniConfigDetails.(map[string]interface{})
		cniConfigDetailsFormatted := make(map[string]string)
		for key, value := range cniConfigDetailsMap {
			cniConfigDetailsFormatted[key] = value.(string)
		}
		p.SetCniconfigdetails(cniConfigDetailsFormatted)
	}

	log.Printf("[DEBUG] Creating Kubernetes Cluster %s", name)
	r, err := cs.Kubernetes.CreateKubernetesCluster(p)
	if err != nil {
		cluster, _, errg := cs.Kubernetes.GetKubernetesClusterByName(
			name,
			cloudstack.WithProject(d.Get("project").(string)),
		)
		if errg == nil {
			d.SetId(cluster.Id)
		}
		return err
	}

	log.Printf("[DEBUG] Kubernetes Cluster %s successfully created", name)
	d.SetId(r.Id)

	if _, ok := d.GetOk("autoscaling_enabled"); ok {
		err = autoscaleKubernetesCluster(d, meta)
		if err != nil {
			return err
		}
	}

	return resourceCloudStackKubernetesClusterRead(d, meta)
}

func resourceCloudStackKubernetesClusterRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	log.Printf("[DEBUG] Retrieving Kubernetes Cluster %s", d.Get("name").(string))

	// Get the Kubernetes Cluster details
	cluster, count, err := cs.Kubernetes.GetKubernetesClusterByID(
		d.Id(),
		cloudstack.WithProject(d.Get("project").(string)),
	)
	if err != nil {
		if count == 0 {
			log.Printf("[DEBUG] Kubernetes Cluster %s does not longer exist", d.Get("name").(string))
			d.SetId("")
			return nil
		}

		return err
	}

	// Update the config
	d.SetId(cluster.Id)
	d.Set("name", cluster.Name)
	d.Set("description", cluster.Description)
	d.Set("control_nodes_size", cluster.Controlnodes)
	d.Set("size", cluster.Size)
	d.Set("autoscaling_enabled", cluster.Autoscalingenabled)
	d.Set("min_size", cluster.Minsize)
	d.Set("max_size", cluster.Maxsize)
	d.Set("keypair", cluster.Keypair)
	d.Set("network_id", cluster.Networkid)
	d.Set("ip_address", cluster.Ipaddress)
	d.Set("state", cluster.State)
	d.Set("account", cluster.Account)
	d.Set("domain_id", cluster.Domainid)
	d.Set("etcd_nodes_size", cluster.Etcdnodes)
	d.Set("cni_configuration_id", cluster.Cniconfigurationid)

	setValueOrID(d, "kubernetes_version", cluster.Kubernetesversionname, cluster.Kubernetesversionid)
	setValueOrID(d, "service_offering", cluster.Serviceofferingname, cluster.Serviceofferingid)
	setValueOrID(d, "project", cluster.Project, cluster.Projectid)
	setValueOrID(d, "zone", cluster.Zonename, cluster.Zoneid)

	return nil
}

func autoscaleKubernetesCluster(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	p := cs.Kubernetes.NewScaleKubernetesClusterParams(d.Id())
	p.SetAutoscalingenabled(d.Get("autoscaling_enabled").(bool))
	p.SetMinsize(int64(d.Get("min_size").(int)))
	p.SetMaxsize(int64(d.Get("max_size").(int)))
	_, err := cs.Kubernetes.ScaleKubernetesCluster(p)
	return err
}

func resourceCloudStackKubernetesClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	if d.HasChange("service_offering") || d.HasChange("size") || d.HasChange("node_offerings") {
		p := cs.Kubernetes.NewScaleKubernetesClusterParams(d.Id())
		serviceOfferingID, e := retrieveID(cs, "service_offering", d.Get("service_offering").(string))
		if e != nil {
			return e.Error()
		}
		p.SetServiceofferingid(serviceOfferingID)
		p.SetSize(int64(d.Get("size").(int)))

		// Handle node offerings if they changed
		if nodeOfferings, ok := d.GetOk("node_offerings"); ok {
			nodeOfferingsMap := nodeOfferings.(map[string]interface{})
			nodeOfferingsFormatted := make(map[string]string)
			for nodeType, offeringName := range nodeOfferingsMap {
				// Retrieve the offering ID
				offeringID, e := retrieveID(cs, "service_offering", offeringName.(string))
				if e != nil {
					return e.Error()
				}
				nodeOfferingsFormatted[nodeType] = offeringID
			}
			p.SetNodeofferings(nodeOfferingsFormatted)
		}

		_, err := cs.Kubernetes.ScaleKubernetesCluster(p)
		if err != nil {
			return fmt.Errorf(
				"Error Scaling Kubernetes Cluster %s: %s", d.Id(), err)
		}
	}

	if d.HasChange("autoscaling_enabled") || d.HasChange("min_size") || d.HasChange("max_size") {
		err := autoscaleKubernetesCluster(d, meta)
		if err != nil {
			return err
		}
	}

	if d.HasChange("kubernetes_version") {
		kubernetesVersionID, e := retrieveID(cs, "kubernetes_version", d.Get("kubernetes_version").(string))
		if e != nil {
			return e.Error()
		}
		p := cs.Kubernetes.NewUpgradeKubernetesClusterParams(d.Id(), kubernetesVersionID)
		_, err := cs.Kubernetes.UpgradeKubernetesCluster(p)
		if err != nil {
			return fmt.Errorf(
				"Error Upgrading Kubernetes Cluster %s: %s", d.Id(), err)
		}
	}

	if d.HasChange("state") {
		state := d.Get("state").(string)
		switch state {
		case "Running":
			p := cs.Kubernetes.NewStartKubernetesClusterParams(d.Id())
			_, err := cs.Kubernetes.StartKubernetesCluster(p)
			if err != nil {
				return fmt.Errorf(
					"Error Starting Kubernetes Cluster %s: %s", d.Id(), err)
			}
		case "Stopped":
			p := cs.Kubernetes.NewStopKubernetesClusterParams(d.Id())
			_, err := cs.Kubernetes.StopKubernetesCluster(p)
			if err != nil {
				return fmt.Errorf(
					"Error Stopping Kubernetes Cluster %s: %s", d.Id(), err)
			}
		default:
			return fmt.Errorf("State must either be 'Running' or 'Stopped'")
		}
	}

	return resourceCloudStackKubernetesClusterRead(d, meta)
}

func resourceCloudStackKubernetesClusterDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.Kubernetes.NewDeleteKubernetesClusterParams(d.Id())

	// Delete the Kubernetes Cluster
	_, err := cs.Kubernetes.DeleteKubernetesCluster(p)
	if err != nil {
		// This is a very poor way to be told the ID does no longer exist :(
		if strings.Contains(err.Error(), fmt.Sprintf(
			"Invalid parameter id value=%s due to incorrect long value format, "+
				"or entity does not exist", d.Id())) {
			return nil
		}

		return fmt.Errorf("Error deleting Kubernetes Cluster: %s", err)
	}

	return nil
}
