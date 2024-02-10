---
layout: default
page_title: "CloudStack: cloudstack_kubernetes_version"
sidebar_current: "docs-cloudstack-resource-kubernetes_version"
description: |-
    Creates a Kubernetes Version
---

# CloudStack: cloudstack_kubernetes_version

A `cloudstack_kubernetes_version` resource manages a Kubernetes version within CloudStack.

## Example Usage

```hcl
resource "cloudstack_kubernetes_version" "example" {
    semantic_version = "1.19.0"
    url = "https://example.com/k8s/1.19.0.tar.gz"
    min_cpu = 2
    min_memory = 2048
}
```

## Argument Reference

The following arguments are supported:

* `semantic_version` - (Required) The semantic version of the Kubernetes version.
* `url` - (Required) The URL to download the Kubernetes version package.
* `min_cpu` - (Required) The minimum CPU requirement for the Kubernetes version.
* `min_memory` - (Required) The minimum memory requirement for the Kubernetes version.
* `name` - (Optional) The name of the Kubernetes version.
* `zone` - (Optional) The zone in which the Kubernetes version should be added.
* `checksum` - (Optional) The checksum of the Kubernetes version package.
* `state` - (Optional) The state of the Kubernetes version. Defaults to "Enabled".

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Kubernetes version.
* `semantic_version` - The semantic version of the Kubernetes version.
* `name` - The name of the Kubernetes version.
* `min_cpu` - The minimum CPU requirement for the Kubernetes version.
* `min_memory` - The minimum memory requirement for the Kubernetes version.
* `state` - The state of the Kubernetes version.

## Import

Kubernetes versions can be imported using the ID of the resource; use `<KUBERNETESVERSIONID>` as the import ID. For example:

```shell
$ terraform import cloudstack_kubernetes_version.example <KUBERNETESVERSIONID>
```
