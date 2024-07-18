package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/servicemanager"
	"github.com/SAP/terraform-provider-btp/internal/validation/uuidvalidator"
)

func newSubaccountServiceBrokerResource() resource.Resource {
	return &subaccountServiceBrokerResource{}
}

type subaccountServiceBrokerResourceType struct {
	/* INPUT */
	SubaccountId types.String `tfsdk:"subaccount_id"`
	Id           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Url          types.String `tfsdk:"url"`
	Username     types.String `tfsdk:"username"`
	Password     types.String `tfsdk:"password"`
	//Labels		 types.Map `tfsdk:"labels"` // not implemented because of NGPBUG-397042

	/* OUTPUT */
	Ready        types.Bool   `tfsdk:"ready"`
	CreatedDate  types.String `tfsdk:"created_date"`
	LastModified types.String `tfsdk:"last_modified"`
}

type subaccountServiceBrokerResource struct {
	cli *btpcli.ClientFacade
}

func (rs *subaccountServiceBrokerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_service_broker", req.ProviderTypeName)
}

func (rs *subaccountServiceBrokerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	rs.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (rs *subaccountServiceBrokerResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Registers a service service broker in a subaccount.

__Tip:__
You must be assigned to the admin role of the subaccount.`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the service broker.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the service broker.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the service broker.",
				Optional:            true,
				Computed:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "The URL of the service broker.",
				Required:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "The username for basic authentication against the service broker.",
				Required:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The password for basic authentication against the service broker.",
				Required:            true,
				Sensitive:           true,
			},
			"created_date": schema.StringAttribute{
				MarkdownDescription: "The date and time when the resource was created in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_modified": schema.StringAttribute{
				MarkdownDescription: "The date and time when the resource was last modified in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
				Computed:            true,
			},
			"ready": schema.BoolAttribute{
				MarkdownDescription: "Shows whether the service broker is ready.",
				Computed:            true,
			},
		},
	}
}

func (rs *subaccountServiceBrokerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state subaccountServiceBrokerResourceType
	diags := req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, rawRes, err := rs.cli.Services.Broker.GetById(ctx, state.SubaccountId.ValueString(), state.Id.ValueString())
	if err != nil {
		handleReadErrors(ctx, rawRes, resp, err, "Resource Service Broker (Subaccount)")
		return
	}

	newState := subaccountServiceBrokerValueFrom(ctx, cliRes)
	newState.SubaccountId = state.SubaccountId
	newState.Name = state.Name
	newState.Username = state.Username
	newState.Password = state.Password

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountServiceBrokerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan subaccountServiceBrokerResourceType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliReq := btpcli.SubaccountServiceBrokerRegisterInput{
		Subaccount:  plan.SubaccountId.ValueString(),
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		URL:         plan.Url.ValueString(),
		User:        plan.Username.ValueString(),
		Password:    plan.Password.ValueString(),
	}

	cliRes, _, err := rs.cli.Services.Broker.Register(ctx, cliReq)
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Resource Service Broker (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	state := subaccountServiceBrokerValueFrom(ctx, cliRes)
	state.SubaccountId = plan.SubaccountId
	state.Name = plan.Name
	state.Username = plan.Username
	state.Password = plan.Password

	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountServiceBrokerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan subaccountServiceBrokerResourceType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliReq := btpcli.SubaccountServiceBrokerUpdateInput{
		Subaccount:  plan.SubaccountId.ValueString(),
		Id:          plan.Id.ValueString(),
		NewName:     plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		URL:         plan.Url.ValueString(),
		User:        plan.Username.ValueString(),
		Password:    plan.Password.ValueString(),
	}

	cliRes, _, err := rs.cli.Services.Broker.Update(ctx, cliReq)
	if err != nil {
		resp.Diagnostics.AddError("API Error Updating Resource Service Broker (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	newState := subaccountServiceBrokerValueFrom(ctx, cliRes)
	newState.SubaccountId = plan.SubaccountId
	newState.Name = plan.Name
	newState.Username = plan.Username
	newState.Password = plan.Password

	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, newState)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountServiceBrokerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state subaccountServiceBrokerResourceType
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := rs.cli.Services.Broker.Unregister(ctx, state.SubaccountId.ValueString(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Service Broker (Subaccount)", fmt.Sprintf("%s", err))
		return
	}
}

func (rs *subaccountServiceBrokerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: subaccount_id,id. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("subaccount_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
}

func subaccountServiceBrokerValueFrom(ctx context.Context, broker servicemanager.ServiceBrokerResponseObject) subaccountServiceBrokerResourceType {
	return subaccountServiceBrokerResourceType{
		Id:           types.StringValue(broker.Id),
		Name:         types.StringValue(broker.Name),
		Description:  types.StringValue(broker.Description),
		Url:          types.StringValue(broker.BrokerUrl),
		Ready:        types.BoolValue(broker.Ready),
		CreatedDate:  timeToValue(broker.CreatedAt),
		LastModified: timeToValue(broker.UpdatedAt),
	}
}
