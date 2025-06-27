package provider

import (
	"context"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/cis"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type directoryHierarchyType struct {
	ID           types.String `tfsdk:"id"`
	CreatedBy    types.String `tfsdk:"created_by"`
	CreatedDate  types.String `tfsdk:"created_date"`
	Directories  types.List   `tfsdk:"directories"`
	Features     types.Set    `tfsdk:"features"`
	ModifiedDate types.String `tfsdk:"last_modified"`
	Name         types.String `tfsdk:"name"`
	ParentID     types.String `tfsdk:"parent_id"`
	ParentName   types.String `tfsdk:"parent_name"`
	ParentType   types.String `tfsdk:"parent_type"`
	State        types.String `tfsdk:"state"`
	Subaccounts  types.List   `tfsdk:"subaccounts"`
	Subdomain    types.String `tfsdk:"subdomain"`
	Type         types.String `tfsdk:"type"`
}

func directoryObjectType(level int) types.ObjectType {

	if level > 1 {
		return types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":           types.StringType,
				"created_by":   types.StringType,
				"created_date": types.StringType,
				"directories": types.ListType{
					ElemType: directoryObjectType(level - 1),
				},
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

	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":           types.StringType,
			"created_by":   types.StringType,
			"created_date": types.StringType,
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

func directoryHierarchyValueFrom(ctx context.Context, dirRes cis.DirectoryResponseObject, parentName types.String, parentType types.String) (directoryHierarchyType, diag.Diagnostics) {

	var summary, diags diag.Diagnostics

	directory := directoryHierarchyType{
		ID:           types.StringValue(dirRes.Guid),
		CreatedBy:    types.StringValue(dirRes.CreatedBy),
		CreatedDate:  timeToValue(dirRes.CreatedDate.Time()),
		ModifiedDate: timeToValue(dirRes.ModifiedDate.Time()),
		Name:         types.StringValue(dirRes.DisplayName),
		ParentID:     types.StringValue(dirRes.ParentGUID),
		ParentName:   parentName,
		ParentType:   parentType,
		State:        types.StringValue(dirRes.EntityState),
		Subdomain:    types.StringValue(dirRes.Subdomain),
		Type:         types.StringValue("Directory"),
	}

	directory.Features, diags = types.SetValueFrom(ctx, types.StringType, dirRes.DirectoryFeatures)
	summary.Append(diags...)

	if len(dirRes.Subaccounts) > 0 {
		subaccounts := subaccountsHierarchyValueFrom(ctx, dirRes.Subaccounts, directory.Name, directory.Type)
		directory.Subaccounts, diags = types.ListValueFrom(ctx, subaccountObjectType, subaccounts)
		summary.Append(diags...)
	} else {
		directory.Subaccounts = types.ListNull(subaccountObjectType)
	}

	return directory, summary
}

func directoriesHierarchyValueFrom(ctx context.Context, dirResponses []cis.DirectoryResponseObject, parentName types.String, parentType types.String, level int) ([]directoryHierarchyType, diag.Diagnostics) {

	var directories = []directoryHierarchyType{}
	var summary diag.Diagnostics

	for _, dirRes := range dirResponses {
		var diags diag.Diagnostics

		directory, diags := directoryHierarchyValueFrom(ctx, dirRes, parentName, parentType)

		if len(dirRes.Children) > 0 && level > 1 {
			//Fetch the values of the sub-directories for the particular level
			subDirectories, diag := directoriesHierarchyValueFrom(ctx, dirRes.Children, directory.Name, directory.Type, level-1)
			diags.Append(diag...)
			directory.Directories, diag = types.ListValueFrom(ctx, directoryObjectType(level-1), subDirectories)
			diags.Append(diag...)
		} else if level > 1 {
			//Set the directories as null for the rest of the levels
			directory.Directories = types.ListNull(directoryObjectType(level - 1))
		}

		directories = append(directories, directory)
		summary.Append(diags...)
	}

	return directories, summary
}
