package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/provisioning"
)

type globalaccountResourceProviderType struct {
	Provider      types.String `tfsdk:"provider_type"`
	TechnicalName types.String `tfsdk:"technical_name"`
	Id            types.String `tfsdk:"id"`
	DisplayName   types.String `tfsdk:"display_name"`
	Description   types.String `tfsdk:"description"`
	Configuration types.String `tfsdk:"configuration"`
}

func globalaccountResourceProviderValueFrom(ctx context.Context, value provisioning.ResourceProviderResponseObject) (globalaccountResourceProviderType, diag.Diagnostics) {
	resourceProvider := globalaccountResourceProviderType{
		Provider:      types.StringValue(value.ResourceProvider),
		TechnicalName: types.StringValue(value.TechnicalName),
		Id:            types.StringValue(value.TechnicalName),
		DisplayName:   types.StringValue(value.DisplayName),
		Description:   types.StringValue(value.Description),
	}

	if value.AdditionalInfo == nil {
		resourceProvider.Configuration = types.StringNull()
	} else {
		resourceProvider.Configuration = types.StringValue(string(*value.AdditionalInfo))
	}

	return resourceProvider, diag.Diagnostics{}
}
