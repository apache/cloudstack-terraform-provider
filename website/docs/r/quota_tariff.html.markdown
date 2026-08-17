---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_quota_tariff"
sidebar_current: "docs-cloudstack-resource-quota-tariff"
description: |-
  Creates and manages CloudStack quota tariffs.
---

# cloudstack_quota_tariff

Provides a CloudStack quota tariff resource. This can be used to create, modify, and delete quota tariffs.

## Example Usage

```hcl
# Create a CPU usage tariff
resource "cloudstack_quota_tariff" "cpu_tariff" {
  name        = "CPU Usage Tariff"
  usage_type  = 1
  value       = 0.05
  description = "Tariff for CPU usage per hour"
}

# Create a memory usage tariff with date range
resource "cloudstack_quota_tariff" "memory_tariff" {
  name           = "Memory Usage Tariff"
  usage_type     = 2
  value          = 0.01
  description    = "Tariff for memory usage per GB per hour"
  start_date     = "2024-01-01"
  end_date       = "2024-12-31"
  activation_rule = "account.type == 'user'"
}

# Create a storage tariff
resource "cloudstack_quota_tariff" "storage_tariff" {
  name        = "Primary Storage Tariff"
  usage_type  = 6
  value       = 0.1
  description = "Tariff for primary storage usage per GB per month"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the quota tariff.

* `usage_type` - (Required) The usage type for the quota tariff. This cannot be changed after creation.

* `value` - (Required) The monetary value of the quota tariff.

* `description` - (Optional) A description of the quota tariff.

* `start_date` - (Optional) The start date for the quota tariff in yyyy-MM-dd format.

* `end_date` - (Optional) The end date for the quota tariff in yyyy-MM-dd format.

* `activation_rule` - (Optional) The activation rule that determines when this tariff applies.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the quota tariff.

* `currency` - The currency used for the tariff.

* `effective_date` - The effective date when the tariff becomes active.

* `usage_name` - The human-readable name of the usage type.

* `usage_unit` - The unit of measurement for the usage.

* `position` - The position/priority of the tariff.

* `removed` - Whether the tariff has been marked as removed.

## Import

Quota tariffs can be imported using the tariff ID:

```
$ terraform import cloudstack_quota_tariff.cpu_tariff 12345678-1234-1234-1234-123456789abc
```