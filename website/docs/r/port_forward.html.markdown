---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_port_forward"
sidebar_current: "docs-cloudstack-resource-port-forward"
description: |-
  Creates port forwards.
---

# cloudstack_port_forward

Creates port forwards.

## Example Usage

### Basic Port Forward

```hcl
resource "cloudstack_port_forward" "default" {
  ip_address_id = "30b21801-d4b3-4174-852b-0c0f30bdbbfb"

  forward {
    protocol           = "tcp"
    private_port       = 80
    public_port        = 8080
    virtual_machine_id = "f8141e2f-4e7e-4c63-9362-986c908b7ea7"
  }
}
```

### Port Forward with Automatic Project Inheritance

```hcl
# Create a VPC in a project
resource "cloudstack_vpc" "project_vpc" {
  name         = "project-vpc"
  cidr         = "10.0.0.0/16"
  vpc_offering = "Default VPC offering"
  project      = "my-project"
  zone         = "zone-1"
}

# IP address automatically inherits project from VPC
resource "cloudstack_ipaddress" "project_ip" {
  vpc_id = cloudstack_vpc.project_vpc.id
  zone   = "zone-1"
}

# Instance in the project
resource "cloudstack_instance" "web" {
  name             = "web-server"
  service_offering = "Small Instance"
  network_id       = "your-network-id"
  template         = "your-template-id"
  zone             = "zone-1"
  project          = "my-project"
}

# Port forward automatically inherits project from IP address
resource "cloudstack_port_forward" "project_forward" {
  ip_address_id = cloudstack_ipaddress.project_ip.id
  # project is automatically inherited from the IP address

  forward {
    protocol           = "tcp"
    private_port       = 80
    public_port        = 8080
    virtual_machine_id = cloudstack_instance.web.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `ip_address_id` - (Required) The IP address ID for which to create the port
    forwards. Changing this forces a new resource to be created.

* `managed` - (Optional) USE WITH CAUTION! If enabled all the port forwards for
    this IP address will be managed by this resource. This means it will delete
    all port forwards that are not in your config! (defaults false)

* `project` - (Optional) The name or ID of the project to deploy this
    resource to. Changing this forces a new resource to be created. If not
    specified, the project will be automatically inherited from the IP address.

* `forward` - (Required) Can be specified multiple times. Each forward block supports
    fields documented below.

The `forward` block supports:

* `protocol` - (Required) The name of the protocol to allow. Valid options are:
    `tcp` and `udp`.

* `private_port` - (Required) The starting port of port forwarding rule's private port range.

* `private_end_port` - (Optional) The ending port of port forwarding rule's private port range.
    If not specified, the private port will be used as the end port.

* `public_port` - (Required) The starting port of port forwarding rule's public port range.

* `public_end_port` - (Optional) The ending port of port forwarding rule's public port range.
    If not specified, the public port will be used as the end port.

* `virtual_machine_id` - (Required) The ID of the virtual machine to forward to.

* `vm_guest_ip` - (Optional) The virtual machine IP address for the port
    forwarding rule (useful when the virtual machine has secondairy NICs
    or IP addresses).

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the IP address for which the port forwards are created.
* `vm_guest_ip` - The IP address of the virtual machine that is used
    for the port forwarding rule.
