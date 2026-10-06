package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
)

func newGlobalaccountRoleCollectionAssignmentResource() resource.Resource {
	return &globalaccountRoleCollectionAssignmentResource{}
}

type globalaccountRoleCollectionAssignmentType struct {
	Id                 types.String `tfsdk:"id"`
	RoleCollectionName types.String `tfsdk:"role_collection_name"`
	Username           types.String `tfsdk:"user_name"`
	Groupname          types.String `tfsdk:"group_name"`
	AttributeName      types.String `tfsdk:"attribute_name"`
	AttributeValue     types.String `tfsdk:"attribute_value"`
	Origin             types.String `tfsdk:"origin"`
}

type globalaccountRoleCollectionAssignmentResource struct {
	cli *btpcli.ClientFacade
}

func (rs *globalaccountRoleCollectionAssignmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_globalaccount_role_collection_assignment", req.ProviderTypeName)
}

func (rs *globalaccountRoleCollectionAssignmentResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	rs.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (rs *globalaccountRoleCollectionAssignmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Assigns a user or a group to a role collection on global account level.

__Tip:__
You must be assigned to the admin role of the global account.`,
		Attributes: map[string]schema.Attribute{
			"role_collection_name": schema.StringAttribute{
				MarkdownDescription: "The name of the role collection.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"id": schema.StringAttribute{ // required by hashicorps terraform plugin testing framework
				DeprecationMessage:  "Use the `role_collection_name` attribute instead",
				MarkdownDescription: "The combined unique ID of the role collection.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_name": schema.StringAttribute{
				MarkdownDescription: "The name of the user to assign.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("user_name"), path.MatchRoot("group_name"), path.MatchRoot("attribute_name")),
					stringvalidator.LengthBetween(1, 256),
				},
			},
			"group_name": schema.StringAttribute{
				MarkdownDescription: "The name of the group to assign.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.MatchRoot("origin")),
					stringvalidator.LengthAtLeast(1),
				},
			},
			"attribute_name": schema.StringAttribute{
				MarkdownDescription: "The name of the attribute to assign.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.MatchRoot("attribute_value")),
					stringvalidator.AlsoRequires(path.MatchRoot("origin")),
					stringvalidator.LengthAtLeast(1),
				},
			},
			"attribute_value": schema.StringAttribute{
				MarkdownDescription: "The value of the attribute to assign.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.MatchRoot("attribute_name")),
					stringvalidator.LengthAtLeast(1),
				},
			},
			"origin": schema.StringAttribute{
				MarkdownDescription: "The identity provider that hosts the user or group. Only needed for custom identity provider.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("ldap"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

type globalaccountRoleCollectionAssignmentIdentityModel struct {
	RoleCollectionName types.String `tfsdk:"role_collection_name"`
	Username           types.String `tfsdk:"user_name"`
	Groupname          types.String `tfsdk:"group_name"`
	AttributeName      types.String `tfsdk:"attribute_name"`
	AttributeValue     types.String `tfsdk:"attribute_value"`
	Origin             types.String `tfsdk:"origin"`
}

func (rs *globalaccountRoleCollectionAssignmentResource) IdentitySchema(_ context.Context, _ resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"role_collection_name": identityschema.StringAttribute{
				RequiredForImport: true,
			},
			"user_name": identityschema.StringAttribute{
				OptionalForImport: true,
			},
			"group_name": identityschema.StringAttribute{
				OptionalForImport: true,
			},
			"attribute_name": identityschema.StringAttribute{
				OptionalForImport: true,
			},
			"attribute_value": identityschema.StringAttribute{
				OptionalForImport: true,
			},
			"origin": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func (rs *globalaccountRoleCollectionAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state globalaccountRoleCollectionAssignmentType

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !state.Username.IsNull() {
		users, _, err := rs.cli.Security.RoleCollection.GetUserAssignmentsByGlobalAccount(ctx, state.RoleCollectionName.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("API Error Reading Resource Role Collection Assignment (Global Account)", fmt.Sprintf("%s", err))
			return
		}
		for _, u := range users {
			if (u.Username == state.Username.ValueString() || u.Email == state.Username.ValueString()) && originMatches(u.Origin, state.Origin.ValueString()) {
				if state.Id.IsNull() || state.Id.IsUnknown() {
					state.Id = types.StringValue(fmt.Sprintf("%s,%s,%s", state.RoleCollectionName.ValueString(), state.Username.ValueString(), state.Origin.ValueString()))
				}
				diags = resp.State.Set(ctx, &state)
				resp.Diagnostics.Append(diags...)
				diags = resp.Identity.Set(ctx, globalaccountRoleCollectionAssignmentIdentityModel{
					RoleCollectionName: state.RoleCollectionName,
					Username:           state.Username,
					Origin:             state.Origin,
				})
				resp.Diagnostics.Append(diags...)
				return
			}
		}
		resp.Diagnostics.Append(resp.Identity.Set(ctx, globalaccountRoleCollectionAssignmentIdentityModel{
			RoleCollectionName: state.RoleCollectionName,
			Username:           state.Username,
			Origin:             state.Origin,
		})...)
		resp.State.RemoveResource(ctx)
		return
	}

	cliRes, _, err := rs.cli.Security.RoleCollection.GetByGlobalAccountWithAttributeMapings(ctx, state.RoleCollectionName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Role Collection Assignment (Global Account)", fmt.Sprintf("%s", err))
		return
	}

	for _, am := range cliRes.SamlAttributeAssignment {
		if !state.Groupname.IsNull() {
			if am.AttributeName == "Groups" && am.AttributeValue == state.Groupname.ValueString() && samlOriginMatches(am.IdentityProvider, am.SamlEntityId, state.Origin.ValueString()) {
				if state.Id.IsNull() || state.Id.IsUnknown() {
					state.Id = types.StringValue(fmt.Sprintf("%s,group:%s,%s", state.RoleCollectionName.ValueString(), state.Groupname.ValueString(), state.Origin.ValueString()))
				}
				diags = resp.State.Set(ctx, &state)
				resp.Diagnostics.Append(diags...)
				diags = resp.Identity.Set(ctx, globalaccountRoleCollectionAssignmentIdentityModel{
					RoleCollectionName: state.RoleCollectionName,
					Groupname:          state.Groupname,
					Origin:             state.Origin,
				})
				resp.Diagnostics.Append(diags...)
				return
			}
		} else {
			if am.AttributeName == state.AttributeName.ValueString() && am.AttributeValue == state.AttributeValue.ValueString() && samlOriginMatches(am.IdentityProvider, am.SamlEntityId, state.Origin.ValueString()) {
				if state.Id.IsNull() || state.Id.IsUnknown() {
					state.Id = types.StringValue(fmt.Sprintf("%s,attribute:%s/%s,%s", state.RoleCollectionName.ValueString(), state.AttributeName.ValueString(), state.AttributeValue.ValueString(), state.Origin.ValueString()))
				}
				diags = resp.State.Set(ctx, &state)
				resp.Diagnostics.Append(diags...)
				diags = resp.Identity.Set(ctx, globalaccountRoleCollectionAssignmentIdentityModel{
					RoleCollectionName: state.RoleCollectionName,
					AttributeName:      state.AttributeName,
					AttributeValue:     state.AttributeValue,
					Origin:             state.Origin,
				})
				resp.Diagnostics.Append(diags...)
				return
			}
		}
	}
	if !state.Groupname.IsNull() {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, globalaccountRoleCollectionAssignmentIdentityModel{
			RoleCollectionName: state.RoleCollectionName,
			Groupname:          state.Groupname,
			Origin:             state.Origin,
		})...)
	} else {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, globalaccountRoleCollectionAssignmentIdentityModel{
			RoleCollectionName: state.RoleCollectionName,
			AttributeName:      state.AttributeName,
			AttributeValue:     state.AttributeValue,
			Origin:             state.Origin,
		})...)
	}
	resp.State.RemoveResource(ctx)
}

func (rs *globalaccountRoleCollectionAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan globalaccountRoleCollectionAssignmentType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var err error
	if !plan.Username.IsNull() {
		// assign user
		_, _, err = rs.cli.Security.RoleCollection.AssignUserByGlobalaccount(ctx, plan.RoleCollectionName.ValueString(), plan.Username.ValueString(), plan.Origin.ValueString())
	} else if !plan.Groupname.IsNull() {
		// assign group
		_, _, err = rs.cli.Security.RoleCollection.AssignGroupByGlobalaccount(ctx, plan.RoleCollectionName.ValueString(), plan.Groupname.ValueString(), plan.Origin.ValueString())
	} else {
		// assign attribute
		_, _, err = rs.cli.Security.RoleCollection.AssignAttributeByGlobalaccount(ctx, plan.RoleCollectionName.ValueString(), plan.AttributeName.ValueString(), plan.AttributeValue.ValueString(), plan.Origin.ValueString())
	}

	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Resource Role Collection Assignment (Global Account)", fmt.Sprintf("%s", err))
		return
	}

	// Setting ID of state - required by hashicorps terraform plugin testing framework for Create. See issue https://github.com/hashicorp/terraform-plugin-testing/issues/84
	plan.Id = types.StringValue(fmt.Sprintf("%s,%s", plan.RoleCollectionName.ValueString(), plan.Username.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	diags = resp.Identity.Set(ctx, globalaccountRoleCollectionAssignmentIdentityModel{
		RoleCollectionName: plan.RoleCollectionName,
		Username:           plan.Username,
		Groupname:          plan.Groupname,
		AttributeName:      plan.AttributeName,
		AttributeValue:     plan.AttributeValue,
		Origin:             plan.Origin,
	})
	resp.Diagnostics.Append(diags...)
}

func (rs *globalaccountRoleCollectionAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan globalaccountRoleCollectionAssignmentType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// since all the attributes are marked to be replaced in case of update, this should never be reached.
	resp.Diagnostics.AddError("API Error Updating Role Collection Assignment (Global Account)", "This resource is not supposed to be updated")

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (rs *globalaccountRoleCollectionAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state globalaccountRoleCollectionAssignmentType
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var err error
	if !state.Username.IsNull() {
		// unassign user
		_, _, err = rs.cli.Security.RoleCollection.UnassignUserByGlobalaccount(ctx, state.RoleCollectionName.ValueString(), state.Username.ValueString(), state.Origin.ValueString())
	} else if !state.Groupname.IsNull() {
		// unassign group
		_, _, err = rs.cli.Security.RoleCollection.UnassignGroupByGlobalaccount(ctx, state.RoleCollectionName.ValueString(), state.Groupname.ValueString(), state.Origin.ValueString())
	} else {
		// unassign attribute
		_, _, err = rs.cli.Security.RoleCollection.UnassignAttributeByGlobalaccount(ctx, state.RoleCollectionName.ValueString(), state.AttributeName.ValueString(), state.AttributeValue.ValueString(), state.Origin.ValueString())
	}

	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Role Collection Assignment (Global Account)", fmt.Sprintf("%s", err))
		return
	}
}

func (rs *globalaccountRoleCollectionAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID != "" {
		idParts := strings.Split(req.ID, ",")

		if len(idParts) != 3 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" {
			resp.Diagnostics.AddError(
				"Unexpected Import Identifier",
				fmt.Sprintf("Expected import identifier with format: role_collection_name,user_name,origin (for user assignments) or role_collection_name,group:<group_name>,origin (for group assignments) or role_collection_name,attribute:<attr_name>/<attr_value>,origin (for attribute assignments). Got: %q", req.ID),
			)
			return
		}

		roleCollectionName := idParts[0]
		assignmentPart := idParts[1]
		origin := idParts[2]

		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("role_collection_name"), roleCollectionName)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("origin"), origin)...)

		if after, ok := strings.CutPrefix(assignmentPart, "group:"); ok {
			groupName := after
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("group_name"), groupName)...)
		} else if after, ok := strings.CutPrefix(assignmentPart, "attribute:"); ok {
			attrPart := after
			before, after, ok := strings.Cut(attrPart, "/")
			if !ok {
				resp.Diagnostics.AddError(
					"Unexpected Import Identifier",
					fmt.Sprintf("Expected attribute assignment in format attribute:<attr_name>/<attr_value>. Got: %q", assignmentPart),
				)
				return
			}
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("attribute_name"), before)...)
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("attribute_value"), after)...)
		} else {
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("user_name"), assignmentPart)...)
		}
		return
	}

	var identityData globalaccountRoleCollectionAssignmentIdentityModel
	resp.Diagnostics.Append(req.Identity.Get(ctx, &identityData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("role_collection_name"), identityData.RoleCollectionName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("origin"), identityData.Origin)...)

	if !identityData.Username.IsNull() && identityData.Username.ValueString() != "" {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("user_name"), identityData.Username)...)
	} else if !identityData.Groupname.IsNull() && identityData.Groupname.ValueString() != "" {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("group_name"), identityData.Groupname)...)
	} else {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("attribute_name"), identityData.AttributeName)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("attribute_value"), identityData.AttributeValue)...)
	}
}
