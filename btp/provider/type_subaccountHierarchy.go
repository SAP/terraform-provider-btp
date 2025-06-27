package provider

import (
	"context"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/cis"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type subaccountHierarchyType struct {
	ID           types.String `tfsdk:"id"`
	CreatedBy    types.String `tfsdk:"created_by"`
	CreatedDate  types.String `tfsdk:"created_date"`
	ModifiedDate types.String `tfsdk:"last_modified"`
	Name         types.String `tfsdk:"name"`
	ParentID     types.String `tfsdk:"parent_id"`
	ParentName   types.String `tfsdk:"parent_name"`
	ParentType   types.String `tfsdk:"parent_type"`
	Region       types.String `tfsdk:"region"`
	State        types.String `tfsdk:"state"`
	Subdomain    types.String `tfsdk:"subdomain"`
	Type         types.String `tfsdk:"type"`
}

var subaccountObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":            types.StringType,
		"created_by":    types.StringType,
		"created_date":  types.StringType,
		"last_modified": types.StringType,
		"name":          types.StringType,
		"parent_id":     types.StringType,
		"parent_name":   types.StringType,
		"parent_type":   types.StringType,
		"region":        types.StringType,
		"state":         types.StringType,
		"subdomain":     types.StringType,
		"type":          types.StringType,
	},
}

func subaccountHierarchyValueFrom(ctx context.Context, subRes cis.SubaccountResponseObject) subaccountHierarchyType {
	subaccount := subaccountHierarchyType{
		ID:           types.StringValue(subRes.Guid),
		CreatedBy:    types.StringValue(subRes.CreatedBy),
		CreatedDate:  timeToValue(subRes.CreatedDate.Time()),
		ModifiedDate: timeToValue(subRes.ModifiedDate.Time()),
		Name:         types.StringValue(subRes.DisplayName),
		ParentID:     types.StringValue(subRes.ParentGUID),
		Region:       types.StringValue(subRes.Region),
		State:        types.StringValue(subRes.State),
		Subdomain:    types.StringValue(subRes.Subdomain),
		Type:         types.StringValue("Subaccount"),
	}
	return subaccount
}

func subaccountsHierarchyValueFrom(ctx context.Context, subResponses []cis.SubaccountResponseObject, parentName types.String, parentType types.String) []subaccountHierarchyType {

	var subaccounts = []subaccountHierarchyType{}

	for _, subRes := range subResponses {
		subaccount := subaccountHierarchyValueFrom(ctx, subRes)
		subaccount.ParentName = parentName
		subaccount.ParentType = parentType
		subaccounts = append(subaccounts, subaccount)
	}

	return subaccounts
}
