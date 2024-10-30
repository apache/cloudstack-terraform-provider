// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package cloudstack

import (
	"context"
	"strings"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ------------------------------------------------------------------------------------------------------------------------------
// Common update methods
// -
func (state *serviceOfferingCommonResourceModel) commonUpdate(ctx context.Context, cs *cloudstack.UpdateServiceOfferingResponse) {
	if cs.Displaytext != "" {
		state.DisplayText = types.StringValue(cs.Displaytext)
	}
	if cs.Domainid != "" {
		state.DomainIds, _ = types.SetValueFrom(ctx, types.StringType, strings.Split(cs.Domainid, ","))
	}
	if cs.Hosttags != "" {
		state.HostTags = types.StringValue(cs.Hosttags)
	}
	if cs.Name != "" {
		state.Name = types.StringValue(cs.Name)
	}
	if cs.Zoneid != "" {
		state.ZoneIds, _ = types.SetValueFrom(ctx, types.StringType, strings.Split(cs.Zoneid, ","))
	}
}

func (plan *serviceOfferingCommonResourceModel) commonUpdateParams(ctx context.Context, p *cloudstack.UpdateServiceOfferingParams) *cloudstack.UpdateServiceOfferingParams {
	if !plan.DisplayText.IsNull() {
		p.SetDisplaytext(plan.DisplayText.ValueString())
	}
	if !plan.DomainIds.IsNull() {
		p.SetDomainid(plan.DomainIds.String())
	}
	if !plan.HostTags.IsNull() {
		p.SetHosttags(plan.HostTags.ValueString())
	}
	if !plan.Name.IsNull() {
		p.SetName(plan.Name.ValueString())
	}
	if !plan.ZoneIds.IsNull() && len(plan.ZoneIds.Elements()) > 0 {
		p.SetZoneid(plan.ZoneIds.String())
	} else {
		p.SetZoneid("all")
	}

	return p

}

// ------------------------------------------------------------------------------------------------------------------------------
// common Read methods
// -
func (state *serviceOfferingCommonResourceModel) commonRead(ctx context.Context, cs *cloudstack.ServiceOffering) {
	state.Id = types.StringValue(cs.Id)

	if cs.Deploymentplanner != "" {
		state.DeploymentPlanner = types.StringValue(cs.Deploymentplanner)
	}
	if cs.Diskofferingid != "" {
		state.DiskOfferingId = types.StringValue(cs.Diskofferingid)
	}
	if cs.Displaytext != "" {
		state.DisplayText = types.StringValue(cs.Displaytext)
	}
	if cs.Domainid != "" {
		state.DomainIds, _ = types.SetValueFrom(ctx, types.StringType, strings.Split(cs.Domainid, ","))
	}
	if cs.Hosttags != "" {
		state.HostTags = types.StringValue(cs.Hosttags)
	}
	if cs.Name != "" {
		state.Name = types.StringValue(cs.Name)
	}
	if cs.Networkrate > 0 {
		state.NetworkRate = types.Int32Value(int32(cs.Networkrate))
	}
	if cs.Zoneid != "" {
		state.ZoneIds, _ = types.SetValueFrom(ctx, types.StringType, strings.Split(cs.Zoneid, ","))
	}

	state.DynamicScalingEnabled = types.BoolValue(cs.Dynamicscalingenabled)
	state.IsVolatile = types.BoolValue(cs.Isvolatile)
	state.LimitCpuUse = types.BoolValue(cs.Limitcpuuse)
	state.OfferHa = types.BoolValue(cs.Offerha)

}

func (state *ServiceOfferingDiskQosHypervisor) commonRead(ctx context.Context, cs *cloudstack.ServiceOffering) {
	if cs.DiskBytesReadRate > 0 {
		state.DiskBytesReadRate = types.Int64Value(cs.DiskBytesReadRate)
	}
	if cs.DiskBytesReadRateMax > 0 {
		state.DiskBytesReadRateMax = types.Int64Value(cs.DiskBytesReadRateMax)
	}
	if cs.DiskBytesReadRateMaxLength > 0 {
		state.DiskBytesReadRateMaxLength = types.Int64Value(cs.DiskBytesReadRateMaxLength)
	}
	if cs.DiskBytesWriteRate > 0 {
		state.DiskBytesWriteRate = types.Int64Value(cs.DiskBytesWriteRate)
	}
	if cs.DiskBytesWriteRateMax > 0 {
		state.DiskBytesWriteRateMax = types.Int64Value(cs.DiskBytesWriteRateMax)
	}
	if cs.DiskBytesWriteRateMaxLength > 0 {
		state.DiskBytesWriteRateMaxLength = types.Int64Value(cs.DiskBytesWriteRateMaxLength)
	}

}

func (state *ServiceOfferingDiskOffering) commonRead(ctx context.Context, cs *cloudstack.ServiceOffering) {

	if cs.CacheMode != "" {
		state.CacheMode = types.StringValue(cs.CacheMode)
	}
	if cs.Diskofferingstrictness {
		state.DiskOfferingStrictness = types.BoolValue(cs.Diskofferingstrictness)
	}
	if cs.Provisioningtype != "" {
		state.ProvisionType = types.StringValue(cs.Provisioningtype)
	}
	if cs.Rootdisksize > 0 {
		state.RootDiskSize = types.Int64Value(cs.Rootdisksize)
	}
	if cs.Storagetype != "" {
		state.StorageType = types.StringValue(cs.Storagetype)
	}
	if cs.Storagetags != "" {
		state.StorageTags = types.StringValue(cs.Storagetags)
	}
}

func (state *ServiceOfferingDiskQosStorage) commonRead(ctx context.Context, cs *cloudstack.ServiceOffering) {
	if cs.Iscustomizediops {
		state.CustomizedIops = types.BoolValue(cs.Iscustomizediops)
	}
	if cs.Hypervisorsnapshotreserve > 0 {
		state.HypervisorSnapshotReserve = types.Int32Value(int32(cs.Hypervisorsnapshotreserve))
	}
	if cs.Maxiops > 0 {
		state.MaxIops = types.Int64Value(cs.Maxiops)
	}
	if cs.Miniops > 0 {
		state.MinIops = types.Int64Value(cs.Miniops)
	}

}

// ------------------------------------------------------------------------------------------------------------------------------
// common Create methods
// -
func (plan *serviceOfferingCommonResourceModel) commonCreateParams(ctx context.Context, p *cloudstack.CreateServiceOfferingParams) *cloudstack.CreateServiceOfferingParams {
	if !plan.DeploymentPlanner.IsNull() && !plan.DeploymentPlanner.IsUnknown() {
		p.SetDeploymentplanner(plan.DeploymentPlanner.ValueString())
	} else {
		plan.DeploymentPlanner = types.StringNull()
	}
	if !plan.DiskOfferingId.IsNull() {
		p.SetDiskofferingid(plan.DiskOfferingId.ValueString())
	}
	if !plan.DomainIds.IsNull() {
		domainids := make([]string, len(plan.DomainIds.Elements()))
		plan.DomainIds.ElementsAs(ctx, &domainids, false)
		p.SetDomainid(domainids)
	}
	if !plan.DynamicScalingEnabled.IsNull() {
		p.SetDynamicscalingenabled(plan.DynamicScalingEnabled.ValueBool())
	}
	if !plan.HostTags.IsNull() {
		p.SetHosttags(plan.HostTags.ValueString())
	}
	if !plan.IsVolatile.IsNull() {
		p.SetIsvolatile(plan.IsVolatile.ValueBool())
	}
	if !plan.LimitCpuUse.IsNull() {
		p.SetLimitcpuuse(plan.LimitCpuUse.ValueBool())
	}
	if !plan.NetworkRate.IsNull() {
		p.SetNetworkrate(int(plan.NetworkRate.ValueInt32()))
	}
	if !plan.OfferHa.IsNull() {
		p.SetOfferha(plan.OfferHa.ValueBool())
	}
	if !plan.ZoneIds.IsNull() {
		zoneIds := make([]string, len(plan.ZoneIds.Elements()))
		plan.ZoneIds.ElementsAs(ctx, &zoneIds, false)
		p.SetZoneid(zoneIds)
	}

	return p

}
func (plan *ServiceOfferingDiskQosHypervisor) commonCreateParams(ctx context.Context, p *cloudstack.CreateServiceOfferingParams) *cloudstack.CreateServiceOfferingParams {
	if !plan.DiskBytesReadRate.IsNull() {
		p.SetBytesreadrate(plan.DiskBytesReadRate.ValueInt64())
	}
	if !plan.DiskBytesReadRateMax.IsNull() {
		p.SetBytesreadratemax(plan.DiskBytesReadRateMax.ValueInt64())
	}
	if !plan.DiskBytesReadRateMaxLength.IsNull() {
		p.SetBytesreadratemaxlength(plan.DiskBytesReadRateMaxLength.ValueInt64())
	}
	if !plan.DiskBytesWriteRate.IsNull() {
		p.SetByteswriterate(plan.DiskBytesWriteRate.ValueInt64())
	}
	if !plan.DiskBytesWriteRateMax.IsNull() {
		p.SetByteswriteratemax(plan.DiskBytesWriteRateMax.ValueInt64())
	}
	if !plan.DiskBytesWriteRateMaxLength.IsNull() {
		p.SetByteswriteratemaxlength(plan.DiskBytesWriteRateMaxLength.ValueInt64())
	}

	return p
}

func (plan *ServiceOfferingDiskOffering) commonCreateParams(ctx context.Context, p *cloudstack.CreateServiceOfferingParams) *cloudstack.CreateServiceOfferingParams {

	if !plan.CacheMode.IsNull() {
		p.SetCachemode(plan.CacheMode.ValueString())
	}
	if !plan.DiskOfferingStrictness.IsNull() {
		p.SetDiskofferingstrictness(plan.DiskOfferingStrictness.ValueBool())
	}
	if !plan.ProvisionType.IsNull() {
		p.SetProvisioningtype(plan.ProvisionType.ValueString())
	}
	if !plan.RootDiskSize.IsNull() {
		p.SetRootdisksize(plan.RootDiskSize.ValueInt64())
	}
	if !plan.StorageType.IsNull() {
		p.SetStoragetype(plan.StorageType.ValueString())
	}
	if !plan.StorageTags.IsNull() {
		p.SetTags(plan.StorageTags.ValueString())
	}

	return p

}

func (plan *ServiceOfferingDiskQosStorage) commonCreateParams(ctx context.Context, p *cloudstack.CreateServiceOfferingParams) *cloudstack.CreateServiceOfferingParams {
	if !plan.CustomizedIops.IsNull() {
		p.SetCustomizediops(plan.CustomizedIops.ValueBool())
	}
	if !plan.HypervisorSnapshotReserve.IsNull() {
		p.SetHypervisorsnapshotreserve(int(plan.HypervisorSnapshotReserve.ValueInt32()))
	}
	if !plan.MaxIops.IsNull() {
		p.SetMaxiops(int64(plan.MaxIops.ValueInt64()))
	}
	if !plan.MinIops.IsNull() {
		p.SetMiniops((plan.MinIops.ValueInt64()))
	}

	return p
}
