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

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCloudStackSecondaryStorage() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackSecondaryStorageCreate,
		Read:   resourceCloudStackSecondaryStorageRead,
		Delete: resourceCloudStackSecondaryStorageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"storage_provider": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"url": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCloudStackSecondaryStorageCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.ImageStore.NewAddImageStoreParams(d.Get("storage_provider").(string))
	if v, ok := d.GetOk("name"); ok {
		p.SetName(v.(string))
	}
	if v, ok := d.GetOk("url"); ok {
		p.SetUrl(v.(string))
	}
	if v, ok := d.GetOk("zone_id"); ok {
		p.SetZoneid(v.(string))
	}

	r, err := cs.ImageStore.AddImageStore(p)
	if err != nil {
		return err
	}

	d.SetId(r.Id)

	return resourceCloudStackSecondaryStorageRead(d, meta)
}

func resourceCloudStackSecondaryStorageRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	r, _, err := cs.ImageStore.GetImageStoreByID(d.Id())
	if err != nil {
		return err
	}

	d.Set("name", r.Name)
	d.Set("storage_provider", r.Providername)
	d.Set("url", r.Url)
	d.Set("zone_id", r.Zoneid)

	return nil
}

func resourceCloudStackSecondaryStorageDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	_, err := cs.ImageStore.DeleteImageStore(cs.ImageStore.NewDeleteImageStoreParams(d.Id()))
	if err != nil {
		return fmt.Errorf("Error deleting secondary storage: %s", err)
	}

	return nil
}
