package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

type ExtractKymaKubeconfigUrlFunction struct{}

var _ function.Function = &ExtractKymaKubeconfigUrlFunction{}

func NewExtractKymaKubeconfigUrlFunction() function.Function {
	return &ExtractKymaKubeconfigUrlFunction{}
}

func (f *ExtractKymaKubeconfigUrlFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "extract_kyma_kubeconfig_url"
}

func (f *ExtractKymaKubeconfigUrlFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "extract_kyma_kubeconfig_url",
		Description:         "Parses the label string of a Kyma environment instance and returns the value of the kubeconfig URL.",
		MarkdownDescription: "Parses the label string of a Kyma environment instance and returns the value of the kubeconfig URL.",

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

func (f *ExtractKymaKubeconfigUrlFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var label string

	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &label))

	if resp.Error != nil {
		return
	}

	kubeconfigUrl, err := ExtractLabelValue(label, EnvironmentLabelKeyKymaKubeconfigUrl)

	if err != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError(err.Error()))
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, kubeconfigUrl))
}
