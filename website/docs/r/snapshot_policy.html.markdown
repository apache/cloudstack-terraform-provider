---
subcategory: "Snapshot"
layout: "cloudstack"
page_title: "CloudStack: cloudstack_snapshot_policy"
sidebar_current: "docs-cloudstack-resource-snapshot-policy"
description: |-
  Creates and manages snapshot policies for volumes.
---

# cloudstack_snapshot_policy

Provides a CloudStack snapshot policy resource. This can be used to create, modify, and delete snapshot policies for volumes.

## Example Usage

### Basic Snapshot Policy

```hcl
resource "cloudstack_snapshot_policy" "daily" {
  volume_id     = cloudstack_disk.data.id
  interval_type = "DAILY"
  max_snaps     = 7
  schedule      = "02:30"
  timezone      = "UTC"
  zone_ids      = [data.cloudstack_zone.zone1.id]

  tags = {
    Environment = "production"
    Purpose     = "backup"
  }
}
```

### Hourly Snapshot Policy

```hcl
resource "cloudstack_snapshot_policy" "hourly" {
  volume_id     = cloudstack_disk.database.id
  interval_type = "HOURLY"
  max_snaps     = 6
  schedule      = "0"  # Top of every hour
  timezone      = "America/New_York"
  zone_ids      = [data.cloudstack_zone.zone1.id]
  
  custom_id = "hourly-db-backup"
}
```

### Multiple Zone Snapshot Policy

```hcl
resource "cloudstack_snapshot_policy" "multi_zone" {
  volume_id     = cloudstack_disk.shared.id
  interval_type = "WEEKLY"
  max_snaps     = 4
  schedule      = "1:03:00"  # Monday at 3:00 AM
  timezone      = "UTC"
  zone_ids      = [
    data.cloudstack_zone.zone1.id,
    data.cloudstack_zone.zone2.id
  ]
}
```

### Monthly Archive Policy

```hcl
resource "cloudstack_snapshot_policy" "monthly_archive" {
  volume_id     = cloudstack_disk.archive.id
  interval_type = "MONTHLY"
  max_snaps     = 12
  schedule      = "1:01:00"  # 1st day of month at 1:00 AM
  timezone      = "UTC"
  zone_ids      = [data.cloudstack_zone.zone1.id]

  tags = {
    Type        = "archive"
    Retention   = "1-year"
    Environment = "production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `volume_id` - (Required) The ID of the volume for which the snapshot policy is being created.

* `interval_type` - (Required) The interval type for the snapshot policy. Valid values are:
  * `HOURLY` - Take snapshots every hour
  * `DAILY` - Take snapshots daily
  * `WEEKLY` - Take snapshots weekly
  * `MONTHLY` - Take snapshots monthly

* `max_snaps` - (Required) Maximum number of snapshots to retain. When this limit is reached, older snapshots are automatically deleted.

* `schedule` - (Required) The schedule for taking snapshots. The format depends on the interval type:
  * **HOURLY**: Minute (0-59), e.g., `"30"` for 30 minutes past every hour
  * **DAILY**: Time in HH:MM format, e.g., `"02:30"` for 2:30 AM daily
  * **WEEKLY**: Day and time in D:HH:MM format, e.g., `"1:02:30"` for Monday at 2:30 AM (1=Monday, 7=Sunday)
  * **MONTHLY**: Day and time in DD:HH:MM format, e.g., `"15:02:30"` for 15th day at 2:30 AM

* `timezone` - (Required) The timezone for the schedule. Use standard timezone names like `UTC`, `America/New_York`, `Europe/London`, etc.

* `zone_ids` - (Optional) List of zone IDs where the snapshot policy should be applied. If not specified, the policy applies to all zones.

* `custom_id` - (Optional) A custom ID for the snapshot policy. This is useful for identification purposes and cannot be changed after creation.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the snapshot policy.

## Schedule Format Examples

### Hourly Schedules
```hcl
schedule = "0"   # Top of every hour (XX:00)
schedule = "15"  # 15 minutes past every hour (XX:15)
schedule = "30"  # 30 minutes past every hour (XX:30)
schedule = "45"  # 45 minutes past every hour (XX:45)
```

### Daily Schedules
```hcl
schedule = "01:00"  # 1:00 AM daily
schedule = "02:30"  # 2:30 AM daily
schedule = "14:00"  # 2:00 PM daily
schedule = "23:59"  # 11:59 PM daily
```

### Weekly Schedules
```hcl
schedule = "1:02:00"  # Monday at 2:00 AM
schedule = "2:03:30"  # Tuesday at 3:30 AM
schedule = "7:01:00"  # Sunday at 1:00 AM
```

### Monthly Schedules
```hcl
schedule = "1:02:00"   # 1st day of month at 2:00 AM
schedule = "15:03:30"  # 15th day of month at 3:30 AM
schedule = "28:01:00"  # 28th day of month at 1:00 AM
```

## Import

Snapshot policies can be imported using the policy ID:

```shell
terraform import cloudstack_snapshot_policy.example 12345678-1234-1234-1234-123456789012
```

## Notes

* Snapshot policies will continue to execute regardless of the volume's state (attached/detached).
* Snapshots are incremental and space-efficient in most CloudStack configurations.
* You can have multiple snapshot policies for the same volume with different intervals.
* Consider storage costs when setting `max_snaps` values, especially for frequently changing volumes.
* The actual snapshot creation time may vary slightly from the scheduled time depending on system load.
* When a snapshot policy is deleted, existing snapshots created by that policy are not automatically deleted.
