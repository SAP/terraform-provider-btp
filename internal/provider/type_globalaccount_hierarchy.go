package provider

import (
	"context"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/cis"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type globalAccountHierarchyType struct {
	ID           types.String `tfsdk:"id"`
	CreatedDate  types.String `tfsdk:"created_date"`
	Directories  types.List   `tfsdk:"directories"`
	ModifiedDate types.String `tfsdk:"last_modified"`
	Name         types.String `tfsdk:"name"`
	State        types.String `tfsdk:"state"`
	Subaccounts  types.List   `tfsdk:"subaccounts"`
	Subdomain    types.String `tfsdk:"subdomain"`
	Type         types.String `tfsdk:"type"`
}

//get global account value
//for directories it has to call get directories value
//same for subaccounts

func globalAccountHierarchyValueFrom(ctx context.Context, value cis.GlobalAccountHierarchyResponseObject) (globalAccountHierarchyType, diag.Diagnostics) {
	globalAccount := globalAccountHierarchyType{
		ID:           types.StringValue(value.Guid),
		CreatedDate:  timeToValue(value.CreatedDate.Time()),
		ModifiedDate: timeToValue(value.ModifiedDate.Time()),
		Name:         types.StringValue(value.DisplayName),
		State:        types.StringValue(value.EntityState),
		Subdomain:    types.StringValue(value.Subdomain),
		Type:         types.StringValue("Global Account"),
	}

	var summary, diags diag.Diagnostics
	var dirs []directoryHierarchyType

	if len(value.Children) > 0 {
		dirs, diags = directoriesHierarchyValueFrom(ctx, value.Children, globalAccount.Name, globalAccount.Type, 2)
		summary.Append(diags...)

		globalAccount.Directories, diags = types.ListValueFrom(ctx, directoriesObjectType(2), dirs)
		summary.Append(diags...)
	}

	if len(value.Subaccounts) > 0 {
		subaccounts := subaccountsHierarchyValueFrom(ctx, value.Subaccounts, globalAccount.Name, globalAccount.Type)
		globalAccount.Subaccounts, diags = types.ListValueFrom(ctx, subaccountObjectType, subaccounts)
		summary.Append(diags...)
	} else {
		globalAccount.Subaccounts = types.ListNull(subaccountObjectType)
	}

	return globalAccount, summary
}
