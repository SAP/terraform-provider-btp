package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

type ExtractKymaApiServerUrlFunction struct{}

var _ function.Function = &ExtractKymaApiServerUrlFunction{}

func NewExtractKymaApiServerUrlFunction() function.Function {
	return &ExtractKymaApiServerUrlFunction{}
}

func (f *ExtractKymaApiServerUrlFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "extract_kyma_api_server_url"
}

func (f *ExtractKymaApiServerUrlFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "extract_kyma_api_server_url",
		Description:         "Parses the label string of a Kyma environment instance and returns the value of the API Server URL.",
		MarkdownDescription: "Parses the label string of a Kyma environment instance and returns the value of the API Server URL.",

		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "kyma_label",
				Description:         "Label string of a Kyma environment instance",
				MarkdownDescription: "Label string of a Kyma environment instance",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f *ExtractKymaApiServerUrlFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var label string

	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &label))

	if resp.Error != nil {
		return
	}

	apiServerUrl, err := ExtractLabelValue(label, EnvironmentLabelKeyKymaApiServerUrl)

	if err != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError(err.Error()))
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, apiServerUrl))
}
