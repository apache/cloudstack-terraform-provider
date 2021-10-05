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
	"fmt"
	"log"
	"strings"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCloudStackSSHKeyPair() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackSSHKeyPairCreate,
		Read:   resourceCloudStackSSHKeyPairRead,
		Delete: resourceCloudStackSSHKeyPairDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"public_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"private_key": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCloudStackSSHKeyPairCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	name := d.Get("name").(string)
	publicKey := d.Get("public_key").(string)

	if publicKey != "" {
		// Register supplied key
		p := cs.SSH.NewRegisterSSHKeyPairParams(name, publicKey)

		// If there is a project supplied, we retrieve and set the project id
		if err := setProjectid(p, cs, d); err != nil {
			return err
		}

		_, err := cs.SSH.RegisterSSHKeyPair(p)
		if err != nil {
			return err
		}
	} else {
		// No key supplied, must create one and return the private key
		p := cs.SSH.NewCreateSSHKeyPairParams(name)

		// If there is a project supplied, we retrieve and set the project id
		if err := setProjectid(p, cs, d); err != nil {
			return err
		}

		r, err := cs.SSH.CreateSSHKeyPair(p)
		if err != nil {
			return err
		}
		d.Set("private_key", r.Privatekey)
	}

	log.Printf("[DEBUG] Key pair successfully generated at Cloudstack")
	d.SetId(name)

	return resourceCloudStackSSHKeyPairRead(d, meta)
}

func resourceCloudStackSSHKeyPairRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	log.Printf("[DEBUG] looking for key pair with name %s", d.Id())

	p := cs.SSH.NewListSSHKeyPairsParams()
	p.SetName(d.Id())

	// If there is a project supplied, we retrieve and set the project id
	if err := setProjectid(p, cs, d); err != nil {
		return err
	}

	r, err := cs.SSH.ListSSHKeyPairs(p)
	if err != nil {
		return err
	}
	if r.Count == 0 {
		log.Printf("[DEBUG] Key pair %s does not exist", d.Id())
		d.SetId("")
		return nil
	}

	//SSHKeyPair name is unique in a cloudstack account so dont need to check for multiple
	d.Set("name", r.SSHKeyPairs[0].Name)
	d.Set("fingerprint", r.SSHKeyPairs[0].Fingerprint)

	return nil
}

func resourceCloudStackSSHKeyPairDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.SSH.NewDeleteSSHKeyPairParams(d.Id())

	// If there is a project supplied, we retrieve and set the project id
	if err := setProjectid(p, cs, d); err != nil {
		return err
	}

	// Remove the SSH Keypair
	_, err := cs.SSH.DeleteSSHKeyPair(p)
	if err != nil {
		// This is a very poor way to be told the ID does no longer exist :(
		if strings.Contains(err.Error(), fmt.Sprintf(
			"A key pair with name '%s' does not exist for account", d.Id())) {
			return nil
		}

		return fmt.Errorf("Error deleting key pair: %s", err)
	}

	return nil
}
