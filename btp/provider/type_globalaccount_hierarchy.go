package provider

import (
	"context"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/cis"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type globalaccountHierarchyType struct {
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

func globalaccountHierarchyValueFrom(ctx context.Context, value cis.GlobalAccountResponseObject) (globalaccountHierarchyType, diag.Diagnostics) {
	globalaccount := globalaccountHierarchyType{
		ID:           types.StringValue(value.Guid),
		CreatedDate:  timeToValue(value.CreatedDate.Time()),
		ModifiedDate: timeToValue(value.ModifiedDate.Time()),
		Name:         types.StringValue(value.DisplayName),
		State:        types.StringValue(value.EntityState),
		Subdomain:    types.StringValue(value.Subdomain),
		Type:         types.StringValue("Global Account"),
	}

	var summary, diags diag.Diagnostics

	if len(value.Children) > 0 {
		//The dirctory level is mentioned as 6 inorder to align with the schema strcuture defined as per the provider.
		directories, diags := directoriesHierarchyValueFrom(ctx, value.Children, globalaccount.Name, globalaccount.Type, 6)
		summary.Append(diags...)
		globalaccount.Directories, diags = types.ListValueFrom(ctx, directoryObjectType(6), directories)
		summary.Append(diags...)
	} else {
		globalaccount.Directories = types.ListNull(directoryObjectType(6))
	}

	if len(value.Subaccounts) > 0 {
		subaccounts := subaccountsHierarchyValueFrom(ctx, value.Subaccounts, globalaccount.Name, globalaccount.Type)
		globalaccount.Subaccounts, diags = types.ListValueFrom(ctx, subaccountObjectType, subaccounts)
		summary.Append(diags...)
	} else {
		globalaccount.Subaccounts = types.ListNull(subaccountObjectType)
	}

	return globalaccount, summary
}
