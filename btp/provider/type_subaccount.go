package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/cis"
)

const EntitlementFeature = "ENTITLEMENTS"
const AuthorizationFeature = "AUTHORIZATIONS"

var subdomainRegex = regexp.MustCompile("^[a-z0-9](?:[a-z0-9|-]{0,61}[a-z0-9])?$")

type subaccountDataSourceType struct {
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
	ContractStatus types.String `tfsdk:"contract_status"`
}

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
	ContractStatus      types.String `tfsdk:"contract_status"`
}

func subaccountValueFrom(ctx context.Context, value cis.SubaccountResponseObject) (subaccountType, diag.Diagnostics) {
	subaccount := subaccountType{
		ID:             types.StringValue(value.Guid),
		BetaEnabled:    types.BoolValue(value.BetaEnabled),
		CreatedBy:      types.StringValue(value.CreatedBy),
		CreatedDate:    timeToValue(value.CreatedDate.Time()),
		Description:    types.StringValue(value.Description),
		LastModified:   timeToValue(value.ModifiedDate.Time()),
		Name:           types.StringValue(value.DisplayName),
		ParentID:       types.StringValue(value.ParentGUID),
		Region:         types.StringValue(value.Region),
		State:          types.StringValue(value.State),
		Subdomain:      types.StringValue(value.Subdomain),
		Usage:          types.StringValue(value.UsedForProduction),
		ContractStatus: types.StringValue(value.ContractStatus),
	}

	var diags, diagnostics diag.Diagnostics

	subaccount.Labels, diags = types.MapValueFrom(ctx, types.SetType{ElemType: types.StringType}, value.Labels)
	diagnostics.Append(diags...)

	subaccount.ParentFeatures, diags = types.SetValueFrom(ctx, types.StringType, value.ParentFeatures)
	diagnostics.Append(diags...)

	return subaccount, diagnostics
}

func subaccountListValueFrom(ctx context.Context, value cis.SubaccountResponseObject) (subaccountType, diag.Diagnostics) {
	subaccount := subaccountType{
		ID:             types.StringValue(value.Guid),
		BetaEnabled:    types.BoolValue(value.BetaEnabled),
		CreatedBy:      types.StringValue(value.CreatedBy),
		CreatedDate:    timeToValue(value.CreatedDate.Time()),
		Description:    types.StringValue(value.Description),
		LastModified:   timeToValue(value.ModifiedDate.Time()),
		Name:           types.StringValue(value.DisplayName),
		ParentID:       types.StringValue(value.ParentGUID),
		Region:         types.StringValue(value.Region),
		State:          types.StringValue(value.State),
		Subdomain:      types.StringValue(value.Subdomain),
		Usage:          types.StringValue(value.UsedForProduction),
		ContractStatus: types.StringValue(value.ContractStatus),
	}

	var diags, diagnostics diag.Diagnostics

	subaccount.Labels, diags = types.MapValueFrom(ctx, types.SetType{ElemType: types.StringType}, value.Labels)
	diagnostics.Append(diags...)

	subaccount.ParentFeatures, diags = types.SetValueFrom(ctx, types.StringType, value.ParentFeatures)
	diagnostics.Append(diags...)

	return subaccount, diagnostics
}

func subaccountDataSourceValueFrom(ctx context.Context, value cis.SubaccountResponseObject) (subaccountDataSourceType, diag.Diagnostics) {
	subaccount := subaccountDataSourceType{
		ID:             types.StringValue(value.Guid),
		BetaEnabled:    types.BoolValue(value.BetaEnabled),
		CreatedBy:      types.StringValue(value.CreatedBy),
		CreatedDate:    timeToValue(value.CreatedDate.Time()),
		Description:    types.StringValue(value.Description),
		LastModified:   timeToValue(value.ModifiedDate.Time()),
		Name:           types.StringValue(value.DisplayName),
		ParentID:       types.StringValue(value.ParentGUID),
		Region:         types.StringValue(value.Region),
		State:          types.StringValue(value.State),
		Subdomain:      types.StringValue(value.Subdomain),
		Usage:          types.StringValue(value.UsedForProduction),
		ContractStatus: types.StringValue(value.ContractStatus),
	}

	var diags, diagnostics diag.Diagnostics

	subaccount.Labels, diags = types.MapValueFrom(ctx, types.SetType{ElemType: types.StringType}, value.Labels)
	diagnostics.Append(diags...)

	subaccount.ParentFeatures, diags = types.SetValueFrom(ctx, types.StringType, value.ParentFeatures)
	diagnostics.Append(diags...)

	return subaccount, diagnostics
}

func determineParentIdByFeature(cli *btpcli.ClientFacade, ctx context.Context, parentIdToVerify string, featureType string, globalAccountGUID string) (parentId string, isParentGlobalaccount bool, err error) {
	if !isParentIdValid(parentIdToVerify) {
		return "", false, fmt.Errorf("invalid parent ID: %s", parentIdToVerify)
	}

	if globalAccountGUID == parentIdToVerify {
		// The parentIdToVerify is the same as the global account GUID, so the parent is the global account.
		return globalAccountGUID, true, nil
	}

	// The parent ID is not the global account, we must traverse the hierarchy to determine the correct parent based on the feature.
	// Flow: Get the parent data, check the features of the parent for the required feature, and traverse up the tree if necessary.
	dataDirectory, _, err := cli.Accounts.Directory.Get(ctx, parentIdToVerify, "")
	if err != nil {
		// The parent is a unamanged directory, and the directory above is a managed directory.
		// We determine the adminDirectoryID of the unmanaged directory
		// we must fall back to the rate-limited flow to determine the information by global account hierarchy.
		return determineParentIdByFeatureByHierarchy(cli, ctx, parentIdToVerify, featureType)
	}

	if hasFeature(dataDirectory.DirectoryFeatures, featureType) {
		//The parent has the required feature
		return parentIdToVerify, false, nil
	}

	//The parent does not have the required feature, so we must traverse up the hierarchy to find the correct parent
	return determineParentIdByFeature(cli, ctx, dataDirectory.ParentGUID, featureType, globalAccountGUID)
}

