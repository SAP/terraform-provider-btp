# Create a secret type API credential at the globalaccount level
resource "btp_globalaccount_api_credential" "with-secret" {
  name = "globalaccount-api-credential-with-secret"
  read_only = false
}

// This datasource runs a go script to dynamically generate a PEM certificate which is used in the resource below
data "external" "values" {
  program = ["go","run","../certificate.go"]
}

# Create a certificate type API credential at the globalaccount level
resource "btp_globalaccount_api_credential" "with-certificate" {
  name = "globalaccount-api-credential-with-certificate"
  certificate_passed = data.external.values.result["certificate"]
  read_only = false
}

# Create a secret type, read-only API credential at the globalaccount level
resource "btp_globalaccount_api_credential" "read-only" {
  name = "read-only-globalaccount-api-credential"
  read_only = true
}