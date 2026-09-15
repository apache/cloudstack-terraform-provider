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
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func diskUpdateData(t *testing.T) *schema.ResourceData {
	sm := schema.InternalMap(resourceCloudStackDisk().Schema)
	state := &terraform.InstanceState{
		ID: "vol-1",
		Attributes: map[string]string{
			"name":               "disk-1",
			"virtual_machine_id": "vm-1",
			"attach":             "true",
		},
	}
	diff := &terraform.InstanceDiff{
		Attributes: map[string]*terraform.ResourceAttrDiff{
			"virtual_machine_id": {Old: "vm-1", New: "vm-2"},
		},
	}
	d, err := sm.Data(state, diff)
	if err != nil {
		t.Fatalf("building resource data: %s", err)
	}
	return d
}

// The disk's attached VM is stored under virtual_machine_id, but the update
// detected a VM change with HasChange("virtual_machine"), a key that does not
// exist, so the branch never ran. A virtual_machine_id change must reach the
// detach path (rather than falling through to the attach path).
func TestDiskUpdateDetachesWhenVirtualMachineChanges(t *testing.T) {
	cs := cloudstack.NewClient("http://127.0.0.1:1", "key", "secret", false)

	if err := resourceCloudStackDiskUpdate(diskUpdateData(t), cs); err == nil {
		t.Fatal("expected an error, got nil")
	} else if !strings.Contains(err.Error(), "detaching") {
		t.Fatalf("a virtual_machine_id change should trigger a detach; got: %s", err)
	}
}

// When a plain detach fails, the fallback stops the VM currently holding the
// disk before retrying. During a virtual_machine_id change that VM is the old
// value, not the new one, so the fallback must stop the source VM.
func TestDiskUpdateDetachFallbackStopsSourceVM(t *testing.T) {
	var stoppedVM string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		r.ParseForm()
		switch r.Form.Get("command") {
		case "listVolumes":
			// The volume is still attached to the old VM.
			w.Write([]byte(`{"listvolumesresponse":{"count":1,"volume":[{"id":"vol-1","virtualmachineid":"vm-1"}]}}`))
		case "detachVolume":
			// Fail so the stop/retry fallback runs.
			w.WriteHeader(http.StatusExpectationFailed)
			w.Write([]byte(`{"detachvolumeresponse":{"errorcode":431,"errortext":"cannot detach"}}`))
		case "stopVirtualMachine":
			stoppedVM = r.Form.Get("id")
			w.WriteHeader(http.StatusExpectationFailed)
			w.Write([]byte(`{"stopvirtualmachineresponse":{"errorcode":431,"errortext":"stop failed"}}`))
		default:
			w.Write([]byte(`{"queryasyncjobresultresponse":{"jobstatus":2,"jobresult":{"errorcode":431,"errortext":"failed"}}}`))
		}
	}))
	defer server.Close()

	cs := cloudstack.NewClient(server.URL, "key", "secret", false)

	_ = resourceCloudStackDiskUpdate(diskUpdateData(t), cs)

	if stoppedVM != "vm-1" {
		t.Fatalf("detach fallback should stop the source VM vm-1, stopped %q", stoppedVM)
	}
}
