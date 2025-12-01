package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

type ExtractCfOrgIdFunction struct{}

var _ function.Function = &ExtractCfOrgIdFunction{}

func NewEExtractCfOrgIdFunction() function.Function {
	return &ExtractCfOrgIdFunction{}
}

func (f *ExtractCfOrgIdFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "extract_cf_org_id"
}

func (f *ExtractCfOrgIdFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "extract_cf_org_id",
		Description:         "Parses the label string of a Cloud Foundry environment instance and returns the value of the CF Org ID.",
		MarkdownDescription: "Parses the label string of a Cloud Foundry environment instance and returns the value of the CF Org ID.",

		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "cf_label",
				Description:         "Label string of a Cloud Foundry environment instance",
				MarkdownDescription: "Label string of a Cloud Foundry environment instance",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f *ExtractCfOrgIdFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var label string

	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &label))

	if resp.Error != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError("Error reading input: "+resp.Error.Error()))
		return
	}

	org_id, err := ExtractLabelValue(label, EnvironmentLabelKeyCfOrgId)

	if err != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError(err.Error()))
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, org_id))
}
