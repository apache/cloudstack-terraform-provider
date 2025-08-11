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
	"strings"
	"time"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudStackProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackProjectCreate,
		Read:   resourceCloudStackProjectRead,
		Update: resourceCloudStackProjectUpdate,
		Delete: resourceCloudStackProjectDelete,
		Importer: &schema.ResourceImporter{
			State: importStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"display_text": {
				Type:     schema.TypeString,
				Required: true, // Required for API version 4.18 and lower. TODO: Make this optional when support for API versions older than 4.18 is dropped.
			},

			"domain": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"account": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"accountid": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"userid": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceCloudStackProjectCreate(d *schema.ResourceData, meta any) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Get the name and display_text
	name := d.Get("name").(string)
	displaytext := d.Get("display_text").(string)

	// Get domain if provided
	var domain string
	domainSet := false
	if domainParam, ok := d.GetOk("domain"); ok {
		domain = domainParam.(string)
		domainSet = true
	}

	// Only check for an existing project if domain is set
	if domainSet {
		existingProject, err := getProjectByName(cs, name, domain)
		if err == nil {
			// Project with this name and domain already exists
			log.Printf("[DEBUG] Project with name %s and domain %s already exists, using existing project with ID: %s", name, domain, existingProject.Id)
			d.SetId(existingProject.Id)

			// Set the basic attributes to match the existing project
			d.Set("name", existingProject.Name)
			d.Set("display_text", existingProject.Displaytext)
			d.Set("domain", existingProject.Domain)

			return resourceCloudStackProjectRead(d, meta)
		} else if !strings.Contains(err.Error(), "not found") {
			// If we got an error other than "not found", return it
			return fmt.Errorf("error checking for existing project: %s", err)
		}
	}

	// Project doesn't exist, create a new one

	// The CloudStack API parameter order differs between versions:
	// - In API 4.18 and lower: displaytext is the first parameter and name is the second
	// - In API 4.19 and higher: name is the first parameter and displaytext is optional
	// The CloudStack Go SDK uses the API 4.18 parameter order
	p := cs.Project.NewCreateProjectParams(displaytext, name)

	// Set the domain if provided
	if domain != "" {
		domainid, e := retrieveID(cs, "domain", domain)
		if e != nil {
			return fmt.Errorf("error retrieving domain ID: %v", e)
		}
		p.SetDomainid(domainid)
	}

	// Set the account if provided
	if account, ok := d.GetOk("account"); ok {
		p.SetAccount(account.(string))
	}

	// Set the accountid if provided
	if accountid, ok := d.GetOk("accountid"); ok {
		p.SetAccountid(accountid.(string))
	}

	// Set the userid if provided
	if userid, ok := d.GetOk("userid"); ok {
		p.SetUserid(userid.(string))
	}

	log.Printf("[DEBUG] Creating project %s", name)
	r, err := cs.Project.CreateProject(p)
	if err != nil {
		return fmt.Errorf("error creating project %s: %s", name, err)
	}

	d.SetId(r.Id)
	log.Printf("[DEBUG] Project created with ID: %s", r.Id)

	// Wait for the project to be available
	// Use a longer timeout to ensure project creation completes
	ctx := context.Background()

	err = retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		project, err := getProjectByID(cs, d.Id(), domain)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				log.Printf("[DEBUG] Project %s not found yet, retrying...", d.Id())
				return retry.RetryableError(fmt.Errorf("project not yet created: %s", err))
			}
			return retry.NonRetryableError(fmt.Errorf("Error retrieving project: %s", err))
		}

		log.Printf("[DEBUG] Project %s found with name %s", d.Id(), project.Name)
		return nil
	})

	// Even if the retry times out, we should still try to read the resource
	// since it might have been created successfully
	if err != nil {
		log.Printf("[WARN] Timeout waiting for project %s to be available: %s", d.Id(), err)
	}

	// Read the resource state
	return resourceCloudStackProjectRead(d, meta)
}

