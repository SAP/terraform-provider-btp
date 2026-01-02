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

type subaccountDestinationResourceType struct {
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
type subaccountDestinationType struct {
	SubaccountID            types.String         `tfsdk:"subaccount_id"`
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

func BuildDestinationConfigurationJSON(destination subaccountDestinationResourceType) (string, error) {
	config := map[string]any{}
	if !destination.Name.IsNull() && destination.Name.ValueString() != "" {
		config["Name"] = destination.Name.ValueString()
	}
	if !destination.Type.IsNull() && destination.Type.ValueString() != "" {
		config["Type"] = destination.Type.ValueString()
	}
	if !destination.ProxyType.IsNull() && destination.ProxyType.ValueString() != "" {
		config["ProxyType"] = destination.ProxyType.ValueString()
	}
	if !destination.URL.IsNull() && destination.URL.ValueString() != "" {
		config["URL"] = destination.URL.ValueString()
	}
	if !destination.Authentication.IsNull() && destination.Authentication.ValueString() != "" {
		config["Authentication"] = destination.Authentication.ValueString()
	}
	if !destination.Description.IsNull() && destination.Description.ValueString() != "" {
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

func destinationDatasourceValueFrom(value connectivity.DestinationResponse, subaccountID types.String, serviceInstanceID types.String) (subaccountDestinationType, diag.Diagnostics) {
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

func destinationResourceValueFrom(value connectivity.DestinationResponse, subaccountID types.String, serviceInstanceID types.String) (subaccountDestinationResourceType, diag.Diagnostics) {
	creationTimeString, err := strconv.ParseInt(value.SystemMetadata.CreationTime, 10, 64)
	if err != nil {
		diagnostics := diag.Diagnostics{
			diag.NewErrorDiagnostic("failed to convert creation time", err.Error()),
		}
		return subaccountDestinationResourceType{}, diagnostics
	}
	creationTime := time.UnixMilli(creationTimeString).UTC().Format(time.RFC3339)
	modificationTimeString, err := strconv.ParseInt(value.SystemMetadata.ModificationTime, 10, 64)
	if err != nil {
		diagnostics := diag.Diagnostics{
			diag.NewErrorDiagnostic("failed to convert modification time", err.Error()),
		}
		return subaccountDestinationResourceType{}, diagnostics
	}
	modifyTime := time.UnixMilli(modificationTimeString).UTC().Format(time.RFC3339)
	destination := subaccountDestinationResourceType{
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
	desc, url, proxy, auth := extract("Description"), extract("URL"), extract("ProxyType"), extract("Authentication")
	destination.CreationTime = types.StringValue(creationTime)
	destination.ModificationTime = types.StringValue(modifyTime)
	destination.Name = types.StringValue(extract("Name"))
	destination.Type = types.StringValue(extract("Type"))
	if desc == "" {
		destination.Description = types.StringNull()
	} else {
		destination.Description = types.StringValue(desc)
	}
	if url == "" {
		destination.URL = types.StringNull()
	} else {
		destination.URL = types.StringValue(url)
	}
	if proxy == "" {
		destination.ProxyType = types.StringNull()
	} else {
		destination.ProxyType = types.StringValue(proxy)
	}
	if auth == "" {
		destination.Authentication = types.StringNull()
	} else {
		destination.Authentication = types.StringValue(auth)
	}
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
			return subaccountDestinationResourceType{}, diagnostics
		}
		destination.AdditionalConfiguration = jsontypes.NewNormalizedValue(string(additionalJSON))
	}
	var diagnostics diag.Diagnostics

	return destination, diagnostics
}

// This function add the masked fields which are not fetched in read operation
func MergeAdditionalConfig(plannedConfig jsontypes.Normalized, responseConfig jsontypes.Normalized) (jsontypes.Normalized, error) {
	plannedMap := make(map[string]string)
	responseMap := make(map[string]string)

	if !plannedConfig.IsNull() && !plannedConfig.IsUnknown() {
		if err := json.Unmarshal([]byte(plannedConfig.ValueString()), &plannedMap); err != nil {
			return jsontypes.Normalized{}, err
		}
	}

	if !responseConfig.IsNull() && !responseConfig.IsUnknown() {
		if err := json.Unmarshal([]byte(responseConfig.ValueString()), &responseMap); err != nil {
			return jsontypes.Normalized{}, err
		}
	}

	for k, plannedVal := range plannedMap {
		responseVal, exists := responseMap[k]
		if !exists || responseVal == "" {
			responseMap[k] = plannedVal
		}
	}

	mergedJSON, err := json.Marshal(responseMap)
	if err != nil {
		return jsontypes.Normalized{}, err
	}

	return jsontypes.NewNormalizedValue(string(mergedJSON)), nil
}
