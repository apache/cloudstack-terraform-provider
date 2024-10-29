package cloudstack

import "github.com/hashicorp/terraform-plugin-framework/types"

type serviceOfferingConstrainedResourceModel struct {
	CpuSpeed     types.Int32 `tfsdk:"cpu_speed"`
	MaxCpuNumber types.Int32 `tfsdk:"max_cpu_number"`
	MaxMemory    types.Int32 `tfsdk:"max_memory"`
	MinCpuNumber types.Int32 `tfsdk:"min_cpu_number"`
	MinMemory    types.Int32 `tfsdk:"min_memory"`
	serviceOfferingCommonResourceModel
}

type serviceOfferingUnconstrainedResourceModel struct {
	serviceOfferingCommonResourceModel
}

type serviceOfferingFixedResourceModel struct {
	CpuNumber types.Int32 `tfsdk:"cpu_number"`
	CpuSpeed  types.Int32 `tfsdk:"cpu_speed"`
	Memory    types.Int32 `tfsdk:"memory"`
	serviceOfferingCommonResourceModel
}

type serviceOfferingCommonResourceModel struct {
	DeploymentPlanner                types.String `tfsdk:"deployment_planner"`
	DiskOfferingId                   types.String `tfsdk:"disk_offering_id"`
	DisplayText                      types.String `tfsdk:"display_text"`
	DomainIds                        types.Set    `tfsdk:"domain_ids"`
	DynamicScalingEnabled            types.Bool   `tfsdk:"dynamic_scaling_enabled"`
	HostTags                         types.String `tfsdk:"host_tags"`
	Id                               types.String `tfsdk:"id"`
	IsVolatile                       types.Bool   `tfsdk:"is_volatile"`
	LimitCpuUse                      types.Bool   `tfsdk:"limit_cpu_use"`
	Name                             types.String `tfsdk:"name"`
	NetworkRate                      types.Int32  `tfsdk:"network_rate"`
	OfferHa                          types.Bool   `tfsdk:"offer_ha"`
	StorageTags                      types.String `tfsdk:"storage_tags"`
	ZoneIds                          types.Set    `tfsdk:"zone_ids"`
	ServiceOfferingDiskQosHypervisor types.Object `tfsdk:"disk_hypervisor"`
	ServiceOfferingDiskOffering      types.Object `tfsdk:"disk_offering"`
	ServiceOfferingDiskQosStorage    types.Object `tfsdk:"disk_storage"`
}

type ServiceOfferingDiskQosHypervisor struct {
	DiskBytesReadRate           types.Int64 `tfsdk:"bytes_read_rate"`
	DiskBytesReadRateMax        types.Int64 `tfsdk:"bytes_read_rate_max"`
	DiskBytesReadRateMaxLength  types.Int64 `tfsdk:"bytes_read_rate_max_length"`
	DiskBytesWriteRate          types.Int64 `tfsdk:"bytes_write_rate"`
	DiskBytesWriteRateMax       types.Int64 `tfsdk:"bytes_write_rate_max"`
	DiskBytesWriteRateMaxLength types.Int64 `tfsdk:"bytes_write_rate_max_length"`
}

type ServiceOfferingDiskOffering struct {
	CacheMode              types.String `tfsdk:"cache_mode"`
	DiskOfferingStrictness types.Bool   `tfsdk:"disk_offering_strictness"`
	ProvisionType          types.String `tfsdk:"provisioning_type"`
	RootDiskSize           types.Int64  `tfsdk:"root_disk_size"`
	StorageType            types.String `tfsdk:"storage_type"`
}

type ServiceOfferingDiskQosStorage struct {
	CustomizedIops            types.Bool  `tfsdk:"customized_iops"`
	HypervisorSnapshotReserve types.Int32 `tfsdk:"hypervisor_snapshot_reserve"`
	MaxIops                   types.Int64 `tfsdk:"max_iops"`
	MinIops                   types.Int64 `tfsdk:"min_iops"`
}
