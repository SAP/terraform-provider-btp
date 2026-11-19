package provider

import (
	"testing"
)

func TestValidateKymaKubeconfigURLSyntax_Allowed(t *testing.T) {
	t.Parallel()
	allowed := []string{
		"https://kyma-env-broker.cp.kyma.cloud.sap/kubeconfig/ABC-123",
		"https://kyma-env-broker.cp.stage.kyma.cloud.sap/kubeconfig/ABC-123",
		"https://kyma-env-broker.cp.dev.kyma.cloud.sap/kubeconfig/ABC-123",
		"https://any-host.example.com/kubeconfig/ABC",
	}
	for _, u := range allowed {
		u := u
		t.Run(u, func(t *testing.T) {
			t.Parallel()
			if err := validateKymaKubeconfigURLSyntax(u); err != nil {
				t.Errorf("expected URL %q to pass, got: %v", u, err)
			}
		})
	}
}

func TestValidateKymaKubeconfigURLSyntax_BlockedScheme(t *testing.T) {
	t.Parallel()
	blocked := []string{
		"http://kyma-env-broker.cp.kyma.cloud.sap/kubeconfig/ABC",
		"ftp://kyma-env-broker.cp.kyma.cloud.sap/kubeconfig/ABC",
		"file:///etc/passwd",
		"",
	}
	for _, u := range blocked {
		u := u
		t.Run(u, func(t *testing.T) {
			t.Parallel()
			if err := validateKymaKubeconfigURLSyntax(u); err == nil {
				t.Errorf("expected URL %q to fail scheme check, but it passed", u)
			}
		})
	}
}

func TestDenyPrivateIPs(t *testing.T) {
	t.Parallel()
	blocked := []string{
		"127.0.0.1",
		"10.0.0.1",
		"172.16.0.1",
		"192.168.1.1",
		"169.254.169.254", // cloud metadata (AWS/GCP/Azure)
	}
	for _, ip := range blocked {
		ip := ip
		t.Run(ip, func(t *testing.T) {
			t.Parallel()
			if err := denyPrivateIPs(ip); err == nil {
				t.Errorf("expected IP %q to be denied, but denyPrivateIPs allowed it", ip)
			}
		})
	}
}
