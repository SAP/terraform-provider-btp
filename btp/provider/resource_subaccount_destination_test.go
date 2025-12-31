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

func TestResourceSubaccountDestination(t *testing.T) {
	t.Parallel()
	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_destination_with_additional_variables_with_service_instance")
		defer stopQuietly(rec)
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBelow(tfversion.Version1_12_0),
			},
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDestinationWithServiceInstance("res1", "res1", "HTTP", "Internet", "https://myservice.example.com", "NoAuthentication", "testing resource for destination", "integration-test-destination", "servtest", map[string]string{"Abc": "test"}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_destination.res1", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_destination.res1", "service_instance_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount_destination.res1", "authentication", "NoAuthentication"),
						resource.TestMatchResourceAttr("btp_subaccount_destination.res1", "creation_time", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount_destination.res1", "description", "testing resource for destination"),
						resource.TestMatchResourceAttr("btp_subaccount_destination.res1", "modification_time", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount_destination.res1", "name", "res1"),
						resource.TestCheckResourceAttr("btp_subaccount_destination.res1", "proxy_type", "Internet"),
						resource.TestCheckResourceAttr("btp_subaccount_destination.res1", "type", "HTTP"),
						resource.TestCheckResourceAttr("btp_subaccount_destination.res1", "additional_configuration", "{\"Abc\":\"test\"}"),
					),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectIdentity("btp_subaccount_destination.res1", map[string]knownvalue.Check{
							"subaccount_id":       knownvalue.StringRegexp(regexpValidUUID),
							"name":                knownvalue.StringExact("res1"),
							"service_instance_id": knownvalue.StringRegexp(regexpValidUUID),
						}),
					},
				},
				{
					ResourceName:    "btp_subaccount_destination.res1",
					ImportState:     true,
					ImportStateKind: resource.ImportBlockWithResourceIdentity,
				},
			},
		})
	})
	t.Run("happy path with additional variables without service instance", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_destination_with_additional_variables_without_service_instance")
		defer stopQuietly(rec)
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDestination("res2", "res2", "HTTP", "Internet", "https://myservice.example.com", "NoAuthentication", "testing resource for destination", "integration-test-destination", map[string]string{"Abc": "test"}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_destination.res2", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount_destination.res2", "authentication", "NoAuthentication"),
						resource.TestMatchResourceAttr("btp_subaccount_destination.res2", "creation_time", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount_destination.res2", "description", "testing resource for destination"),
						resource.TestMatchResourceAttr("btp_subaccount_destination.res2", "modification_time", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount_destination.res2", "name", "res2"),
						resource.TestCheckResourceAttr("btp_subaccount_destination.res2", "proxy_type", "Internet"),
						resource.TestCheckResourceAttr("btp_subaccount_destination.res2", "type", "HTTP"),
						resource.TestCheckResourceAttr("btp_subaccount_destination.res2", "additional_configuration", "{\"Abc\":\"test\"}"),
					),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectIdentity("btp_subaccount_destination.res2", map[string]knownvalue.Check{
							"subaccount_id":       knownvalue.StringRegexp(regexpValidUUID),
							"name":                knownvalue.StringExact("res2"),
							"service_instance_id": knownvalue.Null(),
						}),
					},
				},
				{
					ResourceName:    "btp_subaccount_destination.res2",
					ImportState:     true,
					ImportStateKind: resource.ImportBlockWithResourceIdentity,
				},
			},
		})
	})
	t.Run("happy path without serv instance with additional variables for import id", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_destination_without_service_instance_with_additional_variables_for_import_id")
		defer stopQuietly(rec)
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDestination("res3", "res3", "HTTP", "Internet", "https://myservice.example.com", "NoAuthentication", "testing resource for destination", "integration-test-destination", map[string]string{"Abc": "test"}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_destination.res3", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount_destination.res3", "authentication", "NoAuthentication"),
						resource.TestMatchResourceAttr("btp_subaccount_destination.res3", "creation_time", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount_destination.res3", "description", "testing resource for destination"),
						resource.TestMatchResourceAttr("btp_subaccount_destination.res3", "modification_time", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount_destination.res3", "name", "res3"),
						resource.TestCheckResourceAttr("btp_subaccount_destination.res3", "proxy_type", "Internet"),
						resource.TestCheckResourceAttr("btp_subaccount_destination.res3", "type", "HTTP"),
					),
				},
				{
					ResourceName:      "btp_subaccount_destination.res3",
					ImportStateIdFunc: getIdForSubaccounDestinationFragmentImportId("btp_subaccount_destination.res3", "res3"),
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})
	t.Run("happy path TCP destination with additional variables", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_destination_tcp_with_additional_variables")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) +
						hclResourceDestinationTCP(
							"res_tcp",
							"res_tcp",
							"TCP",
							"Internet",
							"TCP destination test",
							"integration-test-destination",
							map[string]string{
								"Address": "host:1234",
							},
						),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_subaccount_destination.res_tcp", "type", "TCP"),
						resource.TestCheckResourceAttr(
							"btp_subaccount_destination.res_tcp",
							"additional_configuration",
							"{\"Address\":\"host:1234\"}",
						),
					),
				},
			},
		})
	})
	t.Run("happy path LDAP destination with additional configuration", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_destination_ldap_with_additional_configuration")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) +
						hclResourceDestinationLDAP(
							"res_ldap",
							"res_ldap",
							"LDAP",
							"integration-test-destination",
							map[string]string{
								"ldap.url":            "ldap://ldap.example.com:389",
								"ldap.proxyType":      "Internet",
								"ldap.description":    "LDAP destination test",
								"ldap.authentication": "BasicAuthentication",
								"ldap.user":           "abc",
								"ldap.password":       "abc",
							},
						),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_subaccount_destination.res_ldap", "type", "LDAP"),
						resource.TestMatchResourceAttr(
							"btp_subaccount_destination.res_ldap",
							"additional_configuration",
							regexp.MustCompile(`"ldap.url":"ldap://ldap.example.com:389"`),
						),
						resource.TestMatchResourceAttr(
							"btp_subaccount_destination.res_ldap",
							"additional_configuration",
							regexp.MustCompile(`"ldap.authentication":"BasicAuthentication"`),
						),
						resource.TestMatchResourceAttr(
							"btp_subaccount_destination.res_ldap",
							"additional_configuration",
							regexp.MustCompile(`"ldap.user":"abc"`),
						),
					),
				},
			},
		})
	})
	t.Run("happy path MAIL destination with additional configuration", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_destination_mail_with_additional_configuration")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) +
						hclResourceDestinationMAIL(
							"res_mail",
							"res_mail",
							"MAIL",
							"Internet",
							"BasicAuthentication",
							"integration-test-destination",
							map[string]string{
								"mail.url":         "mail://mail.example.com:389",
								"mail.description": "MAIL destination test",
								"mail.user":        "abc",
								"mail.password":    "abc",
							},
						),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_subaccount_destination.res_mail", "type", "MAIL"),
						resource.TestMatchResourceAttr(
							"btp_subaccount_destination.res_mail",
							"additional_configuration",
							regexp.MustCompile(`"mail.url":"mail://mail.example.com:389"`),
						),
						resource.TestMatchResourceAttr(
							"btp_subaccount_destination.res_mail",
							"additional_configuration",
							regexp.MustCompile(`"mail.user":"abc"`),
						),
					),
				},
			},
		})
	})
	t.Run("happy path RFC destination with additional configuration", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_destination_rfc_with_additional_configuration")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) +
						hclResourceDestinationRFC(
							"res_rfc",
							"res_rfc",
							"RFC",
							"integration-test-destination",
							map[string]string{
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
						resource.TestCheckResourceAttr(
							"btp_subaccount_destination.res_rfc",
							"type",
							"RFC",
						),
						resource.TestMatchResourceAttr(
							"btp_subaccount_destination.res_rfc",
							"additional_configuration",
							regexp.MustCompile(`"jco.client.ashost":"va4hci"`),
						),
						resource.TestMatchResourceAttr(
							"btp_subaccount_destination.res_rfc",
							"additional_configuration",
							regexp.MustCompile(`"jco.client.client":"001"`),
						),
						resource.TestMatchResourceAttr(
							"btp_subaccount_destination.res_rfc",
							"additional_configuration",
							regexp.MustCompile(`"jco.client.user":"SAPIPS"`),
						),
						resource.TestMatchResourceAttr(
							"btp_subaccount_destination.res_rfc",
							"additional_configuration",
							regexp.MustCompile(`"jco.destination.proxy_type":"OnPremise"`),
						),
					),
				},
			},
		})
	})

	t.Run("happy path update", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_destination_update")
		defer stopQuietly(rec)
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDestination("res4", "res4", "HTTP", "Internet", "https://myservice.example.com", "NoAuthentication", "testing resource for destination", "integration-test-destination", map[string]string{"Abc": "test"}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_destination.res4", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount_destination.res4", "authentication", "NoAuthentication"),
						resource.TestMatchResourceAttr("btp_subaccount_destination.res4", "creation_time", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount_destination.res4", "description", "testing resource for destination"),
						resource.TestMatchResourceAttr("btp_subaccount_destination.res4", "modification_time", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount_destination.res4", "name", "res4"),
						resource.TestCheckResourceAttr("btp_subaccount_destination.res4", "proxy_type", "Internet"),
						resource.TestCheckResourceAttr("btp_subaccount_destination.res4", "type", "HTTP"),
						resource.TestCheckResourceAttr("btp_subaccount_destination.res4", "additional_configuration", "{\"Abc\":\"test\"}"),
					),
				},
				{
					Config: hclProviderFor(user) + hclResourceDestination("res4", "res4", "HTTP", "Internet", "https://myservice.example.com", "NoAuthentication", "testing resource for destination update", "integration-test-destination", map[string]string{"Abc": "test"}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_destination.res4", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount_destination.res4", "authentication", "NoAuthentication"),
						resource.TestMatchResourceAttr("btp_subaccount_destination.res4", "creation_time", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount_destination.res4", "description", "testing resource for destination update"),
						resource.TestMatchResourceAttr("btp_subaccount_destination.res4", "modification_time", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount_destination.res4", "name", "res4"),
						resource.TestCheckResourceAttr("btp_subaccount_destination.res4", "proxy_type", "Internet"),
						resource.TestCheckResourceAttr("btp_subaccount_destination.res4", "type", "HTTP"),
						resource.TestCheckResourceAttr("btp_subaccount_destination.res4", "additional_configuration", "{\"Abc\":\"test\"}"),
					),
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
					Config:      `resource "btp_subaccount_destination" "res6" {}`,
					ExpectError: regexp.MustCompile(`The argument "subaccount_id" is required, but no definition was found.`),
				},
			},
		})
	})
	t.Run("error path - name must not contain slashes", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclResourceDestination("res7", "res/7", "HTTP", "Internet", "https://myservice.example.com", "NoAuthentication", "testing resource for destination", "integration-test-destination", map[string]string{"Abc": "test"}),
					ExpectError: regexp.MustCompile(`Attribute name must not contain '/', not be empty and not exceed 255`),
				},
			},
		})
	})
}
func hclResourceDestination(resourceName string, displayName string, destType string, proxyType string, url string, authentication string, description string, subaccountName string, additionalConfig map[string]string) string {

	var configBlock string
	if len(additionalConfig) > 0 {
		var configBuilder strings.Builder
		for k, v := range additionalConfig {
			configBuilder.WriteString(fmt.Sprintf("		%s = \"%s\"\n", k, v))
		}

		configBlock = fmt.Sprintf(`additional_configuration = jsonencode({
%s  	})`, configBuilder.String())
	}

	template := ` 
	data "btp_subaccounts" "all" {}
	resource "btp_subaccount_destination" "%s" {
	name           = "%s"
	type           = "%s"
	proxy_type     = "%s"
	url            = "%s"
	authentication = "%s"
	description    = "%s"
	subaccount_id     = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
	%s
}`

	return fmt.Sprintf(template, resourceName, displayName, destType, proxyType, url, authentication, description, subaccountName, configBlock)
}
func hclResourceDestinationTCP(resourceName string, displayName string, destType string, proxyType string, description string, subaccountName string, additionalConfig map[string]string) string {

	var configBlock string
	if len(additionalConfig) > 0 {
		var configBuilder strings.Builder
		for k, v := range additionalConfig {
			configBuilder.WriteString(fmt.Sprintf("		%s = \"%s\"\n", k, v))
		}

		configBlock = fmt.Sprintf(`additional_configuration = jsonencode({
%s  	})`, configBuilder.String())
	}

	template := ` 
	data "btp_subaccounts" "all" {}
	resource "btp_subaccount_destination" "%s" {
	name           = "%s"
	type           = "%s"
	proxy_type     = "%s"
	description    = "%s"
	subaccount_id     = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
	%s
}`

	return fmt.Sprintf(template, resourceName, displayName, destType, proxyType, description, subaccountName, configBlock)
}

