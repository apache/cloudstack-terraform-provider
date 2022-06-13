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
	"testing"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCloudStackKubernetesVersion_basic(t *testing.T) {
	var version cloudstack.KubernetesSupportedVersion

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackKubernetesVersionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackKubernetesVersion_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackKubernetesVersionExists("cloudstack_kubernetes_version.foo", &version),
					testAccCheckCloudStackKubernetesVersionAttributes(&version),
				),
			},
		},
	})
}

func TestAccCloudStackKubernetesVersion_update(t *testing.T) {
	var version cloudstack.KubernetesSupportedVersion

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackKubernetesVersionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackKubernetesVersion_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackKubernetesVersionExists("cloudstack_kubernetes_version.foo", &version),
					testAccCheckCloudStackKubernetesVersionAttributes(&version),
					resource.TestCheckResourceAttr(
						"cloudstack_kubernetes_version.foo", "state", "Enabled"),
				),
			},

			{
				Config: testAccCloudStackKubernetesVersion_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackKubernetesVersionExists("cloudstack_kubernetes_version.foo", &version),
					testAccCheckCloudStackKubernetesVersionAttributes(&version),
					resource.TestCheckResourceAttr(
						"cloudstack_kubernetes_version.foo", "state", "Disabled"),
				),
			},
		},
	})
}

func testAccCheckCloudStackKubernetesVersionExists(
	n string, version *cloudstack.KubernetesSupportedVersion) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No kubernetes version ID is set")
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		ver, _, err := cs.Kubernetes.GetKubernetesSupportedVersionByID(rs.Primary.ID)

		if err != nil {
			return err
		}

		if ver.Id != rs.Primary.ID {
			return fmt.Errorf("Kubernetes Version not found")
		}

		*version = *ver

		return nil
	}
}

func testAccCheckCloudStackKubernetesVersionAttributes(
	kubernetesVersion *cloudstack.KubernetesSupportedVersion) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if kubernetesVersion.Semanticversion != "1.23.3" {
			return fmt.Errorf("Bad semantic version: %s", kubernetesVersion.Name)
		}

		if kubernetesVersion.Mincpunumber != 2 {
			return fmt.Errorf("Bad min cpu: %d", kubernetesVersion.Mincpunumber)
		}

		if kubernetesVersion.Minmemory != 2048 {
			return fmt.Errorf("Bad min memory: %d", kubernetesVersion.Minmemory)
		}

		return nil
	}
}

func testAccCheckCloudStackKubernetesVersionDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudstack_kubernetes_version" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No kubernetes version ID is set")
		}

		_, _, err := cs.Kubernetes.GetKubernetesSupportedVersionByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Kubernetes Version %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

const testAccCloudStackKubernetesVersion_basic = `
resource "cloudstack_kubernetes_version" "foo" {
  semantic_version      = "1.23.3"
  url                   = "http://download.cloudstack.org/cks/setup-1.23.3.iso"
  min_cpu               = 2
  min_memory            = 2048
}`

const testAccCloudStackKubernetesVersion_update = `
resource "cloudstack_kubernetes_version" "foo" {
  semantic_version      = "1.23.3"
  url                   = "http://download.cloudstack.org/cks/setup-1.23.3.iso"
  min_cpu               = 2
  min_memory            = 2048
  state                 = "Disabled"
}`
