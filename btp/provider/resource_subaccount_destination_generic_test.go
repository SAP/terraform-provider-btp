package provider

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestResourceSubaccountDestinationGeneric(t *testing.T) {
	t.Parallel()
	t.Run("happy path HTTP destination with service instance", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_destination_generic_http_with_service_instance")
		defer stopQuietly(rec)
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBelow(tfversion.Version1_12_0),
			},
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDestinationGenericWithServiceInstance("res1", "integration-test-destination", "servtest", map[string]string{
						"ProxyType":      "Internet",
						"URL":            "https://myservice.example.com",
						"Authentication": "NoAuthentication",
						"Description":    "trial destination of basic usecase",
						"Name":           "res1",
						"Type":           "HTTP",
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_destination_generic.res1", "service_instance_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_destination_generic.res1", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_destination_generic.res1", "creation_time", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_destination_generic.res1", "modification_time", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount_destination_generic.res1", "destination_configuration", "{\"Authentication\":\"NoAuthentication\",\"Description\":\"trial destination of basic usecase\",\"Name\":\"res1\",\"ProxyType\":\"Internet\",\"Type\":\"HTTP\",\"URL\":\"https://myservice.example.com\"}"),
					),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectIdentity("btp_subaccount_destination_generic.res1", map[string]knownvalue.Check{
							"subaccount_id":       knownvalue.StringRegexp(regexpValidUUID),
							"name":                knownvalue.StringExact("res1"),
							"service_instance_id": knownvalue.StringRegexp(regexpValidUUID),
						}),
					},
				},
				{
					ResourceName:    "btp_subaccount_destination_generic.res1",
					ImportState:     true,
					ImportStateKind: resource.ImportBlockWithResourceIdentity,
				},
			},
		})
	})
	t.Run("happy path HTTP destination", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_destination_generic_http")
		defer stopQuietly(rec)
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDestinationGeneric("res2", "integration-test-destination", map[string]string{
						"ProxyType":      "Internet",
						"URL":            "https://myservice.example.com",
						"Authentication": "NoAuthentication",
						"Description":    "trial destination of basic usecase",
						"Name":           "res2",
						"Type":           "HTTP",
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_destination_generic.res2", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_destination_generic.res2", "creation_time", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_destination_generic.res2", "modification_time", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount_destination_generic.res2", "destination_configuration", "{\"Authentication\":\"NoAuthentication\",\"Description\":\"trial destination of basic usecase\",\"Name\":\"res2\",\"ProxyType\":\"Internet\",\"Type\":\"HTTP\",\"URL\":\"https://myservice.example.com\"}"),
					),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectIdentity("btp_subaccount_destination_generic.res2", map[string]knownvalue.Check{
							"subaccount_id":       knownvalue.StringRegexp(regexpValidUUID),
							"name":                knownvalue.StringExact("res2"),
							"service_instance_id": knownvalue.Null(),
						}),
					},
				},
				{
					ResourceName:    "btp_subaccount_destination_generic.res2",
					ImportState:     true,
					ImportStateKind: resource.ImportBlockWithResourceIdentity,
				},
			},
		})
	})

	t.Run("happy path HTTP destination update", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_destination_generic_http_update")
		defer stopQuietly(rec)
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDestinationGeneric("res4", "integration-test-destination", map[string]string{
						"ProxyType":      "Internet",
						"URL":            "https://myservice.example.com",
						"Authentication": "NoAuthentication",
						"Description":    "trial destination of basic usecase",
						"Name":           "res4",
						"Type":           "HTTP",
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_destination_generic.res4", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_destination_generic.res4", "creation_time", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_destination_generic.res4", "modification_time", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount_destination_generic.res4", "destination_configuration", "{\"Authentication\":\"NoAuthentication\",\"Description\":\"trial destination of basic usecase\",\"Name\":\"res4\",\"ProxyType\":\"Internet\",\"Type\":\"HTTP\",\"URL\":\"https://myservice.example.com\"}"),
					),
				},
				{
					Config: hclProviderFor(user) + hclResourceDestinationGeneric("res4", "integration-test-destination", map[string]string{
						"ProxyType":      "Internet",
						"URL":            "https://myservice.example.com",
						"Authentication": "NoAuthentication",
						"Description":    "trial destination of basic usecase update",
						"Name":           "res4",
						"Type":           "HTTP",
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_destination_generic.res4", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_destination_generic.res4", "creation_time", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_destination_generic.res4", "modification_time", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount_destination_generic.res4", "destination_configuration", "{\"Authentication\":\"NoAuthentication\",\"Description\":\"trial destination of basic usecase update\",\"Name\":\"res4\",\"ProxyType\":\"Internet\",\"Type\":\"HTTP\",\"URL\":\"https://myservice.example.com\"}"),
					),
				},
			},
		})
	})
	t.Run("happy path TCP destination", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_destination_generic_tcp")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) +
						hclResourceDestinationGenericTCP(
							"res_tcp",
							"integration-test-destination",
							map[string]string{
								"Name":        "res_tcp",
								"Type":        "TCP",
								"Address":     "host:1234",
								"ProxyType":   "OnPremise",
								"Description": "TCP destination example",
							},
						),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(
							"btp_subaccount_destination_generic.res_tcp",
							"destination_configuration",
							"{\"Address\":\"host:1234\",\"Description\":\"TCP destination example\",\"Name\":\"res_tcp\",\"ProxyType\":\"OnPremise\",\"Type\":\"TCP\"}",
						),
					),
				},
			},
		})
	})
	t.Run("happy path LDAP destination", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_destination_generic_ldap")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) +
						hclResourceDestinationGenericLDAP(
							"res_ldap",
							"integration-test-destination",
							map[string]string{
								"Name":                "res_ldap",
								"Type":                "LDAP",
								"ldap.url":            "ldap://ldap.example.com:389",
								"ldap.proxyType":      "Internet",
								"ldap.description":    "LDAP destination test",
								"ldap.authentication": "BasicAuthentication",
								"ldap.user":           "abc",
								"ldap.password":       "abc",
							},
						),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr(
							"btp_subaccount_destination_generic.res_ldap",
							"destination_configuration",
							regexp.MustCompile(`"ldap.url":"ldap://ldap.example.com:389"`),
						),
						resource.TestMatchResourceAttr(
							"btp_subaccount_destination_generic.res_ldap",
							"destination_configuration",
							regexp.MustCompile(`"ldap.authentication":"BasicAuthentication"`),
						),
						resource.TestMatchResourceAttr(
							"btp_subaccount_destination_generic.res_ldap",
							"destination_configuration",
							regexp.MustCompile(`"ldap.user":"abc"`),
						),
					),
				},
			},
		})
	})
	t.Run("happy path MAIL destination", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_destination_generic_mail")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) +
						hclResourceDestinationGenericMAIL(
							"res_mail",
							"integration-test-destination",
							map[string]string{
								"Name":             "mail_dest",
								"Type":             "MAIL",
								"Authentication":   "BasicAuthentication",
								"ProxyType":        "OnPremise",
								"mail.url":         "mail://mail.example.com:389",
								"mail.description": "MAIL destination test",
								"mail.user":        "abc",
								"mail.password":    "abc",
							},
						),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr(
							"btp_subaccount_destination_generic.res_mail",
							"destination_configuration",
							regexp.MustCompile(`"mail.url":"mail://mail.example.com:389"`),
						),
						resource.TestMatchResourceAttr(
							"btp_subaccount_destination_generic.res_mail",
							"destination_configuration",
							regexp.MustCompile(`"mail.user":"abc"`),
						),
					),
				},
			},
		})
	})
	t.Run("happy path RFC destination", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_destination_generic_rfc")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) +
						hclResourceDestinationGenericRFC(
							"res_rfc",
							"integration-test-destination",
							map[string]string{
								"Name":                                  "rfc_dest",
								"Type":                                  "RFC",
								"jco.client.ashost":                     "va4hci",
								"jco.client.client":                     "001",
								"jco.client.delta":                      "1",
								"jco.client.network":                    "LAN",
								"jco.client.passwd":                     "Welcome1",
								"jco.client.serialization_format":       "rowBased",
								"jco.client.sysnr":                      "00",
								"jco.client.trace":                      "0",
								"jco.client.user":                       "SAPIPS",
								"jco.destination.auth_type":             "CONFIGURED_USER",
								"jco.destination.pool_check_connection": "0",
								"jco.destination.proxy_type":            "OnPremise",
							},
						),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr(
							"btp_subaccount_destination_generic.res_rfc",
							"destination_configuration",
							regexp.MustCompile(`"jco.client.ashost":"va4hci"`),
						),
						resource.TestMatchResourceAttr(
							"btp_subaccount_destination_generic.res_rfc",
							"destination_configuration",
							regexp.MustCompile(`"jco.client.client":"001"`),
						),
						resource.TestMatchResourceAttr(
							"btp_subaccount_destination_generic.res_rfc",
							"destination_configuration",
							regexp.MustCompile(`"jco.client.user":"SAPIPS"`),
						),
						resource.TestMatchResourceAttr(
							"btp_subaccount_destination_generic.res_rfc",
							"destination_configuration",
							regexp.MustCompile(`"jco.destination.proxy_type":"OnPremise"`),
						),
					),
				},
			},
		})
	})
	t.Run("error path HTTP destination invalid url", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_destination_generic_http_url_error")
		defer stopQuietly(rec)
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDestinationGeneric("res2", "integration-test-destination", map[string]string{
						"ProxyType":      "Internet",
						"URL":            "htt://myservice.example.com",
						"Authentication": "NoAuthentication",
						"Description":    "trial destination of basic usecase",
						"Name":           "res2",
						"Type":           "HTTP",
					}),
					ExpectError: regexp.MustCompile(`HTTP URL must start with http:// or https://`),
				},
			},
		})
	})

	t.Run("error path - subaccount not provided", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      `resource "btp_subaccount_destination_generic" "res6" {}`,
					ExpectError: regexp.MustCompile(`The argument "subaccount_id" is required, but no definition was found.`),
				},
			},
		})
	})
}
func hclResourceDestinationGeneric(resourceName string, subaccountName string, destinationConfig map[string]string) string {

	var configBlock string
	if len(destinationConfig) > 0 {
		var configBuilder strings.Builder
		for k, v := range destinationConfig {
			configBuilder.WriteString(fmt.Sprintf("		%s = \"%s\"\n", k, v))
		}

		configBlock = fmt.Sprintf(`destination_configuration = jsonencode({
%s  	})`, configBuilder.String())
	}

	template := `
	data "btp_subaccounts" "all" {}
	resource "btp_subaccount_destination_generic" "%s" {
	subaccount_id     = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
	%s
}`

	return fmt.Sprintf(template, resourceName, subaccountName, configBlock)
}

