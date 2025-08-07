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
	"strconv"
	"strings"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// resourceTypeMap maps string resource types to their integer values
var resourceTypeMap = map[string]int{
	"instance":         0,
	"ip":               1,
	"volume":           2,
	"snapshot":         3,
	"template":         4,
	"project":          5,
	"network":          6,
	"vpc":              7,
	"cpu":              8,
	"memory":           9,
	"primarystorage":   10,
	"secondarystorage": 11,
}

func resourceCloudStackLimits() *schema.Resource {
	return &schema.Resource{
		Read:   resourceCloudStackLimitsRead,
		Update: resourceCloudStackLimitsUpdate,
		Create: resourceCloudStackLimitsCreate,
		Delete: resourceCloudStackLimitsDelete,
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"instance", "ip", "volume", "snapshot", "template", "project", "network", "vpc",
					"cpu", "memory", "primarystorage", "secondarystorage",
				}, false), // false disables case-insensitive matching
				Description: "The type of resource to update the limits. Available types are: " +
					"instance, ip, volume, snapshot, template, project, network, vpc, cpu, memory, " +
					"primarystorage, secondarystorage",
			},
			"account": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Update resource for a specified account. Must be used with the domainid parameter.",
			},
			"domainid": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Update resource limits for all accounts in specified domain. If used with the account parameter, updates resource limits for a specified account in specified domain.",
			},
			"max": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Maximum resource limit. Use -1 for unlimited resource limit. A value of 0 means zero resources are allowed, though the CloudStack API may return -1 for a limit set to 0.",
			},
			"projectid": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Update resource limits for project.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// getResourceType gets the resource type from the type field
func getResourceType(d *schema.ResourceData) (int, error) {
	// Check if type is set
	if v, ok := d.GetOk("type"); ok {
		typeStr := v.(string)
		if resourcetype, ok := resourceTypeMap[typeStr]; ok {
			return resourcetype, nil
		}
		return 0, fmt.Errorf("invalid type value: %s", typeStr)
	}

	return 0, fmt.Errorf("type must be specified")
}

func resourceCloudStackLimitsCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	resourcetype, err := getResourceType(d)
	if err != nil {
		return err
	}

	account := d.Get("account").(string)
	domainid := d.Get("domainid").(string)
	projectid := d.Get("projectid").(string)

	// Validate account and domain parameters
	if account != "" && domainid == "" {
		return fmt.Errorf("domainid is required when account is specified")
	}

	// Create a new parameter struct
	p := cs.Limit.NewUpdateResourceLimitParams(resourcetype)
	if account != "" {
		p.SetAccount(account)
	}
	if domainid != "" {
		p.SetDomainid(domainid)
	}
	if maxVal, ok := d.GetOk("max"); ok {
		maxIntVal := maxVal.(int)
		log.Printf("[DEBUG] Setting max value to %d", maxIntVal)
		p.SetMax(int64(maxIntVal))
	}
	if projectid != "" {
		p.SetProjectid(projectid)
	}

	log.Printf("[DEBUG] Updating Resource Limit for type %d", resourcetype)
	_, err = cs.Limit.UpdateResourceLimit(p)

	if err != nil {
		return fmt.Errorf("Error creating resource limit: %s", err)
	}

	// Generate a unique ID based on the parameters
	id := generateResourceID(resourcetype, account, domainid, projectid)
	d.SetId(id)

	return resourceCloudStackLimitsRead(d, meta)
}

// generateResourceID creates a unique ID for the resource based on its parameters
func generateResourceID(resourcetype int, account, domainid, projectid string) string {
	if projectid != "" {
		return fmt.Sprintf("%d-project-%s", resourcetype, projectid)
	}

	if account != "" && domainid != "" {
		return fmt.Sprintf("%d-account-%s-%s", resourcetype, account, domainid)
	}

	if domainid != "" {
		return fmt.Sprintf("%d-domain-%s", resourcetype, domainid)
	}

	return fmt.Sprintf("%d", resourcetype)
}

func resourceCloudStackLimitsRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Get the resourcetype from the type field
	resourcetype, err := getResourceType(d)
	if err != nil {
		// If there's an error getting the type, try to extract it from the ID
		idParts := strings.Split(d.Id(), "-")
		if len(idParts) > 0 {
			if rt, err := strconv.Atoi(idParts[0]); err == nil {
				resourcetype = rt
				// Find the string representation for this numeric type
				for typeStr, typeVal := range resourceTypeMap {
					if typeVal == rt {
						d.Set("type", typeStr)
						break
					}
				}
			}
		}
	}

	account := d.Get("account").(string)
	domainid := d.Get("domainid").(string)
	projectid := d.Get("projectid").(string)

	// Create a new parameter struct
	p := cs.Limit.NewListResourceLimitsParams()
	p.SetResourcetype(resourcetype)
	if account != "" {
		p.SetAccount(account)
	}
	if domainid != "" {
		p.SetDomainid(domainid)
	}
	if projectid != "" {
		p.SetProjectid(projectid)
	}

	// Retrieve the resource limits
	l, err := cs.Limit.ListResourceLimits(p)
	if err != nil {
		return fmt.Errorf("error retrieving resource limits: %s", err)
	}

	if l.Count == 0 {
		log.Printf("[DEBUG] Resource limit not found")
		d.SetId("")
		return nil
	}

	// Update the config
	for _, limit := range l.ResourceLimits {
		if limit.Resourcetype == fmt.Sprintf("%d", resourcetype) {
			log.Printf("[DEBUG] Retrieved max value from API: %d", limit.Max)

			// If the user set max to 0 but the API returned -1, keep it as 0 in the state
			if limit.Max == -1 && d.Get("max").(int) == 0 {
				log.Printf("[DEBUG] API returned -1 for a limit set to 0, keeping it as 0 in state")
				d.Set("max", 0)
			} else {
				d.Set("max", limit.Max)
			}

			// Only set the type field if it was originally specified in the configuration
			if v, ok := d.GetOk("type"); ok {
				// Preserve the original case of the type parameter
				d.Set("type", v.(string))
			}

			return nil
		}
	}

	return fmt.Errorf("resource limit not found")
}

func resourceCloudStackLimitsUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	resourcetype, err := getResourceType(d)
	if err != nil {
		return err
	}

	account := d.Get("account").(string)
	domainid := d.Get("domainid").(string)
	projectid := d.Get("projectid").(string)

	// Create a new parameter struct
	p := cs.Limit.NewUpdateResourceLimitParams(resourcetype)
	if account != "" {
		p.SetAccount(account)
	}
	if domainid != "" {
		p.SetDomainid(domainid)
	}
	if maxVal, ok := d.GetOk("max"); ok {
		maxIntVal := maxVal.(int)
		log.Printf("[DEBUG] Setting max value to %d", maxIntVal)
		p.SetMax(int64(maxIntVal))
	}
	if projectid != "" {
		p.SetProjectid(projectid)
	}

	log.Printf("[DEBUG] Updating Resource Limit for type %d", resourcetype)
	_, err = cs.Limit.UpdateResourceLimit(p)

	if err != nil {
		return fmt.Errorf("Error updating resource limit: %s", err)
	}

	return resourceCloudStackLimitsRead(d, meta)
}

func resourceCloudStackLimitsDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	resourcetype, err := getResourceType(d)
	if err != nil {
		return err
	}

	account := d.Get("account").(string)
	domainid := d.Get("domainid").(string)
	projectid := d.Get("projectid").(string)

	// Create a new parameter struct
	p := cs.Limit.NewUpdateResourceLimitParams(resourcetype)
	if account != "" {
		p.SetAccount(account)
	}
	if domainid != "" {
		p.SetDomainid(domainid)
	}
	if projectid != "" {
		p.SetProjectid(projectid)
	}
	p.SetMax(-1) // Set to -1 to remove the limit

	log.Printf("[DEBUG] Removing Resource Limit for type %d", resourcetype)
	_, err = cs.Limit.UpdateResourceLimit(p)

	if err != nil {
		return fmt.Errorf("Error removing Resource Limit: %s", err)
	}

	d.SetId("")

	return nil
}
