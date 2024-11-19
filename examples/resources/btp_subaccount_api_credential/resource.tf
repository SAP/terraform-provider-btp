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