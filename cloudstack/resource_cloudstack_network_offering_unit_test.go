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
	"testing"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// domain_id is a schema.TypeList, which the SDK returns as []interface{}, so
// asserting v.([]string) on create panics when domain_id is set. The create
// function must reach the API call (and return its error) rather than panic.
func TestNetworkOfferingDomainIdDoesNotPanic(t *testing.T) {
	cs := cloudstack.NewClient("http://127.0.0.1:1", "key", "secret", false)

	d := schema.TestResourceDataRaw(t, resourceCloudStackNetworkOffering().Schema, map[string]interface{}{
		"name":          "no-1",
		"display_text":  "no-1",
		"guest_ip_type": "Isolated",
		"traffic_type":  "Guest",
		"domain_id":     []interface{}{"domain-1", "domain-2"},
	})

	err := resourceCloudStackNetworkOfferingCreate(d, cs)
	if err == nil {
		t.Fatal("expected an API error from the unreachable endpoint, got nil")
	}
}
