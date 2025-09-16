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
			if typeField.Default != "L3" {
				t.Errorf("Type field default should be 'L3', got: %v", typeField.Default)
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
