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
	"github.com/hashicorp/terraform/helper/schema"
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

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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

	// If there is a project supplied, we retrieve and set the project id
	if err := setProjectid(p, cs, d); err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating Kubernetes Cluster %s", name)
	r, err := cs.Kubernetes.CreateKubernetesCluster(p)
	if err != nil {
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
	d.Partial(true)

	if d.HasChange("service_offering") || d.HasChange("size") {
		p := cs.Kubernetes.NewScaleKubernetesClusterParams(d.Id())
		serviceOfferingID, e := retrieveID(cs, "service_offering", d.Get("service_offering").(string))
		if e != nil {
			return e.Error()
		}
		p.SetServiceofferingid(serviceOfferingID)
		p.SetSize(int64(d.Get("size").(int)))
		_, err := cs.Kubernetes.ScaleKubernetesCluster(p)
		if err != nil {
			return fmt.Errorf(
				"Error Scaling Kubernetes Cluster %s: %s", d.Id(), err)
		}
		d.SetPartial("service_offering")
		d.SetPartial("size")
	}

	if d.HasChange("autoscaling_enabled") || d.HasChange("min_size") || d.HasChange("max_size") {
		err := autoscaleKubernetesCluster(d, meta)
		if err != nil {
			return err
		}
		d.SetPartial("autoscaling_enabled")
		d.SetPartial("min_size")
		d.SetPartial("max_size")
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
		d.SetPartial("kubernetes_version")
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
		d.SetPartial("state")
	}

	d.Partial(false)
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