func hclResourceDestinationGenericWithServiceInstance(resourceName string, subaccountName string, serviceInstanceName string, destinationConfig map[string]string) string {

	var configBlock string
	if len(destinationConfig) > 0 {
		var configBuilder strings.Builder
		for k, v := range destinationConfig {
			configBuilder.WriteString(fmt.Sprintf("		%s = \"%s\"\n", k, v))
		}

		configBlock = fmt.Sprintf(`destination_configuration = jsonencode({
%s  	})`, configBuilder.String())
	}

	template := `
	data "btp_subaccounts" "all" {}
	data "btp_subaccount_service_instance" "dest" {
  		subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
  		name          = "%s"
	}
	resource "btp_subaccount_destination_generic" "%s" {
	subaccount_id     = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
	service_instance_id = data.btp_subaccount_service_instance.dest.id
	%s
}`

	return fmt.Sprintf(template, subaccountName, serviceInstanceName, resourceName, subaccountName, configBlock)
}

func hclResourceDestinationGenericTCP(resourceName string, subaccountName string, destinationConfig map[string]string) string {

	var configBlock string
	if len(destinationConfig) > 0 {
		var configBuilder strings.Builder
		for k, v := range destinationConfig {
			configBuilder.WriteString(fmt.Sprintf("		%s = \"%s\"\n", k, v))
		}

		configBlock = fmt.Sprintf(`destination_configuration = jsonencode({
%s  	})`, configBuilder.String())
	}

	template := `
	data "btp_subaccounts" "all" {}
	resource "btp_subaccount_destination_generic" "%s" {
	subaccount_id     = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
	%s
}`

	return fmt.Sprintf(template, resourceName, subaccountName, configBlock)
}

