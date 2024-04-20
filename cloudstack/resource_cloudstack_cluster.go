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

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCloudStackCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackClusterCreate,
		Read:   resourceCloudStackClusterRead,
		Update: resourceCloudStackClusterUpdate,
		Delete: resourceCloudStackClusterDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					validOptions := []string{
						"CloudManaged",
						"ExternalManaged",
					}
					err := validateOptions(validOptions, v.(string), k)
					if err != nil {
						errors = append(errors, err)
					}

					return
				},
			},
			"hypervisor": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					validOptions := []string{
						"XenServer",
						"KVM",
						"VMware",
						"Hyperv",
						"BareMetal",
						"Simulator",
						"Ovm3",
					}
					err := validateOptions(validOptions, v.(string), k)
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},
			"pod_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"allocation_state": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Enabled",
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					validOptions := []string{
						"Enabled",
						"Disabled",
					}
					err := validateOptions(validOptions, v.(string), k)
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},
			"zone_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pod_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCloudStackClusterCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	clusterName := d.Get("name").(string)
	clusterType := d.Get("type").(string)
	hypervisor := d.Get("hypervisor").(string)
	podID := d.Get("pod_id").(string)
	zoneID := d.Get("zone_id").(string)

	// Create a new parameter struct
	p := cs.Cluster.NewAddClusterParams(
		clusterName,
		clusterType,
		hypervisor,
		podID,
		zoneID,
	)

	if allocationState, ok := d.GetOk("allocation_state"); ok {
		p.SetAllocationstate(allocationState.(string))
	}

	log.Printf("[DEBUG] Creating Cluster %s", clusterName)

	c, err := cs.Cluster.AddCluster(p)

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Cluster %s successfully created", clusterName)

	d.SetId(c.Id)

	return resourceCloudStackClusterRead(d, meta)
}

func resourceCloudStackClusterRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	log.Printf("[DEBUG] Retrieving Cluster %s", d.Get("name").(string))

	// Get the Cluster details
	c, count, err := cs.Cluster.GetClusterByName(d.Get("name").(string))

	if err != nil {
		if count == 0 {
			log.Printf("[DEBUG] Cluster %s does no longer exist", d.Get("name").(string))
			d.SetId("")
			return nil
		}
		return err
	}
	log.Printf("[DEBUG] Cluster %+v ", c)
	d.SetId(c.Id)
	d.Set("name", c.Name)
	d.Set("type", c.Clustertype)
	d.Set("hypervisor", c.Hypervisortype)
	d.Set("pod_id", c.Podid)
	d.Set("zone_id", c.Zoneid)
	d.Set("allocation_state", c.Allocationstate)
	d.Set("pod_name", c.Podname)
	d.Set("zone_name", c.Zonename)

	return nil
}

func resourceCloudStackClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	name := d.Get("name").(string)

	if d.HasChange("allocation_state") {
		log.Printf("[DEBUG] allocationState changed for cluster %s, starting update", name)

		p := cs.Cluster.NewUpdateClusterParams(d.Id())

		p.SetAllocationstate(d.Get("allocation_state").(string))

		_, err := cs.Cluster.UpdateCluster(p)
		if err != nil {
			return fmt.Errorf(
				"Error updating the allocation_state for cluster %s: %s", name, err)
		}
	}

	if d.HasChange("name") {
		log.Printf("[DEBUG] name changed for cluster %s, starting update", name)

		p := cs.Cluster.NewUpdateClusterParams(d.Id())

		p.SetClustername(d.Get("name").(string))

		_, err := cs.Cluster.UpdateCluster(p)
		if err != nil {
			return fmt.Errorf(
				"Error updating the name for cluster %s: %s", name, err)
		}
	}

	if d.HasChange("type") {
		log.Printf("[DEBUG] type changed for cluster %s, starting update", name)

		p := cs.Cluster.NewUpdateClusterParams(d.Id())

		p.SetClustertype(d.Get("type").(string))

		_, err := cs.Cluster.UpdateCluster(p)
		if err != nil {
			return fmt.Errorf(
				"Error updating the type for cluster %s: %s", name, err)
		}
	}

	if d.HasChange("hypervisor") {
		log.Printf("[DEBUG] hypervisor changed for cluster %s, starting update", name)

		p := cs.Cluster.NewUpdateClusterParams(d.Id())

		p.SetHypervisor(d.Get("hypervisor").(string))

		_, err := cs.Cluster.UpdateCluster(p)
		if err != nil {
			return fmt.Errorf(
				"Error updating the hypervisor for cluster %s: %s", name, err)
		}
	}
	return resourceCloudStackClusterRead(d, meta)
}

func resourceCloudStackClusterDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.Cluster.NewDeleteClusterParams(d.Id())
	_, err := cs.Cluster.DeleteCluster(p)

	if err != nil {
		return fmt.Errorf("Error deleting Pod: %s", err)
	}

	return nil
}
