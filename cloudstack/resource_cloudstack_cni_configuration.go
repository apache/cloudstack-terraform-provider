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

func resourceCloudStackCniConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackCniConfigurationCreate,
		Read:   resourceCloudStackCniConfigurationRead,
		Delete: resourceCloudStackCniConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: importStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the CNI configuration",
			},

			"cni_config": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "CNI Configuration content to be registered",
			},

			"account": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "An optional account for the CNI configuration. Must be used with domain_id.",
			},

			"domain_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "An optional domain ID for the CNI configuration. If the account parameter is used, domain_id must also be used.",
			},

			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "An optional project for the CNI configuration",
			},

			"params": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Description: "List of variables declared in CNI configuration content",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceCloudStackCniConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	name := d.Get("name").(string)
	log.Printf("[DEBUG] Creating CNI configuration: %s", name)

	p := cs.Configuration.NewRegisterCniConfigurationParams(name)

	if v, ok := d.GetOk("cni_config"); ok {
		cniConfig := v.(string)
		log.Printf("[DEBUG] CNI config data length: %d bytes", len(cniConfig))
		p.SetCniconfig(cniConfig)
	} else {
		return fmt.Errorf("CNI configuration content is required but not provided")
	}

	if account := d.Get("account").(string); account != "" {
		log.Printf("[DEBUG] Setting account: %s", account)
		p.SetAccount(account)
	}

	if domainID := d.Get("domain_id").(string); domainID != "" {
		log.Printf("[DEBUG] Setting domain ID: %s", domainID)
		p.SetDomainid(domainID)
	}

	if projectID := d.Get("project_id").(string); projectID != "" {
		log.Printf("[DEBUG] Setting project ID: %s", projectID)
		p.SetProjectid(projectID)
	}

	if params, ok := d.GetOk("params"); ok {
		paramsList := []string{}
		for _, param := range params.(*schema.Set).List() {
			paramsList = append(paramsList, param.(string))
		}
		if len(paramsList) > 0 {
			paramsStr := strings.Join(paramsList, ",")
			log.Printf("[DEBUG] Setting params: %s", paramsStr)
			p.SetParams(paramsStr)
		}
	}

	resp, err := cs.Configuration.RegisterCniConfiguration(p)
	if err != nil {
		return fmt.Errorf("Error creating CNI configuration %s: %s", name, err)
	}

	log.Printf("[DEBUG] CNI configuration creation response: %+v", resp)

	// List configurations to find the created one by name since direct ID access is not available
	listParams := cs.Configuration.NewListCniConfigurationParams()
	listParams.SetName(name)

	// Add context parameters if available
	if account := d.Get("account").(string); account != "" {
		listParams.SetAccount(account)
	}
	if domainID := d.Get("domain_id").(string); domainID != "" {
		listParams.SetDomainid(domainID)
	}
	if projectID := d.Get("project_id").(string); projectID != "" {
		listParams.SetProjectid(projectID)
	}

	listResp, err := cs.Configuration.ListCniConfiguration(listParams)
	if err != nil {
		return fmt.Errorf("Error listing CNI configurations after creation: %s", err)
	}

	if listResp.Count == 0 {
		return fmt.Errorf("CNI configuration %s was created but could not be found", name)
	}

	// Use the first (and should be only) result
	config := listResp.CniConfiguration[0]
	d.SetId(config.Id)
	log.Printf("[DEBUG] CNI configuration %s successfully created with ID: %s", name, d.Id())

	return resourceCloudStackCniConfigurationRead(d, meta)
}

func resourceCloudStackCniConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	log.Printf("[DEBUG] Reading CNI configuration: %s", d.Id())

	p := cs.Configuration.NewListCniConfigurationParams()
	p.SetId(d.Id())

	config, err := cs.Configuration.ListCniConfiguration(p)
	if err != nil {
		return fmt.Errorf("Error listing CNI configuration: %s", err)
	}
	if config.Count == 0 {
		log.Printf("[DEBUG] CNI configuration %s no longer exists", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("name", config.CniConfiguration[0].Name)
	d.Set("cni_config", config.CniConfiguration[0].Userdata)
	d.Set("account", config.CniConfiguration[0].Account)
	d.Set("domain_id", config.CniConfiguration[0].Domainid)
	d.Set("project_id", config.CniConfiguration[0].Projectid)

	if config.CniConfiguration[0].Params != "" {
		paramsList := strings.Split(config.CniConfiguration[0].Params, ",")
		d.Set("params", paramsList)
	}

	return nil
}

func resourceCloudStackCniConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	log.Printf("[DEBUG] Deleting CNI configuration: %s", d.Id())

	p := cs.Configuration.NewDeleteCniConfigurationParams(d.Id())

	_, err := cs.Configuration.DeleteCniConfiguration(p)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") ||
			strings.Contains(err.Error(), "not found") {
			log.Printf("[DEBUG] CNI configuration %s already deleted", d.Id())
			return nil
		}
		return fmt.Errorf("Error deleting CNI configuration %s: %s", d.Id(), err)
	}

	log.Printf("[DEBUG] CNI configuration %s deleted", d.Id())
	return nil
}
