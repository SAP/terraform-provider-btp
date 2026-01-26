package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/validation/jsonvalidator"
)

func newSubaccountDestinationGenericResource() resource.Resource {
	return &subaccountDestinationGenericResource{}
}

type subaccountDestinationGenericIdentityModel struct {
	SubaccountID      types.String `tfsdk:"subaccount_id"`
	Name              types.String `tfsdk:"name"`
	ServiceInstanceID types.String `tfsdk:"service_instance_id"`
}

type subaccountDestinationGenericResource struct {
	cli *btpcli.ClientFacade
}

func (rs *subaccountDestinationGenericResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_destination_generic", req.ProviderTypeName)
}

func (rs *subaccountDestinationGenericResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	rs.cli = req.ProviderData.(*btpcli.ClientFacade)
}

// Schema defines the schema for the resource.
func (rs *subaccountDestinationGenericResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Manages a destination in a SAP BTP subaccount or in the scope of a specific service instance.

__Tip:__
You must have the appropriate connectivity and destination permissions, such as:

Subaccount Administrator
Destination Administrator
Connectivity and Destination Administrator
__Scope:__
- **Subaccount-level destination**: Specify only the 'subaccount_id' and 'name' attribute.
- **Service instance-level destination**: Specify the 'subaccount_id', 'service_instance_id' and 'name' attributes.

__Notes:__
- 'service_instance_id' is optional. When omitted, the destination is created at the subaccount level.`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
			},
			"id": schema.StringAttribute{ // required by hashicorps terraform plugin testing framework for imports
				DeprecationMessage:  "Use the `subaccount_id,name,service_instance_id` attribute instead",
				MarkdownDescription: "The ID of the destination used for import operations.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The descriptive name of the destination for subaccount",
				Computed:            true,
			},
			"creation_time": schema.StringAttribute{
				MarkdownDescription: "The date and time when the resource was created",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"etag": schema.StringAttribute{
				MarkdownDescription: "The etag for the destination resource",
				Computed:            true,
			},
			"modification_time": schema.StringAttribute{
				MarkdownDescription: "The date and time when the resource was modified",
				Computed:            true,
			},
			"service_instance_id": schema.StringAttribute{
				MarkdownDescription: "The service instance that becomes part of the path used to access the destination of the subaccount.",
				Optional:            true,
			},
			"destination_configuration": schema.StringAttribute{
				MarkdownDescription: "The configuration parameters for the destination.",
				Required:            true,
				Sensitive:           true,
				CustomType:          jsontypes.NormalizedType{},
				Validators: []validator.String{
					jsonvalidator.ValidJSON(),
				},
			},
		},
	}
}

func (rs *subaccountDestinationGenericResource) IdentitySchema(_ context.Context, _ resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"subaccount_id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
			"name": identityschema.StringAttribute{
				RequiredForImport: true,
			},
			"service_instance_id": identityschema.StringAttribute{
				OptionalForImport: true,
			},
		},
	}
}

func (rs *subaccountDestinationGenericResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data subaccountDestinationGenericResourceType
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	planDestinationConfiguration := data.DestinationConfiguration

	cliRes, rawRes, err := rs.cli.Connectivity.Destination.GetBySubaccount(ctx, data.SubaccountID.ValueString(), data.Name.ValueString(), data.ServiceInstanceID.ValueString())
	if err != nil {
		handleReadErrors(ctx, rawRes, cliRes, resp, err, "Resource Destination Generic (Subaccount)")
		return
	}

	data, diags = destinationGenericResourceValueFrom(cliRes, data.SubaccountID, data.ServiceInstanceID, data.Name.ValueString())
	resp.Diagnostics.Append(diags...)

	data.DestinationConfiguration, err = MergeDestinationConfig(planDestinationConfiguration, data.DestinationConfiguration)
	if err != nil {
		resp.Diagnostics.AddError(ErrApiMergingDestinationConfiguration, fmt.Sprintf("%s", err))
		return
	}

	id := data.SubaccountID.ValueString() + "," + data.Name.ValueString() + "," + data.ServiceInstanceID.ValueString()
	data.ID = types.StringValue(id)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)

	var identity subaccountDestinationGenericIdentityModel

	diags = req.Identity.Get(ctx, &identity)
	if diags.HasError() {
		identity = subaccountDestinationGenericIdentityModel{
			SubaccountID:      data.SubaccountID,
			Name:              data.Name,
			ServiceInstanceID: data.ServiceInstanceID,
		}
	}

	diags = resp.Identity.Set(ctx, identity)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountDestinationGenericResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan subaccountDestinationGenericResourceType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	planDestinationConfiguration := plan.DestinationConfiguration

	destinationData, name, err := BuildDestinationGenericConfigurationJSON(plan)
	if err != nil {
		resp.Diagnostics.AddError("Error generating Resource Destination body", fmt.Sprintf("%s", err))
		return
	}

	_, _, err = rs.cli.Connectivity.Destination.CreateBySubaccount(ctx, plan.SubaccountID.ValueString(), destinationData, plan.ServiceInstanceID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Resource Destination", fmt.Sprintf("%s", err))
		return
	}

	cliRes, _, err := rs.cli.Connectivity.Destination.GetBySubaccount(ctx, plan.SubaccountID.ValueString(), name, plan.ServiceInstanceID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(ErrApiReadingDestination, fmt.Sprintf("%s", err))
		return
	}

	plan, diags = destinationGenericResourceValueFrom(cliRes, plan.SubaccountID, plan.ServiceInstanceID, name)
	resp.Diagnostics.Append(diags...)
	plan.DestinationConfiguration, err = MergeDestinationConfig(planDestinationConfiguration, plan.DestinationConfiguration)
	if err != nil {
		resp.Diagnostics.AddError(ErrApiMergingDestinationConfiguration, fmt.Sprintf("%s", err))
		return
	}

	id := plan.SubaccountID.ValueString() + "," + name + "," + plan.ServiceInstanceID.ValueString()
	plan.ID = types.StringValue(id)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	identity := subaccountDestinationGenericIdentityModel{
		SubaccountID:      plan.SubaccountID,
		Name:              types.StringValue(name),
		ServiceInstanceID: plan.ServiceInstanceID,
	}

	diags = resp.Identity.Set(ctx, identity)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountDestinationGenericResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan subaccountDestinationGenericResourceType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	planDestinationConfiguration := plan.DestinationConfiguration
	destinationData, name, err := BuildDestinationGenericConfigurationJSON(plan)
	if err != nil {
		resp.Diagnostics.AddError("Error generating Resource Destination body", fmt.Sprintf("%s", err))
		return
	}

	_, _, err = rs.cli.Connectivity.Destination.UpdateBySubaccount(ctx, plan.SubaccountID.ValueString(), destinationData, plan.ServiceInstanceID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Updating Resource Destination", fmt.Sprintf("%s", err))
		return
	}

	cliRes, _, err := rs.cli.Connectivity.Destination.GetBySubaccount(ctx, plan.SubaccountID.ValueString(), name, plan.ServiceInstanceID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(ErrApiReadingDestination, fmt.Sprintf("%s", err))
		return
	}

	plan, diags = destinationGenericResourceValueFrom(cliRes, plan.SubaccountID, plan.ServiceInstanceID, name)
	resp.Diagnostics.Append(diags...)
	plan.DestinationConfiguration, err = MergeDestinationConfig(planDestinationConfiguration, plan.DestinationConfiguration)
	if err != nil {
		resp.Diagnostics.AddError(ErrApiMergingDestinationConfiguration, fmt.Sprintf("%s", err))
		return
	}

	id := plan.SubaccountID.ValueString() + "," + name + "," + plan.ServiceInstanceID.ValueString()
	plan.ID = types.StringValue(id)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	identity := subaccountDestinationGenericIdentityModel{
		SubaccountID:      plan.SubaccountID,
		Name:              types.StringValue(name),
		ServiceInstanceID: plan.ServiceInstanceID,
	}

	diags = resp.Identity.Set(ctx, identity)
	resp.Diagnostics.Append(diags...)

}

func (rs *subaccountDestinationGenericResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state subaccountDestinationGenericResourceType
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	_, name, err := BuildDestinationGenericConfigurationJSON(state)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving name", fmt.Sprintf("%s", err))
		return
	}

	_, _, err = rs.cli.Connectivity.Destination.DeleteBySubaccount(ctx, state.SubaccountID.ValueString(), name, state.ServiceInstanceID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Destination", fmt.Sprintf("%s", err))
		return
	}
}

func (rs *subaccountDestinationGenericResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID != "" {

		idParts := strings.Split(req.ID, ",")

		switch len(idParts) {
		case 2:
			if idParts[0] == "" || idParts[1] == "" {
				resp.Diagnostics.AddError(
					ErrUnexpectedImportIdentifier,
					fmt.Sprintf("Expected import identifier with format: subaccount_id, name. Got: %q", req.ID),
				)
				return
			}

			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("subaccount_id"), idParts[0])...)
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[1])...)
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_instance_id"), types.StringNull())...)

			return

		case 3:
			if idParts[0] == "" || idParts[1] == "" || idParts[2] == "" {
				resp.Diagnostics.AddError(
					ErrUnexpectedImportIdentifier,
					fmt.Sprintf("Expected import identifier with format: subaccount_id, name, service_instance_id. Got: %q", req.ID),
				)
				return
			}

			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("subaccount_id"), idParts[0])...)
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[1])...)
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_instance_id"), idParts[2])...)

			return

		default:
			resp.Diagnostics.AddError(
				ErrUnexpectedImportIdentifier,
				fmt.Sprintf(
					"Expected one of:\n  - subaccount_id,name\n  - subaccount_id,name,service_instance_id\nGot: %q",
					req.ID,
				),
			)
			return
		}
	}

	var identity subaccountDestinationGenericIdentityModel
	diags := resp.Identity.Get(ctx, &identity)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("subaccount_id"), identity.SubaccountID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), identity.Name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_instance_id"), identity.ServiceInstanceID)...)
}
