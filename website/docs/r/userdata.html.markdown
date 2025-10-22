---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_userdata"
sidebar_current: "docs-cloudstack-resource-userdata"
description: |-
  Registers and manages user data in CloudStack for VM initialization.
---

# cloudstack_userdata

Registers user data in CloudStack that can be used to initialize virtual machines during deployment. User data typically contains scripts, configuration files, or other initialization data that should be executed when a VM starts.

## Example Usage

### Basic User Data

```hcl
resource "cloudstack_userdata" "web_init" {
  name = "web-server-init"
  
  userdata = base64encode(<<-EOF
    #!/bin/bash
    apt-get update
    apt-get install -y nginx
    systemctl enable nginx
    systemctl start nginx
  EOF
  )
}
```

### Parameterized User Data

```hcl
resource "cloudstack_userdata" "app_init" {
  name = "app-server-init"
  
  userdata = base64encode(<<-EOF
    #!/bin/bash
    apt-get update
    apt-get install -y nginx
    
    # Use parameters passed from instance deployment
    echo "<h1>Welcome to $${app_name}!</h1>" > /var/www/html/index.html
    echo "<p>Environment: $${environment}</p>" >> /var/www/html/index.html
    echo "<p>Debug Mode: $${debug_mode}</p>" >> /var/www/html/index.html
    
    systemctl enable nginx
    systemctl start nginx
  EOF
  )
  
  # Define parameters that can be passed during instance deployment
  params = ["app_name", "environment", "debug_mode"]
}
```

### Project-Scoped User Data

```hcl
resource "cloudstack_userdata" "project_init" {
  name       = "project-specific-init"
  project_id = "12345678-1234-1234-1234-123456789012"
  
  userdata = base64encode(<<-EOF
    #!/bin/bash
    # Project-specific initialization
    echo "Initializing project environment..."
  EOF
  )
}
```

## Argument Reference

The following arguments are supported:

### Required Arguments

* `name` - (Required) The name of the user data. Must be unique within the account/project scope.
* `userdata` - (Required) The user data content to be registered. Should be base64 encoded. This is typically a cloud-init script or other initialization data.

### Optional Arguments

* `account` - (Optional) The account name for the user data. Must be used together with `domain_id`. If not specified, uses the current account.
* `domain_id` - (Optional) The domain ID for the user data. Required when `account` is specified.
* `project_id` - (Optional) The project ID to create this user data for. Cannot be used together with `account`/`domain_id`.
* `params` - (Optional) A list of parameter names that are declared in the user data content. These parameters can be passed values during instance deployment using `userdata_details`.

## Attributes Reference

The following attributes are exported:

* `id` - The user data ID.
* `name` - The name of the user data.
* `userdata` - The registered user data content.
* `account` - The account name owning the user data.
* `domain_id` - The domain ID where the user data belongs.
* `project_id` - The project ID if the user data is project-scoped.
* `params` - The list of parameter names defined in the user data.

## Usage with Templates and Instances

User data can be used in multiple ways:

### 1. Linked to Templates

```hcl
resource "cloudstack_template" "web_template" {
  name         = "web-server-template"
  # ... other template arguments ...
  
  userdata_link {
    userdata_id     = cloudstack_userdata.app_init.id
    userdata_policy = "ALLOWOVERRIDE"  # Allow instance to override
  }
}
```

### 2. Direct Instance Usage

```hcl
resource "cloudstack_instance" "web_server" {
  name     = "web-server-01"
  # ... other instance arguments ...
  
  userdata_id = cloudstack_userdata.app_init.id
  
  # Pass parameter values to the userdata script
  userdata_details = {
    "app_name"    = "My Web Application"
    "environment" = "production"
    "debug_mode"  = "false"
  }
}
```

## Import

User data can be imported using the user data ID:

```
terraform import cloudstack_userdata.example 12345678-1234-1234-1234-123456789012
```

## Notes

* User data content should be base64 encoded before registration
* Parameter substitution in user data uses the format `${parameter_name}`
* Parameters must be declared in the `params` list to be usable
* User data is immutable after creation - changes require resource recreation
* Maximum user data size depends on CloudStack configuration (typically 32KB)
