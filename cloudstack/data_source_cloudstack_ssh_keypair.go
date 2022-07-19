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
	"log"
	"regexp"
	"strings"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceCloudstackSSHKeyPair() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudstackSSHKeyPairRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),

			//Computed values
			"fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCloudstackSSHKeyPairRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	p := cs.SSH.NewListSSHKeyPairsParams()
	csSshKeyPairs, err := cs.SSH.ListSSHKeyPairs(p)

	if err != nil {
		return fmt.Errorf("Failed to list ssh key pairs: %s", err)
	}
	filters := d.Get("filter")
	var sshKeyPair *cloudstack.SSHKeyPair

	for _, k := range csSshKeyPairs.SSHKeyPairs {
		match, err := applySshKeyPairsFilters(k, filters.(*schema.Set))
		if err != nil {
			return err
		}
		if match {
			sshKeyPair = k
		}
	}

	if sshKeyPair == nil {
		return fmt.Errorf("No ssh key pair is matching with the specified regex")
	}
	log.Printf("[DEBUG] Selected ssh key pair: %s\n", sshKeyPair.Name)

	return sshKeyPairDescriptionAttributes(d, sshKeyPair)
}

func sshKeyPairDescriptionAttributes(d *schema.ResourceData, sshKeyPair *cloudstack.SSHKeyPair) error {
	d.SetId(sshKeyPair.Name)
	d.Set("fingerprint", sshKeyPair.Fingerprint)
	d.Set("name", sshKeyPair.Name)

	return nil
}

func applySshKeyPairsFilters(sshKeyPair *cloudstack.SSHKeyPair, filters *schema.Set) (bool, error) {
	var sshKeyPairJSON map[string]interface{}
	k, _ := json.Marshal(sshKeyPair)
	err := json.Unmarshal(k, &sshKeyPairJSON)
	if err != nil {
		return false, err
	}

	for _, f := range filters.List() {
		m := f.(map[string]interface{})
		r, err := regexp.Compile(m["value"].(string))
		if err != nil {
			return false, fmt.Errorf("Invalid regex: %s", err)
		}
		updatedName := strings.ReplaceAll(m["name"].(string), "_", "")
		sshKeyPairField := sshKeyPairJSON[updatedName].(string)
		if !r.MatchString(sshKeyPairField) {
			return false, nil
		}

	}
	return true, nil
}
