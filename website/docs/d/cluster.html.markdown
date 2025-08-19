---
subcategory: "Cluster"
layout: "cloudstack"
page_title: "CloudStack: cloudstack_cluster"
description: |-
  Gets information about a cluster.
---

# cloudstack_cluster

Use this data source to get information about a cluster for use in other resources.

## Example Usage

```hcl
data "cloudstack_cluster" "cluster" {
  filter {
    name = "name"
    value = "cluster-1"
  }
}

output "cluster_id" {
  value = data.cloudstack_cluster.cluster.id
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Required) One or more name/value pairs to filter off of. See detailed documentation below.

### Filter Arguments

* `name` - (Required) The name of the field to filter on. This can be any of the fields returned by the CloudStack API.
* `value` - (Required) The value to filter on. This should be a regular expression.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the cluster.
* `name` - The name of the cluster.
* `cluster_type` - Type of the cluster: CloudManaged, ExternalManaged.
* `hypervisor` - Hypervisor type of the cluster: XenServer, KVM, VMware, Hyperv, BareMetal, Simulator, Ovm3.
* `pod_id` - The Pod ID for the cluster.
* `pod_name` - The name of the pod where the cluster is created.
* `zone_id` - The Zone ID for the cluster.
* `zone_name` - The name of the zone where the cluster is created.
* `allocation_state` - The allocation state of the cluster.
* `managed_state` - The managed state of the cluster.
* `cpu_overcommit_ratio` - The CPU overcommit ratio of the cluster.
* `memory_overcommit_ratio` - The memory overcommit ratio of the cluster.
* `arch` - The CPU arch of the cluster.
* `ovm3vip` - Ovm3 vip used for pool (and cluster).
* `capacity` - The capacity information of the cluster. See Capacity below for more details.

### Capacity

The `capacity` attribute supports the following:

* `capacity_allocated` - The capacity allocated.
* `capacity_total` - The total capacity.
* `capacity_used` - The capacity used.
* `cluster_id` - The ID of the cluster.
* `cluster_name` - The name of the cluster.
* `name` - The name of the capacity.
* `percent_used` - The percentage of capacity used.
* `pod_id` - The ID of the pod.
* `pod_name` - The name of the pod.
* `type` - The type of the capacity.
* `zone_id` - The ID of the zone.
* `zone_name` - The name of the zone.