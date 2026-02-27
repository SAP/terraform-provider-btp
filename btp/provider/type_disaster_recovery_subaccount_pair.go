package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/cdr"
)

/*
type DisasterRecoverySubaccountPairType struct {
	ID                  types.String `tfsdk:"id"`
	BetaEnabled         types.Bool   `tfsdk:"beta_enabled"`
	CreatedBy           types.String `tfsdk:"created_by"`
	CreatedDate         types.String `tfsdk:"created_date"`
	Description         types.String `tfsdk:"description"`
	Labels              types.Map    `tfsdk:"labels"`
	LastModified        types.String `tfsdk:"last_modified"`
	Name                types.String `tfsdk:"name"`
	ParentID            types.String `tfsdk:"parent_id"`
	ParentFeatures      types.Set    `tfsdk:"parent_features"`
	Region              types.String `tfsdk:"region"`
	SkipAutoEntitlement types.Bool   `tfsdk:"skip_auto_entitlement"`
	State               types.String `tfsdk:"state"`
	Subdomain           types.String `tfsdk:"subdomain"`
	Usage               types.String `tfsdk:"usage"`
}

func DisasterRecoverySubaccountPairValueFrom(ctx context.Context, value cis.SubaccountResponseObject) (DisasterRecoverySubaccountPairType, diag.Diagnostics) {
	subaccount := DisasterRecoverySubaccountPairType{
		ID:           types.StringValue(value.Guid),
		BetaEnabled:  types.BoolValue(value.BetaEnabled),
		CreatedBy:    types.StringValue(value.CreatedBy),
		CreatedDate:  timeToValue(value.CreatedDate.Time()),
		Description:  types.StringValue(value.Description),
		LastModified: timeToValue(value.ModifiedDate.Time()),
		Name:         types.StringValue(value.DisplayName),
		ParentID:     types.StringValue(value.ParentGUID),
		Region:       types.StringValue(value.Region),
		State:        types.StringValue(value.State),
		Subdomain:    types.StringValue(value.Subdomain),
		Usage:        types.StringValue(value.UsedForProduction),
	}

	var diags, diagnostics diag.Diagnostics

	subaccount.Labels, diags = types.MapValueFrom(ctx, types.SetType{ElemType: types.StringType}, value.Labels)
	diagnostics.Append(diags...)

	subaccount.ParentFeatures, diags = types.SetValueFrom(ctx, types.StringType, value.ParentFeatures)
	diagnostics.Append(diags...)

	return subaccount, diagnostics
}
*/

type SubaccountDrMetadataDataSourceType struct {
	SubaccountId types.String `tfsdk:"subaccount_id"`
	Region       types.String `tfsdk:"region"`
	Subdomain    types.String `tfsdk:"subdomain"`
}

type DisasterRecoverySubaccountPairDataSourceType struct {
	SubaccountId    types.String `tfsdk:"subaccount_id"`
	PairId          types.String `tfsdk:"pair_id"`
	CreatedAt       types.String `tfsdk:"created_at"`
	CreatedBy       types.String `tfsdk:"created_by"`
	GlobalAccountId types.String `tfsdk:"global_account_id"`
	Subaccounts     types.List   `tfsdk:"subaccounts"`
}

func SubaccountPairDataSourceValueFrom(ctx context.Context, subaccountId types.String, value cdr.GetSubaccountPairResponse) (DisasterRecoverySubaccountPairDataSourceType, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	pairData := DisasterRecoverySubaccountPairDataSourceType{
		SubaccountId:    subaccountId,
		PairId:          types.StringValue(value.Id),
		CreatedAt:       types.StringValue(time.Unix(value.CreatedAt, 0).Format(time.RFC3339)),
		CreatedBy:       types.StringValue(value.CreatedBy),
		GlobalAccountId: types.StringValue(value.GlobalAccountId),
	}

	// Convert subaccounts to list
	subaccounts := make([]SubaccountDrMetadataDataSourceType, 0, len(value.Subaccounts))
	for _, sa := range value.Subaccounts {
		subaccounts = append(subaccounts, SubaccountDrMetadataDataSourceType{
			SubaccountId: types.StringValue(sa.Id),
			Region:       types.StringValue(sa.Region),
			Subdomain:    types.StringValue(sa.Subdomain),
		})
	}

	var diags diag.Diagnostics
	pairData.Subaccounts, diags = types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"subaccount_id": types.StringType,
			"region":        types.StringType,
			"subdomain":     types.StringType,
		},
	}, subaccounts)
	diagnostics.Append(diags...)

	return pairData, diagnostics
}
