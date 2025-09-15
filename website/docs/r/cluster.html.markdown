---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_cluster"
sidebar_current: "docs-cloudstack-resource-cluster"
description: |-
  Adds a new cluster
---

# cloudstack_cluster

Adds a new cluster

## Example Usage

Basic usage:

```hcl
resource "cloudstack_cluster" "example" {
	cluster_name = "example"
	cluster_type = "CloudManaged"
	hypervisor   = "KVM"
	pod_id       = cloudstack_pod.example.id
	zone_id      = cloudstack_zone.example.id
}
```

## Argument Reference

The following arguments are supported:

* `allocation_state` - (Optional) Allocation state of this cluster for allocation of new resources.
* `cluster_name` - (Required) the cluster name.
* `cluster_type` - (Required) type of the cluster: CloudManaged, ExternalManaged.
* `guest_vswitch_name` - (Optional) Name of virtual switch used for guest traffic in the cluster. This would override zone wide traffic label setting..
* `guest_vswitch_type` - (Optional) Type of virtual switch used for guest traffic in the cluster. Allowed values are, vmwaresvs (for VMware standard vSwitch) and vmwaredvs (for VMware distributed vSwitch).
* `hypervisor` - (Required) hypervisor type of the cluster: XenServer,KVM,VMware,Hyperv,BareMetal,Simulator,Ovm3.
* `ovm3_cluster` - (Optional) Ovm3 native OCFS2 clustering enabled for cluster.
* `ovm3_pool` - (Optional) Ovm3 native pooling enabled for cluster.
* `ovm3_vip` - (Optional) Ovm3 vip to use for pool (and cluster).
* `password` - (Optional) the password for the host.
* `public_vswitch_name` - (Optional) Name of virtual switch used for public traffic in the cluster. This would override zone wide traffic label setting..
* `public_vswitch_type` - (Optional) Type of virtual switch used for public traffic in the cluster. Allowed values are, vmwaresvs (for VMware standard vSwitch) and vmwaredvs (for VMware distributed vSwitch).
* `pod_id` - (Required) the Pod ID for the host.
* `url` - (Optional) the URL.
* `username` - (Optional) the username for the cluster.
* `vsm_ip_address` - (Optional) the ipaddress of the VSM associated with this cluster.
* `vsm_password` - (Optional) the password for the VSM associated with this cluster.
* `vsm_username` - (Optional) the username for the VSM associated with this cluster.
* `zone_id` - (Required) the Zone ID for the cluster.


## Attributes Reference

The following attributes are exported:

* `id` - The instance ID.



## Import

Clusters can be imported; use `<CLUSTER ID>` as the import ID. For
example:

```shell
terraform import cloudstack_cluster.example 5cf69677-7e4b-4bf4-b868-f0b02bb72ee0
```
