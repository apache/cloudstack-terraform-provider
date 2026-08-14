---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_network_acl_ruleset"
sidebar_current: "docs-cloudstack-resource-network-acl-ruleset"
description: |-
  Manages a complete set of network ACL rules for a given network ACL.
---

# cloudstack_network_acl_ruleset

Manages a complete set of network ACL rules for a given network ACL. This resource is designed
for managing all rules in an ACL as a single unit, with efficient handling of rule insertions
and deletions.

~> **Note:** This resource is recommended over `cloudstack_network_acl_rule` when you need to
manage multiple rules and frequently insert or remove rules. It provides better change management
by identifying rules by their `rule_number` rather than position in a list.

## Example Usage

### Basic Example

```hcl
resource "cloudstack_network_acl_ruleset" "web_server" {
  acl_id = cloudstack_network_acl.example.id

  rule {
    rule_number  = 10
    action       = "allow"
    cidr_list    = ["0.0.0.0/0"]
    protocol     = "tcp"
    port         = "80"
    traffic_type = "ingress"
    description  = "Allow HTTP"
  }

  rule {
    rule_number  = 20
    action       = "allow"
    cidr_list    = ["0.0.0.0/0"]
    protocol     = "tcp"
    port         = "443"
    traffic_type = "ingress"
    description  = "Allow HTTPS"
  }

  rule {
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

### Example with ICMP

```hcl
resource "cloudstack_network_acl_ruleset" "icmp_example" {
  acl_id = cloudstack_network_acl.example.id

  rule {
    rule_number  = 10
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

### Example with Managed Mode

When `managed = true`, the provider will delete any rules not defined in your configuration.
This is useful for ensuring complete control over the ACL.

```hcl
resource "cloudstack_network_acl_ruleset" "managed_example" {
  acl_id  = cloudstack_network_acl.example.id
  managed = true

  rule {
    rule_number  = 10
    action       = "allow"
    cidr_list    = ["10.0.0.0/8"]
    protocol     = "tcp"
    port         = "22"
    traffic_type = "ingress"
    description  = "Allow SSH"
  }
}
```

### Example with Port Range

```hcl
resource "cloudstack_network_acl_ruleset" "port_range" {
  acl_id = cloudstack_network_acl.example.id

  rule {
    rule_number  = 10
    action       = "allow"
    cidr_list    = ["192.168.1.0/24"]
    protocol     = "tcp"
    port         = "8000-8010"
    traffic_type = "ingress"
    description  = "Allow port range"
  }
}
```

### Example with All Protocols

```hcl
resource "cloudstack_network_acl_ruleset" "all_protocols" {
  acl_id = cloudstack_network_acl.example.id

  rule {
    rule_number  = 100
    action       = "allow"
    cidr_list    = ["0.0.0.0/0"]
    protocol     = "all"
    traffic_type = "egress"
    description  = "Allow all outbound traffic"
  }
}
```

## Argument Reference

The following arguments are supported:

* `acl_id` - (Required) The network ACL ID for which to create the rules.
    Changing this forces a new resource to be created.

* `managed` - (Optional) USE WITH CAUTION! If enabled all the ACL rules for
    this network ACL will be managed by this resource. This means it will delete
    all ACL rules that are not in your config! (defaults false)

* `rule` - (Required) Can be specified multiple times. Each rule block supports
    fields documented below.

* `project` - (Optional) The name or ID of the project to deploy this
    instance to. Changing this forces a new resource to be created.

The `rule` block supports:

* `rule_number` - (Required) The number of the ACL item used to order the ACL rules.
    The ACL rule with the lowest number has the highest priority. Each rule_number
    must be unique within the ruleset.

* `action` - (Optional) The action for the rule. Valid options are: `allow` and
    `deny` (defaults allow).

* `cidr_list` - (Required) A CIDR list to allow access to the given ports.

* `protocol` - (Required) The name of the protocol to allow. Valid options are:
    `tcp`, `udp`, `icmp`, or `all`.

* `icmp_type` - (Optional) The ICMP type to allow, or `-1` to allow `any`. This
    can only be specified if the protocol is ICMP. If not specified when protocol
    is ICMP, defaults to `-1` (all types). See the ICMP Types section below for
    valid values.

* `icmp_code` - (Optional) The ICMP code to allow, or `-1` to allow `any`. This
    can only be specified if the protocol is ICMP. If not specified when protocol
    is ICMP, defaults to `-1` (all codes). See the ICMP Types section below for
    valid codes for each type.

* `port` - (Optional) Port or port range to allow. This can only be specified if
    the protocol is TCP or UDP. Valid formats are:
    - Single port: `"80"`
    - Port range: `"8000-8010"`
    - If not specified for TCP/UDP, allows all ports for that protocol

* `traffic_type` - (Optional) The traffic type for the rule. Valid options are:
    `ingress` or `egress` (defaults ingress).

* `description` - (Optional) A description indicating why the ACL rule is required.

## Attributes Reference

The following attributes are exported:

* `id` - The ACL ID for which the rules are managed.

## Import

Network ACL Rulesets can be imported using the ACL ID. For example:

```shell
terraform import cloudstack_network_acl_ruleset.default e8b5982a-1b50-4ea9-9920-6ea2290c7359
```

When importing into a project you need to prefix the import ID with the project name:

```shell
terraform import cloudstack_network_acl_ruleset.default my-project/e8b5982a-1b50-4ea9-9920-6ea2290c7359
```

## Comparison with cloudstack_network_acl_rule

The `cloudstack_network_acl_ruleset` resource is similar to `cloudstack_network_acl_rule` but
with some key differences:

* **Rule identification**: Uses `rule_number` to identify rules (set-based), rather than position
  in a list. This means inserting a rule in the middle only creates that one rule, without
  triggering updates to other rules.

* **Simpler implementation**: Does not support the deprecated `ports` field or auto-numbering
  of rules. All rules must have an explicit `rule_number`.

* **Better for dynamic rulesets**: If you frequently add or remove rules, this resource will
  generate cleaner Terraform plans with fewer spurious changes.

Use `cloudstack_network_acl_rule` if you need auto-numbering or backward compatibility with
the `ports` field. Use `cloudstack_network_acl_ruleset` for cleaner change management when
managing multiple rules.

## ICMP Types and Codes

When using `protocol = "icmp"`, you can specify `icmp_type` and `icmp_code` to control which
ICMP messages are allowed. If not specified, both default to `-1` (allow all).

### Common ICMP Types

| Type | Name | Common Codes | Description |
|------|------|--------------|-------------|
| -1 | Any | -1 (any) | Allow all ICMP types and codes |
| 0 | Echo Reply | 0 | Response to ping request |
| 3 | Destination Unreachable | 0-15 | Various unreachable conditions |
| 5 | Redirect | 0-3 | Route redirection messages |
| 8 | Echo Request | 0 | Ping request |
| 9 | Router Advertisement | 0, 16 | Router discovery |
| 10 | Router Solicitation | 0 | Router discovery |
| 11 | Time Exceeded | 0-1 | TTL expired or fragment reassembly timeout |
| 12 | Parameter Problem | 0-2 | IP header problems |
| 13 | Timestamp Request | 0 | Time synchronization |
| 14 | Timestamp Reply | 0 | Time synchronization response |

### ICMP Type 3 - Destination Unreachable Codes

| Code | Description |
|------|-------------|
| 0 | Network Unreachable |
| 1 | Host Unreachable |
| 2 | Protocol Unreachable |
| 3 | Port Unreachable |
| 4 | Fragmentation Needed and DF Set |
| 5 | Source Route Failed |
| 6 | Destination Network Unknown |
| 7 | Destination Host Unknown |
| 8 | Source Host Isolated |
| 9 | Network Administratively Prohibited |
| 10 | Host Administratively Prohibited |
| 11 | Network Unreachable for ToS |
| 12 | Host Unreachable for ToS |
| 13 | Communication Administratively Prohibited |
| 14 | Host Precedence Violation |
| 15 | Precedence Cutoff in Effect |

### ICMP Type 5 - Redirect Codes

| Code | Description |
|------|-------------|
| 0 | Redirect for Network |
| 1 | Redirect for Host |
| 2 | Redirect for ToS and Network |
| 3 | Redirect for ToS and Host |

### ICMP Type 11 - Time Exceeded Codes

| Code | Description |
|------|-------------|
| 0 | TTL Exceeded in Transit |
| 1 | Fragment Reassembly Time Exceeded |

### ICMP Type 12 - Parameter Problem Codes

| Code | Description |
|------|-------------|
| 0 | Pointer Indicates Error |
| 1 | Missing Required Option |
| 2 | Bad Length |

### Example: Allow Specific ICMP Types

```hcl
resource "cloudstack_network_acl_ruleset" "icmp_example" {
  acl_id = cloudstack_network_acl.example.id

  # Allow all ICMP (default behavior when icmp_type/icmp_code not specified)
  rule {
    rule_number  = 10
    action       = "allow"
    cidr_list    = ["0.0.0.0/0"]
    protocol     = "icmp"
    traffic_type = "ingress"
    description  = "Allow all ICMP"
  }

  # Allow only ping (echo request)
  rule {
    rule_number  = 20
    action       = "allow"
    cidr_list    = ["10.0.0.0/8"]
    protocol     = "icmp"
    icmp_type    = 8
    icmp_code    = 0
    traffic_type = "ingress"
    description  = "Allow ping from internal network"
  }

  # Allow destination unreachable messages
  rule {
    rule_number  = 30
    action       = "allow"
    cidr_list    = ["0.0.0.0/0"]
    protocol     = "icmp"
    icmp_type    = 3
    icmp_code    = -1  # All unreachable codes
    traffic_type = "ingress"
    description  = "Allow all destination unreachable"
  }

  # Allow time exceeded (for traceroute)
  rule {
    rule_number  = 40
    action       = "allow"
    cidr_list    = ["0.0.0.0/0"]
    protocol     = "icmp"
    icmp_type    = 11
    icmp_code    = 0
    traffic_type = "ingress"
    description  = "Allow TTL exceeded for traceroute"
  }
}
```

For a complete list of ICMP types and codes, refer to the
[IANA ICMP Parameters Registry](https://www.iana.org/assignments/icmp-parameters).

