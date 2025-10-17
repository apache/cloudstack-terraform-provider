---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_quota"
sidebar_current: "docs-cloudstack-datasource-quota"
description: |-
  Gets information about CloudStack quota summaries.
---

# cloudstack_quota

Use this data source to retrieve quota summary information for CloudStack accounts and domains. This provides information about quota values and whether quota is enabled for specific accounts.

## Example Usage

```hcl
# Get quota summary for all accounts
data "cloudstack_quota" "all_quotas" {
}

# Get quota summary for a specific account
data "cloudstack_quota" "account_quota" {
  account   = "myaccount"
  domain_id = "domain-uuid"
}

# Get quota summary for a specific domain
data "cloudstack_quota" "domain_quota" {
  domain_id = "domain-uuid"
}
```

## Argument Reference

The following arguments are supported:

* `account` - (Optional) The name of the account to filter quota summaries.

* `domain_id` - (Optional) The ID of the domain to filter quota summaries.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `quotas` - A list of quota summary objects. Each object contains:
  * `account_id` - The ID of the account.
  * `account` - The name of the account.
  * `domain_id` - The ID of the domain.
  * `domain` - The name of the domain.
  * `quota_value` - The current quota value for the account.
  * `quota_enabled` - Whether quota is enabled for this account.