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
	"testing"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// When a resource has been deleted out of band, GetByID reports count 0 (as
// opposed to a transport error, which reports -1). Read must clear the id and
// return nil so Terraform plans a recreate, instead of returning an error.
func TestZoneReadRemovesDeletedFromState(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"listzonesresponse":{"count":0}}`))
	}))
	defer server.Close()

	cs := cloudstack.NewClient(server.URL, "key", "secret", false)

	d := schema.TestResourceDataRaw(t, resourceCloudStackZone().Schema, map[string]interface{}{})
	d.SetId("deleted-zone")

	if err := resourceCloudStackZoneRead(d, cs); err != nil {
		t.Fatalf("read of a deleted zone should not error, got: %s", err)
	}
	if d.Id() != "" {
		t.Fatalf("read of a deleted zone should clear the id, got: %q", d.Id())
	}
}

// A transport or API error (count -1, as opposed to a genuinely absent
// resource) must be returned and must not clear the id, so a temporary outage
// does not drop the resource from state.
func TestZoneReadPreservesStateOnAPIError(t *testing.T) {
	cs := cloudstack.NewClient("http://127.0.0.1:1", "key", "secret", false)

	d := schema.TestResourceDataRaw(t, resourceCloudStackZone().Schema, map[string]interface{}{})
	d.SetId("zone-1")

	if err := resourceCloudStackZoneRead(d, cs); err == nil {
		t.Fatal("expected an error from the unreachable endpoint, got nil")
	}
	if d.Id() != "zone-1" {
		t.Fatalf("a transport error must not clear the id, got: %q", d.Id())
	}
}
