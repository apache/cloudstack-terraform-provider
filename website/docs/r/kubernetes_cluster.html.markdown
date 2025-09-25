---
layout: default
page_title: "CloudStack: cloudstack_kubernetes_cluster"
sidebar_current: "docs-cloudstack-resource-kubernetes_cluster"
description: |-
    Creates and manages a CloudStack Kubernetes Service (CKS) cluster
---

# CloudStack: cloudstack_kubernetes_cluster

A `cloudstack_kubernetes_cluster` resource manages a CloudStack Kubernetes Service (CKS) cluster within CloudStack. This resource supports advanced features including mixed node types, custom templates, CNI configurations, and autoscaling.

## Example Usage

### Basic Cluster

```hcl
resource "cloudstack_kubernetes_cluster" "basic" {
  name               = "basic-cluster"
  zone               = "zone1"
  kubernetes_version = "1.25.0"
  service_offering   = "Medium Instance"
  size               = 3
  description        = "Basic Kubernetes cluster"
}
```

### Advanced Cluster with CKS Features

```hcl
# Kubernetes version resource
resource "cloudstack_kubernetes_version" "k8s_v1_25" {
  semantic_version = "1.25.0"
  name             = "Kubernetes v1.25.0 with Calico"
  url              = "http://example.com/k8s-setup-v1.25.0.iso"
  min_cpu          = 2
  min_memory       = 2048
  zone             = "zone1"
  state            = "Enabled"
}

# CNI configuration
resource "cloudstack_cni_configuration" "calico" {
  name       = "calico-cni-config"
  cni_config = base64encode(jsonencode({
    "name"       = "k8s-pod-network",
    "cniVersion" = "0.3.1",
    "plugins" = [
      {
        "type"           = "calico",
        "datastore_type" = "kubernetes",
        "nodename"       = "KUBERNETES_NODE_NAME",
        "mtu"            = "CNI_MTU"
      }
    ]
  }))
  
  params = ["KUBERNETES_NODE_NAME", "CNI_MTU"]
}

# Advanced cluster with mixed node types
resource "cloudstack_kubernetes_cluster" "advanced" {
  name               = "production-cluster"
  zone               = "zone1"
  kubernetes_version = cloudstack_kubernetes_version.k8s_v1_25.semantic_version
  service_offering   = "Medium Instance"

  # Cluster configuration
  size               = 3
  control_nodes_size = 3
  etcd_nodes_size    = 3

  # Autoscaling
  autoscaling_enabled = true
  min_size            = 2
  max_size            = 10

  # Node configuration
  noderootdisksize = 50

  # Mixed node offerings
  node_offerings = {
    "control" = "Large Instance"
    "worker"  = "Medium Instance"
    "etcd"    = "Medium Instance"
  }

  # Custom templates
  node_templates = {
    "control" = "ubuntu-20.04-k8s-template"
    "worker"  = "ubuntu-20.04-k8s-template"
  }

  # CNI Configuration
  cni_configuration_id = cloudstack_cni_configuration.calico.id
  cni_config_details = {
    "CNI_MTU"               = "1450"
    "KUBERNETES_NODE_NAME"  = "spec.nodeName"
  }

  description = "Production cluster with mixed node types"
  hypervisor  = "KVM"
}
```


## Argument Reference

The following arguments are supported:

### Required Arguments

* `name` - (Required) The name of the Kubernetes cluster.
* `zone` - (Required) The zone where the Kubernetes cluster will be deployed.
* `kubernetes_version` - (Required) The Kubernetes version for the cluster.
* `service_offering` - (Required) The service offering for the nodes in the cluster.

### Basic Configuration

* `size` - (Optional) The number of worker nodes in the Kubernetes cluster. Defaults to `1`.
* `control_nodes_size` - (Optional) The number of control plane nodes in the cluster. Defaults to `1`.
* `etcd_nodes_size` - (Optional) The number of etcd nodes in the cluster. Defaults to `0` (uses control nodes for etcd).
* `description` - (Optional) A description for the Kubernetes cluster.
* `hypervisor` - (Optional) The hypervisor type for the cluster nodes. Defaults to `"KVM"`.

### Autoscaling Configuration

* `autoscaling_enabled` - (Optional) Whether autoscaling is enabled for the cluster. Defaults to `false`.
* `min_size` - (Optional) The minimum number of worker nodes when autoscaling is enabled.
* `max_size` - (Optional) The maximum number of worker nodes when autoscaling is enabled.

### Node Configuration

* `noderootdisksize` - (Optional) Root disk size in GB for each node. Defaults to `20`.
* `node_offerings` - (Optional) A map of node roles to service offerings. Valid roles are `control`, `worker`, and `etcd`. If not specified, the main `service_offering` is used for all nodes.
* `node_templates` - (Optional) A map of node roles to instance templates. Valid roles are `control`, `worker`, and `etcd`. If not specified, system VM template will be used.

### CNI Configuration

* `cni_configuration_id` - (Optional) The ID of a CNI configuration to use for the cluster. If not specified, the default CNI configuration will be used.
* `cni_config_details` - (Optional) A map of CNI configuration parameter values to substitute in the CNI configuration.

### Network and Security

* `keypair` - (Optional) The SSH key pair to use for the nodes in the cluster.
* `network_id` - (Optional) The network ID to connect the Kubernetes cluster to.

### Project and Domain

* `project` - (Optional) The project to assign the Kubernetes cluster to.
* `domain_id` - (Optional) The domain ID for the cluster.
* `account` - (Optional) The account name for the cluster.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Kubernetes cluster.
* `ip_address` - The IP address of the Kubernetes cluster API server.
* `state` - The current state of the Kubernetes cluster.
* `created` - The timestamp when the cluster was created.
* `zone_id` - The zone ID where the cluster is deployed.
* `zone_name` - The zone name where the cluster is deployed.
* `kubernetes_version_id` - The ID of the Kubernetes version used.
* `service_offering_id` - The ID of the service offering used.
* `master_nodes` - The number of master/control nodes in the cluster.
* `cpu_number` - The number of CPUs allocated to the cluster.
* `memory` - The amount of memory (in MB) allocated to the cluster.

## Import

Kubernetes clusters can be imported; use `<KUBERNETESCLUSTERID>` as the import ID. For example:

```shell
$ terraform import cloudstack_kubernetes_cluster.example <KUBERNETESCLUSTERID>
```
