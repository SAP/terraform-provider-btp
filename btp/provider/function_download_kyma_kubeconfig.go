package provider

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

// privateIPNets covers loopback, link-local, private, and cloud-metadata ranges.
var privateIPNets []*net.IPNet

func init() {
	for _, cidr := range []string{
		"127.0.0.0/8",    // loopback
		"::1/128",        // IPv6 loopback
		"10.0.0.0/8",     // private class A
		"172.16.0.0/12",  // private class B
		"192.168.0.0/16", // private class C
		"169.254.0.0/16", // link-local / cloud metadata (AWS, GCP, Azure)
		"fe80::/10",      // IPv6 link-local
		"fc00::/7",       // IPv6 unique-local
		"100.64.0.0/10",  // CGNAT (shared address space)
		"0.0.0.0/8",      // "this" network
		"::ffff:0:0/96",  // IPv4-mapped IPv6
	} {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err == nil {
			privateIPNets = append(privateIPNets, ipNet)
		}
	}
}

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

	if err := validateKymaKubeconfigURL(kubeconfigUrl); err != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError(err.Error()))
		return
	}

	client := &http.Client{
		Timeout: time.Second * 10,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if err := validateKymaKubeconfigURL(req.URL.String()); err != nil {
				return fmt.Errorf("redirect blocked: %w", err)
			}
			return nil
		},
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

// validateKymaKubeconfigURL checks scheme, host allowlist, and resolved IPs.
func validateKymaKubeconfigURL(rawURL string) error {
	if err := validateKymaKubeconfigURLSyntax(rawURL); err != nil {
		return err
	}
	parsed, _ := url.Parse(rawURL)
	return denyPrivateIPs(parsed.Hostname())
}

// validateKymaKubeconfigURLSyntax checks that the URL uses HTTPS and has a non-empty host.
func validateKymaKubeconfigURLSyntax(rawURL string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	if parsed.Scheme != "https" {
		return fmt.Errorf("URL scheme must be 'https', got %q", parsed.Scheme)
	}

	if parsed.Hostname() == "" {
		return fmt.Errorf("URL contains no host")
	}

	return nil
}

// denyPrivateIPs resolves hostname and rejects any address in a private/reserved range.
func denyPrivateIPs(hostname string) error {
	addrs, err := net.LookupHost(hostname)
	if err != nil {
		return fmt.Errorf("could not resolve host %q: %w", hostname, err)
	}

	for _, addr := range addrs {
		ip := net.ParseIP(addr)
		if ip == nil {
			return fmt.Errorf("could not parse resolved address %q", addr)
		}
		for _, ipNet := range privateIPNets {
			if ipNet.Contains(ip) {
				return fmt.Errorf("URL resolves to a disallowed IP address (%s)", addr)
			}
		}
	}
	return nil
}
