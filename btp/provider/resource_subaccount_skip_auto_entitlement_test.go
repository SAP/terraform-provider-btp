package provider

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/cis"
)

func TestSubaccountSkipAutoEntitlementPlanModifier(t *testing.T) {
	t.Parallel()

	var schemaResp resource.SchemaResponse
	(&subaccountResource{}).Schema(t.Context(), resource.SchemaRequest{}, &schemaResp)
	attribute := schemaResp.Schema.Attributes["skip_auto_entitlement"].(schema.BoolAttribute)
	testSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{"skip_auto_entitlement": attribute},
	}

	testCases := []struct {
		name            string
		state           types.Bool
		plan            types.Bool
		create          bool
		destroy         bool
		requiresReplace bool
	}{
		{name: "unset to false", state: types.BoolNull(), plan: types.BoolValue(false)},
		{name: "false to unset", state: types.BoolValue(false), plan: types.BoolNull()},
		{name: "unset to true", state: types.BoolNull(), plan: types.BoolValue(true), requiresReplace: true},
		{name: "true to unset", state: types.BoolValue(true), plan: types.BoolNull(), requiresReplace: true},
		{name: "false to true", state: types.BoolValue(false), plan: types.BoolValue(true), requiresReplace: true},
		{name: "true to false", state: types.BoolValue(true), plan: types.BoolValue(false), requiresReplace: true},
		{name: "unchanged unset", state: types.BoolNull(), plan: types.BoolNull()},
		{name: "unchanged false", state: types.BoolValue(false), plan: types.BoolValue(false)},
		{name: "unchanged true", state: types.BoolValue(true), plan: types.BoolValue(true)},
		{name: "unset to unknown", state: types.BoolNull(), plan: types.BoolUnknown(), requiresReplace: true},
		{name: "false to unknown", state: types.BoolValue(false), plan: types.BoolUnknown(), requiresReplace: true},
		{name: "true to unknown", state: types.BoolValue(true), plan: types.BoolUnknown(), requiresReplace: true},
		{name: "unknown to false", state: types.BoolUnknown(), plan: types.BoolValue(false), requiresReplace: true},
		{name: "create true", state: types.BoolNull(), plan: types.BoolValue(true), create: true},
		{name: "create unknown", state: types.BoolNull(), plan: types.BoolUnknown(), create: true},
		{name: "destroy true", state: types.BoolValue(true), plan: types.BoolNull(), destroy: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			rawValue := func(value types.Bool, absent bool) tftypes.Value {
				if absent {
					return tftypes.NewValue(testSchema.Type().TerraformType(ctx), nil)
				}
				boolValue, err := value.ToTerraformValue(ctx)
				assert.NoError(t, err)
				return tftypes.NewValue(testSchema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"skip_auto_entitlement": boolValue,
				})
			}
			req := planmodifier.BoolRequest{
				Path:       path.Root("skip_auto_entitlement"),
				State:      tfsdk.State{Schema: testSchema, Raw: rawValue(tc.state, tc.create)},
				Plan:       tfsdk.Plan{Schema: testSchema, Raw: rawValue(tc.plan, tc.destroy)},
				StateValue: tc.state,
				PlanValue:  tc.plan,
			}
			resp := planmodifier.BoolResponse{PlanValue: tc.plan}
			assert.Len(t, attribute.PlanModifiers, 1)
			attribute.PlanModifiers[0].PlanModifyBool(ctx, req, &resp)

			assert.Empty(t, resp.Diagnostics)
			assert.Equal(t, tc.requiresReplace, resp.RequiresReplace)
			assert.Equal(t, tc.plan, resp.PlanValue)
		})
	}
}

func TestSubaccountSkipAutoEntitlementUpdateState(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name  string
		state types.Bool
		plan  types.Bool
	}{
		{name: "unset to false", state: types.BoolNull(), plan: types.BoolValue(false)},
		{name: "false to unset", state: types.BoolValue(false), plan: types.BoolNull()},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			transport := subaccountEntitlementTestTransport(func(req *http.Request) (*http.Response, error) {
				assert.True(t, strings.HasSuffix(req.URL.Path, "/accounts/subaccount"))
				assert.Contains(t, []string{"update", "get"}, req.URL.RawQuery)
				var payload struct {
					ParamValues map[string]string `json:"paramValues"`
				}
				assert.NoError(t, json.NewDecoder(req.Body).Decode(&payload))
				assert.NotContains(t, payload.ParamValues, "skipAutoAssignPlans")
				return &http.Response{
					StatusCode: http.StatusOK,
					Header: http.Header{
						btpcli.HeaderCLIBackendStatus: []string{"200"},
					},
					Body: io.NopCloser(strings.NewReader(`{"guid":"test-subaccount","state":"OK","displayName":"test-subaccount"}`)),
				}, nil
			})
			client := btpcli.NewV2ClientWithHttpClient(
				&http.Client{Transport: transport},
				&url.URL{Scheme: "https", Host: "btp.invalid"},
				&btpcli.RetryConfig{Enabled: false},
			)
			rs := &subaccountResource{cli: btpcli.NewClientFacade(client)}
			var schemaResp resource.SchemaResponse
			rs.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
			var identitySchemaResp resource.IdentitySchemaResponse
			rs.IdentitySchema(ctx, resource.IdentitySchemaRequest{}, &identitySchemaResp)

			data, diags := subaccountValueFrom(ctx, cis.SubaccountResponseObject{Guid: "test-subaccount", State: cis.StateOK})
			assert.Empty(t, diags)
			state := tfsdk.State{Schema: schemaResp.Schema}
			data.SkipAutoEntitlement = tc.state
			assert.Empty(t, state.Set(ctx, &data))
			plan := tfsdk.Plan{Schema: schemaResp.Schema}
			data.SkipAutoEntitlement = tc.plan
			assert.Empty(t, plan.Set(ctx, &data))
			resp := resource.UpdateResponse{
				State:    tfsdk.State{Schema: schemaResp.Schema, Raw: plan.Raw},
				Identity: &tfsdk.ResourceIdentity{Schema: identitySchemaResp.IdentitySchema},
			}

			rs.Update(ctx, resource.UpdateRequest{Plan: plan, State: state}, &resp)
			assert.Empty(t, resp.Diagnostics)
			var actual types.Bool
			assert.Empty(t, resp.State.GetAttribute(ctx, path.Root("skip_auto_entitlement"), &actual))
			assert.Equal(t, tc.plan, actual)

			readResp := resource.ReadResponse{State: resp.State, Identity: resp.Identity}
			rs.Read(ctx, resource.ReadRequest{State: resp.State, Identity: resp.Identity}, &readResp)
			assert.Empty(t, readResp.Diagnostics)
			assert.Empty(t, readResp.State.GetAttribute(ctx, path.Root("skip_auto_entitlement"), &actual))
			assert.Equal(t, tc.plan, actual)
		})
	}
}

type subaccountEntitlementTestTransport func(*http.Request) (*http.Response, error)

func (f subaccountEntitlementTestTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
