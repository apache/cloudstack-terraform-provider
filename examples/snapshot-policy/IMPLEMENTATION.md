# CloudStack Snapshot Policy Support

This directory contains the implementation for CloudStack snapshot policy resources and data sources in the Terraform CloudStack provider.

## Overview

CloudStack snapshot policies allow you to automatically create snapshots of volumes at specified intervals. This implementation provides:

- **Resource**: `cloudstack_snapshot_policy` - Create and manage snapshot policies

## Features

### Interval Types Supported
- **HOURLY**: Take snapshots every hour
- **DAILY**: Take snapshots daily  
- **WEEKLY**: Take snapshots weekly
- **MONTHLY**: Take snapshots monthly

### Key Capabilities
- Flexible scheduling with timezone support
- Multiple zone deployment support
- Configurable retention (max_snaps)
- Tag support for resource management
- Custom ID support for identification

## Files

### Core Implementation
- `resource_cloudstack_snapshot_policy.go` - Main resource implementation

### Tests
- `resource_cloudstack_snapshot_policy_test.go` - Resource tests covering all interval types

### Documentation
- `website/docs/r/snapshot_policy.html.markdown` - Resource documentation

## Usage Examples

### Basic Daily Snapshot Policy

```hcl
resource "cloudstack_snapshot_policy" "daily_backup" {
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
resource "cloudstack_snapshot_policy" "hourly_backup" {
  volume_id     = cloudstack_disk.database.id
  interval_type = "HOURLY"
  max_snaps     = 6
  schedule      = "0"  # Top of every hour
  timezone      = "UTC"
  zone_ids      = [data.cloudstack_zone.zone1.id]
  custom_id     = "db-hourly-backup"
}
```

## Schedule Formats

The schedule format depends on the interval type:

- **HOURLY**: `"MM"` (minute, 0-59)
- **DAILY**: `"HH:MM"` (hour:minute)
- **WEEKLY**: `"D:HH:MM"` (day:hour:minute, where 1=Monday)
- **MONTHLY**: `"DD:HH:MM"` (day:hour:minute, day of month)

## Testing

Run the tests using:

```bash
# Run all snapshot policy tests
go test -v ./cloudstack -run TestAccCloudStackSnapshotPolicy

# Run specific test
go test -v ./cloudstack -run TestAccCloudStackSnapshotPolicy_basic
```

## Implementation Notes

### CloudStack API Considerations

1. **Create Response**: The CloudStack API may not return a policy ID in the create response. The implementation includes fallback logic to find the created policy by listing policies for the volume.

2. **Zone ID Extraction**: The CloudStack SDK returns zone information as `[]interface{}` which requires type assertion to extract zone IDs.

3. **List by Volume**: The primary way to retrieve snapshot policies is by volume ID, not policy ID directly.

### Error Handling

The implementation includes comprehensive error handling for:
- Empty or invalid policy IDs
- API communication failures
- Missing resources
- Type assertion failures

### Debugging

Enable debug logging to troubleshoot issues:
```bash
export TF_LOG=DEBUG
terraform plan
```

The implementation includes extensive debug logging for troubleshooting API responses and data handling.

## CloudStack Compatibility

This implementation is compatible with CloudStack versions that support:
- Snapshot policy management APIs
- Zone-based snapshot policies
- Custom ID fields (where available)

Tested with CloudStack Go SDK v2.17.1+.

## Contributing

When modifying the snapshot policy implementation:

1. Update tests to cover new functionality
2. Update documentation for any new fields or behavior
3. Test against actual CloudStack environment
4. Follow existing code patterns and error handling

## Future Enhancements

Potential improvements:
- Support for storage-specific policies
- Policy validation improvements
- Bulk policy management
