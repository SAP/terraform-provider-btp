package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
)

func newSubaccountDestinationTrustDataSource() datasource.DataSource {
	return &subaccountDestinationTrustDataSource{}
}

type subaccountDestinationTrustType struct {
	/* INPUT */
	SubaccountID types.String `tfsdk:"subaccount_id"`
	Active       types.Bool   `tfsdk:"active"`
	/* OUTPUT */
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	BaseURL             types.String `tfsdk:"base_url"`
	Expiration          types.String `tfsdk:"expiration"`
	Owner               types.Object `tfsdk:"owner"`
	GeneratedOn         types.String `tfsdk:"generated_on"`
	PublicKeyBase64     types.String `tfsdk:"public_key_base64"`
	X509PublicKeyBase64 types.String `tfsdk:"x509_public_key_base64"`
}

func destinationTrustDatasourceOwnerObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"instance_id":   types.StringType,
		"subaccount_id": types.StringType,
	}
}

type subaccountDestinationTrustDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountDestinationTrustDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_destination_trust", req.ProviderTypeName)
}

func (ds *subaccountDestinationTrustDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountDestinationTrustDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets details about a specific subaccount destination trust.

__Tip:__
You must be assigned to the admin or viewer role of the subaccount.`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
			},
			"id": schema.StringAttribute{ // required by hashicorps terraform plugin testing framework
				DeprecationMessage:  "Use the `subaccount_id` attribute instead",
				MarkdownDescription: "The ID of the subaccount.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the destination trust.",
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"base_url": schema.StringAttribute{
				MarkdownDescription: "The base URL of the destination trust.",
				Computed:            true,
			},
			"active": schema.BoolAttribute{
				MarkdownDescription: "Shows whether the destination trust is active or passive.",
				Optional:            true,
				Computed:            true,
			},
			"expiration": schema.StringAttribute{
				MarkdownDescription: "The expiration timestamp of the destination trust.",
				Computed:            true,
			},
			"owner": schema.SingleNestedAttribute{
				MarkdownDescription: "The owner of the destination trust.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"instance_id": schema.StringAttribute{
						MarkdownDescription: "The instance ID of the owner.",
						Computed:            true,
					},
					"subaccount_id": schema.StringAttribute{
						MarkdownDescription: "The subaccount ID of the owner.",
						Computed:            true,
					},
				},
			},
			"generated_on": schema.StringAttribute{
				MarkdownDescription: "The generation timestamp of the destination trust.",
				Computed:            true,
			},
			"public_key_base64": schema.StringAttribute{
				MarkdownDescription: "The public key in base64 format of the destination trust.",
				Computed:            true,
			},
			"x509_public_key_base64": schema.StringAttribute{
				MarkdownDescription: "The x509 public key in base64 format of the destination trust.",
				Computed:            true,
			},
		},
	}
}

func (ds *subaccountDestinationTrustDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountDestinationTrustType
	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	active := true // if not provided, default to active destination trust certificate
	if !data.Active.IsNull() && !data.Active.IsUnknown() {
		active = data.Active.ValueBool()
	}

	destinationTrustDetails, _, err := ds.cli.Connectivity.DestinationTrust.GetBySubaccount(ctx, data.SubaccountID.ValueString(), active)
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Destination Trust", fmt.Sprintf("%s", err))
		return
	}

	data.ID = data.SubaccountID
	data.Name = types.StringValue(destinationTrustDetails.Name)
	data.BaseURL = types.StringValue(destinationTrustDetails.BaseURL)
	data.GeneratedOn = types.StringValue(destinationTrustDetails.GeneratedOn)
	data.PublicKeyBase64 = types.StringValue(destinationTrustDetails.PublicKeyBase64)
	data.X509PublicKeyBase64 = types.StringValue(destinationTrustDetails.X509PublicKeyBase64)
	data.Active = types.BoolValue(active)

	// Owner
	expTime := time.UnixMilli(destinationTrustDetails.Expiration).UTC().Format(time.RFC3339Nano)
	data.Expiration = types.StringValue(expTime)

	if destinationTrustDetails.Owner != nil {
		ownerAttrs := map[string]attr.Value{
			"instance_id":   types.StringValue(destinationTrustDetails.Owner.InstanceID),
			"subaccount_id": types.StringValue(destinationTrustDetails.Owner.SubaccountID),
		}
		data.Owner, _ = types.ObjectValue(destinationTrustDatasourceOwnerObjectType(), ownerAttrs)
	} else {
		data.Owner, _ = types.ObjectValue(destinationTrustDatasourceOwnerObjectType(), map[string]attr.Value{
			"instance_id":   types.StringNull(),
			"subaccount_id": types.StringNull(),
		})
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
