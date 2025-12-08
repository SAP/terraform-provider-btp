# Create a secret type API credential at the globalaccount level
resource "btp_globalaccount_api_credential" "with-secret" {
  name      = "globalaccount-api-credential-with-secret"
  read_only = false
}

# Create a certificate type API credential at the globalaccount level
resource "btp_globalaccount_api_credential" "with-certificate" {
  name               = "globalaccount-api-credential-with-certificate"
  certificate_passed = "-----BEGIN CERTIFICATE-----\n-not-a-valid-certificate-\n-----END CERTIFICATE----\n"
  read_only          = false
}

# Create a secret type, read-only API credential at the globalaccount level
resource "btp_globalaccount_api_credential" "read-only" {
  name      = "read-only-globalaccount-api-credential"
  read_only = true
}