// Helper function to get a project by ID
func getProjectByID(cs *cloudstack.CloudStackClient, id string, domain ...string) (*cloudstack.Project, error) {
	p := cs.Project.NewListProjectsParams()
	p.SetId(id)

	// If domain is provided, use it to narrow the search
	if len(domain) > 0 && domain[0] != "" {
		log.Printf("[DEBUG] Looking up project with ID: %s in domain: %s", id, domain[0])
		domainID, err := retrieveID(cs, "domain", domain[0])
		if err != nil {
			log.Printf("[WARN] Error retrieving domain ID for domain %s: %v", domain[0], err)
			// Continue without domain ID, but log the warning
		} else {
			p.SetDomainid(domainID)
		}
	} else {
		log.Printf("[DEBUG] Looking up project with ID: %s (no domain specified)", id)
	}

	l, err := cs.Project.ListProjects(p)
	if err != nil {
		log.Printf("[ERROR] Error calling ListProjects with ID %s: %v", id, err)
		return nil, err
	}

	log.Printf("[DEBUG] ListProjects returned Count: %d for ID: %s", l.Count, id)

	if l.Count == 0 {
		return nil, fmt.Errorf("project with id %s not found", id)
	}

	// Add validation to ensure the returned project ID matches the requested ID
	if l.Projects[0].Id != id {
		log.Printf("[WARN] Project ID mismatch - requested: %s, got: %s", id, l.Projects[0].Id)
		// Continue anyway to see if this is the issue
	}

	log.Printf("[DEBUG] Found project with ID: %s, Name: %s", l.Projects[0].Id, l.Projects[0].Name)
	return l.Projects[0], nil
}

// Helper function to get a project by name
func getProjectByName(cs *cloudstack.CloudStackClient, name string, domain string) (*cloudstack.Project, error) {
	p := cs.Project.NewListProjectsParams()
	p.SetName(name)

	// If domain is provided, use it to narrow the search
	if domain != "" {
		domainID, err := retrieveID(cs, "domain", domain)
		if err != nil {
			return nil, fmt.Errorf("error retrieving domain ID: %v", err)
		}
		p.SetDomainid(domainID)
	}

	log.Printf("[DEBUG] Looking up project with name: %s", name)
	l, err := cs.Project.ListProjects(p)
	if err != nil {
		return nil, err
	}

	if l.Count == 0 {
		return nil, fmt.Errorf("project with name %s not found", name)
	}

	// If multiple projects with the same name exist, log a warning and return the first one
	if l.Count > 1 {
		log.Printf("[WARN] Multiple projects found with name %s, using the first one", name)
	}

	log.Printf("[DEBUG] Found project %s with ID: %s", name, l.Projects[0].Id)
	return l.Projects[0], nil
}

