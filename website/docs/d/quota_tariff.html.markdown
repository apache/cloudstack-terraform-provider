---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_quota_tariff"
sidebar_current: "docs-cloudstack-datasource-quota-tariff"
description: |-
  Gets information about CloudStack quota tariffs.
---

# cloudstack_quota_tariff

Use this data source to retrieve information about quota tariffs in CloudStack. Quota tariffs define the pricing for different resource usage types.

## Example Usage

```hcl
# Get all quota tariffs
data "cloudstack_quota_tariff" "all_tariffs" {
}

# Get tariffs by name
data "cloudstack_quota_tariff" "cpu_tariffs" {
  name = "CPU Tariff"
}

# Get tariffs by usage type
data "cloudstack_quota_tariff" "compute_tariffs" {
  usage_type = 1
}

# Output tariff information
output "tariff_details" {
  value = [
    for tariff in data.cloudstack_quota_tariff.all_tariffs.tariffs : {
      name  = tariff.name
      value = tariff.tariff_value
      unit  = tariff.usage_unit
    }
  ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the quota tariff to filter results.

* `usage_type` - (Optional) The usage type to filter tariffs by.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `tariffs` - A list of quota tariff objects. Each object contains:
  * `id` - The ID of the tariff.
  * `name` - The name of the tariff.
  * `description` - The description of the tariff.
  * `usage_type` - The usage type ID.
  * `usage_name` - The human-readable name of the usage type.
  * `usage_unit` - The unit of measurement for the usage.
  * `tariff_value` - The monetary value of the tariff.
  * `end_date` - The end date of the tariff.
  * `effective_date` - The effective date when the tariff becomes active.
  * `activation_rule` - The rule that determines when this tariff is activated.
  * `removed` - Whether the tariff has been marked as removed.
  * `currency` - The currency used for the tariff.
  * `position` - The position/priority of the tariff.