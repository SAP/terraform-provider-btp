package provider

import (
	"testing"
)

func TestExtractLabelValue_Success_CfApiUrl(t *testing.T) {
	label := `{"API Endpoint":"https://api.cf.example.com","Org Name":"example","Org ID":"8d818824-394a-abcd-0815-7a3c8ce93e57"}`
	val, err := ExtractLabelValue(label, EnvironmentLabelKeyCfApiUrl)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "https://api.cf.example.com"
	if val != want {
		t.Fatalf("got %q, want %q", val, want)
	}
}

func TestExtractLabelValue_Success_CfOrgId(t *testing.T) {
	label := `{"API Endpoint":"https://api.cf.example.com","Org Name":"example","Org ID":"8d818824-394a-abcd-0815-7a3c8ce93e57"}`
	val, err := ExtractLabelValue(label, EnvironmentLabelKeyCfOrgId)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "8d818824-394a-abcd-0815-7a3c8ce93e57"
	if val != want {
		t.Fatalf("got %q, want %q", val, want)
	}
}

func TestExtractLabelValue_Success_KymaApiServerUrl(t *testing.T) {
	label := `{"APIServerURL":"https://api.kyma.example.com","KubeconfigURL":"https://kyma.example.com/kubeconfig/ABC", "Name":"test-terraform-kyma"}`
	val, err := ExtractLabelValue(label, EnvironmentLabelKeyKymaApiServerUrl)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "https://api.kyma.example.com"
	if val != want {
		t.Fatalf("got %q, want %q", val, want)
	}
}

func TestExtractLabelValue_Success_KymaKubeconfigUrl(t *testing.T) {
	label := `{"APIServerURL":"https://api.kyma.example.com","KubeconfigURL":"https://kyma.example.com/kubeconfig/ABC"}`
	val, err := ExtractLabelValue(label, EnvironmentLabelKeyKymaKubeconfigUrl)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "https://kyma.example.com/kubeconfig/ABC"
	if val != want {
		t.Fatalf("got %q, want %q", val, want)
	}
}

func TestExtractLabelValue_Error_EmptyLabel(t *testing.T) {
	_, err := ExtractLabelValue("", EnvironmentLabelKeyCfApiUrl)
	if err == nil {
		t.Fatalf("expected error for empty label, got nil")
	}
}

func TestExtractLabelValue_Error_InvalidJSON(t *testing.T) {
	_, err := ExtractLabelValue("{invalid}", EnvironmentLabelKeyCfApiUrl)
	if err == nil {
		t.Fatalf("expected error for invalid JSON, got nil")
	}
}

func TestExtractLabelValue_Error_NoJSON(t *testing.T) {
	_, err := ExtractLabelValue("invalid", EnvironmentLabelKeyCfApiUrl)
	if err == nil {
		t.Fatalf("expected error for invalid JSON, got nil")
	}
}
func TestExtractLabelValue_Error_MissingKey(t *testing.T) {
	label := `{"Org Name":"example"}`
	_, err := ExtractLabelValue(label, EnvironmentLabelKeyCfApiUrl)
	if err == nil {
		t.Fatalf("expected error for missing key, got nil")
	}
}

func TestExtractLabelValue_Error_NonStringValue(t *testing.T) {
	label := `{"API Endpoint":123}`
	_, err := ExtractLabelValue(label, EnvironmentLabelKeyCfApiUrl)
	if err == nil {
		t.Fatalf("expected error for non-string value, got nil")
	}
}

func TestExtractLabelValue_Error_UnsupportedKey(t *testing.T) {
	label := `{"API Endpoint":"https://api.cf.example.com"}`
	_, err := ExtractLabelValue(label, EnvironmentLabelKey("UnknownKey"))
	if err == nil {
		t.Fatalf("expected error for unsupported key, got nil")
	}
}
