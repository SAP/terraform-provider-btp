package provider

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	// "github.com/hashicorp/terraform-plugin-framework/diag"
	// "github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	// "github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	// "github.com/hashicorp/terraform-plugin-framework/tfsdk"

	// "github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/validation/uuidvalidator"
)

func newSubaccountApiCredentialResource() resource.Resource {
	return &subaccountApiCredentialResource{}
}

type subaccountApiCredentialResource struct {
	cli *btpcli.ClientFacade
}

func (rs *subaccountApiCredentialResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_api_credential", req.ProviderTypeName)
}

func (rs *subaccountApiCredentialResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	rs.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (rs *subaccountApiCredentialResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Assigns the entitlement plan of a service, multitenant application, or environment, to a directory. Note that some environments, such as Cloud Foundry, are available by default to all global accounts and their directorys, and therefore are not made available as entitlements.

__Tip:__
You must be assigned to the global account admin or viewer role.

__Further documentation:__
<https://help.sap.com/docs/btp/sap-business-technology-platform/entitlements-and-quotas>`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"name" : schema.StringAttribute{
				MarkdownDescription: "The name for the api-credential.",
				Optional: 			 true,
				Computed: 			 true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z\d-]+$`), "can contain only alphanumberic values and dashes."),
				},
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "A unique ID associated with the api-credential.",
				Computed:            true,
			},
			"credential_type": schema.StringAttribute{
				MarkdownDescription: "The supported credential types are Secrets (Default) or Certificates.",
				Computed: 			 true,
			},
			"certificate_passed": schema.StringAttribute{
				MarkdownDescription:  "If the user prefers to use a certificate, they must provide the certificate value in PEM format \"----BEGIN CERTIFICATE-----...-----END CERTIFICATE-----\".",
				Optional: true,
			},
			"certificate_received": schema.StringAttribute{
				MarkdownDescription:  "The certificate that is computed based on the one passed by the user.",
				Computed: true,
			},
			"client_secret": schema.StringAttribute{
				MarkdownDescription: "If the certificate is omitted, then a unique secret is generated for the api-credential.",
				Computed: 			 true,
			},
			"key": schema.StringAttribute{
				MarkdownDescription: "RSA key generated if the api-credential is created with a certificate.",
				Computed: true,
			},
			"read_only": schema.BoolAttribute{
				MarkdownDescription: "Access restriction placed on the api-credential. If set to true, the resource has only read-only access. ",
				Optional:            true,
				Computed: 			 true,
			},
			"token_url": schema.StringAttribute{
				MarkdownDescription: "The URL that must used to fetch the access token to make use of the APIs.",
				Computed:            true,
			},
			"api_url": schema.StringAttribute{
				MarkdownDescription: "The URL that must be used to make the API calls.",
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

func (rs *subaccountApiCredentialResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan apiCredentialType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := rs.cli.Security.ApiCredential.CreateBySubaccount(ctx, &btpcli.ApiCredentialInput{
		SubaccountId:     plan.SubaccountId.ValueString(),
		Name:             plan.Name.ValueString(),
		Certificate: 	  plan.CertificatePassed.ValueString(),
		ReadOnly:		  plan.ReadOnly.ValueBool(),
	})
	
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("%s", err))
		return
	}

	updatedPlan, diags := subaccountApiCredentialFromValue(ctx, cliRes)
	resp.Diagnostics.Append(diags...)

	updatedPlan.CertificatePassed = plan.CertificatePassed

	diags = resp.State.Set(ctx, &updatedPlan)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountApiCredentialResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state apiCredentialType
	diags := req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError(){
		return
	}

	cliRes, rawRes, err := rs.cli.Security.ApiCredential.GetBySubaccount(ctx, &btpcli.ApiCredentialInput{
		SubaccountId: state.SubaccountId.ValueString(),
		Name:		  state.Name.ValueString(),
	})  
	if err!=nil {
		handleReadErrors(ctx, rawRes, resp, err, "Resource Api Credential (Subaccount)")
		return
	}

	newState, diags := subaccountApiCredentialFromValue(ctx, cliRes)
	resp.Diagnostics.Append(diags...)

	newState.SubaccountId = state.SubaccountId
	newState.CertificatePassed = state.CertificatePassed

	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountApiCredentialResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
// There is currently no API call that supports the update of the Api credentials
}

func (rs *subaccountApiCredentialResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state apiCredentialType
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, _, err := rs.cli.Security.ApiCredential.DeleteBySubaccount(ctx, &btpcli.ApiCredentialInput{
		SubaccountId: state.SubaccountId.ValueString(),
		Name:		  state.Name.ValueString(),
	})

	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Api Credential (Subaccount)", fmt.Sprintf("%s", err))
		return
	}
}

func (rs *subaccountApiCredentialResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: subaccount_id, name. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("subaccount_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[1])...)
}
