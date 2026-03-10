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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestParseCIDRv6_DefaultGateway(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceCloudStackNetwork().Schema, map[string]interface{}{
		"ip6cidr": "2001:db8::/64",
	})

	result, err := parseCIDRv6(d, false)
	if err != nil {
		t.Fatalf("parseCIDRv6 failed: %v", err)
	}

	// Default gateway should be network address + 1
	expectedGateway := "2001:db8::1"
	if result["ip6gateway"] != expectedGateway {
		t.Errorf("Expected gateway %s, got %s", expectedGateway, result["ip6gateway"])
	}

	// When specifyiprange is false, startipv6 and endipv6 should not be set
	if _, ok := result["startipv6"]; ok {
		t.Errorf("startipv6 should not be set when specifyiprange is false")
	}
	if _, ok := result["endipv6"]; ok {
		t.Errorf("endipv6 should not be set when specifyiprange is false")
	}
}

func TestParseCIDRv6_CustomGateway(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceCloudStackNetwork().Schema, map[string]interface{}{
		"ip6cidr":    "2001:db8::/64",
		"ip6gateway": "2001:db8::1",
	})

	result, err := parseCIDRv6(d, false)
	if err != nil {
		t.Fatalf("parseCIDRv6 failed: %v", err)
	}

	expectedGateway := "2001:db8::1"
	if result["ip6gateway"] != expectedGateway {
		t.Errorf("Expected gateway %s, got %s", expectedGateway, result["ip6gateway"])
	}
}

func TestParseCIDRv6_WithIPRange(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceCloudStackNetwork().Schema, map[string]interface{}{
		"ip6cidr": "2001:db8::/64",
	})

	result, err := parseCIDRv6(d, true)
	if err != nil {
		t.Fatalf("parseCIDRv6 failed: %v", err)
	}

	// Check gateway (should be network address + 1)
	expectedGateway := "2001:db8::1"
	if result["ip6gateway"] != expectedGateway {
		t.Errorf("Expected gateway %s, got %s", expectedGateway, result["ip6gateway"])
	}

	// Check start IP (should be network address + 2)
	expectedStartIP := "2001:db8::2"
	if result["startipv6"] != expectedStartIP {
		t.Errorf("Expected start IP %s, got %s", expectedStartIP, result["startipv6"])
	}

	// Check end IP (should be the last address in the /64 range)
	expectedEndIP := "2001:db8::ffff:ffff:ffff:ffff"
	if result["endipv6"] != expectedEndIP {
		t.Errorf("Expected end IP %s, got %s", expectedEndIP, result["endipv6"])
	}
}

func TestParseCIDRv6_CustomIPRange(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceCloudStackNetwork().Schema, map[string]interface{}{
		"ip6cidr":   "2001:db8:1::/64",
		"startipv6": "2001:db8:1::100",
		"endipv6":   "2001:db8:1::200",
	})

	result, err := parseCIDRv6(d, true)
	if err != nil {
		t.Fatalf("parseCIDRv6 failed: %v", err)
	}

	// Check that custom values are used
	if result["startipv6"] != "2001:db8:1::100" {
		t.Errorf("Expected custom start IP 2001:db8:1::100, got %s", result["startipv6"])
	}
	if result["endipv6"] != "2001:db8:1::200" {
		t.Errorf("Expected custom end IP 2001:db8:1::200, got %s", result["endipv6"])
	}
}

func TestParseCIDRv6_SmallerPrefix(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceCloudStackNetwork().Schema, map[string]interface{}{
		"ip6cidr": "2001:db8::/48",
	})

	result, err := parseCIDRv6(d, true)
	if err != nil {
		t.Fatalf("parseCIDRv6 failed: %v", err)
	}

	// For a /48, the end IP should have the last 80 bits set to 1
	expectedEndIP := "2001:db8:0:ffff:ffff:ffff:ffff:ffff"
	if result["endipv6"] != expectedEndIP {
		t.Errorf("Expected end IP %s, got %s", expectedEndIP, result["endipv6"])
	}
}

