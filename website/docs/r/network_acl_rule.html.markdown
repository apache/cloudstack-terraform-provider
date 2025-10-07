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

```hcl
rresource "cloudstack_network_acl_rule" "default" {
  acl_id = "f3843ce0-334c-4586-bbd3-0c2e2bc946c6"

  rule {
    action       = "allow"
    cidr_list    = ["10.0.0.0/8"]
    protocol     = "tcp"
    port         = "80-443" # preferred, string, supports single port or range
    traffic_type = "ingress"
  }
}
```
# Deprecated example (do not use in new configs)
resource "cloudstack_network_acl_rule" "deprecated" {
  acl_id = "f3843ce0-334c-4586-bbd3-0c2e2bc946c6"

  rule {
    action       = "allow"
    cidr_list    = ["10.0.0.0/8"]
    protocol     = "tcp"
    ports        = ["80", "1000-2000"] # deprecated, use 'port' instead
    traffic_type = "ingress"
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

* `port` - (Optional, string) The port or port range to allow. Preferred for new configs.
    Use a single port (e.g. "80") or a range (e.g. "1000-2000"). Required for tcp or udp protocols. Cannot be used with ports.

* `ports` - (Optional, Deprecated) List of ports and/or port ranges to allow. This can only
    be specified if the protocol is TCP, UDP, ALL or a valid protocol number.
    **Deprecated**: Use port (string) instead. ports will be removed in a future version.

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

### Deprecation Notice:
The `ports` attribute is deprecated and will be removed in a future version. Use `port` (string) instead for all new configurations.
