package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
)

func newSubaccountDestinationCertificateResource() resource.Resource {
	return &subaccountDestinationCertificateResource{}
}

type subaccountDestinationCertificateResource struct {
	cli *btpcli.ClientFacade
}

func (rs *subaccountDestinationCertificateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_destination_certificate", req.ProviderTypeName)
}

func (rs *subaccountDestinationCertificateResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	rs.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (rs *subaccountDestinationCertificateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Creates a destination certificate in a subaccount.
		
		__Tip:__
		You must be assigned the Destination Admin or the Destination Certificate Administrator role in the subaccount.

		__Further Information:__
		<https://help.sap.com/docs/connectivity/sap-btp-connectivity-cf/use-destination-certificates>`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount which contains the certificate.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"certificate_name": schema.StringAttribute{
				MarkdownDescription: "The name of the certificate with a valid certificate extension. Supported certificate types include .pem, .p12, .jks and .pfx",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"certificate_content": schema.StringAttribute{
				MarkdownDescription: "The content of the certificate in base64 format.",
				Required:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"service_instance_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the service instance associated with this certificate.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"certificate_nodes": schema.ListNestedAttribute{
				MarkdownDescription: "List of certificate nodes containing details about the certificate and private key components.",
				Computed:            true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							MarkdownDescription: "Denotes the type of the node i.e., 'private_key' or 'x509_certificate'.",
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"format": schema.StringAttribute{
							MarkdownDescription: "The format of the certificate or key.",
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"algorithm": schema.StringAttribute{
							MarkdownDescription: "The cryptographic algorithm used in the key.",
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"alias": schema.StringAttribute{
							MarkdownDescription: "An identifier used for the key.",
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"subject": schema.StringAttribute{
							MarkdownDescription: "The certificate subject which identifies the owner of the certificate.",
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"issuer": schema.StringAttribute{
							MarkdownDescription: "The certificate issuer which identifies the certificate authority that signed the certificate.",
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"common_name": schema.StringAttribute{
							MarkdownDescription: "The common name (CN) extracted from the certificate subject. May be null if not specified in the certificate.",
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"not_before": schema.StringAttribute{
							MarkdownDescription: "The start date and time (in ISO 8601 format) from which the certificate is valid.",
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"not_after": schema.StringAttribute{
							MarkdownDescription: "The expiration date and time (in ISO 8601 format) after which the certificate is no longer valid.",
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"certificate": schema.StringAttribute{
							MarkdownDescription: "The complete X.509 certificate in PEM format, encoded as base64.",
							Computed:            true,
							Sensitive:           true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
					},
				},
			},
			"certification_creation_details": schema.SingleNestedAttribute{
				MarkdownDescription: "Details about how the destination certificate was created and its configuration settings.",
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"generation_method": schema.StringAttribute{
						MarkdownDescription: "Specifies the method used to create the certificate.",
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"common_name": schema.StringAttribute{
						MarkdownDescription: "The common name (CN) associated with the certificate.",
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"has_password": schema.BoolAttribute{
						MarkdownDescription: "Indicates whether the certificate is protected with a password.",
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"auto_renew": schema.BoolAttribute{
						MarkdownDescription: "Specifies whether the certificate is automatically renewed before it expires.",
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"validity_duration": schema.StringAttribute{
						MarkdownDescription: "The numeric duration for which the certificate is valid.",
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"validity_time_units": schema.StringAttribute{
						MarkdownDescription: "The time unit associated with the validity duration, such as `DAYS`, `MONTHS`, or `YEARS`.",
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
		},
	}
}

func (rs *subaccountDestinationCertificateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var data subaccountDestinationCertificateResourceType

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := rs.cli.Connectivity.DestinationCertificate.Get(ctx, &btpcli.DestinationCertificateGetInput{
		SubaccountId:      data.SubaccountId.ValueString(),
		ServiceInstanceId: data.ServiceInstanceId.ValueString(),
		CertificateName:   data.CertificateName.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Destination Certificate", fmt.Sprintf("%s", err))
		return
	}

	state, diags := subaccountDestinationCertificateValueFrom(ctx, cliRes)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.SubaccountId = data.SubaccountId
	state.ServiceInstanceId = data.ServiceInstanceId
	state.CertificateContent = data.CertificateContent

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

}

func (rs *subaccountDestinationCertificateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var plan subaccountDestinationCertificateResourceType

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := rs.cli.Connectivity.DestinationCertificate.Create(ctx, &btpcli.DestinationCertificateCreateInput{
		SubaccountId:      plan.SubaccountId.ValueString(),
		ServiceInstanceId: plan.ServiceInstanceId.ValueString(),
		Certificate: btpcli.FileInput{
			Filename:           plan.CertificateName.ValueString(),
			CertificateContent: plan.CertificateContent.ValueString(),
		},
	})

	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Resource Destination Certificate", fmt.Sprintf("%s", err))
		return
	}

	state, diags := subaccountDestinationCertificateValueFrom(ctx, cliRes)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.SubaccountId = plan.SubaccountId
	state.ServiceInstanceId = plan.ServiceInstanceId
	state.CertificateContent = plan.CertificateContent

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountDestinationCertificateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	resp.Diagnostics.AddError("Resource Destination Certificate does not support updates", "Terraform will destroy and recreate the resource if any of the user configurable parameters are modified")
}

func (rs *subaccountDestinationCertificateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data subaccountDestinationCertificateResourceType

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := rs.cli.Connectivity.DestinationCertificate.Delete(ctx, &btpcli.DestinationCertificateGetInput{
		SubaccountId:      data.SubaccountId.ValueString(),
		ServiceInstanceId: data.ServiceInstanceId.ValueString(),
		CertificateName:   data.CertificateName.ValueString(),
	})

	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Destination Certificate", fmt.Sprintf("%s", err))
		return
	}
}
