package provider

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/connectivity"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type subaccountDestinationType struct {
	SubaccountID            types.String         `tfsdk:"subaccount_id"`
	ID                      types.String         `tfsdk:"id"`
	CreationTime            types.String         `tfsdk:"creation_time"`
	Etag                    types.String         `tfsdk:"etag"`
	ModificationTime        types.String         `tfsdk:"modification_time"`
	Name                    types.String         `tfsdk:"name"`
	Type                    types.String         `tfsdk:"type"`
	ProxyType               types.String         `tfsdk:"proxy_type"`
	URL                     types.String         `tfsdk:"url"`
	Authentication          types.String         `tfsdk:"authentication"`
	Description             types.String         `tfsdk:"description"`
	ServiceInstanceID       types.String         `tfsdk:"service_instance_id"`
	AdditionalConfiguration jsontypes.Normalized `tfsdk:"additional_configuration"`
}

type subaccountDestinationName struct {
	Name types.String `tfsdk:"name"`
}

func BuildDestinationConfigurationJSON(destination subaccountDestinationType) (string, error) {
	config := map[string]any{}

	if !destination.Name.IsNull() {
		config["Name"] = destination.Name.ValueString()
	}
	if !destination.Type.IsNull() {
		config["Type"] = destination.Type.ValueString()
	}
	if !destination.ProxyType.IsNull() {
		config["ProxyType"] = destination.ProxyType.ValueString()
	}
	if !destination.URL.IsNull() {
		config["URL"] = destination.URL.ValueString()
	}
	if !destination.Authentication.IsNull() {
		config["Authentication"] = destination.Authentication.ValueString()
	}
	if !destination.Description.IsNull() {
		config["Description"] = destination.Description.ValueString()
	}

	if !destination.AdditionalConfiguration.IsNull() {
		var extra map[string]any
		err := json.Unmarshal([]byte(destination.AdditionalConfiguration.ValueString()), &extra)
		if err != nil {
			return "", err
		}
		for k, v := range extra {
			config[k] = v
		}
	}

	jsonBytes, err := json.Marshal(config)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func destinationValueFrom(value connectivity.DestinationResponse, subaccountID types.String, serviceInstanceID types.String) (subaccountDestinationType, diag.Diagnostics) {
	creationTimeString, err := strconv.ParseInt(value.SystemMetadata.CreationTime, 10, 64)
	if err != nil {
		diagnostics := diag.Diagnostics{
			diag.NewErrorDiagnostic("failed to convert creation time", err.Error()),
		}
		return subaccountDestinationType{}, diagnostics
	}
	creationTime := time.UnixMilli(creationTimeString).UTC().Format(time.RFC3339)
	modificationTimeString, err := strconv.ParseInt(value.SystemMetadata.ModificationTime, 10, 64)
	if err != nil {
		diagnostics := diag.Diagnostics{
			diag.NewErrorDiagnostic("failed to convert modification time", err.Error()),
		}
		return subaccountDestinationType{}, diagnostics
	}
	modifyTime := time.UnixMilli(modificationTimeString).UTC().Format(time.RFC3339)
	destination := subaccountDestinationType{
		Etag:         types.StringValue(value.SystemMetadata.Etag),
		SubaccountID: subaccountID,
	}

	tmp := make(map[string]string)
	for k, v := range value.DestinationConfiguration {
		tmp[k] = v
	}

	extract := func(key string) string {
		if v, ok := tmp[key]; ok {
			delete(tmp, key)
			return v
		}
		return ""
	}
	destination.CreationTime = types.StringValue(creationTime)
	destination.ModificationTime = types.StringValue(modifyTime)
	destination.Name = types.StringValue(extract("Name"))
	destination.Type = types.StringValue(extract("Type"))
	destination.Description = types.StringValue(extract("Description"))
	destination.URL = types.StringValue(extract("URL"))
	destination.ProxyType = types.StringValue(extract("ProxyType"))
	destination.Authentication = types.StringValue(extract("Authentication"))
	if serviceInstanceID.ValueString() == "" {
		destination.ServiceInstanceID = types.StringNull()
	} else {
		destination.ServiceInstanceID = serviceInstanceID
	}

	if len(tmp) == 0 {
		destination.AdditionalConfiguration = jsontypes.NewNormalizedNull()
	} else {
		additionalJSON, err := json.Marshal(tmp)
		if err != nil {
			diagnostics := diag.Diagnostics{
				diag.NewErrorDiagnostic("failed to marshal additional configuration", err.Error()),
			}
			return subaccountDestinationType{}, diagnostics
		}
		destination.AdditionalConfiguration = jsontypes.NewNormalizedValue(string(additionalJSON))
	}
	var diagnostics diag.Diagnostics

	return destination, diagnostics
}
