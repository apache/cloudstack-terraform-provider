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
	"strconv"
	"strings"
	"testing"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

var testAccMuxProvider map[string]func() (tfprotov6.ProviderServer, error)

var cloudStackTemplateURL = os.Getenv("CLOUDSTACK_TEMPLATE_URL")

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"cloudstack": testAccProvider,
	}

	testAccMuxProvider = map[string]func() (tfprotov6.ProviderServer, error){
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
		ProtoV6ProviderFactories: testAccMuxProvider,
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

// parseCloudStackVersion parses a CloudStack version string (e.g., "4.22.0.0")
// and returns a numeric value for comparison (e.g., 4.22 -> 4022).
// The numeric value is calculated as: major * 1000 + minor.
// Returns 0 if the version string cannot be parsed.
func parseCloudStackVersion(version string) int {
	parts := strings.Split(version, ".")
	if len(parts) < 2 {
		return 0
	}

	major := 0
	minor := 0

	// Parse major version - extract first numeric part
	majorStr := regexp.MustCompile(`^\d+`).FindString(parts[0])
	if majorStr != "" {
		major, _ = strconv.Atoi(majorStr)
	}

	// Parse minor version - extract first numeric part
	minorStr := regexp.MustCompile(`^\d+`).FindString(parts[1])
	if minorStr != "" {
		minor, _ = strconv.Atoi(minorStr)
	}

	return major*1000 + minor
}

// requireMinimumCloudStackVersion checks if the CloudStack version meets the minimum requirement.
// If the version is below the minimum, it skips the test with an appropriate message.
// The minVersion parameter should be in the format returned by parseCloudStackVersion (e.g., 4022 for 4.22.0).
func requireMinimumCloudStackVersion(t *testing.T, minVersion int, featureName string) {
	t.Helper()
	version := getCloudStackVersion(t)
	if version == "" {
		t.Skipf("Unable to determine CloudStack version, skipping %s test", featureName)
		return
	}

	versionNum := parseCloudStackVersion(version)
	if versionNum < minVersion {
		// Convert minVersion back to readable format (e.g., 4022 -> "4.22")
		major := minVersion / 1000
		minor := minVersion % 1000
		t.Skipf("%s not supported in CloudStack version %s (requires %d.%d+)", featureName, version, major, minor)
	}
}

// testAccPreCheckStaticRouteNexthop checks if the CloudStack version supports
// the nexthop parameter for static routes (requires 4.22.0+)
func testAccPreCheckStaticRouteNexthop(t *testing.T) {
	testAccPreCheck(t)

	const minVersionNum = 4022 // 4.22.0
	requireMinimumCloudStackVersion(t, minVersionNum, "Static route nexthop parameter")
}

// newTestClient creates a CloudStack client from environment variables for use in test PreCheck functions.
// This is needed because PreCheck functions run before the test framework configures the provider,
// so testAccProvider.Meta() is nil at that point.
func newTestClient(t *testing.T) *cloudstack.CloudStackClient {
	t.Helper()
	testAccPreCheck(t)

	cfg := Config{
		APIURL:      os.Getenv("CLOUDSTACK_API_URL"),
		APIKey:      os.Getenv("CLOUDSTACK_API_KEY"),
		SecretKey:   os.Getenv("CLOUDSTACK_SECRET_KEY"),
		HTTPGETOnly: true,
		Timeout:     60,
	}
	cs, err := cfg.NewClient()
	if err != nil {
		t.Fatalf("Failed to create CloudStack client: %v", err)
	}
	return cs
}

// getCloudStackVersion returns the CloudStack version from the API
func getCloudStackVersion(t *testing.T) string {
	t.Helper()
	cs := newTestClient(t)

	p := cs.Configuration.NewListCapabilitiesParams()
	r, err := cs.Configuration.ListCapabilities(p)
	if err != nil {
		t.Fatalf("Failed to get CloudStack capabilities: %v", err)
	}

	if r.Capabilities != nil {
		return r.Capabilities.Cloudstackversion
	}

	return ""
}
