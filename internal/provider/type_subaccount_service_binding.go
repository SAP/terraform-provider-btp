package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/servicemanager"
)

type subaccountServiceBindingType struct {
	SubaccountId      types.String `tfsdk:"subaccount_id"`
	ServiceInstanceId types.String `tfsdk:"service_instance_id"`
	Name              types.String `tfsdk:"name"`
	Parameters        types.String `tfsdk:"parameters"`
	Id                types.String `tfsdk:"id"`
	Ready             types.Bool   `tfsdk:"ready"`
	Context           types.Map    `tfsdk:"context"`
	BindResource      types.Map    `tfsdk:"bind_resource"`
	Credentials       types.String `tfsdk:"credentials"`
	State             types.String `tfsdk:"state"`
	CreatedDate       types.String `tfsdk:"created_date"`
	LastModified      types.String `tfsdk:"last_modified"`
	Labels            types.Map    `tfsdk:"labels"`
}

func subaccountServiceBindingValueFrom(ctx context.Context, value servicemanager.ServiceBindingResponseObject) (subaccountServiceBindingType, diag.Diagnostics) {
	serviceBinding := subaccountServiceBindingType{
		SubaccountId:      types.StringValue(value.SubaccountId),
		Id:                types.StringValue(value.Id),
		Name:              types.StringValue(value.Name),
		Ready:             types.BoolValue(value.Ready),
		ServiceInstanceId: types.StringValue(value.ServiceInstanceId),
		Credentials:       types.StringValue(fmt.Sprintf("%s", value.Credentials)),
		State:             types.StringValue(value.LastOperation.State),
		CreatedDate:       timeToValue(value.CreatedAt),
		LastModified:      timeToValue(value.UpdatedAt),
	}

	var diags, diagnostics diag.Diagnostics

	serviceBinding.Context, diags = types.MapValueFrom(ctx, types.StringType, value.Context)
	diagnostics.Append(diags...)

	serviceBinding.BindResource, diags = types.MapValueFrom(ctx, types.StringType, value.BindResource)
	diagnostics.Append(diags...)

	serviceBinding.Labels, diags = types.MapValueFrom(ctx, types.SetType{ElemType: types.StringType}, value.Labels)
	diagnostics.Append(diags...)

	return serviceBinding, diagnostics
}
