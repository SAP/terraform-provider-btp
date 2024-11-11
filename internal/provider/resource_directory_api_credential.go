package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/validation/uuidvalidator"
)

func newDirectoryApiCredentialResource() resource.Resource {
	return &directoryApiCredentialResource{}
}

type directoryApiCredentialResource struct {
	cli *btpcli.ClientFacade
}

func (rs *directoryApiCredentialResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_directory_api_credential", req.ProviderTypeName)
}

func (rs *directoryApiCredentialResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	rs.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (rs *directoryApiCredentialResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Manage API Credentials at the Directory level. These credentials will enable you to consume
		the REST APIs of the SAP Authorization and Trust Management service (XSUAA).
		With the client ID and client secret, or certificate, you can request an access token for the APIs in the targeted directory.

__Tip:__
You must be assigned to directory admin or viewer role.

__Further documentation:__
<https://help.sap.com/docs/btp/sap-business-technology-platform-internal/managing-api-credentials-for-calling-rest-apis-of-sap-authorization-and-trust-management-service>`,
		Attributes: map[string]schema.Attribute{
			"directory_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the directory.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name for the API credential.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z\d-]+$`), "can contain only alphanumberic values and dashes."),
				},
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "A unique ID associated with the API credential.",
				Computed:            true,
			},
			"credential_type": schema.StringAttribute{
				MarkdownDescription: "The supported credential types are Secrets (Default) or Certificates.",
				Computed:            true,
			},
			"certificate_passed": schema.StringAttribute{
				MarkdownDescription: "If the user prefers to use a certificate, they must provide the certificate value in PEM format \"----BEGIN CERTIFICATE-----...-----END CERTIFICATE-----\".",
				Optional:            true,
			},
			"certificate_received": schema.StringAttribute{
				MarkdownDescription: "The certificate that is computed based on the one passed by the user.",
				Computed:            true,
			},
			"client_secret": schema.StringAttribute{
				MarkdownDescription: "If the certificate is omitted, then a unique secret is generated for the API credential.",
				Computed:            true,
			},
			"key": schema.StringAttribute{
				MarkdownDescription: "RSA key generated if the API credential is created with a certificate.",
				Computed:            true,
			},
			"read_only": schema.BoolAttribute{
				MarkdownDescription: "Access restriction placed on the API credential. If set to true, the resource has only read-only access.",
				Optional:            true,
				Computed:            true,
			},
			"token_url": schema.StringAttribute{
				MarkdownDescription: "The URL to be used to fetch the access token to make use of the XSUAA REST APIs.",
				Computed:            true,
			},
			"api_url": schema.StringAttribute{
				MarkdownDescription: "The URL to be used to make the API calls.",
				Computed:            true,
			},
		},
	}
}

func (rs *directoryApiCredentialResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan directoryApiCredentialType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := rs.cli.Security.ApiCredential.CreateByDirectoryorSubaccount(ctx, &btpcli.ApiCredentialInput{
		Directory:   plan.DirectoryId.ValueString(),
		Name:        plan.Name.ValueString(),
		Certificate: plan.CertificatePassed.ValueString(),
		ReadOnly:    plan.ReadOnly.ValueBool(),
	})

	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("%s", err))
		return
	}

	updatedPlan, diags := directoryApiCredentialFromValue(ctx, cliRes)
	resp.Diagnostics.Append(diags...)

	updatedPlan.CertificatePassed = plan.CertificatePassed

	diags = resp.State.Set(ctx, &updatedPlan)
	resp.Diagnostics.Append(diags...)
}

func (rs *directoryApiCredentialResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state directoryApiCredentialType
	diags := req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, rawRes, err := rs.cli.Security.ApiCredential.GetByDirectoryorSubaccount(ctx, &btpcli.ApiCredentialInput{
		Directory: state.DirectoryId.ValueString(),
		Name:      state.Name.ValueString(),
	})
	if err != nil {
		handleReadErrors(ctx, rawRes, resp, err, "Resource Api Credential (Directory)")
		return
	}

	newState, diags := directoryApiCredentialFromValue(ctx, cliRes)
	resp.Diagnostics.Append(diags...)

	//The below parameters are not returned by the get call to the Api Credential
	newState.DirectoryId = state.DirectoryId
	if !state.CertificatePassed.IsUnknown() {
		newState.CertificatePassed = state.CertificatePassed
		newState.Key = state.Key
	} else {
		newState.ClientSecret = state.ClientSecret
	}

	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
}

func (rs *directoryApiCredentialResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// There is currently no API call that supports the update of the Api credentials
}

func (rs *directoryApiCredentialResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state directoryApiCredentialType
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, _, err := rs.cli.Security.ApiCredential.DeleteByDirectoryorSubaccount(ctx, &btpcli.ApiCredentialInput{
		Directory: state.DirectoryId.ValueString(),
		Name:      state.Name.ValueString(),
	})

	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Api Credential (Directory)", fmt.Sprintf("%s", err))
		return
	}
}
