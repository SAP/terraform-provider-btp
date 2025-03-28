---
page_title: "btp_subaccount_api_credential Resource - terraform-provider-btp"
subcategory: ""
description: |-
  Manage API Credentials at the Subaccount level. These credentials will enable you to consume
  the REST APIs of the SAP Authorization and Trust Management service (XSUAA).
  With the client ID and client secret, or certificate, you can request an access token for the APIs in the targeted subaccount.
  Tip:
  You must be assigned to the subaccount admin or viewer role.
  Further documentation:
  https://help.sap.com/docs/btp/sap-business-technology-platform/managing-api-credentials-for-calling-rest-apis-of-sap-authorization-and-trust-management-service
---

# btp_subaccount_api_credential (Resource)

Manage API Credentials at the Subaccount level. These credentials will enable you to consume
		the REST APIs of the SAP Authorization and Trust Management service (XSUAA).
		With the client ID and client secret, or certificate, you can request an access token for the APIs in the targeted subaccount.

__Tip:__
You must be assigned to the subaccount admin or viewer role.

__Further documentation:__
<https://help.sap.com/docs/btp/sap-business-technology-platform/managing-api-credentials-for-calling-rest-apis-of-sap-authorization-and-trust-management-service>

## Example Usage

```terraform
# Create a secret type API credential at the subaccount level
resource "btp_subaccount_api_credential" "with-secret" {
  name = "subaccount-api-credential-with-secret"
  subaccount_id = "77395f6a-a601-4c9e-8cd0-c1fcefc7f60f"
  read_only = false
}


# Create a certificate type API credential at the subaccount level
resource "btp_subaccount_api_credential" "with-certificate" {
  name = "subaccount-api-credential-with-certificate"
  subaccount_id = "77395f6a-a601-4c9e-8cd0-c1fcefc7f60f"
  certificate_passed = "-----BEGIN CERTIFICATE-----\n-not-a-valid-certificate-\n-----END CERTIFICATE----\n"
  read_only = false
}

# Create a secret type, read-only API credential at the subaccount level
resource "btp_subaccount_api_credential" "read-only" {
  name = "read-only-subaccount-api-credential"
  subaccount_id = "77395f6a-a601-4c9e-8cd0-c1fcefc7f60f"
  read_only = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `subaccount_id` (String) The ID of the subaccount.

### Optional

- `certificate_passed` (String) If the user prefers to use a certificate, they must provide the certificate value in PEM format "----BEGIN CERTIFICATE-----...-----END CERTIFICATE-----".
- `name` (String) The name for the API credential.
- `read_only` (Boolean) Access restriction placed on the API credential. If set to true, the resource has only read-only access.

### Read-Only

- `api_url` (String) The URL to be used to make the API calls.
- `certificate_received` (String) The certificate that is computed based on the one passed by the user.
- `client_id` (String) A unique ID associated with the API credential.
- `client_secret` (String) If the certificate is omitted, then a unique secret is generated for the API credential.
- `credential_type` (String) The supported credential types are Secrets (Default) or Certificates.
- `key` (String) RSA key generated if the API credential is created with a certificate.
- `token_url` (String) The URL to be used to fetch the access token to make use of the XSUAA REST APIs.


