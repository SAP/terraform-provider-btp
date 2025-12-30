package provider

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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

const ErrUnexpectedImportIdentifier = "Unexpected Import Identifier"
const ErrApiReadingDestination = "API Error Reading destination"

func newSubaccountDestinationResource() resource.Resource {
	return &subaccountDestinationResource{}
}

type subaccountDestinationIdentityModel struct {
	SubaccountID      types.String `tfsdk:"subaccount_id"`
	Name              types.String `tfsdk:"name"`
	ServiceInstanceID types.String `tfsdk:"service_instance_id"`
}

type subaccountDestinationResource struct {
	cli *btpcli.ClientFacade
}

func (rs *subaccountDestinationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_destination", req.ProviderTypeName)
}

func (rs *subaccountDestinationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	rs.cli = req.ProviderData.(*btpcli.ClientFacade)
}

// Schema defines the schema for the resource.
func (rs *subaccountDestinationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"creation_time": schema.StringAttribute{
				MarkdownDescription: "The date and time when the resource was created in",
				Computed:            true,
			},
			"etag": schema.StringAttribute{
				MarkdownDescription: "The etag for the destination resource",
				Computed:            true,
			},
			"modification_time": schema.StringAttribute{
				MarkdownDescription: "The date and time when the resource was modified",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The descriptive name of the destination for subaccount",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^[^\/]{1,255}$`), "must not contain '/', not be empty and not exceed 255 characters"),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of request from destination.",
				Required:            true,
			},
			"proxy_type": schema.StringAttribute{
				MarkdownDescription: "The proxytype of the destination.",
				Optional:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "The url of the destination.",
				Optional:            true,
			},
			"authentication": schema.StringAttribute{
				MarkdownDescription: "The authentication of the destination.",
				Optional:            true,
			},
			"service_instance_id": schema.StringAttribute{
				MarkdownDescription: "The service instance that becomes part of the path used to access the destination of the subaccount.",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the destination.",
				Optional:            true,
			},
			"additional_configuration": schema.StringAttribute{
				MarkdownDescription: "The additional configuration parameters for the destination.",
				Optional:            true,
				CustomType:          jsontypes.NormalizedType{},
				Validators: []validator.String{
					jsonvalidator.ValidJSON(),
				},
			},
		},
	}
}

func (rs *subaccountDestinationResource) IdentitySchema(_ context.Context, _ resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
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

func (rs *subaccountDestinationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data subaccountDestinationResourceType
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	oldAdditionalConfiguration := data.AdditionalConfiguration

	cliRes, _, err := rs.cli.Connectivity.Destination.GetBySubaccount(ctx, data.SubaccountID.ValueString(), data.Name.ValueString(), data.ServiceInstanceID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(ErrApiReadingDestination, fmt.Sprintf("%s", err))
		return
	}

	data, diags = destinationResourceValueFrom(cliRes, data.SubaccountID, data.ServiceInstanceID)
	resp.Diagnostics.Append(diags...)

	data.AdditionalConfiguration, err = MergeAdditionalConfig(oldAdditionalConfiguration, data.AdditionalConfiguration)
	if err != nil {
		log.Fatal(err)
	}

	id := data.SubaccountID.ValueString() + "," + data.Name.ValueString() + "," + data.ServiceInstanceID.ValueString()
	data.ID = types.StringValue(id)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)

	identity := subaccountDestinationIdentityModel{
		SubaccountID:      data.SubaccountID,
		Name:              data.Name,
		ServiceInstanceID: data.ServiceInstanceID,
	}

	diags = resp.Identity.Set(ctx, identity)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountDestinationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan subaccountDestinationResourceType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	oldAdditionalConfiguration := plan.AdditionalConfiguration

	destinationData, err := BuildDestinationConfigurationJSON(plan)
	if err != nil {
		resp.Diagnostics.AddError("Error generating Resource Destination body", fmt.Sprintf("%s", err))
		return
	}

	_, _, err = rs.cli.Connectivity.Destination.CreateBySubaccount(ctx, plan.SubaccountID.ValueString(), destinationData, plan.ServiceInstanceID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Resource Destination", fmt.Sprintf("%s", err))
		return
	}

	cliRes, _, err := rs.cli.Connectivity.Destination.GetBySubaccount(ctx, plan.SubaccountID.ValueString(), plan.Name.ValueString(), plan.ServiceInstanceID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(ErrApiReadingDestination, fmt.Sprintf("%s", err))
		return
	}

	plan, diags = destinationResourceValueFrom(cliRes, plan.SubaccountID, plan.ServiceInstanceID)
	resp.Diagnostics.Append(diags...)
	plan.AdditionalConfiguration, err = MergeAdditionalConfig(oldAdditionalConfiguration, plan.AdditionalConfiguration)
	if err != nil {
		log.Fatal(err)
	}

	id := plan.SubaccountID.ValueString() + "," + plan.Name.ValueString() + "," + plan.ServiceInstanceID.ValueString()
	plan.ID = types.StringValue(id)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	identity := subaccountDestinationIdentityModel{
		SubaccountID:      plan.SubaccountID,
		Name:              plan.Name,
		ServiceInstanceID: plan.ServiceInstanceID,
	}

	diags = resp.Identity.Set(ctx, identity)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountDestinationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan subaccountDestinationResourceType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	oldAdditionalConfiguration := plan.AdditionalConfiguration
	destinationData, err := BuildDestinationConfigurationJSON(plan)
	if err != nil {
		resp.Diagnostics.AddError("Error generating Resource Destination body", fmt.Sprintf("%s", err))
		return
	}

	_, _, err = rs.cli.Connectivity.Destination.UpdateBySubaccount(ctx, plan.SubaccountID.ValueString(), destinationData, plan.ServiceInstanceID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Updating Resource Destination", fmt.Sprintf("%s", err))
		return
	}

	cliRes, _, err := rs.cli.Connectivity.Destination.GetBySubaccount(ctx, plan.SubaccountID.ValueString(), plan.Name.ValueString(), plan.ServiceInstanceID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(ErrApiReadingDestination, fmt.Sprintf("%s", err))
		return
	}

	plan, diags = destinationResourceValueFrom(cliRes, plan.SubaccountID, plan.ServiceInstanceID)
	resp.Diagnostics.Append(diags...)
	plan.AdditionalConfiguration, err = MergeAdditionalConfig(oldAdditionalConfiguration, plan.AdditionalConfiguration)
	if err != nil {
		log.Fatal(err)
	}

	id := plan.SubaccountID.ValueString() + "," + plan.Name.ValueString() + "," + plan.ServiceInstanceID.ValueString()
	plan.ID = types.StringValue(id)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	identity := subaccountDestinationIdentityModel{
		SubaccountID:      plan.SubaccountID,
		Name:              plan.Name,
		ServiceInstanceID: plan.ServiceInstanceID,
	}

	diags = resp.Identity.Set(ctx, identity)
	resp.Diagnostics.Append(diags...)

}

func (rs *subaccountDestinationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state subaccountDestinationResourceType
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, _, err := rs.cli.Connectivity.Destination.DeleteBySubaccount(ctx, state.SubaccountID.ValueString(), state.Name.ValueString(), state.ServiceInstanceID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Destination", fmt.Sprintf("%s", err))
		return
	}
}

func (rs *subaccountDestinationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

	var identity subaccountDestinationIdentityModel
	diags := resp.Identity.Get(ctx, &identity)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("subaccount_id"), identity.SubaccountID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), identity.Name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_instance_id"), identity.ServiceInstanceID)...)
}
