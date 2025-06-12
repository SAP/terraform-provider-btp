package provider

import (
	"context"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/xsuaa_authz"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type globalaccountRoleType struct {
	Name              types.String `tfsdk:"name"`
	RoleTemplateAppId types.String `tfsdk:"app_id"`
	RoleTemplateName  types.String `tfsdk:"role_template_name"`
	Description       types.String `tfsdk:"description"`
	IsReadOnly        types.Bool   `tfsdk:"read_only"`
}

func globalaccountRoleFromValue(ctx context.Context, value xsuaa_authz.Role) (globalaccountRoleType, diag.Diagnostics) {
	var globalaccountRole globalaccountRoleType

	globalaccountRole.Description = types.StringValue(value.Description)
	globalaccountRole.IsReadOnly = types.BoolValue(value.IsReadOnly)
	globalaccountRole.Name = types.StringValue(value.Name)
	globalaccountRole.RoleTemplateName = types.StringValue(value.RoleTemplateName)
	globalaccountRole.RoleTemplateAppId = types.StringValue(value.RoleTemplateAppId)

	return globalaccountRole, diag.Diagnostics{}
}
