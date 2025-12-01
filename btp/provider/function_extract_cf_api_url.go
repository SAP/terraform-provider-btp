package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

type ExtractCfApiUrlFunction struct{}

var _ function.Function = &ExtractCfApiUrlFunction{}

func NewExtractCfApiUrlFunction() function.Function {
	return &ExtractCfApiUrlFunction{}
}

func (f *ExtractCfApiUrlFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "extract_cf_api_url"
}

func (f *ExtractCfApiUrlFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "extract_cf_api_url",
		Description:         "Parses the label string of a Cloud Foundry environment instance and returns the value of the CF API endpoint.",
		MarkdownDescription: "Parses the label string of a Cloud Foundry environment instance and returns the value of the CF API endpoint.",

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

func (f *ExtractCfApiUrlFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var label string

	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &label))

	if resp.Error != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError("Error reading input: "+resp.Error.Error()))
		return
	}

	endpoint, err := ExtractLabelValue(label, EnvironmentLabelKeyCfApiUrl)

	if err != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError(err.Error()))
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, endpoint))
}
