# This feature requires Terraform v1.14.0 or later (Stable as of 2026).
# List resources must be defined in .tfquery.hcl files.

# Generic template for a list block.
list "btp_directory_entitlement" "<label_name>" {
  # (Required) Provider instance to use.
  provider = provider_name

  # (Required) Provider-specific filter.
  config {
    # filter_parameter = value
  }
}

# List block to discover all entitlements available within a specific BTP Directory.
# Returns only the resource identities by default.
list "btp_directory_entitlement" "all" {
  provider = btp

  # Required
  config {
    directory_id = "<directory_id>"
  }
}

# List block to discover all entitlements available within a specific BTP Directory with full resource details.
# Setting include_resource = true returns full resource objects (e.g., plan_name, plan_unique_identifier).
list "btp_directory_entitlement" "with_resource" {
  provider         = btp
  include_resource = true

  # Required
  config {
    directory_id   = "<directory_id>"
  }
}