//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

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
	StorageTags            types.String `tfsdk:"storage_tags"`
}

type ServiceOfferingDiskQosStorage struct {
	CustomizedIops            types.Bool  `tfsdk:"customized_iops"`
	HypervisorSnapshotReserve types.Int32 `tfsdk:"hypervisor_snapshot_reserve"`
	MaxIops                   types.Int64 `tfsdk:"max_iops"`
	MinIops                   types.Int64 `tfsdk:"min_iops"`
}
