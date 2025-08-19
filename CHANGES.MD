# CloudStack Terraform Provider - Recent Changes

This document summarizes the latest enhancements and new features added to the CloudStack Terraform Provider.

## 🆕 New Data Sources

### `cloudstack_domain`

- Query domain information from CloudStack
- Support for filtering domains based on various criteria
- Retrieve domain hierarchy and resource information

### `cloudstack_project`

- Query project information and metadata
- Filter projects by name, display text, and other attributes
- Access project resource limits and usage statistics

### `cloudstack_user_data`

- Query registered User Data resources
- Filter and retrieve User Data configurations
- Support for reusable User Data management

## 🆕 New Resources

### `cloudstack_project` (Resource)

- **Full CRUD operations** for CloudStack projects
- **Resource management**: CPU, memory, storage limits
- **Multi-tenancy support** with project isolation
- **Import functionality** for existing projects

### `cloudstack_resource_limit`

- Manage resource limits at domain or account level
- Support for all CloudStack resource types:
  - CPU cores, memory, storage
  - Network resources (IPs, VPCs, networks)
  - Template and snapshot limits
- **Granular control** over resource allocation

### `cloudstack_user_data` (Resource)

- Manage User Data as a separate, reusable resource
- **Template linking** capabilities
- **Versioning support** for User Data configurations
- **Centralized management** of cloud-init scripts

### `cloudstack_user_data_template_link`

- Link User Data resources to specific templates
- **Automated User Data injection** during instance creation
- **Template-specific configurations** support

### `cloudstack_user_keys`

- Manage user API keys and access credentials
- **Secure key rotation** capabilities
- **Access control** for CloudStack API usage

## 🔧 Enhanced Resources

### `cloudstack_instance`

- **New `user_data_id` parameter**: Reference registered User Data resources
- **Mutual exclusivity**: Choose between inline `user_data` or `user_data_id`
- **Improved User Data management** with better reusability

### `cloudstack_domain` (Resource)

- **Complete CRUD support** with full lifecycle management
- **Enhanced attributes**:
  - Resource limits (CPU, memory, storage, networks)
  - Resource consumption and availability metrics
  - Domain hierarchy information (parent/child relationships)
  - Creation timestamps and path information
- **Import functionality** for existing domains
- **Resource limit management** integration

### `cloudstack_private_gateway`

- **New `bypass_vlan_overlap_check` parameter**
- **Improved VLAN management** with conflict resolution
- **Enhanced network isolation** capabilities

## 🚀 Provider Improvements

### Data Source Enhancements

- **Retry mechanisms** for better handling of eventual consistency
- **Server-side filtering** for improved performance and reduced API calls
- **Enhanced error handling** with more descriptive error messages

### Code Organization

- **Alphabetical sorting** of data sources and resources in provider configuration
- **Improved maintainability** with better code structure
- **Consistent naming conventions** across all resources

## 🎯 Key Benefits

### Enterprise Features

- **Multi-tenancy support** through projects and domains
- **Resource governance** with comprehensive limit management
- **Centralized User Data management** for better operational efficiency

### Operational Improvements

- **Better resource reusability** with User Data as a first-class resource
- **Enhanced domain management** for complex organizational structures
- **Improved API key management** for secure access control

### Developer Experience

- **More intuitive resource relationships** with proper data source support
- **Better import capabilities** for existing infrastructure
- **Enhanced error handling** for easier troubleshooting

## 📋 Migration Notes

### Instance User Data

- Existing `user_data` inline configurations continue to work
- Consider migrating to `user_data_id` for better reusability
- No breaking changes to existing configurations

### Domain Management

- Enhanced domain resource now supports full CRUD operations
- Existing domain references remain compatible
- New resource limit management capabilities available

## 🔗 Documentation

New documentation files have been added:

- `website/docs/d/user_data.html.markdown` - User Data data source
- `website/docs/r/user_data.html.markdown` - User Data resource
- `website/docs/r/user_data_template_link.md` - Template linking guide

---

These changes significantly expand the provider's capabilities for enterprise CloudStack environments, providing better resource management, multi-tenancy support, and operational efficiency improvements.
