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
	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudStackServiceOfferingUnConstrained() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackServiceOfferingUnConstrainedCreate,
		Read:   resourceCloudStackServiceOfferingRead,
		Update: resourceCloudStackServiceOfferingUpdate,
		Delete: resourceCloudStackServiceOfferingDelete,
		Schema: serviceOfferingMergeCommonSchema(map[string]*schema.Schema{}),
	}
}

func resourceCloudStackServiceOfferingUnConstrainedCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.ServiceOffering.NewCreateServiceOfferingParams(d.Get("display_text").(string), d.Get("name").(string))

	// set common params
	serviceOfferingCreateParams(p, d)

	// Set unconstrained params
	p.SetCustomized(true)

	// request
	r, err := cs.ServiceOffering.CreateServiceOffering(p)
	if err != nil {
		return err
	}

	d.SetId(r.Id)

	return resourceCloudStackServiceOfferingRead(d, meta)
}
