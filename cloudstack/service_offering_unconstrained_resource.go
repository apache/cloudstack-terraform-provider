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

import (
	"context"
	"fmt"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ resource.Resource              = &serviceOfferingUnconstrainedResource{}
	_ resource.ResourceWithConfigure = &serviceOfferingUnconstrainedResource{}
)

func NewserviceOfferingUnconstrainedResource() resource.Resource {
	return &serviceOfferingUnconstrainedResource{}
}

type serviceOfferingUnconstrainedResource struct {
	client *cloudstack.CloudStackClient
}

// Schema defines the schema for the resource.
func (r *serviceOfferingUnconstrainedResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: serviceOfferingMergeCommonSchema(map[string]schema.Attribute{}),
	}
}

func (r *serviceOfferingUnconstrainedResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan serviceOfferingUnconstrainedResourceModel
	var planDiskQosHypervisor ServiceOfferingDiskQosHypervisor
	var planDiskOffering ServiceOfferingDiskOffering
	var planDiskQosStorage ServiceOfferingDiskQosStorage

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if !plan.ServiceOfferingDiskQosHypervisor.IsNull() {
		resp.Diagnostics.Append(plan.ServiceOfferingDiskQosHypervisor.As(ctx, &planDiskQosHypervisor, basetypes.ObjectAsOptions{})...)
	}
	if !plan.ServiceOfferingDiskOffering.IsNull() {
		resp.Diagnostics.Append(plan.ServiceOfferingDiskOffering.As(ctx, &planDiskOffering, basetypes.ObjectAsOptions{})...)
	}
	if !plan.ServiceOfferingDiskQosStorage.IsNull() {
		resp.Diagnostics.Append(plan.ServiceOfferingDiskQosStorage.As(ctx, &planDiskQosStorage, basetypes.ObjectAsOptions{})...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// cloudstack params
	params := r.client.ServiceOffering.NewCreateServiceOfferingParams(plan.DisplayText.ValueString(), plan.Name.ValueString())
	plan.commonCreateParams(ctx, params)
	planDiskQosHypervisor.commonCreateParams(ctx, params)
	planDiskOffering.commonCreateParams(ctx, params)
	planDiskQosStorage.commonCreateParams(ctx, params)

	// create offering
	cs, err := r.client.ServiceOffering.CreateServiceOffering(params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating service offering",
			"Could not create unconstrained offering, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Id = types.StringValue(cs.Id)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

}

func (r *serviceOfferingUnconstrainedResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state serviceOfferingUnconstrainedResourceModel
	var stateDiskQosHypervisor ServiceOfferingDiskQosHypervisor
	var stateDiskOffering ServiceOfferingDiskOffering
	var stateDiskQosStorage ServiceOfferingDiskQosStorage

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if !state.ServiceOfferingDiskQosHypervisor.IsNull() {
		resp.Diagnostics.Append(state.ServiceOfferingDiskQosHypervisor.As(ctx, &stateDiskQosHypervisor, basetypes.ObjectAsOptions{})...)
	}
	if !state.ServiceOfferingDiskOffering.IsNull() {
		resp.Diagnostics.Append(state.ServiceOfferingDiskOffering.As(ctx, &stateDiskOffering, basetypes.ObjectAsOptions{})...)
	}
	if !state.ServiceOfferingDiskQosStorage.IsNull() {
		resp.Diagnostics.Append(state.ServiceOfferingDiskQosStorage.As(ctx, &stateDiskQosStorage, basetypes.ObjectAsOptions{})...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	cs, _, err := r.client.ServiceOffering.GetServiceOfferingByID(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating service offering",
			"Could not read unconstrained service offering, unexpected error: "+err.Error(),
		)
		return
	}

	state.commonRead(ctx, cs)
	stateDiskQosHypervisor.commonRead(ctx, cs)
	stateDiskOffering.commonRead(ctx, cs)
	stateDiskQosStorage.commonRead(ctx, cs)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

}

func (r *serviceOfferingUnconstrainedResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state serviceOfferingUnconstrainedResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := r.client.ServiceOffering.NewUpdateServiceOfferingParams(state.Id.ValueString())
	state.commonUpdateParams(ctx, params)

	cs, err := r.client.ServiceOffering.UpdateServiceOffering(params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating service offering",
			"Could not update unconstrained service offering, unexpected error: "+err.Error(),
		)
		return
	}

	state.commonUpdate(ctx, cs)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

}

func (r *serviceOfferingUnconstrainedResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state serviceOfferingUnconstrainedResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the service offering
	_, err := r.client.ServiceOffering.DeleteServiceOffering(r.client.ServiceOffering.NewDeleteServiceOfferingParams(state.Id.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting service offering",
			"Could not delete unconstrained offering, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *serviceOfferingUnconstrainedResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*cloudstack.CloudStackClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *cloudstack.CloudStackClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *serviceOfferingUnconstrainedResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_offering_unconstrained"
}
