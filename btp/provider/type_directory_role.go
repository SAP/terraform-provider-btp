package provider

import (
	"context"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/xsuaa_authz"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type directoryRoleScope struct {
	Name                         types.String `tfsdk:"name"`
	Description                  types.String `tfsdk:"description"`
	CustomGrantAsAuthorityToApps types.Set    `tfsdk:"custom_grant_as_authority_to_apps"`
	CustomGrantedApps            types.Set    `tfsdk:"custom_granted_apps"`
	GrantAsAuthorityToApps       types.Set    `tfsdk:"grant_as_authority_to_apps"`
	GrantedApps                  types.Set    `tfsdk:"granted_apps"`
}

type directoryRoleType struct {
	/* INPUT */
	DirectoryId       types.String `tfsdk:"directory_id"`
	Id                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	RoleTemplateAppId types.String `tfsdk:"app_id"`
	RoleTemplateName  types.String `tfsdk:"role_template_name"`
	/* OUTPUT */
	Description types.String `tfsdk:"description"`
	IsReadOnly  types.Bool   `tfsdk:"read_only"`
	Scopes      types.List   `tfsdk:"scopes"`
}

func directoryRoleFromValue(ctx context.Context, value xsuaa_authz.Role) (directoryRoleType, diag.Diagnostics) {
	var dirRole directoryRoleType

	dirRole.Description = types.StringValue(value.Description)
	dirRole.IsReadOnly = types.BoolValue(value.IsReadOnly)
	dirRole.Name = types.StringValue(value.Name)
	dirRole.RoleTemplateName = types.StringValue(value.RoleTemplateName)
	dirRole.RoleTemplateAppId = types.StringValue(value.RoleTemplateAppId)

	// dirRole.Scopes = []directoryRoleScope{}

	dirRoleScopes := []directoryRoleScope{}

	var summary, diags diag.Diagnostics

	for _, scope := range value.Scopes {
		scopeVal := directoryRoleScope{
			Name:        types.StringValue(scope.Name),
			Description: types.StringValue(scope.Description),
		}

		scopeVal.CustomGrantAsAuthorityToApps, diags = types.SetValueFrom(ctx, types.StringType, scope.CustomGrantAsAuthorityToApps)
		summary.Append(diags...)

		scopeVal.CustomGrantedApps, diags = types.SetValueFrom(ctx, types.StringType, scope.CustomGrantedApps)
		summary.Append(diags...)

		scopeVal.GrantAsAuthorityToApps, diags = types.SetValueFrom(ctx, types.StringType, scope.GrantAsAuthorityToApps)
		summary.Append(diags...)

		scopeVal.GrantedApps, diags = types.SetValueFrom(ctx, types.StringType, scope.GrantedApps)
		summary.Append(diags...)

		dirRoleScopes = append(dirRoleScopes, scopeVal)
	}

	dirRole.Scopes, diags = types.ListValueFrom(ctx, directoryScopeObjType, dirRoleScopes)
	summary.Append(diags...)

	return dirRole, summary
}
