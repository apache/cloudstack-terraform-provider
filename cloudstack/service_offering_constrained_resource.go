package cloudstack

import (
	"context"
	"fmt"
	"strconv"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ resource.Resource              = &serviceOfferingConstrainedResource{}
	_ resource.ResourceWithConfigure = &serviceOfferingConstrainedResource{}
)

func NewserviceOfferingConstrainedResource() resource.Resource {
	return &serviceOfferingConstrainedResource{}
}

type serviceOfferingConstrainedResource struct {
	client *cloudstack.CloudStackClient
}

func (r *serviceOfferingConstrainedResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: serviceOfferingMergeCommonSchema(map[string]schema.Attribute{
			"cpu_speed": schema.Int32Attribute{
				Description: "VMware and Xen based hypervisors this is the CPU speed of the service offering in MHz. For the KVM hypervisor the values of the parameters cpuSpeed and cpuNumber will be used to calculate the `shares` value. This value is used by the KVM hypervisor to calculate how much time the VM will have access to the host's CPU. The `shares` value does not have a unit, and its purpose is being a weight value for the host to compare between its guest VMs. For more information, see https://libvirt.org/formatdomain.html#cpu-tuning.",
				Required:    true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"max_cpu_number": schema.Int32Attribute{
				Description: "The maximum number of CPUs to be set with Custom Compute Offering",
				Required:    true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"max_memory": schema.Int32Attribute{
				Description: "The maximum memory size of the custom service offering in MB",
				Required:    true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"min_cpu_number": schema.Int32Attribute{
				Description: "The minimum number of CPUs to be set with Custom Compute Offering",
				Required:    true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"min_memory": schema.Int32Attribute{
				Description: "The minimum memory size of the custom service offering in MB",
				Required:    true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
		}),
	}
}

func (r *serviceOfferingConstrainedResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan serviceOfferingConstrainedResourceModel
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

	// common params
	params := r.client.ServiceOffering.NewCreateServiceOfferingParams(plan.DisplayText.ValueString(), plan.Name.ValueString())
	plan.commonCreateParams(ctx, params)
	planDiskQosHypervisor.commonCreateParams(ctx, params)
	planDiskOffering.commonCreateParams(ctx, params)
	planDiskQosStorage.commonCreateParams(ctx, params)

	// resource specific params
	if !plan.CpuSpeed.IsNull() {
		params.SetCpuspeed(int(plan.CpuSpeed.ValueInt32()))
	}
	if !plan.MaxCpuNumber.IsNull() {
		params.SetMaxcpunumber(int(plan.MaxCpuNumber.ValueInt32()))
	}
	if !plan.MaxMemory.IsNull() {
		params.SetMaxmemory(int(plan.MaxMemory.ValueInt32()))
	}
	if !plan.MinCpuNumber.IsNull() {
		params.SetMincpunumber(int(plan.MinCpuNumber.ValueInt32()))
	}
	if !plan.MinMemory.IsNull() {
		params.SetMinmemory(int(plan.MinMemory.ValueInt32()))
	}

	// create offering
	cs, err := r.client.ServiceOffering.CreateServiceOffering(params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating service offering",
			"Could not create constrained offering, unexpected error: "+err.Error(),
		)
		return
	}

	//
	plan.Id = types.StringValue(cs.Id)

	//
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *serviceOfferingConstrainedResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state serviceOfferingConstrainedResourceModel
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
			"Could not read constrained service offering, unexpected error: "+err.Error(),
		)
		return
	}

	// resource specific
	if cs.Cpuspeed > 0 {
		state.CpuSpeed = types.Int32Value(int32(cs.Cpuspeed))
	}
	if v, found := cs.Serviceofferingdetails["maxcpunumber"]; found {
		i, err := strconv.Atoi(v)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading service offering",
				"Could not read constrained service offering max cpu, unexpected error: "+err.Error(),
			)
			return
		}
		state.MaxCpuNumber = types.Int32Value(int32(i))
	}
	if v, found := cs.Serviceofferingdetails["mincpunumber"]; found {
		i, err := strconv.Atoi(v)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading service offering",
				"Could not read constrained service offering min cpu number, unexpected error: "+err.Error(),
			)
			return
		}
		state.MinCpuNumber = types.Int32Value(int32(i))
	}
	if v, found := cs.Serviceofferingdetails["maxmemory"]; found {
		i, err := strconv.Atoi(v)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading service offering",
				"Could not read constrained service offering max memory, unexpected error: "+err.Error(),
			)
			return
		}
		state.MaxMemory = types.Int32Value(int32(i))
	}
	if v, found := cs.Serviceofferingdetails["minmemory"]; found {
		i, err := strconv.Atoi(v)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading service offering",
				"Could not read constrained service offering min memory, unexpected error: "+err.Error(),
			)
			return
		}
		state.MinMemory = types.Int32Value(int32(i))
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
func (r *serviceOfferingConstrainedResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state serviceOfferingConstrainedResourceModel

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
			"Could not update constrained service offering, unexpected error: "+err.Error(),
		)
		return
	}

	state.commonUpdate(ctx, cs)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *serviceOfferingConstrainedResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state serviceOfferingConstrainedResourceModel

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
			"Could not delete constrained offering, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *serviceOfferingConstrainedResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *serviceOfferingConstrainedResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_offering_constrained"
}
