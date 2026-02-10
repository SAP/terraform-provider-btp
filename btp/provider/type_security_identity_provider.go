package provider

import (
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/xsuaa_authz"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type idpModel struct {
	TenantId     types.String `tfsdk:"tenant_id"`
	TenantType   types.String `tfsdk:"tenant_type"`
	DisplayName  types.String `tfsdk:"display_name"`
	CommonHost   types.String `tfsdk:"common_host"`
	Description  types.String `tfsdk:"description"`
	CustomHost   types.String `tfsdk:"custom_host"`
	CustomerName types.String `tfsdk:"customer_name"`
	CostCenterId types.Int64  `tfsdk:"cost_center_id"`
	DataCenterId types.String `tfsdk:"data_center_id"`
	Host         types.String `tfsdk:"host"`
	CustomerId   types.String `tfsdk:"customer_id"`
	Region       types.String `tfsdk:"region"`
	Status       types.String `tfsdk:"status"`
}

type globalaccountIdentityProviderDataSourceModel struct {
	idpModel
}

type subaccountIdentityProviderDataSourceModel struct {
	SubaccountId types.String `tfsdk:"subaccount_id"`
	idpModel
}

type globalaccountIdentityProvidersDataSourceModel struct {
	Values []idpModel `tfsdk:"values"`
}

type subaccountIdentityProvidersDataSourceModel struct {
	SubaccountId types.String `tfsdk:"subaccount_id"`
	Values       []idpModel   `tfsdk:"values"`
}

func mapIdpToModel(val xsuaa_authz.Idp) idpModel {
	return idpModel{
		TenantId:     types.StringValue(val.TenantId),
		TenantType:   types.StringValue(val.TenantType),
		DisplayName:  types.StringPointerValue(val.DisplayName),
		CommonHost:   types.StringValue(val.CommonHost),
		Description:  types.StringValue(val.Description),
		CustomHost:   types.StringPointerValue(val.CustomHost),
		CustomerName: types.StringPointerValue(val.CustomerName),
		CostCenterId: types.Int64Value(int64(val.CostCenterId)),
		DataCenterId: types.StringValue(val.DataCenterId),
		Host:         types.StringValue(val.Host),
		CustomerId:   types.StringPointerValue(val.CustomerId),
		Region:       types.StringValue(val.Region),
		Status:       types.StringValue(val.Status),
	}
}

func globalaccountIdentityProvidersDataSourceValueFrom(value []xsuaa_authz.Idp) (globalaccountIdentityProvidersDataSourceModel, diag.Diagnostics) {
	data := globalaccountIdentityProvidersDataSourceModel{Values: make([]idpModel, 0, len(value))}
	for _, val := range value {
		data.Values = append(data.Values, mapIdpToModel(val))
	}
	return data, nil
}

func subaccountIdentityProvidersDataSourceValueFrom(value []xsuaa_authz.Idp) (subaccountIdentityProvidersDataSourceModel, diag.Diagnostics) {
	data := subaccountIdentityProvidersDataSourceModel{Values: make([]idpModel, 0, len(value))}
	for _, val := range value {
		data.Values = append(data.Values, mapIdpToModel(val))
	}
	return data, nil
}

func globalaccountIdentityProviderDataSourceValueFrom(val xsuaa_authz.Idp) (globalaccountIdentityProviderDataSourceModel, diag.Diagnostics) {
	return globalaccountIdentityProviderDataSourceModel{idpModel: mapIdpToModel(val)}, nil
}

func subaccountIdentityProviderDataSourceValueFrom(val xsuaa_authz.Idp) (subaccountIdentityProviderDataSourceModel, diag.Diagnostics) {
	return subaccountIdentityProviderDataSourceModel{idpModel: mapIdpToModel(val)}, nil
}
