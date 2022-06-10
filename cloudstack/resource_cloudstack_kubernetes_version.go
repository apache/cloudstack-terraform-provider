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

func resourceCloudStackKubernetesVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackKubernetesVersionCreate,
		Read:   resourceCloudStackKubernetesVersionRead,
		Update: resourceCloudStackKubernetesVersionUpdate,
		Delete: resourceCloudStackKubernetesVersionDelete,
		Importer: &schema.ResourceImporter{
			State: importStatePassthrough,
		},

		Schema: map[string]*schema.Schema{

			"semantic_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"url": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"min_cpu": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"min_memory": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			// Optional Params
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"checksum": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceCloudStackKubernetesVersionCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// State is always Enabled when created
	if state, ok := d.GetOk("state"); ok {
		if state.(string) != "Enabled" {
			return fmt.Errorf("State must be 'Enabled' when first adding an ISO")
		}
	}

	semanticVersion := d.Get("semantic_version").(string)
	url := d.Get("url").(string)
	minCpu := d.Get("min_cpu").(int)
	minMemory := d.Get("min_memory").(int)

	p := cs.Kubernetes.NewAddKubernetesSupportedVersionParams(minCpu, minMemory, semanticVersion)
	p.SetUrl(url)

	if name, ok := d.GetOk("name"); ok {
		p.SetName(name.(string))
	}
	if checksum, ok := d.GetOk("checksum"); ok {
		p.SetName(checksum.(string))
	}
	if zone, ok := d.GetOk("zone"); ok {
		zoneID, e := retrieveID(cs, "zone", zone.(string))
		if e != nil {
			return e.Error()
		}
		p.SetZoneid(zoneID)
	}

	log.Printf("[DEBUG] Creating Kubernetes Version %s", semanticVersion)
	r, err := cs.Kubernetes.AddKubernetesSupportedVersion(p)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Kubernetes Version %s successfully created", semanticVersion)
	d.SetId(r.Id)
	return resourceCloudStackKubernetesVersionRead(d, meta)
}

func resourceCloudStackKubernetesVersionRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	log.Printf("[DEBUG] Retrieving Kubernetes Version %s", d.Get("semantic_version").(string))

	// Get the Kubernetes Version details
	version, count, err := cs.Kubernetes.GetKubernetesSupportedVersionByID(
		d.Id(),
	)
	if err != nil {
		if count == 0 {
			log.Printf("[DEBUG] Kubernetes Version %s does not longer exist", d.Get("semantic_version").(string))
			d.SetId("")
			return nil
		}

		return err
	}

	// Update the config
	d.SetId(version.Id)
	d.Set("semantic_version", version.Semanticversion)
	d.Set("name", version.Name)
	d.Set("min_cpu", version.Mincpunumber)
	d.Set("min_memory", version.Minmemory)
	d.Set("state", version.State)

	setValueOrID(d, "zone", version.Zonename, version.Zoneid)
	return nil
}

func resourceCloudStackKubernetesVersionUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	d.Partial(true)

	if d.HasChange("state") {
		p := cs.Kubernetes.NewUpdateKubernetesSupportedVersionParams(d.Id(), d.Get("state").(string))
		_, err := cs.Kubernetes.UpdateKubernetesSupportedVersion(p)
		if err != nil {
			return fmt.Errorf(
				"Error Updating Kubernetes Version %s: %s", d.Id(), err)
		}
		d.SetPartial("state")
	}

	d.Partial(false)
	return resourceCloudStackKubernetesVersionRead(d, meta)
}

func resourceCloudStackKubernetesVersionDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.Kubernetes.NewDeleteKubernetesSupportedVersionParams(d.Id())

	// Delete the Kubernetes Version
	_, err := cs.Kubernetes.DeleteKubernetesSupportedVersion(p)
	if err != nil {
		// This is a very poor way to be told the ID does no longer exist :(
		if strings.Contains(err.Error(), fmt.Sprintf(
			"Invalid parameter id value=%s due to incorrect long value format, "+
				"or entity does not exist", d.Id())) {
			return nil
		}

		return fmt.Errorf("Error deleting Kubernetes Version: %s", err)
	}

	return nil
}
