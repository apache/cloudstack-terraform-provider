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

package cloudstack

import (
	"testing"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dsPanicFilterSet(pairs ...[2]string) *schema.Set {
	hash := func(i interface{}) int {
		m := i.(map[string]interface{})
		return schema.HashString(m["name"].(string) + "|" + m["value"].(string))
	}
	items := make([]interface{}, 0, len(pairs))
	for _, p := range pairs {
		items = append(items, map[string]interface{}{"name": p[0], "value": p[1]})
	}
	return schema.NewSet(hash, items)
}

// A data-source filter referencing an unknown field name, or a non-string field, must not panic.
func TestApplyVolumeFiltersDoesNotPanicOnUnknownOrNonStringField(t *testing.T) {
	vol := &cloudstack.Volume{Name: "vol-a", Size: 5368709120}

	// Unknown filter field: the JSON lookup returns nil; must not panic and must not match.
	if match, err := applyVolumeFilters(vol, dsPanicFilterSet([2]string{"no_such_field", "x"})); err != nil {
		t.Fatalf("unexpected error: %s", err)
	} else if match {
		t.Errorf("an unknown filter field should not match")
	}

	// Non-string (numeric) field: the JSON value is a float64; must not panic.
	if _, err := applyVolumeFilters(vol, dsPanicFilterSet([2]string{"size", "anything"})); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}
