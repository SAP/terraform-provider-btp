package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type directoryHierarchyType struct {
	ID             types.String `tfsdk:"id"`
	CreatedBy      types.String `tfsdk:"created_by"`
	CreatedDate    types.String `tfsdk:"created_date"`
	Directories	   types.List	`tfsdk:"directories"`
	DirectoryType  types.String	`tfsdk:"directory_type"`
	Features       types.Set    `tfsdk:"features"`
	ModifiedDate   types.String `tfsdk:"modified_date"`
	Name           types.String `tfsdk:"name"`
	ParentID       types.String `tfsdk:"parent_id"`
	ParentName 	   types.String `tfsdk:"parent_name"`
	ParentType 	   types.String `tfsdk:"parent_type"`
	State          types.String `tfsdk:"state"`
	Subaccounts    types.List	`tfsdk:"subaccounts"`
	Subdomain      types.String `tfsdk:"subdomain"`
	Type 		   types.String `tfsdk:"type"`
	// Description  types.String `tfsdk:"description"`
	// Labels       types.Map    `tfsdk:"labels"`
	// LastModified types.String `tfsdk:"last_modified"`
}