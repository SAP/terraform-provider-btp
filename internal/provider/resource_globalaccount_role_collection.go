package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
)

func newGlobalaccountRoleCollectionResource() resource.Resource {
	return &globalaccountRoleCollectionResource{}
}

type globalaccountRoleCollectionType struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

type globalaccountRoleCollectionResource struct {
	cli *btpcli.ClientFacade
}

func (rs *globalaccountRoleCollectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_globalaccount_role_collection", req.ProviderTypeName)
}

func (rs *globalaccountRoleCollectionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	rs.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (rs *globalaccountRoleCollectionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Create a role collection in a global account.

__Further documentation__
https://help.sap.com/viewer/65de2977205c403bbc107264b8eccf4b/Cloud/en-US/0039cf082d3d43eba9200fe15647922a.html`,
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the role collection.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Whether the role collection is readonly.",
				Optional:            true,
				Computed:            true,
			},
			/*"role_references": schema.ListNestedAttribute{
			    NestedObject: schema.NestedAttributeObject{
			        Attributes: map[string]schema.Attribute{
			            "role_template_name": schema.StringAttribute{
			                MarkdownDescription: "The name of the referenced role template.",
			                Computed:    true,
			            },
			            "role_template_app_id": schema.StringAttribute{
			                MarkdownDescription: "The name of the referenced template app id",
			                Computed:    true,
			            },
			            "description": schema.StringAttribute{
			                MarkdownDescription: "The description of the referenced role",
			                Computed:    true,
			            },
			            "name": schema.StringAttribute{
			                MarkdownDescription: "The name of the referenced role.",
			                Computed:    true,
			            },
			        },
			    },
			    Computed: true,
			},*/
		},
	}
}

func (rs *globalaccountRoleCollectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state globalaccountRoleCollectionType

	diags := req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := rs.cli.Security.RoleCollection.GetByGlobalAccount(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Role Collection (Global Account)", fmt.Sprintf("%s", err))
		return
	}

	state.Name = types.StringValue(cliRes.Name)
	state.Description = types.StringValue(cliRes.Description)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (rs *globalaccountRoleCollectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan globalaccountRoleCollectionType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := rs.cli.Security.RoleCollection.CreateByGlobalAccount(ctx, plan.Name.ValueString(), plan.Description.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Resource Role Collection (Global Account)", fmt.Sprintf("%s", err))
		return
	}

	plan.Name = types.StringValue(cliRes.Name)
	plan.Description = types.StringValue(cliRes.Description)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (rs *globalaccountRoleCollectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan globalaccountRoleCollectionType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.AddError("API Error Updating Resource Role Collection (Global Account)", "Update is not yet implemented.")

	/*TODO cliRes, err := rs.cli.Execute(ctx, btpcli.Update, rs.command, plan)
	if err != nil {
		resp.Diagnostics.AddError("API Error Updating Resource Role Collection (Global Account)", fmt.Sprintf("%s", err))
		return
	}*/

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (rs *globalaccountRoleCollectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state globalaccountRoleCollectionType
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, _, err := rs.cli.Security.RoleCollection.DeleteByGlobalAccount(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Role Collection (Global Account)", fmt.Sprintf("%s", err))
		return
	}
}
