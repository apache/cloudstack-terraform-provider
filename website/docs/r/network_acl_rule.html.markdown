---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_network_acl_rule"
sidebar_current: "docs-cloudstack-resource-network-acl-rule"
description: |-
  Creates network ACL rules for a given network ACL.
---

# cloudstack_network_acl_rule

Creates network ACL rules for a given network ACL.

## Example Usage

### Basic Example with Port

```hcl
resource "cloudstack_network_acl_rule" "default" {
  acl_id = "f3843ce0-334c-4586-bbd3-0c2e2bc946c6"

  rule {
    action       = "allow"
    cidr_list    = ["10.0.0.0/8"]
    protocol     = "tcp"
    port         = "80"
    traffic_type = "ingress"
  }
}
```

### Example with Port Range

```hcl
resource "cloudstack_network_acl_rule" "port_range" {
  acl_id = "f3843ce0-334c-4586-bbd3-0c2e2bc946c6"

  rule {
    action       = "allow" 
    cidr_list    = ["192.168.1.0/24"]
    protocol     = "tcp"
    port         = "8000-8010"
    traffic_type = "ingress"
  }
}
```

### Example with No Port (Allow All Ports)

```hcl
resource "cloudstack_network_acl_rule" "all_ports" {
  acl_id = "f3843ce0-334c-4586-bbd3-0c2e2bc946c6"

  rule {
    action       = "allow"
    cidr_list    = ["10.0.0.0/16"]
    protocol     = "tcp"
    traffic_type = "ingress"
    description  = "Allow all TCP traffic from internal network"
  }
}
```

### Example with ICMP

```hcl
resource "cloudstack_network_acl_rule" "icmp" {
  acl_id = "f3843ce0-334c-4586-bbd3-0c2e2bc946c6"

  rule {
    action       = "allow"
    cidr_list    = ["0.0.0.0/0"]
    protocol     = "icmp"
    icmp_type    = 8
    icmp_code    = 0
    traffic_type = "ingress"
    description  = "Allow ping"
  }
}
```

### Complete Example with Multiple Rules

```hcl
resource "cloudstack_network_acl_rule" "web_server" {
  acl_id = "f3843ce0-334c-4586-bbd3-0c2e2bc946c6"

  # HTTP traffic
  rule {
    rule_number  = 10
    action       = "allow"
    cidr_list    = ["0.0.0.0/0"]
    protocol     = "tcp"
    port         = "80"
    traffic_type = "ingress"
    description  = "Allow HTTP"
  }

  # HTTPS traffic
  rule {
    rule_number  = 20
    action       = "allow"
    cidr_list    = ["0.0.0.0/0"]
    protocol     = "tcp"
    port         = "443"
    traffic_type = "ingress"
    description  = "Allow HTTPS"
  }

  # SSH from management network
  rule {
    rule_number  = 30
    action       = "allow"
    cidr_list    = ["192.168.100.0/24"]
    protocol     = "tcp"
    port         = "22"
    traffic_type = "ingress"
    description  = "Allow SSH from management"
  }

  # Allow all outbound traffic
  rule {
    rule_number  = 100
    action       = "allow"
    cidr_list    = ["0.0.0.0/0"]
    protocol     = "tcp"
    traffic_type = "egress"
    description  = "Allow all outbound TCP"
  }
}

## Argument Reference

The following arguments are supported:

* `acl_id` - (Required) The network ACL ID for which to create the rules.
    Changing this forces a new resource to be created.

* `managed` - (Optional) USE WITH CAUTION! If enabled all the firewall rules for
    this network ACL will be managed by this resource. This means it will delete
    all firewall rules that are not in your config! (defaults false)

* `rule` - (Optional) Can be specified multiple times. Each rule block supports
    fields documented below. If `managed = false` at least one rule is required!

* `project` - (Optional) The name or ID of the project to deploy this
    instance to. Changing this forces a new resource to be created.

* `parallelism` (Optional) Specifies how much rules will be created or deleted
    concurrently. (defaults 2)

The `rule` block supports:

* `rule_number` - (Optional) The number of the ACL item used to order the ACL rules. The ACL rule with the lowest number has the highest priority. If not specified, the ACL item will be created with a number one greater than the highest numbered rule.

* `action` - (Optional) The action for the rule. Valid options are: `allow` and
    `deny` (defaults allow).

* `cidr_list` - (Required) A CIDR list to allow access to the given ports.

* `protocol` - (Required) The name of the protocol to allow. Valid options are:
    `tcp`, `udp`, `icmp`, `all` or a valid protocol number.

* `icmp_type` - (Optional) The ICMP type to allow, or `-1` to allow `any`. This
    can only be specified if the protocol is ICMP. (defaults 0)

* `icmp_code` - (Optional) The ICMP code to allow, or `-1` to allow `any`. This
    can only be specified if the protocol is ICMP. (defaults 0)

* `port` - (Optional) Port or port range to allow. This can only be specified if 
    the protocol is TCP, UDP, ALL or a valid protocol number. Valid formats are:
    - Single port: `"80"`
    - Port range: `"8000-8010"`
    - If not specified for TCP/UDP, allows all ports for that protocol

* `ports` - (Optional) **DEPRECATED**: Use `port` instead. List of ports and/or 
    port ranges to allow. This field is deprecated and will be removed in a future 
    version. For backward compatibility only.

* `traffic_type` - (Optional) The traffic type for the rule. Valid options are:
    `ingress` or `egress` (defaults ingress).

* `description` - (Optional) A description indicating why the ACL rule is required.

## Attributes Reference

The following attributes are exported:

* `id` - The ACL ID for which the rules are created.

## Import

Network ACL Rules can be imported; use `<NETWORK ACL Rule ID>` as the import ID. For
example:

```shell
terraform import cloudstack_network_acl_rule.default e8b5982a-1b50-4ea9-9920-6ea2290c7359
```

When importing into a project you need to prefix the import ID with the project name:

```shell
terraform import cloudstack_network_acl_rule.default my-project/e8b5982a-1b50-4ea9-9920-6ea2290c7359
```