func hclResourceDestinationGenericLDAP(resourceName string, subaccountName string, destinationConfig map[string]string) string {

	var configBlock string
	if len(destinationConfig) > 0 {
		var configBuilder strings.Builder
		for k, v := range destinationConfig {
			configBuilder.WriteString(fmt.Sprintf(`        "%s" = "%s"`+"\n", k, v))
		}

		configBlock = fmt.Sprintf(`destination_configuration = jsonencode({
%s  	})`, configBuilder.String())
	}

	template := `
	data "btp_subaccounts" "all" {}
	resource "btp_subaccount_destination_generic" "%s" {
	subaccount_id     = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
	%s
}`

	return fmt.Sprintf(template, resourceName, subaccountName, configBlock)
}
func hclResourceDestinationGenericMAIL(resourceName string, subaccountName string, destinationConfig map[string]string) string {

	var configBlock string
	if len(destinationConfig) > 0 {
		var configBuilder strings.Builder
		for k, v := range destinationConfig {
			configBuilder.WriteString(fmt.Sprintf(`        "%s" = "%s"`+"\n", k, v))
		}

		configBlock = fmt.Sprintf(`destination_configuration = jsonencode({
%s  	})`, configBuilder.String())
	}

	template := `
	data "btp_subaccounts" "all" {}
	resource "btp_subaccount_destination_generic" "%s" {
	subaccount_id     = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
	%s
}`

	return fmt.Sprintf(template, resourceName, subaccountName, configBlock)
}

func hclResourceDestinationGenericRFC(resourceName string, subaccountName string, destinationConfig map[string]string) string {

	var configBlock string
	if len(destinationConfig) > 0 {
		var configBuilder strings.Builder
		for k, v := range destinationConfig {
			configBuilder.WriteString(fmt.Sprintf(`        "%s" = "%s"`+"\n", k, v))
		}

		configBlock = fmt.Sprintf(`destination_configuration = jsonencode({
%s  	})`, configBuilder.String())
	}

	template := `
	data "btp_subaccounts" "all" {}
	resource "btp_subaccount_destination_generic" "%s" {
	subaccount_id     = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
	%s
}`

	return fmt.Sprintf(template, resourceName, subaccountName, configBlock)
}
