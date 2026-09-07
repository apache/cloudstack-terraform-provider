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

// For a managed egress firewall, Read collects unknown rules into a dummy rule.
// cidr_list is a schema.TypeSet, so its value must be a *schema.Set. The dummy
// rule stored the uuid string there instead, which corrupts the rule set and
// later panics when createEgressFirewallRule asserts cidr_list as *schema.Set.
func TestEgressFirewallReadManagedDummyRuleCidrListIsSet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"listegressfirewallrulesresponse":{"count":1,"firewallrule":[{"id":"rule-1","protocol":"tcp","cidrlist":"10.0.0.0/8","startport":80,"endport":80}]}}`))
	}))
	defer server.Close()

	cs := cloudstack.NewClient(server.URL, "key", "secret", false)

	d := schema.TestResourceDataRaw(t, resourceCloudStackEgressFirewall().Schema, map[string]interface{}{
		"managed": true,
	})
	d.SetId("net-1")

	if err := resourceCloudStackEgressFirewallRead(d, cs); err != nil {
		t.Fatalf("read of a managed egress firewall should not error, got: %s", err)
	}

	rules := d.Get("rule").(*schema.Set)
	if rules.Len() != 1 {
		t.Fatalf("expected one managed dummy rule, got %d", rules.Len())
	}

	for _, raw := range rules.List() {
		rule := raw.(map[string]interface{})
		if _, ok := rule["cidr_list"].(*schema.Set); !ok {
			t.Fatalf("dummy rule cidr_list must be a *schema.Set, got %T", rule["cidr_list"])
		}
	}
}