func resourceCloudStackProjectRead(d *schema.ResourceData, meta any) error {
	cs := meta.(*cloudstack.CloudStackClient)

	log.Printf("[DEBUG] Retrieving project %s", d.Id())

	// Get project name and domain for potential fallback lookup
	name := d.Get("name").(string)
	var domain string
	if domainParam, ok := d.GetOk("domain"); ok {
		domain = domainParam.(string)
	}

	// Get the project details by ID
	project, err := getProjectByID(cs, d.Id(), domain)

	// If project not found by ID and we have a name, try to find it by name
	if err != nil && name != "" && (strings.Contains(err.Error(), "not found") ||
		strings.Contains(err.Error(), "does not exist") ||
		strings.Contains(err.Error(), "could not be found") ||
		strings.Contains(err.Error(), fmt.Sprintf(
			"Invalid parameter id value=%s due to incorrect long value format, "+
				"or entity does not exist", d.Id()))) {

		log.Printf("[DEBUG] Project %s not found by ID, trying to find by name: %s", d.Id(), name)
		project, err = getProjectByName(cs, name, domain)

		// If project not found by name either, resource doesn't exist
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				log.Printf("[DEBUG] Project with name %s not found either, marking as gone", name)
				d.SetId("")
				return nil
			}
			// For other errors during name lookup, return them
			return fmt.Errorf("error looking up project by name: %s", err)
		}

		// Found by name, update the ID
		log.Printf("[DEBUG] Found project by name %s with ID: %s", name, project.Id)
		d.SetId(project.Id)
	} else if err != nil {
		// For other errors during ID lookup, return them
		return fmt.Errorf("error retrieving project %s: %s", d.Id(), err)
	}

	log.Printf("[DEBUG] Found project %s: %s", d.Id(), project.Name)

	// Set the basic attributes
	d.Set("name", project.Name)
	d.Set("display_text", project.Displaytext)
	d.Set("domain", project.Domain)

	// Handle owner information more safely
	// Only set the account, accountid, and userid if they were explicitly set in the configuration
	// and if the owner information is available
	if _, ok := d.GetOk("account"); ok {
		// Safely handle the case where project.Owner might be nil or empty
		if len(project.Owner) > 0 {
			foundAccount := false
			for _, owner := range project.Owner {
				if account, ok := owner["account"]; ok {
					d.Set("account", account)
					foundAccount = true
					break
				}
			}
			if !foundAccount {
				log.Printf("[DEBUG] Project %s owner information doesn't contain account, keeping original value", d.Id())
			}
		} else {
			// Keep the original account value from the configuration
			// This prevents Terraform from thinking the resource has disappeared
			log.Printf("[DEBUG] Project %s owner information not available yet, keeping original account value", d.Id())
		}
	}

	if _, ok := d.GetOk("accountid"); ok {
		if len(project.Owner) > 0 {
			foundAccountID := false
			for _, owner := range project.Owner {
				if accountid, ok := owner["accountid"]; ok {
					d.Set("accountid", accountid)
					foundAccountID = true
					break
				}
			}
			if !foundAccountID {
				log.Printf("[DEBUG] Project %s owner information doesn't contain accountid, keeping original value", d.Id())
			}
		} else {
			log.Printf("[DEBUG] Project %s owner information not available yet, keeping original accountid value", d.Id())
		}
	}

	if _, ok := d.GetOk("userid"); ok {
		if len(project.Owner) > 0 {
			foundUserID := false
			for _, owner := range project.Owner {
				if userid, ok := owner["userid"]; ok {
					d.Set("userid", userid)
					foundUserID = true
					break
				}
			}
			if !foundUserID {
				log.Printf("[DEBUG] Project %s owner information doesn't contain userid, keeping original value", d.Id())
			}
		} else {
			log.Printf("[DEBUG] Project %s owner information not available yet, keeping original userid value", d.Id())
		}
	}

	return nil
}

func resourceCloudStackProjectUpdate(d *schema.ResourceData, meta any) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Check if the name or display text is changed
	if d.HasChange("name") || d.HasChange("display_text") {
		// Create a new parameter struct
		p := cs.Project.NewUpdateProjectParams(d.Id())

		// Set the name and display_text if they have changed
		// Note: The 'name' parameter is only available in API 4.19 and higher
		// If you're using API 4.18 or lower, the SetName method might not work
		// In that case, you might need to update the display_text only
		if d.HasChange("name") {
			p.SetName(d.Get("name").(string))
		}

		if d.HasChange("display_text") {
			p.SetDisplaytext(d.Get("display_text").(string))
		}

		log.Printf("[DEBUG] Updating project %s", d.Id())
		_, err := cs.Project.UpdateProject(p)
		if err != nil {
			return fmt.Errorf("Error updating project %s: %s", d.Id(), err)
		}
	}

	// Check if the account, accountid, or userid is changed
	if d.HasChange("account") || d.HasChange("accountid") || d.HasChange("userid") {
		// Create a new parameter struct
		p := cs.Project.NewUpdateProjectParams(d.Id())

		// Set swapowner to true to swap ownership with the account/user provided
		p.SetSwapowner(true)

		// Set the account if it has changed
		if d.HasChange("account") {
			p.SetAccount(d.Get("account").(string))
		}

		// Set the userid if it has changed
		if d.HasChange("userid") {
			p.SetUserid(d.Get("userid").(string))
		}

		// Note: accountid is not directly supported by the UpdateProject API,
		// but we can use the account parameter instead if accountid has changed
		if d.HasChange("accountid") && !d.HasChange("account") {
			// If accountid has changed but account hasn't, we need to look up the account name
			// This is a placeholder - in a real implementation, you would need to look up
			// the account name from the accountid
			log.Printf("[WARN] Updating accountid is not directly supported by the API. Please use account instead.")
		}

		log.Printf("[DEBUG] Updating project owner %s", d.Id())
		_, err := cs.Project.UpdateProject(p)
		if err != nil {
			return fmt.Errorf("Error updating project owner %s: %s", d.Id(), err)
		}
	}

	// Wait for the project to be updated
	ctx := context.Background()

	// Get domain if provided
	var domain string
	if domainParam, ok := d.GetOk("domain"); ok {
		domain = domainParam.(string)
	}

	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		project, err := getProjectByID(cs, d.Id(), domain)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				log.Printf("[DEBUG] Project %s not found after update, retrying...", d.Id())
				return retry.RetryableError(fmt.Errorf("project not found after update: %s", err))
			}
			return retry.NonRetryableError(fmt.Errorf("Error retrieving project after update: %s", err))
		}

		// Check if the project has the expected values
		if d.HasChange("name") && project.Name != d.Get("name").(string) {
			log.Printf("[DEBUG] Project %s name not updated yet, retrying...", d.Id())
			return retry.RetryableError(fmt.Errorf("project name not updated yet"))
		}

		if d.HasChange("display_text") && project.Displaytext != d.Get("display_text").(string) {
			log.Printf("[DEBUG] Project %s display_text not updated yet, retrying...", d.Id())
			return retry.RetryableError(fmt.Errorf("project display_text not updated yet"))
		}

		log.Printf("[DEBUG] Project %s updated successfully", d.Id())
		return nil
	})

	// Even if the retry times out, we should still try to read the resource
	// since it might have been updated successfully
	if err != nil {
		log.Printf("[WARN] Timeout waiting for project %s to be updated: %s", d.Id(), err)
	}

	// Read the resource state
	return resourceCloudStackProjectRead(d, meta)
}

