package provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"net"
	"net/url"
	"strconv"
	"strings"
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

// for extracting name and creation of json string from the destination json configuration
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
	rawName, ok := config["Name"]
	if !ok {
		return "", "", errors.New(`destination_configuration JSON must contain a "Name" property`)
	}
	name, ok = rawName.(string)
	if !ok {
		return "", "", fmt.Errorf(`destination_configuration JSON property "Name" must be a string, got %T`, rawName)
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return "", "", errors.New(`destination_configuration JSON property "Name" must be a non-empty string`)
	}
	jsonBytes, err := json.Marshal(config)
	if err != nil {
		return "", "", err
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

func ValidateFromJSON(jsonStr string) error {
	var m map[string]any

	if err := json.Unmarshal([]byte(jsonStr), &m); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	t, ok := m["Type"].(string)
	if !ok || t == "" {
		return errors.New("missing or invalid 'Type'")
	}

	switch strings.ToUpper(t) {

	case "HTTP":
		u, ok := m["URL"].(string)
		if !ok || u == "" {
			return errors.New("missing 'URL' for HTTP type")
		}
		if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
			return errors.New("HTTP URL must start with http:// or https://")
		}
		if _, err := url.ParseRequestURI(u); err != nil {
			return errors.New("invalid HTTP URL")
		}

	case "LDAP":
		u, ok := m["ldap.url"].(string)
		if !ok || u == "" {
			return errors.New("missing 'ldap.url'")
		}
		if !strings.HasPrefix(u, "ldap://") && !strings.HasPrefix(u, "ldaps://") {
			return errors.New("LDAP URL must start with ldap:// or ldaps://")
		}
		if _, err := url.ParseRequestURI(u); err != nil {
			return errors.New("invalid LDAP URL")
		}

	case "TCP":
		addr, ok := m["Address"].(string)
		if !ok || addr == "" {
			return errors.New("missing 'Address' for TCP type")
		}
		host, port, err := net.SplitHostPort(addr)
		if err != nil || host == "" || port == "" {
			return errors.New("TCP Address must be in host:port format")
		}
		p, err := strconv.Atoi(port)
		if err != nil || p <= 0 || p > 65535 {
			return errors.New("TCP port must be a valid number between 1 and 65535")
		}

	default:
		return nil // intentionally no validation for unknown types
	}
	return nil
}

// This function add the masked fields which are not fetched in read operation
func MergeGenericDestinationConfig(plannedConfig jsontypes.Normalized, responseConfig jsontypes.Normalized) (jsontypes.Normalized, error) {
	if plannedConfig.IsNull() {
		return responseConfig, nil
	}

	if responseConfig.IsNull() {
		return plannedConfig, nil
	}

	plannedMap := make(map[string]string)
	responseMap := make(map[string]string)

	if !plannedConfig.IsUnknown() {
		if err := json.Unmarshal([]byte(plannedConfig.ValueString()), &plannedMap); err != nil {
			return jsontypes.Normalized{}, err
		}
	}

	if !responseConfig.IsUnknown() {
		if err := json.Unmarshal([]byte(responseConfig.ValueString()), &responseMap); err != nil {
			return jsontypes.Normalized{}, err
		}
	}

	for k, plannedVal := range plannedMap {
		responseVal, exists := responseMap[k]
		// Some values are not returned by the API or are masked with <redacted>
		// To achieve consistency between state and plan we take the planned value in those cases
		if !exists || responseVal == "" || responseVal == "<redacted>" {
			responseMap[k] = plannedVal
		}
	}

	mergedJSON, err := json.Marshal(responseMap)
	if err != nil {
		return jsontypes.Normalized{}, err
	}

	return jsontypes.NewNormalizedValue(string(mergedJSON)), nil
}
