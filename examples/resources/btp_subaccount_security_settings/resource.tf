resource "btp_subaccount_security_settings" "subaccount" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"

  default_identity_provider = "sap.custom"

  access_token_validity  = 3600
  refresh_token_validity = 3600

  treat_users_with_same_email_as_same_user = true

  custom_email_domains = ["yourdomain.test"]
}
