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
	if err.Error()[:len(expectedError)] != expectedError {
		t.Errorf("Expected error message to start with '%s', got '%s'", expectedError, err.Error())
	}
}
