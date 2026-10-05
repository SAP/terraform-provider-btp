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
	"github.com/SAP/terraform-provider-btp/internal/validation/uuidvalidator"
)

func newDirectoryRoleCollectionAssignmentResource() resource.Resource {
	return &directoryRoleCollectionAssignmentResource{}
}

type directoryRoleCollectionAssignmentType struct {
	Id                 types.String `tfsdk:"id"`
	DirectoryId        types.String `tfsdk:"directory_id"`
	RoleCollectionName types.String `tfsdk:"role_collection_name"`
	Username           types.String `tfsdk:"user_name"`
	Groupname          types.String `tfsdk:"group_name"`
	AttributeName      types.String `tfsdk:"attribute_name"`
	AttributeValue     types.String `tfsdk:"attribute_value"`
	Origin             types.String `tfsdk:"origin"`
}

type directoryRoleCollectionAssignmentResource struct {
	cli *btpcli.ClientFacade
}

func (rs *directoryRoleCollectionAssignmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_directory_role_collection_assignment", req.ProviderTypeName)
}

func (rs *directoryRoleCollectionAssignmentResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	rs.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (rs *directoryRoleCollectionAssignmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Assigns a user to a role collection on a directory level.

__Tip:__
You must be assigned to the admin role of the global account or the directory.`,
		Attributes: map[string]schema.Attribute{
			"directory_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the directory.",
				Required:            true,
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
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
				DeprecationMessage:  "Use the `directory_id` and `role_collection_name` attributes instead",
				MarkdownDescription: "The combined unique ID of the role collection.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_name": schema.StringAttribute{
				MarkdownDescription: "The username of the user to assign.",
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
					stringvalidator.LengthAtLeast(1),
					stringvalidator.AlsoRequires(path.MatchRoot("origin")),
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
				MarkdownDescription: "The identity provider that hosts the user or a group. Only needed for custom identity provider.",
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

type directoryRoleCollectionAssignmentIdentityModel struct {
	DirectoryId        types.String `tfsdk:"directory_id"`
	RoleCollectionName types.String `tfsdk:"role_collection_name"`
	Username           types.String `tfsdk:"user_name"`
	Groupname          types.String `tfsdk:"group_name"`
	AttributeName      types.String `tfsdk:"attribute_name"`
	AttributeValue     types.String `tfsdk:"attribute_value"`
	Origin             types.String `tfsdk:"origin"`
}

func (rs *directoryRoleCollectionAssignmentResource) IdentitySchema(_ context.Context, _ resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"directory_id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
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

func (rs *directoryRoleCollectionAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state directoryRoleCollectionAssignmentType

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !state.Username.IsNull() {
		users, _, err := rs.cli.Security.RoleCollection.GetUserAssignmentsByDirectory(ctx, state.DirectoryId.ValueString(), state.RoleCollectionName.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("API Error Reading Resource Role Collection Assignment (Directory)", fmt.Sprintf("%s", err))
			return
		}
		for _, u := range users {
			if (u.Username == state.Username.ValueString() || u.Email == state.Username.ValueString()) && originMatches(u.Origin, state.Origin.ValueString()) {
				if state.Id.IsNull() || state.Id.IsUnknown() {
					state.Id = types.StringValue(fmt.Sprintf("%s,%s,%s", state.DirectoryId.ValueString(), state.RoleCollectionName.ValueString(), state.Username.ValueString()))
				}
				diags = resp.State.Set(ctx, &state)
				resp.Diagnostics.Append(diags...)
				diags = resp.Identity.Set(ctx, directoryRoleCollectionAssignmentIdentityModel{
					DirectoryId:        state.DirectoryId,
					RoleCollectionName: state.RoleCollectionName,
					Username:           state.Username,
					Origin:             state.Origin,
				})
				resp.Diagnostics.Append(diags...)
				return
			}
		}
		resp.State.RemoveResource(ctx)
		return
	}

	cliRes, _, err := rs.cli.Security.RoleCollection.GetByDirectoryWithAttributeMappings(ctx, state.DirectoryId.ValueString(), state.RoleCollectionName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Role Collection Assignment (Directory)", fmt.Sprintf("%s", err))
		return
	}

	for _, am := range cliRes.SamlAttributeAssignment {
		if !state.Groupname.IsNull() {
			if am.AttributeName == "Groups" && am.AttributeValue == state.Groupname.ValueString() && samlOriginMatches(am.IdentityProvider, am.SamlEntityId, state.Origin.ValueString()) {
				if state.Id.IsNull() || state.Id.IsUnknown() {
					state.Id = types.StringValue(fmt.Sprintf("%s,%s,group:%s", state.DirectoryId.ValueString(), state.RoleCollectionName.ValueString(), state.Groupname.ValueString()))
				}
				diags = resp.State.Set(ctx, &state)
				resp.Diagnostics.Append(diags...)
				diags = resp.Identity.Set(ctx, directoryRoleCollectionAssignmentIdentityModel{
					DirectoryId:        state.DirectoryId,
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
					state.Id = types.StringValue(fmt.Sprintf("%s,%s,attribute:%s/%s", state.DirectoryId.ValueString(), state.RoleCollectionName.ValueString(), state.AttributeName.ValueString(), state.AttributeValue.ValueString()))
				}
				diags = resp.State.Set(ctx, &state)
				resp.Diagnostics.Append(diags...)
				diags = resp.Identity.Set(ctx, directoryRoleCollectionAssignmentIdentityModel{
					DirectoryId:        state.DirectoryId,
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
	resp.State.RemoveResource(ctx)
}

func (rs *directoryRoleCollectionAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan directoryRoleCollectionAssignmentType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var err error
	if !plan.Username.IsNull() {
		// assign user
		_, _, err = rs.cli.Security.RoleCollection.AssignUserByDirectory(ctx, plan.DirectoryId.ValueString(), plan.RoleCollectionName.ValueString(), plan.Username.ValueString(), plan.Origin.ValueString())
	} else if !plan.Groupname.IsNull() {
		// assign group
		_, _, err = rs.cli.Security.RoleCollection.AssignGroupByDirectory(ctx, plan.DirectoryId.ValueString(), plan.RoleCollectionName.ValueString(), plan.Groupname.ValueString(), plan.Origin.ValueString())
	} else {
		// assign attribute
		_, _, err = rs.cli.Security.RoleCollection.AssignAttributeByDirectory(ctx, plan.DirectoryId.ValueString(), plan.RoleCollectionName.ValueString(), plan.AttributeName.ValueString(), plan.AttributeValue.ValueString(), plan.Origin.ValueString())
	}

	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Resource Role Collection Assignment (Directory)", fmt.Sprintf("%s", err))
		return
	}

	// Setting ID of state - required by hashicorps terraform plugin testing framework for Create. See issue https://github.com/hashicorp/terraform-plugin-testing/issues/84
	plan.Id = types.StringValue(fmt.Sprintf("%s,%s,%s", plan.DirectoryId.ValueString(), plan.RoleCollectionName.ValueString(), plan.Username.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	diags = resp.Identity.Set(ctx, directoryRoleCollectionAssignmentIdentityModel{
		DirectoryId:        plan.DirectoryId,
		RoleCollectionName: plan.RoleCollectionName,
		Username:           plan.Username,
		Groupname:          plan.Groupname,
		AttributeName:      plan.AttributeName,
		AttributeValue:     plan.AttributeValue,
		Origin:             plan.Origin,
	})
	resp.Diagnostics.Append(diags...)
}

func (rs *directoryRoleCollectionAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan directoryRoleCollectionAssignmentType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// since all the attributes are marked to be replaced in case of update, this should never be reached.
	resp.Diagnostics.AddError("API Error Updating Resource Role Collection Assignment (Directory)", "This resource is not supposed to be updated")

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (rs *directoryRoleCollectionAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state directoryRoleCollectionAssignmentType
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var err error
	if !state.Username.IsNull() {
		// unassign user
		_, _, err = rs.cli.Security.RoleCollection.UnassignUserByDirectory(ctx, state.DirectoryId.ValueString(), state.RoleCollectionName.ValueString(), state.Username.ValueString(), state.Origin.ValueString())
	} else if !state.Groupname.IsNull() {
		// unassign group
		_, _, err = rs.cli.Security.RoleCollection.UnassignGroupByDirectory(ctx, state.DirectoryId.ValueString(), state.RoleCollectionName.ValueString(), state.Groupname.ValueString(), state.Origin.ValueString())
	} else {
		// unassign attribute
		_, _, err = rs.cli.Security.RoleCollection.UnassignAttributeByDirectory(ctx, state.DirectoryId.ValueString(), state.RoleCollectionName.ValueString(), state.AttributeName.ValueString(), state.AttributeValue.ValueString(), state.Origin.ValueString())
	}

	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Role Collection Assignment (Directory)", fmt.Sprintf("%s", err))
		return
	}
}

func (rs *directoryRoleCollectionAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID != "" {
		idParts := strings.Split(req.ID, ",")

		if len(idParts) != 4 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" || idParts[3] == "" {
			resp.Diagnostics.AddError(
				"Unexpected Import Identifier",
				fmt.Sprintf("Expected import identifier with format: directory_id,role_collection_name,user_name,origin (for user assignments) or directory_id,role_collection_name,group:<group_name>,origin (for group assignments) or directory_id,role_collection_name,attribute:<attr_name>/<attr_value>,origin (for attribute assignments). Got: %q", req.ID),
			)
			return
		}

		directoryID := idParts[0]
		roleCollectionName := idParts[1]
		assignmentPart := idParts[2]
		origin := idParts[3]

		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("directory_id"), directoryID)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("role_collection_name"), roleCollectionName)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("origin"), origin)...)

		if strings.HasPrefix(assignmentPart, "group:") {
			groupName := strings.TrimPrefix(assignmentPart, "group:")
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("group_name"), groupName)...)
		} else if strings.HasPrefix(assignmentPart, "attribute:") {
			attrPart := strings.TrimPrefix(assignmentPart, "attribute:")
			slashIdx := strings.Index(attrPart, "/")
			if slashIdx == -1 {
				resp.Diagnostics.AddError(
					"Unexpected Import Identifier",
					fmt.Sprintf("Expected attribute assignment in format attribute:<attr_name>/<attr_value>. Got: %q", assignmentPart),
				)
				return
			}
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("attribute_name"), attrPart[:slashIdx])...)
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("attribute_value"), attrPart[slashIdx+1:])...)
		} else {
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("user_name"), assignmentPart)...)
		}
		return
	}

	var identityData directoryRoleCollectionAssignmentIdentityModel
	resp.Diagnostics.Append(req.Identity.Get(ctx, &identityData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("directory_id"), identityData.DirectoryId)...)
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
