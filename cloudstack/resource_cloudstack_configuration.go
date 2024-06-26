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

func resourceCloudStackConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackConfigurationCreate,
		Read:   resourceCloudStackConfigurationRead,
		Update: resourceCloudStackConfigurationUpdate,
		Delete: resourceCloudStackConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: importStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"accountid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"clusterid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domainid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"imagestoreuuid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"storeid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"zoneid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// computed
			"category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_dynamic": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scope": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCloudStackConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	p := cs.Configuration.NewListConfigurationsParams()

	// required
	p.SetName(d.Id())

	// optional
	if v, ok := d.GetOk("accountid"); ok {
		p.SetAccountid(v.(string))
	}
	if v, ok := d.GetOk("category"); ok {
		p.SetCategory(v.(string))
	}
	if v, ok := d.GetOk("clusterid"); ok {
		p.SetClusterid(v.(string))
	}
	if v, ok := d.GetOk("domainid"); ok {
		p.SetDomainid(v.(string))
	}
	if v, ok := d.GetOk("imagestoreuuid"); ok {
		p.SetImagestoreuuid(v.(string))
	}
	if v, ok := d.GetOk("storeid"); ok {
		p.SetStorageid(v.(string))
	}
	if v, ok := d.GetOk("zoneid"); ok {
		p.SetZoneid(v.(string))
	}

	cfg, err := cs.Configuration.ListConfigurations(p)
	if err != nil {
		return err
	}

	found := false
	for _, v := range cfg.Configurations {
		if v.Name == d.Id() {
			d.Set("category", v.Category)
			d.Set("description", v.Description)
			d.Set("is_dynamic", v.Isdynamic)
			d.Set("name", v.Name)
			d.Set("value", v.Value)
			d.Set("scope", v.Scope)
			found = true
		}
	}

	if !found {
		return fmt.Errorf("listConfiguration failed. no matching names found %s", d.Id())
	}

	return nil

}

func resourceCloudStackConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	if v, ok := d.GetOk("name"); ok {
		d.SetId(v.(string))
	}

	resourceCloudStackConfigurationUpdate(d, meta)

	return nil

}

func resourceCloudStackConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	p := cs.Configuration.NewUpdateConfigurationParams(d.Id())

	// Optional
	if v, ok := d.GetOk("accountid"); ok {
		p.SetAccountid(v.(string))
	}
	if v, ok := d.GetOk("clusterid"); ok {
		p.SetClusterid(v.(string))
	}
	if v, ok := d.GetOk("domainid"); ok {
		p.SetDomainid(v.(string))
	}
	if v, ok := d.GetOk("imagestoreuuid"); ok {
		p.SetImagestoreuuid(v.(string))
	}
	if v, ok := d.GetOk("storeid"); ok {
		p.SetStorageid(v.(string))
	}
	if v, ok := d.GetOk("value"); ok {
		p.SetValue(v.(string))
	}
	if v, ok := d.GetOk("zoneid"); ok {
		p.SetZoneid(v.(string))
	}

	_, err := cs.Configuration.UpdateConfiguration(p)
	if err != nil {
		return err
	}

	resourceCloudStackConfigurationRead(d, meta)

	return nil
}

func resourceCloudStackConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	p := cs.Configuration.NewResetConfigurationParams(d.Id())

	// Optional
	if v, ok := d.GetOk("accountid"); ok {
		p.SetAccountid(v.(string))
	}
	if v, ok := d.GetOk("clusterid"); ok {
		p.SetClusterid(v.(string))
	}
	if v, ok := d.GetOk("domainid"); ok {
		p.SetDomainid(v.(string))
	}
	if v, ok := d.GetOk("imagestoreuuid"); ok {
		p.SetImagestoreid(v.(string))
	}
	if v, ok := d.GetOk("storeid"); ok {
		p.SetStorageid(v.(string))
	}
	if v, ok := d.GetOk("zoneid"); ok {
		p.SetZoneid(v.(string))
	}

	_, err := cs.Configuration.ResetConfiguration(p)
	if err != nil {
		return err
	}

	return nil
}