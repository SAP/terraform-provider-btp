package provider

import (
	"context"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/xsuaa_authz"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type subaccountRoleType struct {
	SubaccountId      types.String `tfsdk:"subaccount_id"`
	Id                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	RoleTemplateAppId types.String `tfsdk:"app_id"`
	RoleTemplateName  types.String `tfsdk:"role_template_name"`
	Description       types.String `tfsdk:"description"`
	IsReadOnly        types.Bool   `tfsdk:"read_only"`
}

func subaccountRoleFromValue(ctx context.Context, value xsuaa_authz.Role) (subaccountRoleType, diag.Diagnostics) {
	var subaccountRole subaccountRoleType

	subaccountRole.Description = types.StringValue(value.Description)
	subaccountRole.IsReadOnly = types.BoolValue(value.IsReadOnly)
	subaccountRole.Name = types.StringValue(value.Name)
	subaccountRole.RoleTemplateName = types.StringValue(value.RoleTemplateName)
	subaccountRole.RoleTemplateAppId = types.StringValue(value.RoleTemplateAppId)

	return subaccountRole, diag.Diagnostics{}
}
