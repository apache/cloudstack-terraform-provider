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
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

var cloudStackTemplateURL = os.Getenv("CLOUDSTACK_TEMPLATE_URL")

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"cloudstack": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func TestMuxServer(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"cloudstack": func() (tfprotov6.ProviderServer, error) {
				ctx := context.Background()

				upgradedSdkServer, err := tf5to6server.UpgradeServer(
					ctx,
					Provider().GRPCProvider,
				)

				if err != nil {
					return nil, err
				}

				providers := []func() tfprotov6.ProviderServer{
					providerserver.NewProtocol6(New()),
					func() tfprotov6.ProviderServer {
						return upgradedSdkServer
					},
				}

				muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)

				if err != nil {
					return nil, err
				}

				return muxServer.ProviderServer(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config:      testMuxServerConfig_conflict,
				ExpectError: regexp.MustCompile("Invalid Attribute Combination"),
			},
			{
				Config: testMuxServerConfig_basic,
			},
		},
	})
}

const testMuxServerConfig_basic = `
resource "cloudstack_zone" "zone_resource"{
	name       		= "TestZone1"
  	dns1       		= "8.8.8.8"
  	internal_dns1  	=  "172.20.0.1"
  	network_type   	=  "Advanced"
  }

  data "cloudstack_zone" "zone_data_source"{
    filter{
    	name 	=	"name"
    	value	=	cloudstack_zone.zone_resource.name
    }
  }
  `

const testMuxServerConfig_conflict = `
provider "cloudstack" {
	api_url = "http://localhost:8080/client/api"
	api_key = "xxxxx"
	secret_key = "xxxxx"
	config = "cloudstack.ini"
}

data "cloudstack_zone" "zone_data_source"{
    filter{
    	name 	=	"name"
    	value	=	"test"
    }
}
  `

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("CLOUDSTACK_API_URL"); v == "" {
		t.Fatal("CLOUDSTACK_API_URL must be set for acceptance tests")
	}
	if v := os.Getenv("CLOUDSTACK_API_KEY"); v == "" {
		t.Fatal("CLOUDSTACK_API_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("CLOUDSTACK_SECRET_KEY"); v == "" {
		t.Fatal("CLOUDSTACK_SECRET_KEY must be set for acceptance tests")
	}
}
