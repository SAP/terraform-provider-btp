package provider

import (
	"context"
	"encoding/json"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/xsuaa_authz"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type subaccountRoleAttributePlain struct {
	AttributeName        string   `json:"attributeName"`
	AttributeValueOrigin string   `json:"attributeValueOrigin"`
	AttributeValues      []string `json:"attributeValues"`
	ValueRequired        bool     `json:"valueRequired"`
}

type subaccountRoleAttribute struct {
	AttributeName        types.String `tfsdk:"attribute_name"`
	AttributeValueOrigin types.String `tfsdk:"attribute_value_origin"`
	AttributeValues      types.Set    `tfsdk:"attribute_values"`
	ValueRequired        types.Bool   `tfsdk:"value_required"`
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
			ValueRequired:        types.BoolValue(attribute.ValueRequired),
		}

		attributeLine.AttributeValues, _ = types.SetValueFrom(ctx, types.StringType, attribute.AttributeValues)

		subaccountRole.AttributeList = append(subaccountRole.AttributeList, attributeLine)
	}
	return subaccountRole, diag.Diagnostics{}
}

func subaccountAttributeListToJsonString(attributeList []subaccountRoleAttribute) (string, error) {

	var attributeListPlain []subaccountRoleAttributePlain
	var attributeValuePlain []string

	for _, attribute := range attributeList {
		attributeLinePlain := subaccountRoleAttributePlain{
			AttributeName:        attribute.AttributeName.ValueString(),
			AttributeValueOrigin: attribute.AttributeValueOrigin.ValueString(),
			ValueRequired:        attribute.ValueRequired.ValueBool(),
		}

		for _, value := range attribute.AttributeValues.Elements() {
			attributeValuePlain = append(attributeValuePlain, value.(types.String).ValueString())
		}
		attributeLinePlain.AttributeValues = attributeValuePlain

		attributeListPlain = append(attributeListPlain, attributeLinePlain)
	}

	json, err := json.Marshal(attributeListPlain)

	if err != nil {
		return "", err
	}
	return string(json), nil
}
