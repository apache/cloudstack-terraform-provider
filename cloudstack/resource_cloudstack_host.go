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
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCloudStackHost() *schema.Resource {
	return &schema.Resource{
		Read:   resourceCloudStackHostRead,
		Update: resourceCloudStackHostUpdate,
		Create: resourceCloudStackHostCreate,
		Delete: resourceCloudStackHostDelete,
		Schema: map[string]*schema.Schema{
			"hypervisor": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					validHypervisors := []string{"xenserver", "kvm", "vmware", "baremetal", "simulator"}

					sort.Strings(validHypervisors)

					if sort.SearchStrings(validHypervisors, v.(string)) >= len(validHypervisors) {
						errors = append(errors, fmt.Errorf("%q must be one of %v", k, validHypervisors))
					}
					return
				},
				ForceNew: true,
			},
			"pod_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"url": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ConflictsWith: []string{
					"cluster_name",
				},
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ConflictsWith: []string{
					"cluster_id",
				},
			},
			"host_tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"prevent_destroy": {
				Type:        schema.TypeBool,
				Description: "Prevent the host from being destroyed. This is useful when you want to avoid destroy the host in any change.",
				Optional:    true,
				Default:     false,
			},
			"force_destroy": {
				Type:        schema.TypeBool,
				Description: "Force the host to be destroyed.",
				Optional:    true,
				Default:     false,
			},
			"allocation_state": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Enabled",
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// Creating a host can fail if the instance is still being created and Cloudsack
			// user is created during cloud-init. This timeout is used to wait for the host
			// to be created and Cloudstack user to be available.
			"create_timeout": {
				Type:        schema.TypeInt,
				Description: "Timeout in seconds to wait for the host to be created.",
				Optional:    true,
				Default:     300,
			},
			// Destroying a host will put it in Maintenance mode first. If the VMs are still
			// being migrated, the host will be in state PrepareForMaintenance. This timeout
			// is used to wait for the host to be in Maintenance state.
			"destroy_timeout": {
				Type:        schema.TypeInt,
				Description: "Timeout in seconds to wait for the host to be destroyed.",
				Optional:    true,
				Default:     300,
			},
		},
	}
}

func resourceCloudStackHostCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	hypervisor := d.Get("hypervisor").(string)
	pod_id := d.Get("pod_id").(string)
	url := d.Get("url").(string)
	zone_id := d.Get("zone_id").(string)

	p := cs.Host.NewAddHostParams(hypervisor, pod_id, url, zone_id)

	if cluster_id, ok := d.GetOk("cluster_id"); ok {
		p.SetClusterid(cluster_id.(string))
	}

	if cluster_name, ok := d.GetOk("cluster_name"); ok {
		p.SetClustername(cluster_name.(string))
	}

	if host_tags, ok := d.GetOk("host_tags"); ok {
		p.SetHosttags(host_tags.([]string))
	}

	if username, ok := d.GetOk("username"); ok {
		p.SetUsername(username.(string))
	}

	if password, ok := d.GetOk("password"); ok {
		p.SetPassword(password.(string))
	}

	timeout := time.After(time.Duration(d.Get("create_timeout").(int)) * time.Second)
	tick := time.NewTicker(5 * time.Second)
	var err error
	var host *cloudstack.AddHostResponse

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout waiting for Host to be created, with error: %s", err)
		case <-tick.C:
			log.Printf("[DEBUG] Trying to create host %s", d.Get("url").(string))
			host, err = cs.Host.AddHost(p)
			if err != nil {
				log.Printf("[ERROR] Error creating host %s: %s. Will try again...", d.Get("url").(string), err)
				continue
			}

			if host.Id != "" {
				log.Printf("[DEBUG] Host %s successfully created", url)
				d.SetId(host.Id)
				return resourceCloudStackHostRead(d, meta)
			}
		}
	}
}

func resourceCloudStackHostRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	log.Printf("[DEBUG] Retrieving Host %s", d.Get("url").(string))

	h, count, err := cs.Host.GetHostByID(d.Id())

	if err != nil {
		if count == 0 {
			log.Printf("[WARN] Host %s does no longer exist", d.Get("url").(string))
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(h.Id)

	fields := map[string]interface{}{
		"hypervisor":     h.Hypervisor,
		"pod_id":         h.Podid,
		"zone_id":        h.Zoneid,
		"state":          h.State,
		"resource_state": h.Resourcestate,
		"name":           h.Name,
	}

	for k, v := range fields {
		if err := d.Set(k, v); err != nil {
			return err
		}
	}

	if cluster_id := d.Get("cluster_id"); cluster_id != "" {
		d.Set("cluster_id", h.Clusterid)
	} else {
		d.Set("cluster_name", h.Clustername)
	}

	if h.Hosttags != "" {
		d.Set("host_tags", strings.Split(h.Hosttags, ","))
	}

	return nil
}

func resourceCloudStackHostUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Updating Host: %s", d.Id())

	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.Host.NewUpdateHostParams(d.Id())

	if d.HasChange("allocation_state") {
		log.Printf("[DEBUG] Updating Host allocation state: %s", d.Id())
		p.SetAllocationstate(d.Get("allocation_state").(string))
	}

	if d.HasChange("host_tags") {
		log.Printf("[DEBUG] Updating Host tags: %s", d.Id())
		p.SetHosttags(d.Get("host_tags").([]string))
	}

	return resourceCloudStackHostRead(d, meta)
}

func resourceCloudStackHostDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	if d.Get("prevent_destroy").(bool) {
		log.Printf("[INFO] Skipping Host deletion: %s", d.Id())
		return fmt.Errorf("host %s is marked to be protected from deletion", d.Id())
	}

	log.Printf("[INFO] Removing Host: %s", d.Id())
	mm := cs.Host.NewPrepareHostForMaintenanceParams(d.Id())
	_, err := cs.Host.PrepareHostForMaintenance(mm)

	if err != nil {
		return fmt.Errorf("error preparing Host for maintenance: %s", err)
	}

	timeout := time.After(time.Duration(d.Get("destroy_timeout").(int)) * time.Second)
	tick := time.NewTicker(3 * time.Second)

	for {
		select {
		case <-timeout:
			return errors.New("timeout waiting for Host to enter Maintenance state")
		case <-tick.C:
			log.Printf("[DEBUG] Checking Host state: %s", d.Id())
			err = resourceCloudStackHostRead(d, meta)
			if err != nil {
				return fmt.Errorf("error reading Host: %s", err)
			}

			if d.Get("resource_state").(string) == "Maintenance" || d.Get("resource_state").(string) == "Disconnected" {
				log.Printf("[INFO] Deleting Host: %s", d.Id())
				h := cs.Host.NewDeleteHostParams(d.Id())
				_, err = cs.Host.DeleteHost(h)

				if err != nil {
					return fmt.Errorf("error deleting Host: %s", err)
				}
				return nil
			}
		}
	}
}
