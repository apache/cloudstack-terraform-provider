---
layout: default
page_title: "CloudStack: cloudstack_kubernetes_cluster"
sidebar_current: "docs-cloudstack-resource-kubernetes_cluster"
description: |-
    Creates a Kubernetes Cluster
---

# CloudStack: cloudstack_kubernetes_cluster

A `cloudstack_kubernetes_cluster` resource manages a Kubernetes cluster within CloudStack.

## Example Usage

```hcl
resource "cloudstack_kubernetes_cluster" "example" {
    name = "example-cluster"
    zone = "zone-id"
    kubernetes_version = "1.18.6"
    service_offering = "small"
    size = 1
    autoscaling_enabled = true
    min_size = 1
    max_size = 5
    control_nodes_size = 1
    description = "An example Kubernetes cluster"
    keypair = "my-ssh-key"
    network_id = "net-id"
    state = "Running"
    project = "my-project"
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Kubernetes cluster.
* `zone` - (Required) The zone where the Kubernetes cluster will be deployed.
* `kubernetes_version` - (Required) The Kubernetes version for the cluster.
* `service_offering` - (Required) The service offering for the nodes in the cluster.
* `size` - (Optional) The initial size of the Kubernetes cluster. Defaults to `1`.
* `autoscaling_enabled` - (Optional) Whether autoscaling is enabled for the cluster.
* `min_size` - (Optional) The minimum size of the Kubernetes cluster when autoscaling is enabled.
* `max_size` - (Optional) The maximum size of the Kubernetes cluster when autoscaling is enabled.
* `control_nodes_size` - (Optional) The size of the control nodes in the cluster.
* `description` - (Optional) A description for the Kubernetes cluster.
* `keypair` - (Optional) The SSH key pair to use for the nodes in the cluster.
* `network_id` - (Optional) The network ID to connect the Kubernetes cluster to.
* `ip_address` - (Computed) The IP address of the Kubernetes cluster.
* `state` - (Optional) The state of the Kubernetes cluster. Defaults to `"Running"`.
* `project` - (Optional) The project to assign the Kubernetes cluster to.
* `noderootdisksize` - (Optional) root disk size in GB for each node.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Kubernetes cluster.
* `name` - The name of the Kubernetes cluster.
* `description` - The description of the Kubernetes cluster.
* `control_nodes_size` - The size of the control nodes in the cluster.
* `size` - The size of the Kubernetes cluster.
* `autoscaling_enabled` - Whether autoscaling is enabled for the cluster.
* `min_size` - The minimum size of the Kubernetes cluster when autoscaling is enabled.
* `max_size` - The maximum size of the Kubernetes cluster when autoscaling is enabled.
* `keypair` - The SSH key pair used for the nodes in the cluster.
* `network_id` - The network ID connected to the Kubernetes cluster.
* `ip_address` - The IP address of the Kubernetes cluster.
* `state` - The state of the Kubernetes cluster.
* `project` - The project assigned to the Kubernetes cluster.

## Import

Kubernetes clusters can be imported; use `<KUBERNETESCLUSTERID>` as the import ID. For example:

```shell
$ terraform import cloudstack_kubernetes_cluster.example <KUBERNETESCLUSTERID>
```