func TestParseCIDRv6_RejectsIPv4(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceCloudStackNetwork().Schema, map[string]interface{}{
		"ip6cidr": "10.0.0.0/24",
	})

	_, err := parseCIDRv6(d, false)
	if err == nil {
		t.Fatal("parseCIDRv6 should reject IPv4 CIDR")
	}

	expectedError := "ip6cidr must be an IPv6 CIDR, got IPv4"
	if !strings.HasPrefix(err.Error(), expectedError) {
		t.Errorf("Expected error message to start with '%s', got '%s'", expectedError, err.Error())
	}
}

func TestParseCIDRv6_Prefix128_NoIPRange(t *testing.T) {
	// /128 is a single address - should fail even without IP range
	d := schema.TestResourceDataRaw(t, resourceCloudStackNetwork().Schema, map[string]interface{}{
		"ip6cidr": "2001:db8::1/128",
	})

	_, err := parseCIDRv6(d, false)
	if err == nil {
		t.Fatal("parseCIDRv6 should reject /128 prefix (single address)")
	}

	expectedError := "ip6cidr prefix /128 is too small"
	if !strings.HasPrefix(err.Error(), expectedError) {
		t.Errorf("Expected error message to start with '%s', got '%s'", expectedError, err.Error())
	}
}

func TestParseCIDRv6_Prefix127_NoIPRange(t *testing.T) {
	// /127 has 2 addresses - should work without IP range (only needs gateway)
	d := schema.TestResourceDataRaw(t, resourceCloudStackNetwork().Schema, map[string]interface{}{
		"ip6cidr": "2001:db8::/127",
	})

	result, err := parseCIDRv6(d, false)
	if err != nil {
		t.Fatalf("parseCIDRv6 should accept /127 prefix without IP range: %v", err)
	}

	// Should have gateway
	if _, ok := result["ip6gateway"]; !ok {
		t.Error("Expected ip6gateway to be set")
	}

	// Should not have start/end IP
	if _, ok := result["startipv6"]; ok {
		t.Error("startipv6 should not be set when specifyiprange is false")
	}
}

func TestParseCIDRv6_Prefix127_WithIPRange(t *testing.T) {
	// /127 has only 2 addresses - should fail with IP range (needs 3+ addresses)
	d := schema.TestResourceDataRaw(t, resourceCloudStackNetwork().Schema, map[string]interface{}{
		"ip6cidr": "2001:db8::/127",
	})

	_, err := parseCIDRv6(d, true)
	if err == nil {
		t.Fatal("parseCIDRv6 should reject /127 prefix with IP range (only 2 addresses)")
	}

	expectedError := "ip6cidr prefix /127 is too small for automatic IP range generation"
	if !strings.HasPrefix(err.Error(), expectedError) {
		t.Errorf("Expected error message to start with '%s', got '%s'", expectedError, err.Error())
	}
}

func TestParseCIDRv6_Prefix126_WithIPRange(t *testing.T) {
	// /126 has 4 addresses - should work with IP range
	d := schema.TestResourceDataRaw(t, resourceCloudStackNetwork().Schema, map[string]interface{}{
		"ip6cidr": "2001:db8::/126",
	})

	result, err := parseCIDRv6(d, true)
	if err != nil {
		t.Fatalf("parseCIDRv6 should accept /126 prefix with IP range: %v", err)
	}

	// Should have gateway, start, and end
	if _, ok := result["ip6gateway"]; !ok {
		t.Error("Expected ip6gateway to be set")
	}
	if _, ok := result["startipv6"]; !ok {
		t.Error("Expected startipv6 to be set")
	}
	if _, ok := result["endipv6"]; !ok {
		t.Error("Expected endipv6 to be set")
	}

	// Verify the end IP is correct for /126 (last 2 bits set to 1)
	expectedEndIP := "2001:db8::3"
	if result["endipv6"] != expectedEndIP {
		t.Errorf("Expected end IP %s for /126, got %s", expectedEndIP, result["endipv6"])
	}
}
