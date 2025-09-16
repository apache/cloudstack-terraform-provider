package cloudstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestParseCIDR(t *testing.T) {
	networkResource := resourceCloudStackNetwork()
	
	t.Run("L2 network should return empty map", func(t *testing.T) {
		config := map[string]interface{}{
			"type": "L2",
		}
		
		resourceData := schema.TestResourceDataRaw(t, networkResource.Schema, config)
		
		result, err := parseCIDR(resourceData, false)
		if err != nil {
			t.Errorf("Expected no error for L2 network, got: %v", err)
		}
		if result == nil {
			t.Error("Expected non-nil result for L2 network")
		}
		if len(result) != 0 {
			t.Errorf("Expected empty map for L2 network, got: %v", result)
		}
	})
	
	t.Run("L3 network with valid CIDR should parse correctly", func(t *testing.T) {
		config := map[string]interface{}{
			"type": "L3",
			"cidr": "10.0.0.0/16",
		}
		
		resourceData := schema.TestResourceDataRaw(t, networkResource.Schema, config)
		
		result, err := parseCIDR(resourceData, true)
		if err != nil {
			t.Errorf("Expected no error for L3 network with valid CIDR, got: %v", err)
		}
		if result == nil {
			t.Error("Expected non-nil result for L3 network")
		}
		if result["gateway"] == "" {
			t.Error("Expected gateway to be set")
		}
		if result["netmask"] == "" {
			t.Error("Expected netmask to be set")
		}
	})
	
	t.Run("L3 network without CIDR should return error", func(t *testing.T) {
		config := map[string]interface{}{
			"type": "L3",
		}
		
		resourceData := schema.TestResourceDataRaw(t, networkResource.Schema, config)
		
		_, err := parseCIDR(resourceData, true)
		if err == nil {
			t.Error("Expected error for L3 network without CIDR, but got none")
		}
		if err != nil && err.Error() != "cidr is required for L3 networks" {
			t.Errorf("Expected specific error message, got: %v", err)
		}
	})
}
