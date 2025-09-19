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

func intervalTypeToString(intervalType int) string {
	switch intervalType {
	case 0:
		return "HOURLY"
	case 1:
		return "DAILY"
	case 2:
		return "WEEKLY"
	case 3:
		return "MONTHLY"
	default:
		return fmt.Sprintf("%d", intervalType)
	}
}

func resourceCloudStackSnapshotPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudstackSnapshotPolicyCreate,
		Read:   resourceCloudstackSnapshotPolicyRead,
		Update: resourceCloudstackSnapshotPolicyUpdate,
		Delete: resourceCloudstackSnapshotPolicyDelete,

		Schema: map[string]*schema.Schema{
			"volume_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"interval_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"max_snaps": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"schedule": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"timezone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"custom_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceCloudstackSnapshotPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.Snapshot.NewCreateSnapshotPolicyParams(
		d.Get("interval_type").(string),
		d.Get("max_snaps").(int),
		d.Get("schedule").(string),
		d.Get("timezone").(string),
		d.Get("volume_id").(string),
	)

	if v, ok := d.GetOk("zone_ids"); ok && v != nil {
		zoneIDs := []string{}
		for _, id := range v.([]interface{}) {
			zoneIDs = append(zoneIDs, id.(string))
		}
		p.SetZoneids(zoneIDs)
	}

	snapshotPolicy, err := cs.Snapshot.CreateSnapshotPolicy(p)
	if err != nil {
		return fmt.Errorf("Error creating snapshot policy: %s", err)
	}

	log.Printf("[DEBUG] CreateSnapshotPolicy response: %+v", snapshotPolicy)

	if snapshotPolicy.Id == "" {
		log.Printf("[DEBUG] CloudStack returned empty ID, trying to find created policy by volume ID")

		listParams := cs.Snapshot.NewListSnapshotPoliciesParams()
		listParams.SetVolumeid(d.Get("volume_id").(string))

		resp, listErr := cs.Snapshot.ListSnapshotPolicies(listParams)
		if listErr != nil {
			return fmt.Errorf("Error listing snapshot policies to find created policy: %s", listErr)
		}

		if resp.Count == 0 {
			return fmt.Errorf("No snapshot policies found for volume after creation")
		}

		foundPolicy := resp.SnapshotPolicies[resp.Count-1]
		log.Printf("[DEBUG] Found policy with ID: %s", foundPolicy.Id)
		d.SetId(foundPolicy.Id)
	} else {
		d.SetId(snapshotPolicy.Id)
	}

	log.Printf("[DEBUG] Snapshot policy created with ID: %s", d.Id())

	// Set tags if provided
	if err := setTags(cs, d, "SnapshotPolicy"); err != nil {
		return fmt.Errorf("Error setting tags on snapshot policy %s: %s", d.Id(), err)
	}

	return resourceCloudstackSnapshotPolicyRead(d, meta)
}

func resourceCloudstackSnapshotPolicyRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	if d.Id() == "" {
		log.Printf("[DEBUG] Snapshot policy ID is empty")
		return fmt.Errorf("Snapshot policy ID is empty")
	}

	p := cs.Snapshot.NewListSnapshotPoliciesParams()
	p.SetId(d.Id())

	resp, err := cs.Snapshot.ListSnapshotPolicies(p)
	if err != nil {
		return fmt.Errorf("Failed to list snapshot policies: %s", err)
	}

	if resp.Count == 0 {
		log.Printf("[DEBUG] Snapshot policy %s not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	snapshotPolicy := resp.SnapshotPolicies[0]

	d.Set("volume_id", snapshotPolicy.Volumeid)
	d.Set("interval_type", intervalTypeToString(snapshotPolicy.Intervaltype))
	d.Set("max_snaps", snapshotPolicy.Maxsnaps)
	d.Set("schedule", snapshotPolicy.Schedule)
	d.Set("timezone", snapshotPolicy.Timezone)

	if snapshotPolicy.Zone != nil {
		zoneIDs := []string{}
		for _, zone := range snapshotPolicy.Zone {
			if zoneMap, ok := zone.(map[string]interface{}); ok {
				if id, ok := zoneMap["id"].(string); ok {
					zoneIDs = append(zoneIDs, id)
				}
			}
		}
		d.Set("zone_ids", zoneIDs)
	} else {
		d.Set("zone_ids", nil)
	}

	// Handle tags
	tags := make(map[string]interface{})
	for _, tag := range snapshotPolicy.Tags {
		tags[tag.Key] = tag.Value
	}
	d.Set("tags", tags)

	return nil
}

func resourceCloudstackSnapshotPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.Snapshot.NewUpdateSnapshotPolicyParams()
	p.SetId(d.Id())
	if v, ok := d.GetOk("custom_id"); ok {
		p.SetCustomid(v.(string))
	}

	_, err := cs.Snapshot.UpdateSnapshotPolicy(p)
	if err != nil {
		return err
	}

	// Handle tags
	if d.HasChange("tags") {
		if err := updateTags(cs, d, "SnapshotPolicy"); err != nil {
			return fmt.Errorf("Error updating tags on snapshot policy %s: %s", d.Id(), err)
		}
	}

	return resourceCloudstackSnapshotPolicyRead(d, meta)
}

func resourceCloudstackSnapshotPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.Snapshot.NewDeleteSnapshotPoliciesParams()
	p.SetId(d.Id())

	_, err := cs.Snapshot.DeleteSnapshotPolicies(p)
	if err != nil {
		return fmt.Errorf("Failed to delete snapshot policy %s: %s", d.Id(), err)
	}

	d.SetId("")

	return nil
}
