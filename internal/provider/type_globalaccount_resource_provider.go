package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/provisioning"
)

type globalaccountResourceProviderType struct {
	ResourceProvider types.String `tfsdk:"resource_provider"`
	Id               types.String `tfsdk:"id"`
	DisplayName      types.String `tfsdk:"display_name"`
	Description      types.String `tfsdk:"description"`
	Parameters       types.String `tfsdk:"parameters"`
}

func globalaccountResourceProviderValueFrom(ctx context.Context, value provisioning.ResourceProviderResponseObject) (globalaccountResourceProviderType, diag.Diagnostics) {
	resourceProvider := globalaccountResourceProviderType{
		ResourceProvider: types.StringValue(value.ResourceProvider),
		Id:               types.StringValue(value.ResourceTechnicalName),
		DisplayName:      types.StringValue(value.DisplayName),
		Description:      types.StringValue(value.Description),
	}

	if value.AdditionalInfo == nil {
		resourceProvider.Parameters = types.StringNull()
	} else {
		resourceProvider.Parameters = types.StringValue(fmt.Sprintf("%s", *value.AdditionalInfo))
	}

	return resourceProvider, diag.Diagnostics{}
}
