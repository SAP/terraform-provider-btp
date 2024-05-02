package provider

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	// "github.com/hashicorp/terraform-plugin-framework/diag"
	// "github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	// "github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	// "github.com/hashicorp/terraform-plugin-framework/tfsdk"

	// "github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/cis_entitlements"
	"github.com/SAP/terraform-provider-btp/internal/tfutils"
	"github.com/SAP/terraform-provider-btp/internal/validation/uuidvalidator"
)

func newApiCredentialResource() resource.Resource {
	return &apiCredentialResource{}
}

type apiCredentialResource struct {
	cli *btpcli.ClientFacade
}

func (rs *apiCredentialResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_api_credential", req.ProviderTypeName)
}

func (rs *apiCredentialResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	rs.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (rs *apiCredentialResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Assigns the entitlement plan of a service, multitenant application, or environment, to a directory. Note that some environments, such as Cloud Foundry, are available by default to all global accounts and their directorys, and therefore are not made available as entitlements.

__Tip:__
You must be assigned to the global account admin or viewer role.

__Further documentation:__
<https://help.sap.com/docs/btp/sap-business-technology-platform/entitlements-and-quotas>`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the directory.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"name" : schema.StringAttribute{
				MarkdownDescription: "The ID of the directory.",
				Optional: 			 true,
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the secret.",
				Computed:            true,
			},
			"credential_type": schema.StringAttribute{
				MarkdownDescription: "The name of the entitled service.",
				Computed: 			 true,
			},
			"client_secret": schema.StringAttribute{
				MarkdownDescription: "The name of the entitled service plan.",
				Computed: 			 true,
				Optional:			 true,
			},
			"certificate": schema.StringAttribute{
				Optional: true,
			},
			"key": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"read_only": schema.BoolAttribute{
				MarkdownDescription: "The quota assigned to the directory.",
				Optional:            true,
			},
			"token_url": schema.StringAttribute{
				MarkdownDescription: "Determines whether the plans of entitlements that have a numeric quota with the amount specified in `auto_distribute_amount` are automatically allocated to any new subaccount that is added to the directory in the future. For entitlements without a numeric quota, it shows if the plan are assigned to any new subaccount that is added to the directory in the future (`auto_distribute_amount` is not needed). If the `distribute` parameter is set, the same assignment is also made to all subaccounts currently in the directory. Entitlements are subject to available quota in the directory.",
				Computed:            true,
			},
			"api_url": schema.StringAttribute{
				MarkdownDescription: "The quota of the specified plan automatically allocated to any new subaccount that is created in the future in the directory. When applying this option, `auto_assign` and/or `distribute` must also be set. Applies only to entitlements that have a numeric quota.",
				Computed:            true,
			},
			"xsapp_name": schema.StringAttribute{
				Computed: true,
			},
			"service_instance_id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (rs *apiCredentialResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan apiCredentialType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := rs.cli.Security.ApiCredential.CreateBySubaccount(ctx, &btpcli.ApiCredentialCreateInput{
		SubaccountId:     plan.SubaccountId.ValueString(),
		Name:             plan.Name.ValueString(),
		Certificate: 	  plan.Certificate.ValueString(),
		ReadOnly:		  plan.ReadOnly.ValueBool(),
	})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("%s", err))
		return
	}

	updatedPlan, diags := subaccountApiCredentialFromValue(ctx, cliRes)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &updatedPlan)
	resp.Diagnostics.Append(diags...)
}

func (rs *apiCredentialResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (rs *apiCredentialResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (rs *apiCredentialResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state directoryEntitlementType
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var err error
	if !hasPlanQuotaDir(state) {
		_, err = rs.cli.Accounts.Entitlement.DisableInDirectory(ctx, state.DirectoryId.ValueString(), state.ServiceName.ValueString(), state.PlanName.ValueString(), state.Distribute.ValueBool(), state.AutoAssign.ValueBool())
	} else {

		dirAssignmentInput := btpcli.DirectoryAssignmentInput{
			DirectoryId:          state.DirectoryId.ValueString(),
			ServiceName:          state.ServiceName.ValueString(),
			ServicePlanName:      state.PlanName.ValueString(),
			Amount:               0,
			Distribute:           false,
			AutoAssign:           false,
			AutoDistributeAmount: 0,
		}
		_, err = rs.cli.Accounts.Entitlement.AssignToDirectory(ctx, dirAssignmentInput)
	}

	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Entitlement (Directory)", fmt.Sprintf("%s", err))
		return
	}

	deleteStateConf := &tfutils.StateChangeConf{
		Pending: []string{cis_entitlements.StateStarted, cis_entitlements.StateProcessing},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {

			entitlement, _, err := rs.cli.Accounts.Entitlement.GetEntitledByDirectory(ctx, state.DirectoryId.ValueString(), state.ServiceName.ValueString(), state.PlanName.ValueString())

			if reflect.ValueOf(entitlement).IsNil() {
				return entitlement, "DELETED", nil
			}

			if err != nil {
				return entitlement, cis_entitlements.StateProcessingFailed, err
			}

			return entitlement, cis_entitlements.StateProcessing, nil
		},
		Timeout:    10 * time.Minute,
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	_, err = deleteStateConf.WaitForStateContext(ctx)

	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Entitlement (Directory)", fmt.Sprintf("%s", err))
		return
	}
}

func (rs *apiCredentialResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	//Known Gap: The DISTRIBUTE flag cannot be fetched via the platform APIs. Hence, we cannot import the value, it will always be FALSE
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 3 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: directory,service_name,plan_name. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("directory_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), idParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("plan_name"), idParts[2])...)
}

