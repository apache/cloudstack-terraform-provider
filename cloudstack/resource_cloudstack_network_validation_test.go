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
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestResourceCloudStackNetworkSchema(t *testing.T) {
	networkResource := resourceCloudStackNetwork()

	// Test that required fields exist
	t.Run("Schema should have type field", func(t *testing.T) {
		if typeField, ok := networkResource.Schema["type"]; !ok {
			t.Error("Schema should have 'type' field")
		} else {
			if typeField.Type != schema.TypeString {
				t.Errorf("Type field should be TypeString, got: %v", typeField.Type)
			}
			if typeField.Required {
				t.Error("Type field should not be required")
			}
			if typeField.Optional != true {
				t.Error("Type field should be optional")
			}
			if typeField.Computed != true {
				t.Error("Type field should be computed")
			}
			if typeField.Default != nil {
				t.Errorf("Type field should not have a default (it's computed), got: %v", typeField.Default)
			}
		}
	})

	t.Run("Schema should have cidr field as optional", func(t *testing.T) {
		if cidrField, ok := networkResource.Schema["cidr"]; !ok {
			t.Error("Schema should have 'cidr' field")
		} else {
			if cidrField.Required {
				t.Error("CIDR field should not be required")
			}
			if cidrField.Optional != true {
				t.Error("CIDR field should be optional")
			}
		}
	})

	t.Run("Schema should have CustomizeDiff", func(t *testing.T) {
		if networkResource.CustomizeDiff == nil {
			t.Error("Resource should have CustomizeDiff function")
		}
	})
}

func diffNetwork(t *testing.T, state *terraform.InstanceState, raw map[string]interface{}) (*terraform.InstanceDiff, error) {
	t.Helper()
	r := resourceCloudStackNetwork()
	c := terraform.NewResourceConfigRaw(raw)
	return r.Diff(context.Background(), state, c, nil)
}

func TestResourceCloudStackNetworkTypeCidrDiff(t *testing.T) {
	t.Run("new L2 network with explicit type and no cidr plans cleanly", func(t *testing.T) {
		d, err := diffNetwork(t, nil, map[string]interface{}{
			"name":             "terraform-l2-network",
			"type":             "L2",
			"network_offering": "DefaultL2NetworkOffering",
			"zone":             "Sandbox-simulator",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if d == nil {
			t.Fatal("expected a diff for new resource")
		}
		if attr, ok := d.Attributes["type"]; !ok || attr.New != "L2" {
			t.Fatalf("expected planned type L2, got %+v", d.Attributes["type"])
		}
	})

	t.Run("new isolated network without cidr and without type plans cleanly", func(t *testing.T) {
		d, err := diffNetwork(t, nil, map[string]interface{}{
			"name":             "terraform-isolated-no-cidr",
			"network_offering": "DefaultIsolatedNetworkOfferingWithSourceNatService",
			"zone":             "Sandbox-simulator",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// type should stay computed so the value from the API is accepted after apply
		if attr, ok := d.Attributes["type"]; ok && !attr.NewComputed && attr.New != "" {
			t.Fatalf("type should remain computed, got %+v", attr)
		}
		if attr, ok := d.Attributes["cidr"]; ok && !attr.NewComputed && attr.New != "" {
			t.Fatalf("cidr should remain computed, got %+v", attr)
		}
	})

	t.Run("new L3 network with cidr plans cleanly", func(t *testing.T) {
		_, err := diffNetwork(t, nil, map[string]interface{}{
			"name":             "terraform-network",
			"cidr":             "10.1.1.0/24",
			"network_offering": "DefaultIsolatedNetworkOfferingWithSourceNatService",
			"zone":             "Sandbox-simulator",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("L2 with cidr is rejected", func(t *testing.T) {
		_, err := diffNetwork(t, nil, map[string]interface{}{
			"name":             "terraform-l2-network",
			"type":             "L2",
			"cidr":             "10.1.1.0/24",
			"network_offering": "DefaultL2NetworkOffering",
			"zone":             "Sandbox-simulator",
		})
		if err == nil {
			t.Fatal("expected error for L2 network with cidr")
		}
	})

	t.Run("explicit L3 without cidr is rejected", func(t *testing.T) {
		_, err := diffNetwork(t, nil, map[string]interface{}{
			"name":             "terraform-network",
			"type":             "L3",
			"network_offering": "DefaultIsolatedNetworkOfferingWithSourceNatService",
			"zone":             "Sandbox-simulator",
		})
		if err == nil {
			t.Fatal("expected error for explicit L3 network without cidr")
		}
	})

	t.Run("existing isolated network created without cidr shows no drift", func(t *testing.T) {
		// state as it looks after creating an isolated network without cidr
		state := &terraform.InstanceState{
			ID: "net-1",
			Attributes: map[string]string{
				"id":                        "net-1",
				"name":                      "terraform-isolated-no-cidr",
				"display_text":              "terraform-isolated-no-cidr",
				"type":                      "L3",
				"cidr":                      "10.1.1.0/24",
				"gateway":                   "10.1.1.1",
				"ip6cidr":                   "",
				"ip6gateway":                "",
				"startip":                   "",
				"endip":                     "",
				"startipv6":                 "",
				"endipv6":                   "",
				"network_domain":            "cs1cloud.internal",
				"network_offering":          "DefaultIsolatedNetworkOfferingWithSourceNatService",
				"vpc_id":                    "",
				"acl_id":                    "none",
				"project":                   "",
				"source_nat_ip":             "false",
				"source_nat_ip_address":     "",
				"source_nat_ip_id":          "",
				"zone":                      "Sandbox-simulator",
				"bypass_vlan_overlap_check": "false",
				"tags.%":                    "0",
			},
		}
		d, err := diffNetwork(t, state, map[string]interface{}{
			"name":             "terraform-isolated-no-cidr",
			"network_offering": "DefaultIsolatedNetworkOfferingWithSourceNatService",
			"zone":             "Sandbox-simulator",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if d != nil && len(d.Attributes) > 0 {
			t.Fatalf("expected empty diff, got %+v", d.Attributes)
		}
	})

	t.Run("existing L2 network shows no drift", func(t *testing.T) {
		state := &terraform.InstanceState{
			ID: "net-2",
			Attributes: map[string]string{
				"id":                        "net-2",
				"name":                      "terraform-l2-network",
				"display_text":              "terraform-l2-network",
				"type":                      "L2",
				"cidr":                      "",
				"gateway":                   "",
				"ip6cidr":                   "",
				"ip6gateway":                "",
				"startip":                   "",
				"endip":                     "",
				"startipv6":                 "",
				"endipv6":                   "",
				"network_domain":            "cs1cloud.internal",
				"network_offering":          "DefaultL2NetworkOffering",
				"vpc_id":                    "",
				"acl_id":                    "none",
				"project":                   "",
				"source_nat_ip":             "false",
				"source_nat_ip_address":     "",
				"source_nat_ip_id":          "",
				"zone":                      "Sandbox-simulator",
				"bypass_vlan_overlap_check": "false",
				"tags.%":                    "0",
			},
		}
		d, err := diffNetwork(t, state, map[string]interface{}{
			"name":             "terraform-l2-network",
			"type":             "L2",
			"network_offering": "DefaultL2NetworkOffering",
			"zone":             "Sandbox-simulator",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if d != nil && len(d.Attributes) > 0 {
			t.Fatalf("expected empty diff, got %+v", d.Attributes)
		}
	})
}
