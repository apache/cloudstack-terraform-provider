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
* `allocation_state` - Allocation state of this cluster for allocation of new resources.
* `cluster_name` - The cluster name.
* `name` - The name of the cluster.
* `cluster_type` - Type of the cluster: CloudManaged, ExternalManaged.
* `guest_vswitch_name` - Name of virtual switch used for guest traffic in the cluster. This would override zone wide traffic label setting.
* `guest_vswitch_type` - Type of virtual switch used for guest traffic in the cluster. Allowed values are, vmwaresvs (for VMware standard vSwitch) and vmwaredvs (for VMware distributed vSwitch).
* `hypervisor` - Hypervisor type of the cluster: XenServer, KVM, VMware, Hyperv, BareMetal, Simulator, Ovm3.
* `ovm3_cluster` - Ovm3 native OCFS2 clustering enabled for cluster.
* `ovm3_pool` - Ovm3 native pooling enabled for cluster.
* `ovm3_vip` - Ovm3 vip to use for pool (and cluster).
* `ovm3vip` - Ovm3 vip used for pool (and cluster).
* `password` - The password for the host.
* `public_vswitch_name` - Name of virtual switch used for public traffic in the cluster. This would override zone wide traffic label setting.
* `public_vswitch_type` - Type of virtual switch used for public traffic in the cluster. Allowed values are, vmwaresvs (for VMware standard vSwitch) and vmwaredvs (for VMware distributed vSwitch).
* `pod_id` - The Pod ID for the cluster.
* `pod_name` - The name of the pod where the cluster is created.
* `url` - The URL.
* `username` - The username for the cluster.
* `vsm_ip_address` - The ipaddress of the VSM associated with this cluster.
* `vsm_password` - The password for the VSM associated with this cluster.
* `vsm_username` - The username for the VSM associated with this cluster.
* `zone_id` - The Zone ID for the cluster.
* `zone_name` - The name of the zone where the cluster is created.
* `managed_state` - The managed state of the cluster.
* `cpu_overcommit_ratio` - The CPU overcommit ratio of the cluster.
* `memory_overcommit_ratio` - The memory overcommit ratio of the cluster.
* `arch` - The CPU arch of the cluster.
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
