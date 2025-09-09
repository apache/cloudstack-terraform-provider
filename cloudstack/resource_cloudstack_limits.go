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
	"context"
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
				Description: "Update resource for a specified account. Must be used with the domain_id parameter.",
			},
			"domain_id": {
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
			"configured_max": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Internal field to track the originally configured max value to distinguish between 0 and -1 when CloudStack returns -1.",
			},
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Update resource limits for project.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: resourceCloudStackLimitsImport,
		},
	}
}

// resourceCloudStackLimitsImport parses composite import IDs and sets resource fields accordingly.
func resourceCloudStackLimitsImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// Expected formats:
	// - type-account-accountname-domain_id (for account-specific limits)
	// - type-project-projectid (for project-specific limits)
	// - type-domain-domain_id (for domain-specific limits)

	log.Printf("[DEBUG] Importing resource with ID: %s", d.Id())

	// First, extract the resource type which is always the first part
	idParts := strings.SplitN(d.Id(), "-", 2)
	if len(idParts) < 2 {
		return nil, fmt.Errorf("unexpected import ID format (%q), expected type-account-accountname-domain_id, type-domain-domain_id, or type-project-projectid", d.Id())
	}

	// Parse the resource type
	typeInt, err := strconv.Atoi(idParts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid type value in import ID: %s", idParts[0])
	}

	// Find the string representation for this numeric type
	var typeStr string
	for k, v := range resourceTypeMap {
		if v == typeInt {
			typeStr = k
			break
		}
	}
	if typeStr == "" {
		return nil, fmt.Errorf("unknown type value in import ID: %d", typeInt)
	}
	if err := d.Set("type", typeStr); err != nil {
		return nil, err
	}

	// Get the original resource ID from the state
	originalID := d.Id()
	log.Printf("[DEBUG] Original import ID: %s", originalID)

	// Instead of trying to parse the complex ID, let's create a new resource
	// and read it from the API to get the correct values
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct for listing resource limits
	p := cs.Limit.NewListResourceLimitsParams()
	p.SetResourcetype(typeInt)

	// Try to determine the resource scope from the ID format
	remainingID := idParts[1]

	// Extract the resource scope from the ID
	if strings.HasPrefix(remainingID, "domain-") {
		// It's a domain-specific limit
		log.Printf("[DEBUG] Detected domain-specific limit")
		// We'll use the Read function to get the domain ID from the state
		// after setting a temporary ID
		d.SetId(originalID)
		return []*schema.ResourceData{d}, nil
	} else if strings.HasPrefix(remainingID, "project-") {
		// It's a project-specific limit
		log.Printf("[DEBUG] Detected project-specific limit")
		// We'll use the Read function to get the project ID from the state
		// after setting a temporary ID
		d.SetId(originalID)
		return []*schema.ResourceData{d}, nil
	} else if strings.HasPrefix(remainingID, "account-") {
		// It's an account-specific limit
		log.Printf("[DEBUG] Detected account-specific limit")
		// We'll use the Read function to get the account and domain ID from the state
		// after setting a temporary ID
		d.SetId(originalID)
		return []*schema.ResourceData{d}, nil
	} else {
		// For backward compatibility, assume it's a global limit
		log.Printf("[DEBUG] Detected global limit")
		d.SetId(originalID)
		return []*schema.ResourceData{d}, nil
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
	domain_id := d.Get("domain_id").(string)
	projectid := d.Get("project").(string)

	// Validate account and domain parameters
	if account != "" && domain_id == "" {
		return fmt.Errorf("domain_id is required when account is specified")
	}

	// Create a new parameter struct
	p := cs.Limit.NewUpdateResourceLimitParams(resourcetype)
	if account != "" {
		p.SetAccount(account)
	}
	if domain_id != "" {
		p.SetDomainid(domain_id)
	}
	// Check for max value - need to handle zero values explicitly
	maxVal := d.Get("max")
	if maxVal != nil {
		maxIntVal := maxVal.(int)
		log.Printf("[DEBUG] Setting max value to %d", maxIntVal)
		p.SetMax(int64(maxIntVal))

		// Store the original configured value for later reference
		// This helps the Read function distinguish between 0 and -1 when CloudStack returns -1
		if err := d.Set("configured_max", maxIntVal); err != nil {
			return fmt.Errorf("error storing configured max value: %w", err)
		}
	} else {
		log.Printf("[DEBUG] No max value found in configuration during Create")
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
	id := generateResourceID(resourcetype, account, domain_id, projectid)
	d.SetId(id)

	return resourceCloudStackLimitsRead(d, meta)
}

// generateResourceID creates a unique ID for the resource based on its parameters
func generateResourceID(resourcetype int, account, domain_id, projectid string) string {
	if projectid != "" {
		return fmt.Sprintf("%d-project-%s", resourcetype, projectid)
	}

	if account != "" && domain_id != "" {
		return fmt.Sprintf("%d-account-%s-%s", resourcetype, account, domain_id)
	}

	if domain_id != "" {
		return fmt.Sprintf("%d-domain-%s", resourcetype, domain_id)
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
						if err := d.Set("type", typeStr); err != nil {
							return fmt.Errorf("error setting type: %s", err)
						}
						break
					}
				}

				// Handle different ID formats
				if len(idParts) >= 3 {
					if idParts[1] == "domain" {
						// Format: resourcetype-domain-domain_id
						if err := d.Set("domain_id", idParts[2]); err != nil {
							return fmt.Errorf("error setting domain_id: %s", err)
						}
					} else if idParts[1] == "project" {
						// Format: resourcetype-project-projectid
						if err := d.Set("project", idParts[2]); err != nil {
							return fmt.Errorf("error setting project: %s", err)
						}
					} else if idParts[1] == "account" && len(idParts) >= 4 {
						// Format: resourcetype-account-account-domain_id
						if err := d.Set("account", idParts[2]); err != nil {
							return fmt.Errorf("error setting account: %s", err)
						}
						if err := d.Set("domain_id", idParts[3]); err != nil {
							return fmt.Errorf("error setting domain_id: %s", err)
						}
					}
				}
			}
		}
	}

	account := d.Get("account").(string)
	domain_id := d.Get("domain_id").(string)
	projectid := d.Get("project").(string)

	// Create a new parameter struct
	p := cs.Limit.NewListResourceLimitsParams()
	p.SetResourcetype(resourcetype)
	if account != "" {
		p.SetAccount(account)
	}
	if domain_id != "" {
		p.SetDomainid(domain_id)
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

	// Get the first (and should be only) limit from the results
	limit := l.ResourceLimits[0]

	// Handle the max value - CloudStack may return -1 for both unlimited and zero limits
	// We need to preserve the original value from the configuration when possible
	log.Printf("[DEBUG] CloudStack returned max value: %d", limit.Max)
	if limit.Max == -1 {
		// CloudStack returns -1 for both unlimited and zero limits
		// Check if we have the originally configured value stored
		if configuredMax, hasConfiguredMax := d.GetOk("configured_max"); hasConfiguredMax {
			configuredValue := configuredMax.(int)
			log.Printf("[DEBUG] Found configured max value: %d, using it", configuredValue)
			// Use the originally configured value (0 for zero limit, -1 for unlimited)
			if err := d.Set("max", configuredValue); err != nil {
				return fmt.Errorf("error setting max to configured value %d: %w", configuredValue, err)
			}
		} else {
			log.Printf("[DEBUG] No configured max value found, treating -1 as unlimited")
			// If no configured value is stored, treat -1 as unlimited
			if err := d.Set("max", -1); err != nil {
				return fmt.Errorf("error setting max to unlimited (-1): %w", err)
			}
		}
	} else {
		log.Printf("[DEBUG] Using positive max value from API: %d", limit.Max)
		// For any positive value, use it directly from the API
		if err := d.Set("max", int(limit.Max)); err != nil {
			return fmt.Errorf("error setting max: %w", err)
		}
	}

	// Preserve original type configuration if it exists
	if typeValue, ok := d.GetOk("type"); ok {
		if err := d.Set("type", typeValue.(string)); err != nil {
			return fmt.Errorf("error setting type: %w", err)
		}
	}

	return nil
}

func resourceCloudStackLimitsUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	resourcetype, err := getResourceType(d)
	if err != nil {
		return err
	}

	account := d.Get("account").(string)
	domain_id := d.Get("domain_id").(string)
	projectid := d.Get("project").(string)

	// Create a new parameter struct
	p := cs.Limit.NewUpdateResourceLimitParams(resourcetype)
	if account != "" {
		p.SetAccount(account)
	}
	if domain_id != "" {
		p.SetDomainid(domain_id)
	}
	if maxVal, ok := d.GetOk("max"); ok {
		maxIntVal := maxVal.(int)
		log.Printf("[DEBUG] Setting max value to %d", maxIntVal)
		p.SetMax(int64(maxIntVal))

		// Store the original configured value for later reference
		// This helps the Read function distinguish between 0 and -1 when CloudStack returns -1
		log.Printf("[DEBUG] Storing configured max value in update: %d", maxIntVal)
		if err := d.Set("configured_max", maxIntVal); err != nil {
			return fmt.Errorf("error storing configured max value: %w", err)
		}
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
	domain_id := d.Get("domain_id").(string)
	projectid := d.Get("project").(string)

	// Create a new parameter struct
	p := cs.Limit.NewUpdateResourceLimitParams(resourcetype)
	if account != "" {
		p.SetAccount(account)
	}
	if domain_id != "" {
		p.SetDomainid(domain_id)
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
