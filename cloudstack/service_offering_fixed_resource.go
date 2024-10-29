package cloudstack

import (
	"context"
	"fmt"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ resource.Resource              = &serviceOfferingFixedResource{}
	_ resource.ResourceWithConfigure = &serviceOfferingFixedResource{}
)

func NewserviceOfferingFixedResource() resource.Resource {
	return &serviceOfferingFixedResource{}
}

type serviceOfferingFixedResource struct {
	client *cloudstack.CloudStackClient
}

func (r *serviceOfferingFixedResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: serviceOfferingMergeCommonSchema(map[string]schema.Attribute{
			"cpu_number": schema.Int32Attribute{
				Description: "Number of CPU cores",
				Required:    true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"cpu_speed": schema.Int32Attribute{
				Description: "the CPU speed of the service offering in MHz.  This does not apply to KVM.",
				Required:    true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"memory": schema.Int32Attribute{
				Description: "the total memory of the service offering in MB",
				Required:    true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
		}),
	}
}

func (r *serviceOfferingFixedResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan serviceOfferingFixedResourceModel
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

	// resource specific params
	if !plan.CpuNumber.IsNull() {
		params.SetCpunumber(int(plan.CpuNumber.ValueInt32()))
	}
	if !plan.CpuSpeed.IsNull() {
		params.SetCpuspeed(int(plan.CpuSpeed.ValueInt32()))
	}
	if !plan.Memory.IsNull() {
		params.SetMemory(int(plan.Memory.ValueInt32()))
	}

	// create offering
	cs, err := r.client.ServiceOffering.CreateServiceOffering(params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating service offering",
			"Could not create fixed offering, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Id = types.StringValue(cs.Id)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *serviceOfferingFixedResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var state serviceOfferingFixedResourceModel
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
			"Error reading service offering",
			"Could not read fixed service offering, unexpected error: "+err.Error(),
		)
		return
	}

	// resource specific
	if cs.Cpunumber > 0 {
		state.CpuNumber = types.Int32Value(int32(cs.Cpunumber))
	}
	if cs.Cpuspeed > 0 {
		state.CpuSpeed = types.Int32Value(int32(cs.Cpuspeed))
	}
	if cs.Memory > 0 {
		state.Memory = types.Int32Value(int32(cs.Memory))
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

// Update updates the resource and sets the updated Terraform state on success.
func (r *serviceOfferingFixedResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state serviceOfferingFixedResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := r.client.ServiceOffering.NewUpdateServiceOfferingParams(state.Id.ValueString())
	state.commonUpdateParams(ctx, params)

	cs, err := r.client.ServiceOffering.UpdateServiceOffering(params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating fixed service offering",
			"Could not update fixed service offering, unexpected error: "+err.Error(),
		)
		return
	}

	state.commonUpdate(ctx, cs)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *serviceOfferingFixedResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state serviceOfferingFixedResourceModel

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
			"Could not delete fixed offering, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *serviceOfferingFixedResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
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

// Metadata returns the resource type name.
func (r *serviceOfferingFixedResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_offering_fixed"
}
