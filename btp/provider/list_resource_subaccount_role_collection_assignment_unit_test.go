package provider

import (
	"context"
	"testing"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/xsuaa_authz"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/xsuaa_trust"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestSubaccountRoleCollectionAssignmentListResourceConfigure(t *testing.T) {
	t.Parallel()

	r := NewSubaccountRoleCollectionAssignmentListResource().(list.ListResourceWithConfigure)
	resp := &resource.ConfigureResponse{}
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: struct{}{}}, resp)

	assert.True(t, resp.Diagnostics.HasError())
}

func TestSubaccountRoleCollectionAssignmentUserEntries(t *testing.T) {
	t.Parallel()

	users := []xsuaa_authz.UserReference{
		{Username: "api-user-name", Email: "user@example.com", Origin: "sap.default"},
		{Username: "other-user", Email: "other@example.com", Origin: "other-origin"},
	}

	t.Run("filtered", func(t *testing.T) {
		entries := subaccountRoleCollectionAssignmentUserEntries(users, "Subaccount Viewer", "ldap")
		if !assert.Len(t, entries, 1) {
			return
		}
		assert.Equal(t, "api-user-name", entries[0].username)
		assert.Equal(t, "ldap", entries[0].origin)
	})

	t.Run("all origins", func(t *testing.T) {
		entries := subaccountRoleCollectionAssignmentUserEntries(users, "Subaccount Viewer", "")
		if !assert.Len(t, entries, 2) {
			return
		}
		assert.Equal(t, "sap.default", entries[0].origin)
		assert.Equal(t, "other-origin", entries[1].origin)
	})
}

func TestSubaccountRoleCollectionAssignmentSAMLEntries(t *testing.T) {
	t.Parallel()

	assignments := []xsuaa_authz.SamlAttrAssignment{
		{AttributeName: "Groups", AttributeValue: "viewer-group", IdentityProvider: "custom-origin"},
		{AttributeName: "cost_center", AttributeValue: "1234", IdentityProvider: "A display name", SamlEntityId: "https://other.accounts.example"},
	}
	resolver := subaccountRoleCollectionAssignmentOriginResolver{
		trustConfigurations: xsuaa_trust.TrustConfigurationResponseCollectionObject{
			{OriginKey: "custom-origin", IdentityProvider: "custom.accounts.example"},
			{OriginKey: "other-origin", IdentityProvider: "other.accounts.example"},
		},
	}

	entries, err := subaccountRoleCollectionAssignmentSAMLEntries(assignments, "Subaccount Viewer", resolver)

	assert.NoError(t, err)
	if !assert.Len(t, entries, 2) {
		return
	}
	assert.Equal(t, "custom-origin", entries[0].origin)
	assert.Equal(t, "viewer-group", entries[0].groupName)
	assert.Equal(t, "other-origin", entries[1].origin)
	assert.Equal(t, "cost_center", entries[1].attributeName)
	assert.Equal(t, "1234", entries[1].attributeValue)
}

func TestSubaccountRoleCollectionAssignmentSAMLOriginResolutionError(t *testing.T) {
	t.Parallel()

	assignment := xsuaa_authz.SamlAttrAssignment{
		AttributeName:    "Groups",
		AttributeValue:   "viewer-group",
		IdentityProvider: "A display name",
		SamlEntityId:     "https://shared.accounts.example",
	}
	resolver := subaccountRoleCollectionAssignmentOriginResolver{
		trustConfigurations: xsuaa_trust.TrustConfigurationResponseCollectionObject{
			{OriginKey: "first-origin", IdentityProvider: "shared.accounts.example"},
			{OriginKey: "second-origin", IdentityProvider: "shared.accounts.example"},
		},
	}

	_, _, err := resolver.resolveSAMLOrigin(assignment)
	assert.ErrorContains(t, err, "multiple trust configurations")
}

func TestSubaccountRoleCollectionAssignmentListResourceValues(t *testing.T) {
	t.Parallel()

	subaccountID := types.StringValue("subaccount-id")
	tests := map[string]struct {
		entry      subaccountRoleCollectionAssignmentListEntry
		expectedID string
	}{
		"user": {
			entry:      subaccountRoleCollectionAssignmentListEntry{roleCollectionName: "Viewer", origin: "user-origin", username: "api-user"},
			expectedID: "subaccount-id,Viewer,api-user,user-origin",
		},
		"group": {
			entry:      subaccountRoleCollectionAssignmentListEntry{roleCollectionName: "Viewer", origin: "group-origin", groupName: "viewer-group"},
			expectedID: "subaccount-id,Viewer,group:viewer-group,group-origin",
		},
		"attribute": {
			entry: subaccountRoleCollectionAssignmentListEntry{
				roleCollectionName: "Viewer",
				origin:             "attribute-origin",
				attributeName:      "cost_center",
				attributeValue:     "1234",
			},
			expectedID: "subaccount-id,Viewer,attribute:cost_center/1234,attribute-origin",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			value := subaccountRoleCollectionAssignmentListResourceValue(subaccountID, test.entry)
			assert.Equal(t, test.expectedID, value.Id.ValueString())
			assert.Equal(t, test.entry.origin, value.Origin.ValueString())
		})
	}
}
