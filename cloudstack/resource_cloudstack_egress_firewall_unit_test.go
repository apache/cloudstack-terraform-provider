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

func TestIsAllPortsTCPUDP(t *testing.T) {
	tests := []struct {
		protocol string
		start    int
		end      int
		expected bool
		name     string
	}{
		{"tcp", 0, 0, true, "TCP with 0/0"},
		{"TCP", 0, 0, true, "TCP uppercase with 0/0"},
		{"udp", -1, -1, true, "UDP with -1/-1"},
		{"UDP", -1, -1, true, "UDP uppercase with -1/-1"},
		{"tcp", 1, 65535, true, "TCP with 1/65535"},
		{"udp", 1, 65535, true, "UDP with 1/65535"},
		{"tcp", 80, 80, false, "TCP with specific port"},
		{"udp", 53, 53, false, "UDP with specific port"},
		{"icmp", 0, 0, false, "ICMP protocol"},
		{"all", 0, 0, false, "ALL protocol"},
		{"tcp", 1, 1000, false, "TCP with port range"},
		{"tcp", 0, 1, false, "TCP with 0/1"},
		{"tcp", -1, 0, false, "TCP with -1/0"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := isAllPortsTCPUDP(test.protocol, test.start, test.end)
			if result != test.expected {
				t.Errorf("isAllPortsTCPUDP(%q, %d, %d) = %v, expected %v",
					test.protocol, test.start, test.end, result, test.expected)
			}
		})
	}
}

func TestNormalizeRemoteCIDRs(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
		name     string
	}{
		{"", []string{}, "empty string"},
		{"10.0.0.0/8", []string{"10.0.0.0/8"}, "single CIDR"},
		{"10.0.0.0/8,192.168.1.0/24", []string{"10.0.0.0/8", "192.168.1.0/24"}, "two CIDRs"},
		{"10.0.0.0/8, 192.168.1.0/24", []string{"10.0.0.0/8", "192.168.1.0/24"}, "two CIDRs with space"},
		{" 10.0.0.0/8 , 192.168.1.0/24 ", []string{"10.0.0.0/8", "192.168.1.0/24"}, "CIDRs with extra spaces"},
		{"192.168.1.0/24,10.0.0.0/8", []string{"10.0.0.0/8", "192.168.1.0/24"}, "unsorted CIDRs (should be sorted)"},
		{"10.0.0.0/8,,192.168.1.0/24", []string{"10.0.0.0/8", "192.168.1.0/24"}, "empty CIDR in middle"},
		{" , , ", []string{}, "only commas and spaces"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := normalizeRemoteCIDRs(test.input)
			if len(result) != len(test.expected) {
				t.Errorf("normalizeRemoteCIDRs(%q) length = %d, expected %d",
					test.input, len(result), len(test.expected))
				return
			}
			for i, v := range result {
				if v != test.expected[i] {
					t.Errorf("normalizeRemoteCIDRs(%q)[%d] = %q, expected %q",
						test.input, i, v, test.expected[i])
				}
			}
		})
	}
}

func TestNormalizeLocalCIDRs(t *testing.T) {
	tests := []struct {
		input    *schema.Set
		expected []string
		name     string
	}{
		{nil, []string{}, "nil set"},
		{schema.NewSet(schema.HashString, []interface{}{}), []string{}, "empty set"},
		{schema.NewSet(schema.HashString, []interface{}{"10.0.0.0/8"}), []string{"10.0.0.0/8"}, "single CIDR"},
		{schema.NewSet(schema.HashString, []interface{}{"10.0.0.0/8", "192.168.1.0/24"}), []string{"10.0.0.0/8", "192.168.1.0/24"}, "two CIDRs"},
		{schema.NewSet(schema.HashString, []interface{}{"192.168.1.0/24", "10.0.0.0/8"}), []string{"10.0.0.0/8", "192.168.1.0/24"}, "unsorted CIDRs"},
		{schema.NewSet(schema.HashString, []interface{}{" 10.0.0.0/8 ", " 192.168.1.0/24 "}), []string{"10.0.0.0/8", "192.168.1.0/24"}, "CIDRs with spaces"},
		{schema.NewSet(schema.HashString, []interface{}{"10.0.0.0/8", "", "192.168.1.0/24"}), []string{"10.0.0.0/8", "192.168.1.0/24"}, "with empty string"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := normalizeLocalCIDRs(test.input)
			if len(result) != len(test.expected) {
				t.Errorf("normalizeLocalCIDRs() length = %d, expected %d",
					len(result), len(test.expected))
				return
			}
			for i, v := range result {
				if v != test.expected[i] {
					t.Errorf("normalizeLocalCIDRs()[%d] = %q, expected %q",
						i, v, test.expected[i])
				}
			}
		})
	}
}

func TestCidrSetsEqual(t *testing.T) {
	tests := []struct {
		remote   string
		local    *schema.Set
		expected bool
		name     string
	}{
		{"", schema.NewSet(schema.HashString, []interface{}{}), true, "both empty"},
		{"10.0.0.0/8", schema.NewSet(schema.HashString, []interface{}{"10.0.0.0/8"}), true, "single matching CIDR"},
		{"10.0.0.0/8,192.168.1.0/24", schema.NewSet(schema.HashString, []interface{}{"192.168.1.0/24", "10.0.0.0/8"}), true, "multiple CIDRs different order"},
		{"10.0.0.0/8, 192.168.1.0/24", schema.NewSet(schema.HashString, []interface{}{"10.0.0.0/8", "192.168.1.0/24"}), true, "remote with spaces"},
		{"10.0.0.0/8", schema.NewSet(schema.HashString, []interface{}{"192.168.1.0/24"}), false, "different CIDRs"},
		{"10.0.0.0/8,192.168.1.0/24", schema.NewSet(schema.HashString, []interface{}{"10.0.0.0/8"}), false, "different count"},
		{"", schema.NewSet(schema.HashString, []interface{}{"10.0.0.0/8"}), false, "remote empty, local not"},
		{"10.0.0.0/8", schema.NewSet(schema.HashString, []interface{}{}), false, "local empty, remote not"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := cidrSetsEqual(test.remote, test.local)
			if result != test.expected {
				t.Errorf("cidrSetsEqual(%q, %v) = %v, expected %v",
					test.remote, test.local.List(), result, test.expected)
			}
		})
	}
}
