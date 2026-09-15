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
	"strings"
	"testing"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// host_tags is a schema.TypeList, which the SDK returns as []interface{}, so
// asserting host_tags.([]string) panics. Creating a host with host_tags set
// must get past that conversion without panicking.
func TestHostCreateHostTagsDoesNotPanic(t *testing.T) {
	cs := cloudstack.NewClient("http://127.0.0.1:1", "key", "secret", false)

	d := schema.TestResourceDataRaw(t, resourceCloudStackHost().Schema, map[string]interface{}{
		"hypervisor":     "KVM",
		"pod_id":         "pod-1",
		"url":            "http://host-1",
		"zone_id":        "zone-1",
		"host_tags":      []interface{}{"tag1", "tag2"},
		"create_timeout": 1,
	})

	if err := resourceCloudStackHostCreate(d, cs); err == nil {
		t.Fatal("expected an error, got nil")
	}
}

// Updating a host must actually call UpdateHost. The update function used to
// build the params and then return without sending them, so the change was a
// silent no-op. With a real diff (allocation_state changed) the update must
// reach UpdateHost, which here fails against the unreachable endpoint.
func TestHostUpdateCallsUpdateHost(t *testing.T) {
	cs := cloudstack.NewClient("http://127.0.0.1:1", "key", "secret", false)

	sm := schema.InternalMap(resourceCloudStackHost().Schema)
	state := &terraform.InstanceState{
		ID:         "host-1",
		Attributes: map[string]string{"allocation_state": "Enabled"},
	}
	diff := &terraform.InstanceDiff{
		Attributes: map[string]*terraform.ResourceAttrDiff{
			"allocation_state": {Old: "Enabled", New: "Disabled"},
		},
	}
	d, err := sm.Data(state, diff)
	if err != nil {
		t.Fatalf("building resource data: %s", err)
	}

	err = resourceCloudStackHostUpdate(d, cs)
	if err == nil {
		t.Fatal("expected an error from the update request, got nil")
	}
	if !strings.Contains(err.Error(), "updating host") {
		t.Fatalf("expected the error to come from the UpdateHost call, got: %s", err)
	}
}
