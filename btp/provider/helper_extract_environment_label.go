package provider

import (
	"encoding/json"
	"errors"
)

type EnvironmentLabelKey string

const (
	EnvironmentLabelKeyCfApiUrl          EnvironmentLabelKey = "API Endpoint"
	EnvironmentLabelKeyCfOrgId           EnvironmentLabelKey = "Org ID"
	EnvironmentLabelKeyKymaApiServerUrl  EnvironmentLabelKey = "APIServerURL"
	EnvironmentLabelKeyKymaKubeconfigUrl EnvironmentLabelKey = "KubeconfigURL"
)

func (k EnvironmentLabelKey) String() string { return string(k) }

// Validate that the provided key is supported
func isValidEnvironmentLabelKey(k EnvironmentLabelKey) bool {
	switch k {
	case EnvironmentLabelKeyCfApiUrl,
		EnvironmentLabelKeyCfOrgId,
		EnvironmentLabelKeyKymaApiServerUrl,
		EnvironmentLabelKeyKymaKubeconfigUrl:
		return true
	default:
		return false
	}
}

func ExtractLabelValue(label string, key EnvironmentLabelKey) (string, error) {

	if label == "" {
		return "", errors.New("label is empty")
	}

	if key.String() == "" {
		return "", errors.New("key is empty")
	}

	if !isValidEnvironmentLabelKey(key) {
		return "", errors.New("unsupported key: '" + key.String() + "'")
	}

	var data map[string]any
	if err := json.Unmarshal([]byte(label), &data); err != nil {
		return "", errors.New("failed to parse label JSON: " + err.Error())
	}

	val, ok := data[key.String()]
	if !ok {
		return "", errors.New("label does not contain '" + key.String() + "'")
	}

	strVal, ok := val.(string)
	if !ok || strVal == "" {
		return "", errors.New("the value for '" + key.String() + "' is missing or not a string")
	}

	return strVal, nil
}
