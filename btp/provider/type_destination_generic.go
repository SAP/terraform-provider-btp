package provider

import (
	"encoding/json"
	"maps"
	"strconv"
	"time"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/connectivity"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type subaccountDestinationGenericResourceType struct {
	SubaccountID             types.String         `tfsdk:"subaccount_id"`
	ID                       types.String         `tfsdk:"id"`
	CreationTime             types.String         `tfsdk:"creation_time"`
	Etag                     types.String         `tfsdk:"etag"`
	Name                     types.String         `tfsdk:"name"`
	ModificationTime         types.String         `tfsdk:"modification_time"`
	ServiceInstanceID        types.String         `tfsdk:"service_instance_id"`
	DestinationConfiguration jsontypes.Normalized `tfsdk:"destination_configuration"`
}

type subaccountDestinationGenericType struct {
	SubaccountID             types.String         `tfsdk:"subaccount_id"`
	CreationTime             types.String         `tfsdk:"creation_time"`
	Etag                     types.String         `tfsdk:"etag"`
	ModificationTime         types.String         `tfsdk:"modification_time"`
	Name                     types.String         `tfsdk:"name"`
	ServiceInstanceID        types.String         `tfsdk:"service_instance_id"`
	DestinationConfiguration jsontypes.Normalized `tfsdk:"destination_configuration"`
}

// for excracting name and creation of json string from the destination json configuration
func BuildDestinationGenericConfigurationJSON(destination subaccountDestinationGenericResourceType) (string, string, error) {
	config := map[string]any{}
	name := ""
	if !destination.DestinationConfiguration.IsNull() {
		var extra map[string]any
		err := json.Unmarshal([]byte(destination.DestinationConfiguration.ValueString()), &extra)
		if err != nil {
			return "", name, err
		}
		maps.Copy(config, extra)
	}
	name, _ = config["Name"].(string)

	jsonBytes, err := json.Marshal(config)
	if err != nil {
		return "", name, err
	}
	return string(jsonBytes), name, nil
}

func destinationGenericResourceValueFrom(value connectivity.DestinationResponse, subaccountID types.String, serviceInstanceID types.String, name string) (subaccountDestinationGenericResourceType, diag.Diagnostics) {
	creationTimeString, err := strconv.ParseInt(value.SystemMetadata.CreationTime, 10, 64)
	if err != nil {
		diagnostics := diag.Diagnostics{
			diag.NewErrorDiagnostic("failed to convert creation time", err.Error()),
		}
		return subaccountDestinationGenericResourceType{}, diagnostics
	}
	creationTime := time.UnixMilli(creationTimeString).UTC().Format(time.RFC3339)

	modificationTimeString, err := strconv.ParseInt(value.SystemMetadata.ModificationTime, 10, 64)
	if err != nil {
		diagnostics := diag.Diagnostics{
			diag.NewErrorDiagnostic("failed to convert modification time", err.Error()),
		}
		return subaccountDestinationGenericResourceType{}, diagnostics
	}
	modifyTime := time.UnixMilli(modificationTimeString).UTC().Format(time.RFC3339)

	destination := subaccountDestinationGenericResourceType{
		Etag:         types.StringValue(value.SystemMetadata.Etag),
		SubaccountID: subaccountID,
	}

	tmp := make(map[string]string)
	maps.Copy(tmp, value.DestinationConfiguration)

	destination.CreationTime = types.StringValue(creationTime)
	destination.ModificationTime = types.StringValue(modifyTime)
	destination.Name = types.StringValue(name)
	if serviceInstanceID.ValueString() == "" {
		destination.ServiceInstanceID = types.StringNull()
	} else {
		destination.ServiceInstanceID = serviceInstanceID
	}

	destinationJSON, err := json.Marshal(tmp)
	if err != nil {
		diagnostics := diag.Diagnostics{
			diag.NewErrorDiagnostic("failed to marshal additional configuration", err.Error()),
		}
		return subaccountDestinationGenericResourceType{}, diagnostics
	}
	destination.DestinationConfiguration = jsontypes.NewNormalizedValue(string(destinationJSON))

	var diagnostics diag.Diagnostics

	return destination, diagnostics
}

func destinationGenericDatasourceValueFrom(value connectivity.DestinationResponse, subaccountID types.String, serviceInstanceID types.String) (subaccountDestinationGenericType, diag.Diagnostics) {

	creationTimeString, err := strconv.ParseInt(value.SystemMetadata.CreationTime, 10, 64)
	if err != nil {
		diagnostics := diag.Diagnostics{
			diag.NewErrorDiagnostic("failed to convert creation time", err.Error()),
		}
		return subaccountDestinationGenericType{}, diagnostics
	}
	creationTime := time.UnixMilli(creationTimeString).UTC().Format(time.RFC3339)

	modificationTimeString, err := strconv.ParseInt(value.SystemMetadata.ModificationTime, 10, 64)
	if err != nil {
		diagnostics := diag.Diagnostics{
			diag.NewErrorDiagnostic("failed to convert modification time", err.Error()),
		}
		return subaccountDestinationGenericType{}, diagnostics
	}
	modifyTime := time.UnixMilli(modificationTimeString).UTC().Format(time.RFC3339)

	destination := subaccountDestinationGenericType{
		Etag:         types.StringValue(value.SystemMetadata.Etag),
		SubaccountID: subaccountID,
	}

	tmp := make(map[string]string)
	maps.Copy(tmp, value.DestinationConfiguration)

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

	if serviceInstanceID.ValueString() == "" {
		destination.ServiceInstanceID = types.StringNull()
	} else {
		destination.ServiceInstanceID = serviceInstanceID
	}

	destinationJSON, err := json.Marshal(tmp)
	if err != nil {
		diagnostics := diag.Diagnostics{
			diag.NewErrorDiagnostic("failed to marshal additional configuration", err.Error()),
		}
		return subaccountDestinationGenericType{}, diagnostics
	}
	destination.DestinationConfiguration = jsontypes.NewNormalizedValue(string(destinationJSON))

	var diagnostics diag.Diagnostics

	return destination, diagnostics
}
