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