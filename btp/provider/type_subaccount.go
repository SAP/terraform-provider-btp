package provider

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/cis"
)

const EntitlementFeature = "ENTITLEMENTS"
const AuthorizationFeature = "AUTHORIZATIONS"

var subdomainRegex = regexp.MustCompile("^[a-z0-9](?:[a-z0-9|-]{0,61}[a-z0-9])?$")

type subaccountType struct {
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

func determineParentIdByFeature(cli *btpcli.ClientFacade, ctx context.Context, parentIdToVerify string, featureType string) (parentId string, isParentGlobalaccount bool, err error) {
	if parentIdToVerify == "" {
		return "", true, nil
	}

	globalAccountHierarchy, _, err := cli.Accounts.GlobalAccount.GetWithHierarchy(ctx)
	if err != nil {
		return "", false, err
	}

	if parentIdToVerify == globalAccountHierarchy.Guid {
		return globalAccountHierarchy.GlobalAccountGUID, true, nil
	}

	parentId = parentIdToVerify
	parentIdNew := ""

	// Due to the structure of the hierarchy, we will end up at the root which is the global account.
	for parentId != globalAccountHierarchy.Guid {
		var parentFeatures []string
		parentFeatures, parentIdNew = findTargetFeaturesAndParent(parentId, globalAccountHierarchy.Children)
		if hasFeature(parentFeatures, featureType) {
			return parentId, false, nil
		}
		parentId = parentIdNew
	}

	return globalAccountHierarchy.GlobalAccountGUID, true, nil
}

func hasFeature(features []string, featureType string) (featureTypeFound bool) {
	for _, f := range features {
		if f == featureType {
			featureTypeFound = true
		}
	}

	return
}

func findTargetFeaturesAndParent(targetID string, hierarchy []cis.DirectoryResponseObject) (targetFeatures []string, parentId string) {

	for _, child := range hierarchy {
		if child.Guid == targetID {
			return child.DirectoryFeatures, child.ParentGUID
		}
	}

	for _, child := range hierarchy {
		targetFeatures, parentId = findTargetFeaturesAndParent(targetID, child.Children)
		if parentId != "" {
			return
		}
	}

	return
}

func determineParentIdForEntitlement(cli *btpcli.ClientFacade, ctx context.Context, parentIdToVerify string) (parentId string, isParentGlobalaccount bool, err error) {
	return determineParentIdByFeature(cli, ctx, parentIdToVerify, EntitlementFeature)
}

func determineParentIdForAuthorization(cli *btpcli.ClientFacade, ctx context.Context, parentIdToVerify string) (parentId string, isParentGlobalaccount bool, err error) {
	return determineParentIdByFeature(cli, ctx, parentIdToVerify, AuthorizationFeature)
}
