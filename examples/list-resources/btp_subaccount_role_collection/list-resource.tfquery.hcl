# This feature requires Terraform v1.14.0 or later (Stable as of 2026).
# List resources must be defined in .tfquery.hcl files.

# Generic template for a list block.
list "btp_subaccount_role_collection" "<label_name>" {
  # (Required) Provider instance to use.
  provider = provider_name

  # (Required) Provider-specific filter.
  config {
    # filter_parameter = value
  }
}

# List block to discover all role collections available within a specific BTP subaccount.
# Returns only the resource identities by default.
list "btp_subaccount_role_collection" "all" {
  provider = btp

  # Required
  config {
    subaccount_id = "<subaccount_id>"
  }
}

# List block to discover all role collections available within a specific BTP subaccount with full resource details.
# Setting include_resource = true returns full resource objects (e.g., description, name).
list "btp_subaccount_role_collection" "with_resource" {
  provider         = btp
  include_resource = true

  # Required
  config {
    subaccount_id   = "<subaccount_id>"
  }
}
