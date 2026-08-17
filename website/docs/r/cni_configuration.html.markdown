---
layout: default
page_title: "CloudStack: cloudstack_cni_configuration"
sidebar_current: "docs-cloudstack-resource-cni_configuration"
description: |-
    Creates and manages a CloudStack CNI (Container Network Interface) configuration
---

# CloudStack: cloudstack_cni_configuration

A `cloudstack_cni_configuration` resource manages a Container Network Interface (CNI) configuration for CloudStack Kubernetes Service (CKS) clusters. CNI configurations define how network connectivity is provided to Kubernetes pods.

## Example Usage

### Basic Calico CNI Configuration

```hcl
resource "cloudstack_cni_configuration" "calico" {
  name       = "calico-cni-config"
  cni_config = base64encode(jsonencode({
    "name"       = "k8s-pod-network",
    "cniVersion" = "0.3.1",
    "plugins" = [
      {
        "type"           = "calico",
        "log_level"      = "info",
        "datastore_type" = "kubernetes",
        "nodename"       = "KUBERNETES_NODE_NAME",
        "mtu"            = "CNI_MTU",
        "ipam" = {
          "type" = "calico-ipam"
        },
        "policy" = {
          "type" = "k8s"
        },
        "kubernetes" = {
          "kubeconfig" = "KUBECONFIG_FILEPATH"
        }
      },
      {
        "type" = "portmap",
        "snat" = true,
        "capabilities" = { "portMappings" = true }
      }
    ]
  }))
  
  params = [
    "KUBERNETES_NODE_NAME",
    "CNI_MTU",
    "KUBECONFIG_FILEPATH"
  ]
}
```

### Flannel CNI Configuration

```hcl
resource "cloudstack_cni_configuration" "flannel" {
  name       = "flannel-cni-config"
  cni_config = base64encode(jsonencode({
    "name"       = "cbr0",
    "cniVersion" = "0.3.1",
    "plugins" = [
      {
        "type"   = "flannel",
        "delegate" = {
          "hairpinMode" = true,
          "isDefaultGateway" = true
        }
      },
      {
        "type" = "portmap",
        "capabilities" = {
          "portMappings" = true
        }
      }
    ]
  }))
  
  params = ["FLANNEL_NETWORK", "FLANNEL_SUBNET"]
  
  domain_id = "domain-uuid"
  account   = "admin"
}
```

## Argument Reference

The following arguments are supported:

### Required Arguments

* `name` - (Required) The name of the CNI configuration. Must be unique within the account/domain.
* `cni_config` - (Required) The CNI configuration in base64-encoded JSON format. This should contain the complete CNI plugin configuration according to the CNI specification.

### Optional Arguments

* `params` - (Optional) A list of parameter names that can be substituted in the CNI configuration. These parameters can be provided with actual values when creating a Kubernetes cluster using `cni_config_details`.
* `domain_id` - (Optional) The domain ID for the CNI configuration. If not specified, uses the default domain.
* `account` - (Optional) The account name for the CNI configuration. If not specified, uses the account of the authenticated user.
* `project_id` - (Optional) The project ID to assign the CNI configuration to.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the CNI configuration.
* `created` - The timestamp when the CNI configuration was created.
* `domain` - The domain name where the CNI configuration belongs.
* `project` - The project name if the CNI configuration is assigned to a project.

## CNI Configuration Format

The `cni_config` should be a base64-encoded JSON string that follows the CNI specification. The configuration supports parameter substitution using placeholder names that can be defined in the `params` list.

### Parameter Substitution

Parameters in the CNI configuration can be specified as placeholders and will be replaced with actual values when the configuration is used in a Kubernetes cluster:

```json
{
  "name": "k8s-pod-network",
  "cniVersion": "0.3.1",
  "plugins": [
    {
      "type": "calico",
      "nodename": "KUBERNETES_NODE_NAME",
      "mtu": "CNI_MTU"
    }
  ]
}
```

The `KUBERNETES_NODE_NAME` and `CNI_MTU` placeholders will be replaced when creating a cluster using this configuration.

### Supported CNI Plugins

CloudStack supports various CNI plugins including:

* **Calico** - Provides networking and network policy for Kubernetes
* **Flannel** - Simple overlay network for Kubernetes
* **Weave** - Container networking solution
* **Custom plugins** - Any CNI-compliant plugin can be configured

## Usage with Kubernetes Clusters

CNI configurations are used with Kubernetes clusters by referencing the configuration ID:

```hcl
resource "cloudstack_kubernetes_cluster" "example" {
  name                 = "example-cluster"
  zone                 = "zone1"
  kubernetes_version   = "1.25.0"
  service_offering     = "Medium Instance"
  
  cni_configuration_id = cloudstack_cni_configuration.calico.id
  cni_config_details = {
    "CNI_MTU"               = "1450"
    "KUBERNETES_NODE_NAME"  = "spec.nodeName"
    "KUBECONFIG_FILEPATH"   = "/etc/cni/net.d/calico-kubeconfig"
  }
}
```

## Import

CNI configurations can be imported using the configuration ID:

```shell
$ terraform import cloudstack_cni_configuration.example <CNI_CONFIGURATION_ID>
```
