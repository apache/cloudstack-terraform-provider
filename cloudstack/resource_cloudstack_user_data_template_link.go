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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudStackUserDataTemplateLink() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackUserDataTemplateLinkCreate,
		Read:   resourceCloudStackUserDataTemplateLinkRead,
		Update: resourceCloudStackUserDataTemplateLinkUpdate,
		Delete: resourceCloudStackUserDataTemplateLinkDelete,

		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"iso_id"},
				ForceNew:      true,
			},

			"iso_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"template_id"},
				ForceNew:      true,
			},

			"user_data_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"user_data_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ALLOWOVERRIDE",
				ValidateFunc: validateUserDataPolicy,
			},

			// Computed attributes from template response
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"display_text": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"is_ready": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"template_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"user_data_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"user_data_params": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func validateUserDataPolicy(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	validPolicies := []string{"ALLOWOVERRIDE", "APPEND", "DENYOVERRIDE"}

	for _, policy := range validPolicies {
		if value == policy {
			return
		}
	}

	errors = append(errors, fmt.Errorf("user_data_policy must be one of: %v", validPolicies))
	return
}

func resourceCloudStackUserDataTemplateLinkCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create parameter struct
	p := cs.Template.NewLinkUserDataToTemplateParams()

	// Set template or ISO ID
	if templateId, ok := d.GetOk("template_id"); ok {
		p.SetTemplateid(templateId.(string))
		d.SetId(fmt.Sprintf("template-%s", templateId.(string)))
	} else if isoId, ok := d.GetOk("iso_id"); ok {
		p.SetIsoid(isoId.(string))
		d.SetId(fmt.Sprintf("iso-%s", isoId.(string)))
	} else {
		return fmt.Errorf("Either template_id or iso_id must be specified")
	}

	// Set optional parameters
	if userDataId, ok := d.GetOk("user_data_id"); ok {
		p.SetUserdataid(userDataId.(string))
	}

	if userDataPolicy, ok := d.GetOk("user_data_policy"); ok {
		p.SetUserdatapolicy(userDataPolicy.(string))
	}

	log.Printf("[DEBUG] Linking UserData to Template/ISO")
	r, err := cs.Template.LinkUserDataToTemplate(p)
	if err != nil {
		return fmt.Errorf("Error linking UserData to Template/ISO: %s", err)
	}

	// Store the template/ISO ID as resource ID
	if r.Id != "" {
		d.SetId(r.Id)
	}

	return resourceCloudStackUserDataTemplateLinkRead(d, meta)
}

func resourceCloudStackUserDataTemplateLinkRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	var template *cloudstack.Template
	var err error
	var count int

	// Determine if we're dealing with a template or ISO
	if templateId, ok := d.GetOk("template_id"); ok {
		// Get template details
		template, count, err = cs.Template.GetTemplateByID(
			templateId.(string),
			"all",
		)
	} else if isoId, ok := d.GetOk("iso_id"); ok {
		// Get ISO details (ISOs are also handled by the Template service)
		template, count, err = cs.Template.GetTemplateByID(
			isoId.(string),
			"all",
		)
	} else {
		return fmt.Errorf("Either template_id or iso_id must be specified")
	}

	if err != nil {
		if count == 0 {
			log.Printf("[DEBUG] Template/ISO no longer exists")
			d.SetId("")
			return nil
		}
		return err
	}

	// Update computed attributes
	d.Set("name", template.Name)
	d.Set("display_text", template.Displaytext)
	d.Set("is_ready", template.Isready)
	d.Set("template_type", template.Templatetype)
	d.Set("user_data_name", template.Userdataname)
	d.Set("user_data_params", template.Userdataparams)

	return nil
}

func resourceCloudStackUserDataTemplateLinkUpdate(d *schema.ResourceData, meta interface{}) error {
	// If user_data_id or user_data_policy changes, we need to re-link
	if d.HasChange("user_data_id") || d.HasChange("user_data_policy") {
		return resourceCloudStackUserDataTemplateLinkCreate(d, meta)
	}

	return resourceCloudStackUserDataTemplateLinkRead(d, meta)
}

func resourceCloudStackUserDataTemplateLinkDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create parameter struct for unlinking (no userdata id = unlink)
	p := cs.Template.NewLinkUserDataToTemplateParams()

	// Set template or ISO ID
	if templateId, ok := d.GetOk("template_id"); ok {
		p.SetTemplateid(templateId.(string))
	} else if isoId, ok := d.GetOk("iso_id"); ok {
		p.SetIsoid(isoId.(string))
	} else {
		return fmt.Errorf("Either template_id or iso_id must be specified")
	}

	// Don't set userdataid - this will unlink existing userdata

	log.Printf("[DEBUG] Unlinking UserData from Template/ISO")
	_, err := cs.Template.LinkUserDataToTemplate(p)
	if err != nil {
		return fmt.Errorf("Error unlinking UserData from Template/ISO: %s", err)
	}

	return nil
}