func hclResourceDestinationLDAP(resourceName string, displayName string, destType string, subaccountName string, additionalConfig map[string]string) string {

	var configBlock string
	if len(additionalConfig) > 0 {
		var configBuilder strings.Builder
		for k, v := range additionalConfig {
			configBuilder.WriteString(fmt.Sprintf(`        "%s" = "%s"`+"\n", k, v))
		}

		configBlock = fmt.Sprintf(`additional_configuration = jsonencode({
%s  	})`, configBuilder.String())
	}

	template := ` 
	data "btp_subaccounts" "all" {}
	resource "btp_subaccount_destination" "%s" {
	name           = "%s"
	type           = "%s"
	subaccount_id     = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
	%s
}`

	return fmt.Sprintf(template, resourceName, displayName, destType, subaccountName, configBlock)
}
func hclResourceDestinationMAIL(resourceName string, displayName string, destType string, proxyType string, authentication string, subaccountName string, additionalConfig map[string]string) string {

	var configBlock string
	if len(additionalConfig) > 0 {
		var configBuilder strings.Builder
		for k, v := range additionalConfig {
			configBuilder.WriteString(fmt.Sprintf(`        "%s" = "%s"`+"\n", k, v))
		}

		configBlock = fmt.Sprintf(`additional_configuration = jsonencode({
%s  	})`, configBuilder.String())
	}

	template := ` 
	data "btp_subaccounts" "all" {}
	resource "btp_subaccount_destination" "%s" {
	name           = "%s"
	type           = "%s"
	proxy_type     = "%s"
	authentication = "%s"
	subaccount_id     = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
	%s
}`

	return fmt.Sprintf(template, resourceName, displayName, destType, proxyType, authentication, subaccountName, configBlock)
}

