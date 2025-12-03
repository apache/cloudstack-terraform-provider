---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_instance"
sidebar_current: "docs-cloudstack-resource-instance"
description: |-
  Creates and automatically starts a virtual machine based on a service offering, disk offering, and template.
---

# cloudstack_instance

Creates and automatically starts a virtual machine based on a service offering,
disk offering, and template.

## Example Usage

### Basic Instance

```hcl
resource "cloudstack_instance" "web" {
  name             = "server-1"
  service_offering = "small"
  network_id       = "6eb22f91-7454-4107-89f4-36afcdf33021"
  template         = "CentOS 6.5"
  zone             = "zone-1"
}
```

### Instance with Inline User Data

```hcl
resource "cloudstack_instance" "web_with_userdata" {
  name             = "web-server"
  service_offering = "small"
  network_id       = "6eb22f91-7454-4107-89f4-36afcdf33021"
  template         = "Ubuntu 20.04"
  zone             = "zone-1"
  
  user_data = base64encode(<<-EOF
    #!/bin/bash
    apt-get update
    apt-get install -y nginx
    systemctl enable nginx
    systemctl start nginx
  EOF
  )
}
```

### Instance with Registered User Data

```hcl
# First, create registered user data
resource "cloudstack_userdata" "web_init" {
  name = "web-server-init"
  
  userdata = base64encode(<<-EOF
    #!/bin/bash
    apt-get update
    apt-get install -y nginx
    
    # Use parameters
    echo "<h1>Welcome to $${app_name}!</h1>" > /var/www/html/index.html
    echo "<p>Environment: $${environment}</p>" >> /var/www/html/index.html
    
    systemctl enable nginx
    systemctl start nginx
  EOF
  )
  
  params = ["app_name", "environment"]
}

# Deploy instance with parameterized user data
resource "cloudstack_instance" "app_server" {
  name             = "app-server-01"
  service_offering = "medium"
  network_id       = "6eb22f91-7454-4107-89f4-36afcdf33021"
  template         = "Ubuntu 20.04"
  zone             = "zone-1"
  
  userdata_id = cloudstack_userdata.web_init.id
  
  userdata_details = {
    "app_name"    = "My Application"
    "environment" = "production"
  }
}
```

### Instance with Template-Linked User Data

```hcl
# Use a template that has user data pre-linked
resource "cloudstack_instance" "from_template" {
  name             = "template-instance"
  service_offering = "small"
  network_id       = "6eb22f91-7454-4107-89f4-36afcdf33021"
  template         = cloudstack_template.web_template.id  # Template with userdata_link
  zone             = "zone-1"
  
  # Override parameters for the template's linked user data
  userdata_details = {
    "app_name" = "Template-Based App"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the instance.

* `display_name` - (Optional) The display name of the instance.

* `service_offering` - (Required) The name or ID of the service offering used
    for this instance.

* `disk_offering` - (Optional) The name or ID of the disk offering for the virtual machine.
   If the template is of ISO format, the disk offering is for the root disk volume.
   Otherwise this parameter is used to indicate the offering for the data disk volume.
   If the template parameter passed is from a Template object, the disk offering refers
   to a DATA Disk Volume created. If the template parameter passed is from an ISO object,
   the disk offering refers to a ROOT Disk Volume created.

* `override_disk_offering` - (Optional) The name or ID of the disk offering for the virtual
   machine to be used for root volume instead of the disk offering mapped in service offering.
   In case of virtual machine deploying from ISO, then the disk offering specified for root
   volume is ignored and uses this override disk offering.

* `host_id` -  (Optional)  destination Host ID to deploy the VM to - parameter available
   for root admin only

* `pod_id` -  (Optional) destination Pod ID to deploy the VM to - parameter available for root admin only

* `cluster_id` - (Optional) destination Cluster ID to deploy the VM to - parameter available
   for root admin only

* `network_id` - (Optional) The ID of the network to connect this instance
    to. Changing this forces a new resource to be created.

* `ip_address` - (Optional) The IP address to assign to this instance. Changing
    this forces a new resource to be created.

* `template` - (Required) The name or ID of the template used for this
    instance. Changing this forces a new resource to be created.

* `root_disk_size` - (Optional) The size of the root disk in gigabytes. The
    root disk is resized on deploy. Only applies to template-based deployments.
    Changing this forces a new resource to be created.

* `group` - (Optional) The group name of the instance.

* `affinity_group_ids` - (Optional) List of affinity group IDs to apply to this
    instance.

* `affinity_group_names` - (Optional) List of affinity group names to apply to
    this instance.

* `security_group_ids` - (Optional) List of security group IDs to apply to this
    instance. Changing this forces a new resource to be created.

* `security_group_names` - (Optional) List of security group names to apply to
    this instance. Changing this forces a new resource to be created.

* `project` - (Optional) The name or ID of the project to deploy this
    instance to. Changing this forces a new resource to be created.

* `zone` - (Required) The name or ID of the zone where this instance will be
    created. Changing this forces a new resource to be created.

* `start_vm` - (Optional) This determines if the instances is started after it
    is created (defaults true)

* `user_data` - (Optional) The user data to provide when launching the
    instance. This can be either plain text or base64 encoded text.

* `userdata_id` - (Optional) The ID of a registered CloudStack user data to use for this instance.
    Cannot be used together with `user_data`.

* `userdata_details` - (Optional) A map of key-value pairs to pass as parameters to the user data script.
    Only valid when `userdata_id` is specified. Keys must match the parameter names defined in the user data.

* `keypair` - (Optional) The name of the SSH key pair that will be used to
    access this instance. (Mutual exclusive with keypairs)

* `keypairs` - (Optional) A list of SSH key pair names that will be used to
    access this instance. (Mutual exclusive with keypair)

* `expunge` - (Optional) This determines if the instance is expunged when it is
    destroyed (defaults false)

* `uefi` - (Optional) When set, will boot the instance in UEFI/Legacy mode (defaults false)

* `boot_mode` - (Optional) The boot mode of the instance. Can only be specified when uefi is true. Valid options are 'Legacy' and 'Secure'.

* `deleteprotection` - (Optional) Set delete protection for the virtual machine. If true, the instance will be protected from deletion.
    Note: If the instance is managed by another service like autoscaling groups or CKS, delete protection will be ignored.

## Attributes Reference

The following attributes are exported:

* `id` - The instance ID.
* `display_name` - The display name of the instance.

## Import

Instances can be imported; use `<INSTANCE ID>` as the import ID. For
example:

```shell
terraform import cloudstack_instance.default 5cf69677-7e4b-4bf4-b868-f0b02bb72ee0
```

When importing into a project you need to prefix the import ID with the project name:

```shell
terraform import cloudstack_instance.default my-project/5cf69677-7e4b-4bf4-b868-f0b02bb72ee0
```
