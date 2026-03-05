# This feature requires Terraform v1.14.0 or later (Stable as of 2026)
# List resources must be defined in .tfquery.hcl files.

# Generic template for a list block
list "btp_directory" "<label_name>" {
  # (Required) Provider instance to use
  provider = provider_name

   config {
    # Provider specific filters
  }
}

# List block to discover all directories
# Returns only the resource identities by default.
list "btp_directory" "all" {
  provider = btp
}

# List block to discover all directories with full resource details
# Setting include_resource = true returns full resource objects (e.g., description, name..)
list "btp_directory" "with_resource" {
  provider         = btp
  include_resource = true
}
