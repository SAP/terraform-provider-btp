package provider

import (
	"context"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/destinations"
)

type subaccountDestinationCertificatesType struct {
	SubaccountId                types.String `tfsdk:"subaccount_id"`
	ServiceInstanceId           types.String `tfsdk:"service_instance_id"`
	SubaccountCertificates      types.List   `tfsdk:"subaccount_level_certificates"`
	ServiceInstanceCertificates types.List   `tfsdk:"service_instance_level_certificates"`
}

// Used for the certificate(s) data source
type subaccountDestinationCertificatesDataSourceType struct {
	CertificateName types.String `tfsdk:"certificate_name"`
	Nodes           types.List   `tfsdk:"certificate_nodes"`
	Creation        types.Object `tfsdk:"certification_creation_details"`
}

// Used for the certificate resource
type subaccountDestinationCertificateDataSourceType struct {
	SubaccountId      types.String `tfsdk:"subaccount_id"`
	ServiceInstanceId types.String `tfsdk:"service_instance_id"`
	CertificateName   types.String `tfsdk:"certificate_name"`
	Nodes             types.List   `tfsdk:"certificate_nodes"`
	Creation          types.Object `tfsdk:"certification_creation_details"`
}

type subaccountDestinationCertificateResourceType struct {
	SubaccountId       types.String `tfsdk:"subaccount_id"`
	ServiceInstanceId  types.String `tfsdk:"service_instance_id"`
	CertificateContent types.String `tfsdk:"certificate_content"`
	CertificateName    types.String `tfsdk:"certificate_name"`
	Nodes              types.List   `tfsdk:"certificate_nodes"`
	Creation           types.Object `tfsdk:"certification_creation_details"`
}

type DestinationCertificateNodeType struct {
	Type        types.String `tfsdk:"type"`
	Format      types.String `tfsdk:"format"`
	Algorithm   types.String `tfsdk:"algorithm"`
	Alias       types.String `tfsdk:"alias"`
	Subject     types.String `tfsdk:"subject"`
	Issuer      types.String `tfsdk:"issuer"`
	CommonName  types.String `tfsdk:"common_name"`
	NotBefore   types.String `tfsdk:"not_before"`
	NotAfter    types.String `tfsdk:"not_after"`
	Certificate types.String `tfsdk:"certificate"`
}

func subaccountDestinationCertificateValueFrom(ctx context.Context, value destinations.DestinationCertificateResponseObject) (subaccountDestinationCertificateResourceType, diag.Diagnostics) {

	var diagnostics diag.Diagnostics

	destinationCertificate := subaccountDestinationCertificateResourceType{
		CertificateName: types.StringValue(value.Name),
	}

	if len(value.Nodes) > 0 {
		nodes, diags := types.ListValueFrom(ctx, certificateNodeObjType, value.Nodes)
		diagnostics.Append(diags...)

		if diagnostics.HasError() {
			return destinationCertificate, diagnostics
		}

		nodesVal := []DestinationCertificateNodeType{}
		diags = nodes.ElementsAs(ctx, &nodesVal, true)
		diagnostics.Append(diags...)

		if diagnostics.HasError() {
			return destinationCertificate, diagnostics
		}

		// replaces all empty strings in each node with null
		nodes, diags = convertEmptyStringToNull(ctx, nodesVal)
		diagnostics.Append(diags...)

		if diagnostics.HasError() {
			return destinationCertificate, diagnostics
		}

		destinationCertificate.Nodes = nodes
	}

	creationData, diags := types.ObjectValueFrom(ctx, creationDataObjType.AttrTypes, value.Creation)
	diagnostics.Append(diags...)

	if diagnostics.HasError() {
		return destinationCertificate, diagnostics
	}

	destinationCertificate.Creation = creationData

	return destinationCertificate, diagnostics
}

func subaccountDestinationCertificatesValueFrom(ctx context.Context, values map[string][]destinations.DestinationCertificateResponseObject) (subaccountDestinationCertificatesType, diag.Diagnostics) {

	var diagnostics diag.Diagnostics

	allDestinationCertificates := subaccountDestinationCertificatesType{}

	for k, v := range values {

		destinationCertificates := []subaccountDestinationCertificatesDataSourceType{}
		certificates := v

		for _, cert := range certificates {
			certValue, diags := subaccountDestinationCertificateValueFrom(ctx, cert)
			diagnostics.Append(diags...)

			if diagnostics.HasError() {
				return subaccountDestinationCertificatesType{}, diagnostics
			}

			destinationCertificates = append(destinationCertificates, subaccountDestinationCertificatesDataSourceType{
				CertificateName: certValue.CertificateName,
				Nodes:           certValue.Nodes,
				Creation:        certValue.Creation,
			})
		}

		destinationCertificatesList, diags := types.ListValueFrom(ctx, subaccountDestinationCertificateObjType, destinationCertificates)
		diagnostics.Append(diags...)

		if diagnostics.HasError() {
			return subaccountDestinationCertificatesType{}, diagnostics
		}

		switch k {
		case "subaccount":
			if len(destinationCertificates) > 0 {
				allDestinationCertificates.SubaccountCertificates = destinationCertificatesList
			} else {
				allDestinationCertificates.SubaccountCertificates = types.ListNull(subaccountDestinationCertificateObjType)
			}
		case "serviceInstance":
			if len(destinationCertificates) > 0 {
				allDestinationCertificates.ServiceInstanceCertificates = destinationCertificatesList
			} else {
				allDestinationCertificates.ServiceInstanceCertificates = types.ListNull(subaccountDestinationCertificateObjType)
			}
		}
	}

	return allDestinationCertificates, diagnostics
}

func convertEmptyStringToNull(ctx context.Context, nodes []DestinationCertificateNodeType) (types.List, diag.Diagnostics) {

	var diags diag.Diagnostics

	for i, node := range nodes {

		v := reflect.ValueOf(&node).Elem()
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			if field.Type() == reflect.TypeOf(types.String{}) {
				strVal := field.Interface().(types.String)
				if strVal.ValueString() == "" {
					field.Set(reflect.ValueOf(types.StringNull()))
				}
			}
		}

		nodes[i] = node
	}

	nodesList, diags := types.ListValueFrom(ctx, certificateNodeObjType, nodes)

	return nodesList, diags
}
