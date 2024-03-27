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
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

var cloudStackTemplateURL = os.Getenv("CLOUDSTACK_TEMPLATE_URL")

func init() {
	testAccProvider = New()
	testAccProviders = map[string]*schema.Provider{
		"cloudstack": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := New().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = New()
}

func TestMuxServer(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: map[string]func() (tfprotov5.ProviderServer, error){
			"cloudstack": func() (tfprotov5.ProviderServer, error) {
				ctx := context.Background()
				providers := []func() tfprotov5.ProviderServer{
					New().GRPCProvider,
				}

				muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)

				if err != nil {
					return nil, err
				}

				return muxServer.ProviderServer(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testMuxServerConfig_basic,
			},
		},
	})
}

const testMuxServerConfig_basic = `
resource "cloudstack_zone" "zone_resource"{
	name       		= "TestZone"
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
