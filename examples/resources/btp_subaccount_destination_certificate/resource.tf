# Subaccount specific destination certificate
resource "btp_subaccount_destination_certificate" "dest_cert_for_sa" {
  subaccount_id       = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  certificate_name    = "test.pem" // Ensure to include a certificate file extension, refer documentation for the valid extensions
  certificate_content = filebase64("path_to_file")
}

# Service instance specific destination certificate
resource "btp_subaccount_destination_certificate" "dest_cert_for_si" {
  subaccount_id       = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  service_instance_id = "bc8a216f-1184-49dc-b4b4-17cfe2828965"
  certificate_name    = "test.pem" // Ensure to include a certificate file extension, refer documentation for the valid extensions
  certificate_content = filebase64("path_to_file")
}
