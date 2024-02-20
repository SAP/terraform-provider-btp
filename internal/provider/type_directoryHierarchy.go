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
	Directories	   types.List	`tfsdk:"directories"`
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
}


func directoryObjectType( level int ) types.ObjectType{

	if level > 1 {
		return types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":			types.StringType,
				"created_by":   types.StringType,
				"created_date": types.StringType,
				"directory_type": types.StringType,
				"directories":	types.ListType{
					ElemType: directoryObjectType(level-1),
				},
				"features": 	types.SetType{
					ElemType: types.StringType,
				},
				"last_modified": types.StringType,
				"name":          types.StringType,
				"parent_id":     types.StringType,
				"parent_name": 	 types.StringType,
				"parent_type":	 types.StringType,
				"state":         types.StringType,
				"subaccounts":	 types.ListType{
					ElemType: 	subaccountObjectType,
				},
				"subdomain":     types.StringType,
				"type":			 types.StringType,
			},
		}
	}

	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":			types.StringType,
			"created_by":   types.StringType,
			"created_date": types.StringType,
			"directory_type": types.StringType,
			"features": 	types.SetType{
				ElemType: types.StringType,
			},
			"last_modified": types.StringType,
			"name":          types.StringType,
			"parent_id":     types.StringType,
			"parent_name": 	 types.StringType,
			"parent_type":	 types.StringType,
			"state":         types.StringType,
			"subaccounts":	 types.ListType{
				ElemType: 	subaccountObjectType,
			},
			"subdomain":     types.StringType,
			"type":			 types.StringType,
		},
	}

	// return types.ObjectNull()
}

func directoriesObjectType(level int) types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":           types.StringType,
			"created_by":   types.StringType,
			"created_date": types.StringType,
			"directories":	types.ListType{
				ElemType: directoryObjectType(level),
			},
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
				ElemType: subaccountObjectType,
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
	} else {
		dir.Subaccounts = types.ListNull(subaccountObjectType)
	}

	return dir, summary
}

func directoriesHierarchyValueFrom(ctx context.Context, dirResponses []cis.DirectoryHierarchyResponseObject, parentName types.String, parentType types.String, level int) ([]directoryHierarchyType, diag.Diagnostics) {

	var dirs = []directoryHierarchyType{}
	var summary diag.Diagnostics

	for _, dirRes := range dirResponses {
		var diags diag.Diagnostics

		dir, diags := directoryHierarchyValueFrom(ctx, dirRes, parentName, parentType)

		if len(dirRes.Children) > 0 && level>0{
			directories, diag := directoriesHierarchyValueFrom(ctx, dirRes.Children, dir.Name, dir.Type, level-1)
			diags.Append(diag...)
			dir.Directories, diag = types.ListValueFrom(ctx, directoryObjectType(level), directories)
			diags.Append(diag...)
		} else if level>0{
			// nullDirs := []directoryHierarchyNoChildrenType{}
			// nullDirs = append(nullDirs, directoryHierarchyListNull(level-1))
			// dir.Directories, diags = types.ListValueFrom(ctx, directoryObjectType(level-1), nullDirs)

			//only 1 child
			dir.Directories = types.ListNull(directoryObjectType(level))
		}

		dirs = append(dirs, dir)
		summary.Append(diags...)
	}

	return dirs, summary
}

// func directoryHierarchyListNull (level int) directoryHierarchyNoChildrenType {
// 	return directoryHierarchyNoChildrenType{
// 		ID: types.StringNull(),
// 		CreatedBy:   types.StringNull(),
// 		CreatedDate: types.StringNull(),
// 		DirectoryType: types.StringNull(),
// 		Features: types.SetNull(types.StringType),
// 		ModifiedDate: types.StringNull(),
// 		Name: types.StringNull(),
// 		ParentID: types.StringNull(),
// 		ParentName: types.StringNull(),
// 		ParentType: types.StringNull(),
// 		State: types.StringNull(),
// 		Subaccounts : types.ListNull(subaccountObjectType),
// 		Subdomain : types.StringNull(),
// 		Type: types.StringNull(),
// 	}
// }