package servicemanager

import (
	"encoding/json"
	"testing"
)

// The Service Manager returns instance parameters as a flat, top-level JSON
// object (verified against a live cloud-logging instance, issue #278):
//
//	{"backend":{"api_enabled":true,"max_data_nodes":2},"ingest_otlp":{"enabled":true}}
//
// ServiceInstanceParametersPlain must deserialize that whole body into
// Parameters. Previously the field carried `json:"-"`, so it was never
// populated and instance parameters silently vanished.
func TestServiceInstanceParametersPlain_UnmarshalFlatBody(t *testing.T) {
	raw := `{"backend":{"api_enabled":true,"max_data_nodes":2},"ingest_otlp":{"enabled":true}}`

	var p ServiceInstanceParametersPlain
	if err := json.Unmarshal([]byte(raw), &p); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if len(p.Parameters) == 0 {
		t.Fatalf("expected parameters to be populated from flat body, got empty map: %#v", p.Parameters)
	}

	ingest, ok := p.Parameters["ingest_otlp"].(map[string]any)
	if !ok {
		t.Fatalf("expected ingest_otlp object in parameters, got: %#v", p.Parameters["ingest_otlp"])
	}
	if ingest["enabled"] != true {
		t.Errorf("expected ingest_otlp.enabled=true, got %#v", ingest["enabled"])
	}

	backend, ok := p.Parameters["backend"].(map[string]any)
	if !ok {
		t.Fatalf("expected backend object in parameters, got: %#v", p.Parameters["backend"])
	}
	if backend["api_enabled"] != true {
		t.Errorf("expected backend.api_enabled=true, got %#v", backend["api_enabled"])
	}
}

// Instances without parameters must keep working: an empty JSON object, a JSON
// null, and an empty body (io.EOF at the caller) must all deserialize cleanly
// to an empty/nil map — never an error, never a panic. doGet() then treats
// len(Parameters)==0 as "no parameters", exactly as before the fix.
func TestServiceInstanceParametersPlain_UnmarshalWithoutParameters(t *testing.T) {
	cases := map[string]string{
		"empty object": `{}`,
		"json null":    `null`,
	}
	for name, raw := range cases {
		t.Run(name, func(t *testing.T) {
			var p ServiceInstanceParametersPlain
			if err := json.Unmarshal([]byte(raw), &p); err != nil {
				t.Fatalf("unmarshal of %q failed: %v", raw, err)
			}
			if len(p.Parameters) != 0 {
				t.Errorf("expected empty parameters for %q, got: %#v", raw, p.Parameters)
			}
		})
	}
}
