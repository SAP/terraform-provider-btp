package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/cis"
)

type directoryType struct {
	ID           types.String `tfsdk:"id"`
	CreatedBy    types.String `tfsdk:"created_by"`
	CreatedDate  types.String `tfsdk:"created_date"`
	Description  types.String `tfsdk:"description"`
	Features     types.Set    `tfsdk:"features"`
	Labels       types.Map    `tfsdk:"labels"`
	LastModified types.String `tfsdk:"last_modified"`
	Name         types.String `tfsdk:"name"`
	ParentID     types.String `tfsdk:"parent_id"`
	State        types.String `tfsdk:"state"`
	Subdomain    types.String `tfsdk:"subdomain"`
}

func directoryValueFrom(ctx context.Context, value cis.DirectoryResponseObject) (directoryType, diag.Diagnostics) {
	directory := directoryType{
		ID:           types.StringValue(value.Guid),
		CreatedBy:    types.StringValue(value.CreatedBy),
		CreatedDate:  timeToValue(value.CreatedDate.Time()),
		Description:  types.StringValue(value.Description),
		LastModified: timeToValue(value.ModifiedDate.Time()),
		Name:         types.StringValue(value.DisplayName),
		ParentID:     types.StringValue(value.ParentGUID),
		State:        types.StringValue(value.EntityState),
		Subdomain:    types.StringValue(value.Subdomain),
	}

	var summary, diags diag.Diagnostics

	directory.Features, diags = types.SetValueFrom(ctx, types.StringType, value.DirectoryFeatures)
	summary.Append(diags...)

	directory.Labels, diags = types.MapValueFrom(ctx, types.SetType{ElemType: types.StringType}, value.Labels)
	summary.Append(diags...)

	return directory, summary
}
