# This feature requires Terraform v1.14.0 or later (Stable as of 2026)
# List resources must be defined in .tfquery.hcl files.

# Generic template for a list block
list "btp_globalaccount_role" "<label_name>" {
  # (Required) Provider instance to use
  provider = provider_name
}

# List block to discover all roles in global account
# Returns only the resource identities (names) by default.
list "btp_globalaccount_role" "all" {
  provider = btp
}

# List block to discover all roles in global account with full resource details
# Setting include_resource = true returns full resource objects (e.g., description, read_only)
list "btp_globalaccount_role" "with_resource" {
  provider         = btp
  include_resource = true
}
