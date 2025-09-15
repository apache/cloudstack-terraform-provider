---
subcategory: "Cluster"
layout: "cloudstack"
page_title: "CloudStack: cloudstack_cluster"
description: |-
  Creates a cluster.
---

# cloudstack_cluster

Creates a cluster.

## Example Usage

```hcl
resource "cloudstack_cluster" "default" {
  name = "cluster-1"
  cluster_type = "CloudManaged"
  hypervisor = "KVM"
  pod_id = "1"
  zone_id = "1"
}
```

## Argument Reference

The following arguments are supported:

<<<<<<< HEAD
* `name` - (Required) The name of the cluster.
* `cluster_type` - (Required) Type of the cluster: CloudManaged, ExternalManaged.
* `hypervisor` - (Required) Hypervisor type of the cluster: XenServer, KVM, VMware, Hyperv, BareMetal, Simulator, Ovm3.
* `pod_id` - (Required) The Pod ID for the cluster.
* `zone_id` - (Required) The Zone ID for the cluster.
* `allocation_state` - (Optional) Allocation state of this cluster for allocation of new resources.
* `arch` - (Optional) The CPU arch of the cluster. Valid options are: x86_64, aarch64.
* `guest_vswitch_name` - (Optional) Name of virtual switch used for guest traffic in the cluster. This would override zone wide traffic label setting.
* `guest_vswitch_type` - (Optional) Type of virtual switch used for guest traffic in the cluster. Allowed values are: vmwaresvs (for VMware standard vSwitch) and vmwaredvs (for VMware distributed vSwitch).
* `ovm3cluster` - (Optional) Ovm3 native OCFS2 clustering enabled for cluster.
* `ovm3pool` - (Optional) Ovm3 native pooling enabled for cluster.
* `ovm3vip` - (Optional) Ovm3 vip to use for pool (and cluster).
* `password` - (Optional) The password for the host.
* `public_vswitch_name` - (Optional) Name of virtual switch used for public traffic in the cluster. This would override zone wide traffic label setting.
* `public_vswitch_type` - (Optional) Type of virtual switch used for public traffic in the cluster. Allowed values are: vmwaresvs (for VMware standard vSwitch) and vmwaredvs (for VMware distributed vSwitch).
* `url` - (Optional) The URL for the cluster.
* `username` - (Optional) The username for the cluster.
* `vsm_ip_address` - (Optional) The IP address of the VSM associated with this cluster.
* `vsm_password` - (Optional) The password for the VSM associated with this cluster.
* `vsm_username` - (Optional) The username for the VSM associated with this cluster.
=======
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

>>>>>>> apache/main

## Attributes Reference

The following attributes are exported:

<<<<<<< HEAD
* `id` - The ID of the cluster.
* `pod_name` - The name of the pod where the cluster is created.
* `zone_name` - The name of the zone where the cluster is created.
* `managed_state` - The managed state of the cluster.
* `cpu_overcommit_ratio` - The CPU overcommit ratio of the cluster.
* `memory_overcommit_ratio` - The memory overcommit ratio of the cluster.

## Import

Clusters can be imported; use `<CLUSTER ID>` as the import ID. For example:

```shell
terraform import cloudstack_cluster.default 5fb02d7f-9513-4f96-9fbe-b5d167f4e90b
=======
* `id` - The instance ID.



## Import

Clusters can be imported; use `<CLUSTER ID>` as the import ID. For
example:

```shell
terraform import cloudstack_cluster.example 5cf69677-7e4b-4bf4-b868-f0b02bb72ee0
```
>>>>>>> apache/main
