/*
 * Entitlements Service
 *
 * The Entitlements service provides REST APIs that manage the assignments of entitlements and quotas to subaccounts and directories.   Entitlements and their quota are automatically assigned to the global account when a customer order is fulfilled. Use the APIs in this service to manage the distribution of this global quota to your directories and subaccounts.   NOTE: These APIs are relevant only for cloud management tools feature set B. For details and information about whether this applies to your global account, see [Cloud Management Tools - Feature Set Overview](https://help.sap.com/viewer/65de2977205c403bbc107264b8eccf4b/Cloud/en-US/caf4e4e23aef4666ad8f125af393dfb2.html).  See also: * [Authorization](https://help.sap.com/viewer/65de2977205c403bbc107264b8eccf4b/latest/en-US/3670474a58c24ac2b082e76cbbd9dc19.html) * [Rate Limiting](https://help.sap.com/viewer/65de2977205c403bbc107264b8eccf4b/latest/en-US/77b217b3f57a45b987eb7fbc3305ce1e.html) * [Error Response Format](https://help.sap.com/viewer/65de2977205c403bbc107264b8eccf4b/latest/en-US/77fef2fb104b4b1795e2e6cee790e8b8.html) * [Asynchronous Jobs](https://help.sap.com/viewer/65de2977205c403bbc107264b8eccf4b/latest/en-US/0a0a6ab0ad114d72a6611c1c6b21683e.html)
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package cis

type DataCenterResponseObject struct {
	// Descriptive name of the data center for customer-facing UIs.
	DisplayName string `json:"displayName,omitempty"`
	// The domain of the data center
	Domain string `json:"domain,omitempty"`
	// The environment that the data center supports. For example: Kubernetes, Cloud Foundry.
	Environment string `json:"environment,omitempty"`
	// The infrastructure provider for the data center. Valid values: * <b>AWS:</b> Amazon Web Services. * <b>GCP:</b> Google Cloud Platform. * <b>AZURE:</b> Microsoft Azure. * <b>SAP:</b> SAP BTP (Neo). * <b>ALI:</b> Alibaba Cloud. * <b>IBM:</b> IBM Cloud.
	IaasProvider string `json:"iaasProvider,omitempty"`
	// Technical name of the data center. Must be unique within the cloud deployment.
	Name string `json:"name,omitempty"`
	// Provisioning service URL.
	ProvisioningServiceUrl string `json:"provisioningServiceUrl,omitempty"`
	// The region in which the data center is located.
	Region string `json:"region,omitempty"`
	// Saas-Registry service URL.
	SaasRegistryServiceUrl string `json:"saasRegistryServiceUrl,omitempty"`
	// Whether the specified datacenter supports trial accounts.
	SupportsTrial bool `json:"supportsTrial,omitempty"`
}