---
page_title: "btp_directory_api_credential Resource - terraform-provider-btp"
subcategory: ""
description: |-
  Manage API Credentials at the Directory level. These credentials will enable you to consume
  the REST APIs of the SAP Authorization and Trust Management service (XSUAA).
  With the client ID and client secret, or certificate, you can request an access token for the APIs in the targeted directory.
  Tip:
  You must be assigned to directory admin or viewer role.
  Further documentation:
  https://help.sap.com/docs/btp/sap-business-technology-platform/entitlements-and-quotas
---

# btp_directory_api_credential (Resource)

Manage API Credentials at the Directory level. These credentials will enable you to consume
		the REST APIs of the SAP Authorization and Trust Management service (XSUAA).
		With the client ID and client secret, or certificate, you can request an access token for the APIs in the targeted directory.

__Tip:__
You must be assigned to directory admin or viewer role.

__Further documentation:__
<https://help.sap.com/docs/btp/sap-business-technology-platform/entitlements-and-quotas>

## Example Usage

```terraform
# Create a secret type API credential at the directory level
resource "btp_directory_api_credential" "with-secret" {
  name = "directory-api-credential-with-secret"
  directory_id = "d1298936-ddaf-4a82-b1d7-3ad29a732b61"
  read_only = false
}

// This datasource runs a go script to dynamically generate a PEM certificate which is used in the resource below
data "external" "values" {
  program = ["go","run","../certificate.go"]
}

# Create a certificate type API credential at the directory level
resource "btp_directory_api_credential" "with-certificate" {
  name = "directory-api-credential-with-certificate"
  directory_id = "d1298936-ddaf-4a82-b1d7-3ad29a732b61"
  certificate_passed = data.external.values.result["certificate"]
  read_only = false
}

# Create a secret type, read-only API credential at the directory level
resource "btp_directory_api_credential" "read-only" {
  name = "read-only-directory-api-credential"
  directory_id = "d1298936-ddaf-4a82-b1d7-3ad29a732b61"
  read_only = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `directory_id` (String) The ID of the directory.

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

