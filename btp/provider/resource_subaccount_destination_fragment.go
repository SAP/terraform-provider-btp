package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/connectivity"
	"github.com/SAP/terraform-provider-btp/internal/validation/uuidvalidator"
)

func newSubaccountDestinationFragmentResource() resource.Resource {
	return &subaccountDestinationFragmentResource{}
}

type subaccountDestinationFragmentResourceConfig struct {
	// INPUT
	SubaccountID      types.String `tfsdk:"subaccount_id"`
	Name              types.String `tfsdk:"name"`
	ServiceInstanceID types.String `tfsdk:"service_instance_id"`
	// OUTPUT
	DestinationFragment types.Map `tfsdk:"fragment_content"`
}

type subaccountDestinationFragmentIdentityModel struct {
	SubaccountID      types.String `tfsdk:"subaccount_id"`
	Name              types.String `tfsdk:"name"`
	ServiceInstanceID types.String `tfsdk:"service_instance_id"`
}

type subaccountDestinationFragmentResource struct {
	cli *btpcli.ClientFacade
}

func (rs *subaccountDestinationFragmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_destination_fragment", req.ProviderTypeName)
}

func (rs *subaccountDestinationFragmentResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	rs.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (rs *subaccountDestinationFragmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Manages a destination fragment in a SAP BTP subaccount or in the scope of a specific service instance.

__Tip:__
You must have the appropriate connectivity and destination permissions, such as:
- Subaccount Administrator  
- Destination Administrator  
- Connectivity and Destination Administrator

To learn more about these roles, see the SAP Help documentation:  
https://help.sap.com/docs/connectivity/sap-btp-connectivity-cf/user-roles?q=role+collection

__Scope:__
- **Subaccount-level fragment**: Specify only the 'subaccount_id' and 'name' attribute.
- **Service instance-level fragment**: Specify the 'subaccount_id', 'service_instance_id' and 'name' attributes.

__Notes:__
- 'service_instance_id' is optional. When omitted, the fragment is created at the subaccount level.
- The fragment content is defined using the 'fragment_content' map and the API requires the 'FragmentName' field, which is automatically handled by the provider.`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the destination fragment.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"service_instance_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the service instance.",
				Optional:            true,
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"fragment_content": schema.MapAttribute{
				MarkdownDescription: "The content of the destination fragment.",
				ElementType:         types.StringType,
				Optional:            true,
			},
		},
	}
}

func (rs *subaccountDestinationFragmentResource) IdentitySchema(_ context.Context, _ resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
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

func (rs *subaccountDestinationFragmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data subaccountDestinationFragmentResourceConfig
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	hasServiceInstance := !data.ServiceInstanceID.IsNull() && !data.ServiceInstanceID.IsUnknown() && data.ServiceInstanceID.ValueString() != ""

	cliRes, err := rs.readFragment(ctx, data, hasServiceInstance, resp.Diagnostics)
	if err != nil {
		return
	}

	delete(cliRes.Content, "FragmentName")

	if !data.DestinationFragment.IsNull() {
		data.DestinationFragment, diags = types.MapValueFrom(ctx, types.StringType, cliRes.Content)
		resp.Diagnostics.Append(diags...)
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)

	identity := subaccountDestinationFragmentIdentityModel{
		SubaccountID:      data.SubaccountID,
		Name:              data.Name,
		ServiceInstanceID: data.ServiceInstanceID,
	}
	diags = resp.Identity.Set(ctx, identity)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountDestinationFragmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan subaccountDestinationFragmentResourceConfig
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	rawContent, diags := extractFragmentContent(plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := buildPayload(plan.Name.ValueString(), rawContent)

	hasServiceInstance := !plan.ServiceInstanceID.IsNull() && !plan.ServiceInstanceID.IsUnknown() && plan.ServiceInstanceID.ValueString() != ""

	err := rs.createFragment(ctx, plan, payload, hasServiceInstance, resp.Diagnostics)
	if err != nil {
		return
	}

	destinationFragmentDetails, err := rs.readFragment(ctx, plan, hasServiceInstance, resp.Diagnostics)
	if err != nil {
		return
	}

	delete(destinationFragmentDetails.Content, "FragmentName")

	if !plan.DestinationFragment.IsNull() {
		plan.DestinationFragment, diags = types.MapValueFrom(ctx, types.StringType, destinationFragmentDetails.Content)
		resp.Diagnostics.Append(diags...)
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	identity := subaccountDestinationFragmentIdentityModel{
		SubaccountID:      plan.SubaccountID,
		Name:              plan.Name,
		ServiceInstanceID: plan.ServiceInstanceID,
	}

	diags = resp.Identity.Set(ctx, identity)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountDestinationFragmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan subaccountDestinationFragmentResourceConfig
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	rawContent, diags := extractFragmentContent(plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := buildPayload(plan.Name.ValueString(), rawContent)

	hasServiceInstance := !plan.ServiceInstanceID.IsNull() && !plan.ServiceInstanceID.IsUnknown() && plan.ServiceInstanceID.ValueString() != ""

	err := rs.updateFragment(ctx, plan, payload, hasServiceInstance, resp.Diagnostics)
	if err != nil {
		return
	}

	destinationFragmentDetails, err := rs.readFragment(ctx, plan, hasServiceInstance, resp.Diagnostics)
	if err != nil {
		return
	}

	delete(destinationFragmentDetails.Content, "FragmentName")

	if !plan.DestinationFragment.IsNull() {
		plan.DestinationFragment, diags = types.MapValueFrom(ctx, types.StringType, destinationFragmentDetails.Content)
		resp.Diagnostics.Append(diags...)
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	identity := subaccountDestinationFragmentIdentityModel{
		SubaccountID:      plan.SubaccountID,
		Name:              plan.Name,
		ServiceInstanceID: plan.ServiceInstanceID,
	}

	diags = resp.Identity.Set(ctx, identity)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountDestinationFragmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state subaccountDestinationFragmentResourceConfig
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	hasServiceInstance := !state.ServiceInstanceID.IsNull() && !state.ServiceInstanceID.IsUnknown() && state.ServiceInstanceID.ValueString() != ""

	if hasServiceInstance {
		_, _, err := rs.cli.Connectivity.DestinationFragment.DeleteByServiceInstance(ctx, state.SubaccountID.ValueString(), state.Name.ValueString(), state.ServiceInstanceID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("API Error Deleting Destination Fragment at Service Instance Level", fmt.Sprintf("%s", err))
			return
		}
	} else {
		_, _, err := rs.cli.Connectivity.DestinationFragment.DeleteBySubaccount(ctx, state.SubaccountID.ValueString(), state.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("API Error Deleting Destination Fragment at Subaccount Level", fmt.Sprintf("%s", err))
			return
		}
	}
}

func (rs *subaccountDestinationFragmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var identity subaccountDestinationFragmentIdentityModel
	diags := resp.Identity.Get(ctx, &identity)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("subaccount_id"), identity.SubaccountID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), identity.Name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_instance_id"), identity.ServiceInstanceID)...)
}

func extractFragmentContent(plan subaccountDestinationFragmentResourceConfig) (map[string]string, diag.Diagnostics) {
	rawContent := map[string]string{}
	if !plan.DestinationFragment.IsNull() && !plan.DestinationFragment.IsUnknown() {
		diags := plan.DestinationFragment.ElementsAs(context.Background(), &rawContent, false)
		if diags.HasError() {
			return nil, diags
		}
	}
	return rawContent, diag.Diagnostics{}
}

func buildPayload(name string, rawContent map[string]string) map[string]string {
	payload := map[string]string{
		"FragmentName": name,
	}

	for k, v := range rawContent {
		payload[k] = v
	}
	return payload
}

func (rs *subaccountDestinationFragmentResource) createFragment(ctx context.Context, plan subaccountDestinationFragmentResourceConfig, payload map[string]string, hasServiceInstance bool, respDiags diag.Diagnostics) error {
	if hasServiceInstance {
		_, _, err := rs.cli.Connectivity.DestinationFragment.CreateByServiceInstance(ctx, plan.SubaccountID.ValueString(), plan.ServiceInstanceID.ValueString(), payload)
		if err != nil {
			respDiags.AddError("API Error Creating Destination Fragment at Service Instance Level", fmt.Sprintf("%s", err))
		}
		return err
	}

	_, _, err := rs.cli.Connectivity.DestinationFragment.CreateBySubaccount(ctx, plan.SubaccountID.ValueString(), payload)
	if err != nil {
		respDiags.AddError("API Error Creating Destination Fragment at Subaccount Level", fmt.Sprintf("%s", err))

	}
	return err
}

func (rs *subaccountDestinationFragmentResource) updateFragment(ctx context.Context, plan subaccountDestinationFragmentResourceConfig, payload map[string]string, hasServiceInstance bool, respDiags diag.Diagnostics) error {
	if hasServiceInstance {
		_, _, err := rs.cli.Connectivity.DestinationFragment.UpdateByServiceInstance(ctx, plan.SubaccountID.ValueString(), plan.ServiceInstanceID.ValueString(), payload)
		if err != nil {
			respDiags.AddError("API Error Updating Destination Fragment at Service Instance Level", fmt.Sprintf("%s", err))
		}
		return err
	}

	_, _, err := rs.cli.Connectivity.DestinationFragment.UpdateBySubaccount(ctx, plan.SubaccountID.ValueString(), payload)
	if err != nil {
		respDiags.AddError("API Error Updating Destination Fragment at Subaccount Level", fmt.Sprintf("%s", err))

	}
	return err
}

func (rs *subaccountDestinationFragmentResource) readFragment(ctx context.Context, plan subaccountDestinationFragmentResourceConfig, hasServiceInstance bool, respDiags diag.Diagnostics) (connectivity.DestinationFragment, error) {
	if hasServiceInstance {
		destinationFragmentDetails, _, err := rs.cli.Connectivity.DestinationFragment.GetByServiceInstance(ctx, plan.SubaccountID.ValueString(), plan.Name.ValueString(), plan.ServiceInstanceID.ValueString())
		if err != nil {
			respDiags.AddError("API Error Reading Destination Fragment at Service Instance Level", fmt.Sprintf("%s", err))
			return connectivity.DestinationFragment{}, err
		}
		return destinationFragmentDetails, nil
	}

	destinationFragmentDetails, _, err := rs.cli.Connectivity.DestinationFragment.GetBySubaccount(ctx, plan.SubaccountID.ValueString(), plan.Name.ValueString())
	if err != nil {
		respDiags.AddError("API Error Reading Destination Fragment at Subaccount Level", fmt.Sprintf("%s", err))
		return connectivity.DestinationFragment{}, err
	}
	return destinationFragmentDetails, nil
}
