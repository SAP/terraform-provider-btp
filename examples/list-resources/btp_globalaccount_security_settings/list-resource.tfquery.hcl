# This feature requires Terraform v1.14.0 or later (Stable as of 2026)
# List resources must be defined in .tfquery.hcl files.

# Generic template for a list block
list "btp_globalaccount_security_settings" "<label_name>" {
  # (Required) Provider instance to use
  provider = provider_name
}

# List block to discover all security settings for a global account
# Returns only the resource identities by default.
list "btp_globalaccount_security_settings" "all" {
  provider = btp
}

# List block to discover all security settings for a global account with full resource details
# Setting include_resource = true returns full resource objects (e.g., access_token_validity, refresh_token_validity..)
list "btp_globalaccount_security_settings" "with_resource" {
  provider         = btp
  include_resource = true
}