func resourceCloudStackProjectDelete(d *schema.ResourceData, meta any) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Get project name and domain for potential fallback lookup
	name := d.Get("name").(string)
	var domain string
	if domainParam, ok := d.GetOk("domain"); ok {
		domain = domainParam.(string)
	}

	// First check if the project still exists by ID
	log.Printf("[DEBUG] Checking if project %s exists before deleting", d.Id())
	project, err := getProjectByID(cs, d.Id(), domain)

	// If project not found by ID, try to find it by name
	if err != nil && strings.Contains(err.Error(), "not found") {
		log.Printf("[DEBUG] Project %s not found by ID, trying to find by name: %s", d.Id(), name)
		project, err = getProjectByName(cs, name, domain)

		// If project not found by name either, we're done
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				log.Printf("[DEBUG] Project with name %s not found either, nothing to delete", name)
				return nil
			}
			// For other errors during name lookup, return them
			return fmt.Errorf("error looking up project by name: %s", err)
		}

		// Found by name, update the ID
		log.Printf("[DEBUG] Found project by name %s with ID: %s", name, project.Id)
		d.SetId(project.Id)
	} else if err != nil {
		// For other errors during ID lookup, return them
		return fmt.Errorf("error checking project existence before delete: %s", err)
	}

	log.Printf("[DEBUG] Found project %s (%s), proceeding with delete", d.Id(), project.Name)

	// Create a new parameter struct
	p := cs.Project.NewDeleteProjectParams(d.Id())

	log.Printf("[INFO] Deleting project: %s (%s)", d.Id(), project.Name)
	_, err = cs.Project.DeleteProject(p)
	if err != nil {
		// Check for various "not found" or "does not exist" error patterns
		if strings.Contains(err.Error(), "not found") ||
			strings.Contains(err.Error(), "does not exist") ||
			strings.Contains(err.Error(), fmt.Sprintf(
				"Invalid parameter id value=%s due to incorrect long value format, "+
					"or entity does not exist", d.Id())) {
			log.Printf("[DEBUG] Project %s no longer exists after delete attempt", d.Id())
			return nil
		}

		return fmt.Errorf("error deleting project %s: %s", d.Id(), err)
	}

	log.Printf("[DEBUG] Successfully deleted project: %s (%s)", d.Id(), project.Name)
	return nil
}
