package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type subaccountHierarchyType struct {
	ID             types.String `tfsdk:"id"`
	CreatedBy      types.String `tfsdk:"created_by"`
	CreatedDate    types.String `tfsdk:"created_date"`
	ModifiedDate   types.String	`tfdsk:"last_modified"`
	Name           types.String `tfsdk:"name"`
	ParentID       types.String `tfsdk:"parent_id"`
	ParentName	   types.String `tfsdk:"parent_name"`
	ParentType     types.String	`tfsdk:"parent_type"`
	Region         types.String `tfsdk:"region"`
	State          types.String `tfsdk:"state"`
	Subdomain      types.String `tfsdk:"subdomain"`
	Type		   types.String	`tfsdk:"type"`

	// Description    types.String `tfsdk:"description"`
	// Labels         types.Map    `tfsdk:"labels"`
	// LastModified   types.String `tfsdk:"last_modified"`
	// ParentFeatures types.Set    `tfsdk:"parent_features"`
	// Usage          types.String `tfsdk:"usage"`
}