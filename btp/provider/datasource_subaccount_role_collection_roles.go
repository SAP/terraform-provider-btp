package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
)

func newSubaccountRoleCollectionRolesDataSource() datasource.DataSource {
	return &subaccountRoleCollectionRolesDataSource{}
}

// Struct for the inner role objects
type subaccountRoleCollectionRolesValueModel struct {
	RoleName          types.String `tfsdk:"role_name"`
	RoleTemplateName  types.String `tfsdk:"role_template_name"`
	RoleTemplateAppId types.String `tfsdk:"role_template_app_id"`
	Description       types.String `tfsdk:"description"`
}

// Map for type conversion (Required for types.SetValueFrom)
var roleRolesAttrTypes = map[string]attr.Type{
	"role_name":            types.StringType,
	"role_template_name":   types.StringType,
	"role_template_app_id": types.StringType,
	"description":          types.StringType,
}

type subaccountRoleCollectionRolesDataSourceModel struct {
	SubaccountId types.String `tfsdk:"subaccount_id"`
	Name         types.String `tfsdk:"name"`
	Values       types.Set    `tfsdk:"values"` // Computed list of roles
}

type subaccountRoleCollectionRolesDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountRoleCollectionRolesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_role_collection_roles", req.ProviderTypeName)
}

func (ds *subaccountRoleCollectionRolesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountRoleCollectionRolesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Lists all roles assigned to a specific subaccount role collection.`,
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
			"values": schema.SetNestedAttribute{
				MarkdownDescription: "The roles assigned to the role collection.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"role_name": schema.StringAttribute{
							Computed: true,
						},
						"role_template_name": schema.StringAttribute{
							Computed: true,
						},
						"role_template_app_id": schema.StringAttribute{
							Computed: true,
						},
						"description": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (ds *subaccountRoleCollectionRolesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountRoleCollectionRolesDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	roleCollection, _, err := ds.cli.Security.RoleCollection.GetBySubaccount(ctx, data.SubaccountId.ValueString(), data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Role Collection Roles", fmt.Sprintf("%s", err))
		return
	}

	var roles []subaccountRoleCollectionRolesValueModel
	for _, ref := range roleCollection.RoleReferences {
		roles = append(roles, subaccountRoleCollectionRolesValueModel{
			RoleName:          types.StringValue(ref.Name),
			RoleTemplateName:  types.StringValue(ref.RoleTemplateName),
			RoleTemplateAppId: types.StringValue(ref.RoleTemplateAppId),
			Description:       types.StringValue(ref.Description),
		})
	}

	roleSet, diags := types.SetValueFrom(ctx, types.ObjectType{AttrTypes: roleRolesAttrTypes}, roles)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Values = roleSet

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
