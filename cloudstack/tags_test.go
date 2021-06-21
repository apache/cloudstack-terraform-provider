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
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestDiffTags(t *testing.T) {
	cases := []struct {
		Old, New       map[string]interface{}
		Create, Remove map[string]string
	}{
		// Basic add/remove
		{
			Old: map[string]interface{}{
				"foo": "bar",
			},
			New: map[string]interface{}{
				"bar": "baz",
			},
			Create: map[string]string{
				"bar": "baz",
			},
			Remove: map[string]string{
				"foo": "bar",
			},
		},

		// Modify
		{
			Old: map[string]interface{}{
				"foo": "bar",
			},
			New: map[string]interface{}{
				"foo": "baz",
			},
			Create: map[string]string{
				"foo": "baz",
			},
			Remove: map[string]string{
				"foo": "bar",
			},
		},
	}

	for i, tc := range cases {
		r, c := diffTags(tagsFromSchema(tc.Old), tagsFromSchema(tc.New))
		if !reflect.DeepEqual(r, tc.Remove) {
			t.Fatalf("%d: bad remove: %#v", i, r)
		}
		if !reflect.DeepEqual(c, tc.Create) {
			t.Fatalf("%d: bad create: %#v", i, c)
		}
	}
}

// testAccCheckResourceTags is an helper to test tags creation on any resource.
func testAccCheckResourceTags(
	n interface{}) resource.TestCheckFunc {
	res := struct {
		Tags []struct {
			Key   string `json:"key,omitempty"`
			Value string `json:"value,omitempty"`
		} `json:"tags,omitempty"`
	}{}
	return func(s *terraform.State) error {
		b, _ := json.Marshal(n)
		json.Unmarshal(b, &res)
		tags := make(map[string]string)
		for _, tag := range res.Tags {
			tags[tag.Key] = tag.Value
		}
		return testAccCheckTags(tags, "terraform-tag", "true")
	}
}

// testAccCheckTags can be used to check the tags on a resource.
func testAccCheckTags(tags map[string]string, key string, value string) error {
	v, ok := tags[key]
	if !ok {
		return fmt.Errorf("Missing tag: %s", key)
	}

	if v != value {
		return fmt.Errorf("%s: bad value: %s", key, v)
	}

	return nil
}
