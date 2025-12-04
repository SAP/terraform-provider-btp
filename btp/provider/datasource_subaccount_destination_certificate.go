package provider

import (
	"context"
	"fmt"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func newSubaccountDestinationCertificateDataSource() datasource.DataSource {
	return &subaccountDestinationCertificateDataSource{}
}

type subaccountDestinationCertificateDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountDestinationCertificateDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_destination_certificate", req.ProviderTypeName)
}

func (ds *subaccountDestinationCertificateDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountDestinationCertificateDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets the destination certificate in a subaccount.
		
		__Tip:__
		You must be assigned the Destination Admin role in the subaccount.

		__Further Information:__
		<https://help.sap.com/docs/connectivity/sap-btp-connectivity-cf/use-destination-certificates>`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount which contains the certificate.",
				Required:            true,
			},
			"service_instance_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the service instance associated with this certificate.",
				Optional:            true,
			},
			"certificate_name": schema.StringAttribute{
				MarkdownDescription: "The name of the certificate with a valid certificate extension. Supported certificate types include .pem, .p12, .jks and .pfx",
				Required:            true,
			},
			"certificate_nodes": schema.ListNestedAttribute{
				MarkdownDescription: "List of certificate nodes containing details about the certificate and private key components.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							MarkdownDescription: "Denotes the type of the node i.e., 'private_key' or 'x509_certificate'.",
							Computed:            true,
						},
						"format": schema.StringAttribute{
							MarkdownDescription: "The format of the certificate or key.",
							Computed:            true,
						},
						"algorithm": schema.StringAttribute{
							MarkdownDescription: "The cryptographic algorithm used in the key.",
							Computed:            true,
						},
						"alias": schema.StringAttribute{
							MarkdownDescription: "An identifier used for the key.",
							Computed:            true,
						},
						"subject": schema.StringAttribute{
							MarkdownDescription: "The certificate subject which identifies the owner of the certificate.",
							Computed:            true,
						},
						"issuer": schema.StringAttribute{
							MarkdownDescription: "The certificate issuer which identifies the certificate authority that signed the certificate.",
							Computed:            true,
						},
						"common_name": schema.StringAttribute{
							MarkdownDescription: "The common name (CN) extracted from the certificate subject. May be null if not specified in the certificate.",
							Computed:            true,
						},
						"not_before": schema.StringAttribute{
							MarkdownDescription: "The start date and time (in ISO 8601 format) from which the certificate is valid.",
							Computed:            true,
						},
						"not_after": schema.StringAttribute{
							MarkdownDescription: "The expiration date and time (in ISO 8601 format) after which the certificate is no longer valid.",
							Computed:            true,
						},
						"certificate": schema.StringAttribute{
							MarkdownDescription: "The complete X.509 certificate in PEM format, encoded as base64.",
							Computed:            true,
							Sensitive:           true,
						},
					},
				},
			},
			"certification_creation_details": schema.SingleNestedAttribute{
				MarkdownDescription: "Details about how the destination certificate was created and its configuration settings.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"generation_method": schema.StringAttribute{
						MarkdownDescription: "Specifies the method used to create the certificate.",
						Computed:            true,
					},
					"common_name": schema.StringAttribute{
						MarkdownDescription: "The common name (CN) associated with the certificate.",
						Computed:            true,
					},
					"has_password": schema.BoolAttribute{
						MarkdownDescription: "Indicates whether the certificate is protected with a password.",
						Computed:            true,
					},
					"auto_renew": schema.BoolAttribute{
						MarkdownDescription: "Specifies whether the certificate is automatically renewed before it expires.",
						Computed:            true,
					},
					"validity_duration": schema.StringAttribute{
						MarkdownDescription: "The numeric duration for which the certificate is valid.",
						Computed:            true,
					},
					"validity_time_units": schema.StringAttribute{
						MarkdownDescription: "The time unit associated with the validity duration, such as `DAYS`, `MONTHS`, or `YEARS`.",
						Computed:            true,
					},
				},
			},
		},
	}
}

func (ds *subaccountDestinationCertificateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data subaccountDestinationCertificateDataSourceType

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := ds.cli.Connectivity.DestinationCertificate.Get(ctx, &btpcli.DestinationCertificateGetInput{
		SubaccountId:      data.SubaccountId.ValueString(),
		ServiceInstanceId: data.ServiceInstanceId.ValueString(),
		CertificateName:   data.CertificateName.ValueString(),
	})

	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Destination Certificate", fmt.Sprintf("%s", err))
		return
	}

	dataRes, diags := subaccountDestinationCertificateValueFrom(ctx, cliRes)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := subaccountDestinationCertificateDataSourceType{
		SubaccountId:      data.SubaccountId,
		ServiceInstanceId: data.ServiceInstanceId,
		CertificateName:   dataRes.CertificateName,
		Nodes:             dataRes.Nodes,
		Creation:          dataRes.Creation,
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
