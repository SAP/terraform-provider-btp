package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/cdr"
)

type DisasterRecoverySubaccountPairType struct {
	SubaccountId       types.String `tfsdk:"subaccount_id"`
	PairedSubaccountId types.String `tfsdk:"paired_subaccount_id"`
	PairId             types.String `tfsdk:"pair_id"`
	CreatedDate        types.String `tfsdk:"created_date"`
	CreatedBy          types.String `tfsdk:"created_by"`
	GlobalAccountId    types.String `tfsdk:"globalaccount_id"`
}

func disasterRecoverySubaccountPairValueFrom(ctx context.Context, subaccountId types.String, pairedSubaccountId types.String, value cdr.GetSubaccountPairResponse) (DisasterRecoverySubaccountPairType, diag.Diagnostics) {

	return DisasterRecoverySubaccountPairType{
		SubaccountId:       subaccountId,
		PairedSubaccountId: pairedSubaccountId,
		PairId:             types.StringValue(value.Id),
		CreatedDate:        types.StringValue(time.Unix(value.CreatedAt, 0).Format(time.RFC3339)),
		CreatedBy:          types.StringValue(value.CreatedBy),
		GlobalAccountId:    types.StringValue(value.GlobalAccountId),
	}, nil

}

type subaccountDrMetadataDataSourceType struct {
	SubaccountId types.String `tfsdk:"id"`
	Region       types.String `tfsdk:"region"`
	Subdomain    types.String `tfsdk:"subdomain"`
}

type disasterRecoverySubaccountPairDataSourceType struct {
	SubaccountId    types.String                         `tfsdk:"subaccount_id"`
	PairId          types.String                         `tfsdk:"pair_id"`
	CreatedDate     types.String                         `tfsdk:"created_date"`
	CreatedBy       types.String                         `tfsdk:"created_by"`
	GlobalAccountId types.String                         `tfsdk:"globalaccount_id"`
	Subaccounts     []subaccountDrMetadataDataSourceType `tfsdk:"subaccounts"`
}

func disasterRecoverySubaccountPairDataSourceValueFrom(ctx context.Context, subaccountId types.String, value cdr.GetSubaccountPairResponse) (disasterRecoverySubaccountPairDataSourceType, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	pairData := disasterRecoverySubaccountPairDataSourceType{
		SubaccountId:    subaccountId,
		PairId:          types.StringValue(value.Id),
		CreatedDate:     types.StringValue(time.Unix(value.CreatedAt, 0).Format(time.RFC3339)),
		CreatedBy:       types.StringValue(value.CreatedBy),
		GlobalAccountId: types.StringValue(value.GlobalAccountId),
	}

	// Convert subaccounts to list
	subaccounts := make([]subaccountDrMetadataDataSourceType, 0, len(value.Subaccounts))
	for _, sa := range value.Subaccounts {
		subaccounts = append(subaccounts, subaccountDrMetadataDataSourceType{
			SubaccountId: types.StringValue(sa.Id),
			Region:       types.StringValue(sa.Region),
			Subdomain:    types.StringValue(sa.Subdomain),
		})
	}

	pairData.Subaccounts = subaccounts

	return pairData, diagnostics
}
