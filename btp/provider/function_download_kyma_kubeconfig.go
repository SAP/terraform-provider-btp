package provider

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

type DownloadKymaKubeconfigFunction struct{}

var _ function.Function = &DownloadKymaKubeconfigFunction{}

func NewDownloadKymaKubeconfigFunction() function.Function {
	return &DownloadKymaKubeconfigFunction{}
}

func (f *DownloadKymaKubeconfigFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "download_kyma_kubeconfig"
}

func (f *DownloadKymaKubeconfigFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "download_kyma_kubeconfig",
		Description:         "Downloads the kubeconfig of a Kyma environment.",
		MarkdownDescription: "Downloads the kubeconfig of a Kyma environment.",

		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "kyma_kubeconfig_url",
				Description:         "The URL of the Kyma environment kubeconfig",
				MarkdownDescription: "The URL of the Kyma environment kubeconfig",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f *DownloadKymaKubeconfigFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var kubeconfigUrl string

	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &kubeconfigUrl))

	if resp.Error != nil {
		return
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	httpResponse, err := client.Get(kubeconfigUrl)
	if err != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError(fmt.Sprintf("Error making HTTP request: %s", err.Error())))
		return
	}
	defer func() {
		// Ignore error on close intentionally as there is nothing to do when it fails
		_ = httpResponse.Body.Close()
	}()

	kubeconfig, err := io.ReadAll(httpResponse.Body)

	if err != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError(fmt.Sprintf("Error reading response body: %s", err.Error())))
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, string(kubeconfig)))
}
