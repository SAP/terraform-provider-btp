package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/validation/jsonvalidator"
)

func newGlobalaccountResourceProviderResource() resource.Resource {
	return &resourceGlobalaccountProviderResource{}
}

type resourceGlobalaccountProviderResource struct {
	cli *btpcli.ClientFacade
}

func (rs *resourceGlobalaccountProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_globalaccount_resource_provider", req.ProviderTypeName)
}

func (rs *resourceGlobalaccountProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	rs.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (rs *resourceGlobalaccountProviderResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Create a resource provider instance to allow your global account to connect to your provider account on a non-SAP cloud vendor. Through this channel, you can then consume remote service resources that you already own and which are supported by SAP BTP.
For example, if you are subscribed to Amazon Web Services (AWS) and have already purchased services, such as PostgreSQL, you can register the vendor as a resource provider in SAP BTP and consume this service across your subaccounts together with other services offered by SAP.

The use of this functionality is subject to the availability of the supported non-SAP cloud vendors in your country or region.

__Tips__
* You must be assigned to the global account admin role.
* You can create more than one instance of a given resource provider, each with its unique configuration properties. In such cases, the display name and technical name should be descriptive enough so that you and developers can easily differentiate between each instance.
* After you configure a new resource provider instance, its supported services are added as entitlements in your global account.

__Further documentation__
https://help.sap.com/viewer/65de2977205c403bbc107264b8eccf4b/Cloud/en-US/e2c250dc5abd468a81f4f619206157a2.html`,
		Attributes: map[string]schema.Attribute{
			"resource_provider": schema.StringAttribute{
				MarkdownDescription: "Provider of the requested resource. Possible values are: `AWS`, `AZURE`.",
				Required:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique technical name of the resource provider.",
				Required:            true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Descriptive name of the resource provider.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the resource provider.",
				Computed:            true,
			},
			"parameters": schema.StringAttribute{
				MarkdownDescription: "Any relevant information about the resource provider that is not provided by other parameter values.",
				Required:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					jsonvalidator.ValidJSON(),
				},
			},
		},
	}
}

func (rs *resourceGlobalaccountProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state globalaccountResourceProviderType

	diags := req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := rs.cli.Accounts.ResourceProvider.Get(ctx, state.ResourceProvider.ValueString(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Resource Provider (Global Account)", fmt.Sprintf("%s", err))
		return
	}

	state, diags = globalaccountResourceProviderValueFrom(ctx, cliRes)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (rs *resourceGlobalaccountProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan globalaccountResourceProviderType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := rs.cli.Accounts.ResourceProvider.Create(ctx, btpcli.GlobalaccountResourceProviderCreateInput{
		Provider:          plan.ResourceProvider.ValueString(),
		TechnicalName:     plan.Id.ValueString(),
		DisplayName:       plan.DisplayName.ValueString(),
		Description:       plan.Description.ValueString(),
		ConfigurationInfo: plan.Parameters.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Resource Resource Provider (Global Account)", fmt.Sprintf("%s", err))
		return
	}

	plan, diags = globalaccountResourceProviderValueFrom(ctx, cliRes)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (rs *resourceGlobalaccountProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan globalaccountResourceProviderType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.AddError("API Error Updating Resource Resource Provider (Global Account)", "Update is not yet implemented.")

	/* TODO: implementation of UPDATE operation
	cliRes, err := gen.client.Execute(ctx, btpcli.Update, gen.command, plan)
	if err != nil {
		resp.Diagnostics.AddError("API Error Updating Resource Resource Provider (Global Account)", fmt.Sprintf("%s", err))
		return
	}*/

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (rs *resourceGlobalaccountProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state globalaccountResourceProviderType
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, _, err := rs.cli.Accounts.ResourceProvider.Delete(ctx, state.ResourceProvider.ValueString(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Resource Provider (Global Account)", fmt.Sprintf("%s", err))
		return
	}
}

func (rs *resourceGlobalaccountProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: resource_provider,unique_technical_name. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("resource_provider"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
}
