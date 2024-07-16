package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/cis"
)

const EntitlementFeature = "ENTITLEMENTS"
const AuthorizationFeature = "AUTHORIZATIONS"

type subaccountType struct {
	ID             types.String `tfsdk:"id"`
	BetaEnabled    types.Bool   `tfsdk:"beta_enabled"`
	CreatedBy      types.String `tfsdk:"created_by"`
	CreatedDate    types.String `tfsdk:"created_date"`
	Description    types.String `tfsdk:"description"`
	Labels         types.Map    `tfsdk:"labels"`
	LastModified   types.String `tfsdk:"last_modified"`
	Name           types.String `tfsdk:"name"`
	ParentID       types.String `tfsdk:"parent_id"`
	ParentFeatures types.Set    `tfsdk:"parent_features"`
	Region         types.String `tfsdk:"region"`
	State          types.String `tfsdk:"state"`
	Subdomain      types.String `tfsdk:"subdomain"`
	Usage          types.String `tfsdk:"usage"`
}

func subaccountValueFrom(ctx context.Context, value cis.SubaccountResponseObject) (subaccountType, diag.Diagnostics) {
	subaccount := subaccountType{
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

func determineParentIdByFeature(cli *btpcli.ClientFacade, ctx context.Context, parentIdToVerify string, featureType string) (parentId string, isParentGlobalaccount bool) {
	parentData, _, err := cli.Accounts.Directory.Get(ctx, parentIdToVerify)

	// The parent is the global account
	if err != nil {
		isParentGlobalaccount = true
		return
	}

	if hasFeature(parentData.DirectoryFeatures, featureType) {
		// Parent is a directory with entitlements feature enabled
		parentId = parentIdToVerify
	} else {
		// Parent is a directory, but not with entitlements feature enabled -> step up the hierarchy
		parentId, isParentGlobalaccount = determineParentIdByFeature(cli, ctx, parentData.ParentGUID, featureType)
	}

	return
}

func hasFeature(features []string, featureType string) (featureTypeFound bool) {
	for _, f := range features {
		if f == featureType {
			featureTypeFound = true
		}
	}

	return
}

func determineParentIdForEntitlement(cli *btpcli.ClientFacade, ctx context.Context, parentIdToVerify string) (parentId string, isParentGlobalaccount bool) {
	return determineParentIdByFeature(cli, ctx, parentIdToVerify, EntitlementFeature)
}

func determineParentIdForAuthorization(cli *btpcli.ClientFacade, ctx context.Context, parentIdToVerify string) (parentId string, isParentGlobalaccount bool) {
	return determineParentIdByFeature(cli, ctx, parentIdToVerify, AuthorizationFeature)
}
