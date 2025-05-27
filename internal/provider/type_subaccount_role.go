package provider

import (
	"context"
	"encoding/json"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/xsuaa_authz"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type subaccountRoleAttribute struct {
	AttributeName        types.String `tfsdk:"attribute_name" json:"attributeName,omitempty"`
	AttributeValueOrigin types.String `tfsdk:"attribute_value_origin" json:"attributeValueOrigin,omitempty"`
	AttributeValues      types.Set    `tfsdk:"attribute_values" json:"attributeValues,omitempty"`
	Description          types.String `tfsdk:"description" json:"description,omitempty"`      // The description of the role attribute.
	ValueRequired        types.Bool   `tfsdk:"value_required" json:"valueRequired,omitempty"` // Indicates whether the attribute value is required.
}

type subaccountRoleType struct {
	SubaccountId      types.String              `tfsdk:"subaccount_id"`
	Id                types.String              `tfsdk:"id"`
	Name              types.String              `tfsdk:"name"`
	RoleTemplateAppId types.String              `tfsdk:"app_id"`
	RoleTemplateName  types.String              `tfsdk:"role_template_name"`
	Description       types.String              `tfsdk:"description"`
	AttributeList     []subaccountRoleAttribute `tfsdk:"attribute_list"`
	IsReadOnly        types.Bool                `tfsdk:"read_only"`
}

func subaccountRoleFromValue(ctx context.Context, value xsuaa_authz.Role) (subaccountRoleType, diag.Diagnostics) {
	var subaccountRole subaccountRoleType

	subaccountRole.Description = types.StringValue(value.Description)
	subaccountRole.IsReadOnly = types.BoolValue(value.IsReadOnly)
	subaccountRole.Name = types.StringValue(value.Name)
	subaccountRole.RoleTemplateName = types.StringValue(value.RoleTemplateName)
	subaccountRole.RoleTemplateAppId = types.StringValue(value.RoleTemplateAppId)

	for _, attribute := range value.AttributeList {
		attributeLine := subaccountRoleAttribute{
			AttributeName:        types.StringValue(attribute.AttributeName),
			AttributeValueOrigin: types.StringValue(attribute.AttributeValueOrigin),
			Description:          types.StringValue(attribute.Description),
			ValueRequired:        types.BoolValue(attribute.ValueRequired),
		}

		attributeLine.AttributeValues, _ = types.SetValueFrom(ctx, types.StringType, attribute.AttributeValues)

		subaccountRole.AttributeList = append(subaccountRole.AttributeList, attributeLine)
	}
	return subaccountRole, diag.Diagnostics{}
}

func subaccountAttributeListToJsonString(attributeList []subaccountRoleAttribute) (string, error) {
	json, err := json.Marshal(attributeList)

	if err != nil {
		return "", err
	}
	return string(json), nil
}
