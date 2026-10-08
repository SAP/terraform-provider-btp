package provider

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestResolveServiceMetadataCacheTTL(t *testing.T) {
	cases := []struct {
		name, env string
		value     types.String
		expected  time.Duration
		ok        bool
		warnings  int
	}{
		{"default", "", types.StringNull(), 5 * time.Minute, true, 0},
		{"environment", "30s", types.StringNull(), 30 * time.Second, true, 0},
		{"explicit wins", "30s", types.StringValue("1m"), time.Minute, true, 1},
		{"disabled", "", types.StringValue("0s"), 0, true, 0},
		{"negative", "", types.StringValue("-1s"), 0, false, 0},
		{"invalid", "invalid", types.StringNull(), 0, false, 0},
		{"unknown", "", types.StringUnknown(), 0, false, 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("BTP_SERVICE_METADATA_CACHE_TTL", tc.env)
			response := &provider.ConfigureResponse{}
			ttl, ok := resolveServiceMetadataCacheTTL(tc.value, response)
			require.Equal(t, tc.ok, ok)
			require.Equal(t, tc.expected, ttl)
			require.Equal(t, tc.warnings, response.Diagnostics.WarningsCount())
			if !ok && !tc.value.IsUnknown() {
				require.True(t, response.Diagnostics.HasError())
			}
		})
	}
}
