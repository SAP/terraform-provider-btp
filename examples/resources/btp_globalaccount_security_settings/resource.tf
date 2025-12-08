resource "btp_globalaccount_security_settings" "this" {
  default_identity_provider = "sap.custom"

  access_token_validity  = 3600
  refresh_token_validity = 3600

  treat_users_with_same_email_as_same_user = true

  custom_email_domains = ["yourdomain.test"]

  iframe_domains = ["https://yourdomain.test"]
}