func determineParentIdByFeatureByHierarchy(cli *btpcli.ClientFacade, ctx context.Context, parentIdToVerify string, featureType string) (parentId string, isParentGlobalaccount bool, err error) {
	if !isParentIdValid(parentIdToVerify) {
		return "", false, fmt.Errorf("invalid parent ID: %s", parentIdToVerify)
	}

	globalAccountHierarchy, _, err := cli.Accounts.GlobalAccount.GetWithHierarchy(ctx)
	if err != nil {
		return "", false, err
	}

	if parentIdToVerify == globalAccountHierarchy.Guid {
		return globalAccountHierarchy.Guid, true, nil
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

		if parentIdNew == "" {
			// An invalid parentIdToVerify was provided - add error with message that the parentIdToVerify was not found in the hierarchy
			return "", false, fmt.Errorf("the ID %s not found in the global account hierarchy", parentIdToVerify)
		}

		parentId = parentIdNew
	}

	return globalAccountHierarchy.Guid, true, nil
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

func determineParentIdForEntitlement(cli *btpcli.ClientFacade, ctx context.Context, parentIdToVerify string, closestEntitlementManagedParentGUID string) (parentId string, isParentGlobalaccount bool, err error) {

	if !isParentIdValid(parentIdToVerify) {
		return "", false, fmt.Errorf("invalid parent ID: %s", parentIdToVerify)
	}

	if parentIdToVerify == "" {
		// Import sceanrio
		return "", true, nil
	}

	dataGlobalAccount, _, err := cli.Accounts.GlobalAccount.Get(ctx)
	if err != nil {
		return "", false, fmt.Errorf("failed to get global account data: %w", err)
	}

	// The new value for the closest entitlement managed parent GUID takes precedence if it is provided.
	if closestEntitlementManagedParentGUID != "" {

		if closestEntitlementManagedParentGUID == dataGlobalAccount.Guid {
			return closestEntitlementManagedParentGUID, true, nil
		} else {
			return closestEntitlementManagedParentGUID, false, nil
		}
	}

	// Fallback: recursive determination of the parent ID by feature if the closest entitlement managed parent GUID is not provided.
	return determineParentIdByFeature(cli, ctx, parentIdToVerify, EntitlementFeature, dataGlobalAccount.Guid)
}

func determineParentIdForAuthorization(cli *btpcli.ClientFacade, ctx context.Context, parentIdToVerify string, closestEntitlementManagedParentGUID string) (parentId string, isParentGlobalaccount bool, err error) {

	if !isParentIdValid(parentIdToVerify) {
		return "", false, fmt.Errorf("invalid parent ID: %s", parentIdToVerify)
	}

	if parentIdToVerify == "" {
		// Import sceanrio
		return "", true, nil
	}

	dataGlobalAccount, _, err := cli.Accounts.GlobalAccount.Get(ctx)
	if err != nil {
		return "", false, fmt.Errorf("failed to get global account data: %w", err)
	}

	// If the parentIdToVerify matches the global account GUID, authorization is handled on global account level
	if dataGlobalAccount.Guid == parentIdToVerify {
		return dataGlobalAccount.Guid, true, nil
	}

	// If the closest entitlement managed parent GUID is provided, it takes precedence.
	// Prerequiste for the authorization feature is the entitlement feature, so the closest entitlement managed parent GUID is relevant.
	if closestEntitlementManagedParentGUID != "" {
		if closestEntitlementManagedParentGUID == dataGlobalAccount.Guid {
			// The closest entitlement managed parent GUID matches the global account GUID, so authorization is handled at the global account level.
			return closestEntitlementManagedParentGUID, true, nil
		} else {
			// The closest entitlement managed parent GUID does not match the global account GUID, i.e. it is a directory
			// However the authorization feature needs to be checked
			dataDirectory, _, err := cli.Accounts.Directory.Get(ctx, closestEntitlementManagedParentGUID, "")
			if err != nil {
				return "", false, fmt.Errorf("failed to get directory data: %w", err)
			}

			if hasFeature(dataDirectory.DirectoryFeatures, AuthorizationFeature) {
				// The directory has the authorization feature enabled, so it can be used as the parent for authorization.
				return closestEntitlementManagedParentGUID, false, nil
			}

			// the directory has the authorization feature disabled, so the parent for authorization is the global account.
			return dataGlobalAccount.Guid, true, nil
		}
	}

	// Fall back if no closest entitlement managed parent GUID is or can be provided
	return determineParentIdByFeature(cli, ctx, parentIdToVerify, AuthorizationFeature, dataGlobalAccount.Guid)
}

func isParentIdValid(parentId string) bool {
	return parentId != "00000000-0000-0000-0000-000000000000"
}
