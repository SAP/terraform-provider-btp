package provider

import (
	"context"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/cis"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type directoryHierarchyType struct {
	ID          types.String `tfsdk:"id"`
	CreatedBy   types.String `tfsdk:"created_by"`
	CreatedDate types.String `tfsdk:"created_date"`
	// Directories	   types.List	`tfsdk:"directories"`
	DirectoryType types.String `tfsdk:"directory_type"`
	Features      types.Set    `tfsdk:"features"`
	ModifiedDate  types.String `tfsdk:"last_modified"`
	Name          types.String `tfsdk:"name"`
	ParentID      types.String `tfsdk:"parent_id"`
	ParentName    types.String `tfsdk:"parent_name"`
	ParentType    types.String `tfsdk:"parent_type"`
	State         types.String `tfsdk:"state"`
	Subaccounts   types.List   `tfsdk:"subaccounts"`
	Subdomain     types.String `tfsdk:"subdomain"`
	Type          types.String `tfsdk:"type"`
	// Description  types.String `tfsdk:"description"`
	// Labels       types.Map    `tfsdk:"labels"`
	// LastModified types.String `tfsdk:"last_modified"`
}

// func directoryObjectType( level int ) types.ObjectType{

// 	if level > 1 {
// 		return types.ObjectType{
// 			AttrTypes: map[string]attr.Type{
// 				"id":			types.StringType,
// 				"created_by":   types.StringType,
// 				"created_date": types.StringType,
// 				"directory_type": types.StringType,
// 				"directories":	types.ListType{
// 					ElemType: directoryObjectType(level-1),
// 				},
// 				"features": 	types.SetType{
// 					ElemType: types.StringType,
// 				},
// 				"last_modified": types.StringType,
// 				"name":          types.StringType,
// 				"parent_id":     types.StringType,
// 				"parent_name": 	 types.StringType,
// 				"parent_type":	 types.StringType,
// 				"state":         types.StringType,
// 				"subaccounts":	 types.ListType{
// 					ElemType: 	subaccountObjectType,
// 				},
// 				"subdomain":     types.StringType,
// 				"type":			 types.StringType,
// 			},
// 		}
// 	}

// 	return types.ObjectType{
// 		AttrTypes: map[string]attr.Type{
// 			"id":			types.StringType,
// 			"created_by":   types.StringType,
// 			"created_date": types.StringType,
// 			"features": 	types.SetType{
// 				ElemType: types.StringType,
// 			},
// 			"last_modified": types.StringType,
// 			"name":          types.StringType,
// 			"parent_id":     types.StringType,
// 			"parent_name": 	 types.StringType,
// 			"parent_type":	 types.StringType,
// 			"state":         types.StringType,
// 			"subaccounts":	 types.ListType{
// 				ElemType: 	subaccountObjectType,
// 			},
// 			"subdomain":     types.StringType,
// 			"type":			 types.StringType,
// 		},
// 	}
// }

func directoriesObjectType(level int) types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":           types.StringType,
			"created_by":   types.StringType,
			"created_date": types.StringType,
			// "directories":	types.ListType{
			// 	ElemType: directoryObjectType(level),
			// },
			"directory_type": types.StringType,
			"features": types.SetType{
				ElemType: types.StringType,
			},
			"last_modified": types.StringType,
			"name":          types.StringType,
			"parent_id":     types.StringType,
			"parent_name":   types.StringType,
			"parent_type":   types.StringType,
			"state":         types.StringType,
			"subaccounts": types.ListType{
				ElemType: types.ObjectType{ AttrTypes : subaccountObjectType.AttrTypes },
			},
			"subdomain": types.StringType,
			"type":      types.StringType,
		},
	}
}

func directoryHierarchyValueFrom(ctx context.Context, dirRes cis.DirectoryHierarchyResponseObject, parentName types.String, parentType types.String) (directoryHierarchyType, diag.Diagnostics) {

	var summary, diags diag.Diagnostics

	dir := directoryHierarchyType{
		ID:            types.StringValue(dirRes.Guid),
		CreatedBy:     types.StringValue(dirRes.CreatedBy),
		CreatedDate:   timeToValue(dirRes.CreatedDate.Time()),
		DirectoryType: types.StringValue(dirRes.DirectoryType),
		ModifiedDate:  timeToValue(dirRes.ModifiedDate.Time()),
		Name:          types.StringValue(dirRes.DisplayName),
		ParentID:      types.StringValue(dirRes.ParentGUID),
		ParentName:    parentName,
		ParentType:    parentType,
		State:         types.StringValue(dirRes.EntityState),
		Subdomain:     types.StringValue(dirRes.Subdomain),
		Type:          types.StringValue("Directory"),
	}

	dir.Features, diags = types.SetValueFrom(ctx, types.StringType, dirRes.DirectoryFeatures)
	summary.Append(diags...)

	if len(dirRes.Subaccounts) > 0 {
		subaccounts := subaccountsHierarchyValueFrom(ctx, dirRes.Subaccounts, dir.Name, dir.Type)
		dir.Subaccounts, diags = types.ListValueFrom(ctx, subaccountObjectType, subaccounts)
		summary.Append(diags...)
	}

	return dir, summary
}

func directoriesHierarchyValueFrom(ctx context.Context, dirResponses []cis.DirectoryHierarchyResponseObject, parentName types.String, parentType types.String, level int) ([]directoryHierarchyType, diag.Diagnostics) {

	var dirs = []directoryHierarchyType{}
	var summary diag.Diagnostics

	for _, dirRes := range dirResponses {
		var diags diag.Diagnostics

		dir, diags := directoryHierarchyValueFrom(ctx, dirRes, parentName, parentType)

		// if len(dirRes.Children) > 0 {
		// 	directories, diag := directoriesHierarchyValueFrom(ctx, dirRes.Children, dir.Name, dir.Type, level-1)
		// 	diags.Append(diag...)
		// 	dir.Directories, diag = types.ListValueFrom(ctx, directoriesObjectType(level-1), directories)
		// 	diags.Append(diag...)
		// }

		dirs = append(dirs, dir)
		summary.Append(diags...)
	}

	return dirs, summary
}
