# This feature requires Terraform v1.14.0 or later (Stable as of 2026)
# List resources must be defined in .tfquery.hcl files.

# Generic template for a list block
list "btp_directory_role" "<label_name>" {
  # (Required) Provider instance to use
  provider = provider_name

   config {
    # Provider specific filters
  }
}

# List block to discover all directory roles for given directory
# Returns only the resource identities by default.
list "btp_directory_role" "all" {
  provider = btp

  # Required
  config {
    directory_id = "<directory_id>"
  }
}

# List block to discover all directory roles with full resource details
# Setting include_resource = true returns full resource objects (e.g., platform_id, name..)
list "btp_directory_role" "with_resource" {
  provider         = btp
  include_resource = true
  config {
  # Required  
  directory_id = "<directory_id>"
  }
}
