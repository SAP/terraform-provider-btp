package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
)

func newSubaccountRoleCollectionRoleDataSource() datasource.DataSource {
	return &subaccountRoleCollectionRoleDataSource{}
}

// This represents the configuration and state for a SINGLE role lookup
type subaccountRoleCollectionRoleDataSourceConfig struct {
	/* INPUTS */
	SubaccountId types.String `tfsdk:"subaccount_id"`
	Name         types.String `tfsdk:"name"`      // Collection Name
	RoleName     types.String `tfsdk:"role_name"` // The specific Role to find

	/* OUTPUTS */
	RoleTemplateName  types.String `tfsdk:"role_template_name"`
	RoleTemplateAppId types.String `tfsdk:"role_template_app_id"`
	Description       types.String `tfsdk:"description"`
}

type subaccountRoleCollectionRoleDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountRoleCollectionRoleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_role_collection_role", req.ProviderTypeName)
}

func (ds *subaccountRoleCollectionRoleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountRoleCollectionRoleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets details about a specific role assigned to a subaccount role collection.`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the role collection.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"role_name": schema.StringAttribute{
				MarkdownDescription: "The name of the role.",
				Required:            true,
			},
			"role_template_name": schema.StringAttribute{
				MarkdownDescription: "The name of the referenced role template.",
				Computed:            true,
			},
			"role_template_app_id": schema.StringAttribute{
				MarkdownDescription: "The name of the referenced template app id.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the referenced role.",
				Computed:            true,
			},
		},
	}
}

func (ds *subaccountRoleCollectionRoleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountRoleCollectionRoleDataSourceConfig
	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	roleCollection, _, err := ds.cli.Security.RoleCollection.GetBySubaccount(ctx, data.SubaccountId.ValueString(), data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Role Collection", fmt.Sprintf("%s", err))
		return
	}

	found := false
	for _, ref := range roleCollection.RoleReferences {
		if ref.Name == data.RoleName.ValueString() {
			data.RoleTemplateName = types.StringValue(ref.RoleTemplateName)
			data.RoleTemplateAppId = types.StringValue(ref.RoleTemplateAppId)
			data.Description = types.StringValue(ref.Description)
			found = true
			break
		}
	}

	if !found {
		resp.Diagnostics.AddError(
			"Role Not Found",
			fmt.Sprintf("Role '%s' not found in role collection '%s'", data.RoleName.ValueString(), data.Name.ValueString()),
		)
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
