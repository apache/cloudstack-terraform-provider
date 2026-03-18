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
```

### Using `ruleset` for Better Change Management

The `ruleset` field is recommended when you need to insert or remove rules without
triggering unnecessary updates to other rules. Unlike `rule` (which uses a list),
`ruleset` uses a set that identifies rules by their `rule_number` rather than position.

**Key differences:**
- `ruleset` requires `rule_number` on all rules (no auto-numbering)
- Each `rule_number` must be unique within the ruleset; if you define multiple rules with the same `rule_number`, only the last one will be kept (Terraform's TypeSet behavior)
- `ruleset` does not support the deprecated `ports` field (use `port` instead)
- Inserting a rule in the middle only creates that one rule, without updating others

```hcl
resource "cloudstack_network_acl_rule" "web_server_set" {
  acl_id = "f3843ce0-334c-4586-bbd3-0c2e2bc946c6"

  # HTTP traffic
  ruleset {
    rule_number  = 10
    action       = "allow"
    cidr_list    = ["0.0.0.0/0"]
    protocol     = "tcp"
    port         = "80"
    traffic_type = "ingress"
    description  = "Allow HTTP"
  }

  # HTTPS traffic
  ruleset {
    rule_number  = 20
    action       = "allow"
    cidr_list    = ["0.0.0.0/0"]
    protocol     = "tcp"
    port         = "443"
    traffic_type = "ingress"
    description  = "Allow HTTPS"
  }

  # SSH from management network
  ruleset {
    rule_number  = 30
    action       = "allow"
    cidr_list    = ["192.168.100.0/24"]
    protocol     = "tcp"
    port         = "22"
    traffic_type = "ingress"
    description  = "Allow SSH from management"
  }
}
```

**Note:** You cannot use both `rule` and `ruleset` in the same resource. Choose one based on your needs:
- Use `rule` if you want auto-numbering and don't mind Terraform showing updates when inserting rules
- Use `ruleset` if you frequently insert/remove rules and want minimal plan changes

## Argument Reference

The following arguments are supported:

* `acl_id` - (Required) The network ACL ID for which to create the rules.
    Changing this forces a new resource to be created.

* `managed` - (Optional) USE WITH CAUTION! If enabled all the firewall rules for
    this network ACL will be managed by this resource. This means it will delete
    all firewall rules that are not in your config! (defaults false)

* `rule` - (Optional) Can be specified multiple times. Each rule block supports
    fields documented below. If `managed = false` at least one rule or ruleset is required!
    **Cannot be used together with `ruleset`.**

* `ruleset` - (Optional) Can be specified multiple times. Similar to `rule` but uses
    a set instead of a list, which prevents spurious updates when inserting rules.
    Each ruleset block supports the same fields as `rule` (documented below), with these differences:
    - `rule_number` is **required** (no auto-numbering)
    - `ports` field is not supported (use `port` instead)
    **Cannot be used together with `rule`.**

* `project` - (Optional) The name or ID of the project to deploy this
    instance to. Changing this forces a new resource to be created.

* `parallelism` (Optional) Specifies how much rules will be created or deleted
    concurrently. (defaults 2)

The `rule` and `ruleset` blocks support:

* `rule_number` - (Optional for `rule`, **Required** for `ruleset`) The number of the ACL
    item used to order the ACL rules. The ACL rule with the lowest number has the highest
    priority.
    - For `rule`: If not specified, the provider will auto-assign rule numbers starting at 1,
      increasing sequentially in the order the rules are defined and filling any gaps, rather
      than basing the number on the highest existing rule in the ACL.
    - For `ruleset`: Must be specified for all rules (no auto-numbering).

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
    version. For backward compatibility only. **Not available in `ruleset`.**

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