func hclResourceDestinationRFC(resourceName string, displayName string, destType string, subaccountName string, additionalConfig map[string]string) string {

	var configBlock string
	if len(additionalConfig) > 0 {
		var configBuilder strings.Builder
		for k, v := range additionalConfig {
			configBuilder.WriteString(fmt.Sprintf(`        "%s" = "%s"`+"\n", k, v))
		}

		configBlock = fmt.Sprintf(`additional_configuration = jsonencode({
%s  	})`, configBuilder.String())
	}

	template := ` 
	data "btp_subaccounts" "all" {}
	resource "btp_subaccount_destination" "%s" {
	name           = "%s"
	type           = "%s"
	subaccount_id     = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
	%s
}`

	return fmt.Sprintf(template, resourceName, displayName, destType, subaccountName, configBlock)
}

func hclResourceDestinationWithServiceInstance(resourceName string, displayName string, destType string, proxyType string, url string, authentication string, description string, subaccountName string, serviceInstanceName string, additionalConfig map[string]string) string {

	var configBlock string
	if len(additionalConfig) > 0 {
		var configBuilder strings.Builder
		for k, v := range additionalConfig {
			configBuilder.WriteString(fmt.Sprintf("		%s = \"%s\"\n", k, v))
		}

		configBlock = fmt.Sprintf(`additional_configuration = jsonencode({
%s  	})`, configBuilder.String())
	}

	template := ` 
	data "btp_subaccounts" "all" {}
	data "btp_subaccount_service_instance" "dest" {
  		subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
  		name          = "%s"
	}
	resource "btp_subaccount_destination" "%s" {
	name           = "%s"
	type           = "%s"
	proxy_type     = "%s"
	url            = "%s"
	authentication = "%s"
	description    = "%s"
	subaccount_id     = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
	service_instance_id = data.btp_subaccount_service_instance.dest.id
	%s
}`

	return fmt.Sprintf(template, subaccountName, serviceInstanceName, resourceName, displayName, destType, proxyType, url, authentication, description, subaccountName, configBlock)
}
