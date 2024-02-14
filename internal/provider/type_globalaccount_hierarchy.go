package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type globalaccountHierarchyType struct {
	ID             types.String `tfsdk:"id"`
	CreatedBy      types.String `tfsdk:"created_by"`
	CreatedDate    types.String `tfsdk:"created_date"`
	Directories    types.List   `tfsdk:"directories"`
	ModifiedDate   types.String	`tfdsk:"last_modified"`
	Name           types.String `tfsdk:"name"`
	Region         types.String `tfsdk:"region"`
	State          types.String `tfsdk:"state"`
	Subaccounts    types.List	`tfsdk:"subaccounts"`
	Subdomain      types.String `tfsdk:"subdomain"`
	Type		   types.String	`tfsdk:"type"`
}