# This feature requires Terraform v1.14.0 or later (Stable as of 2026)
# List resources must be defined in .tfquery.hcl files.

# Generic template for a list block
list "btp_subaccount" "<label_name>" {
  # (Required) Provider instance to use
  provider = provider_name

   config {
    # Provider specific filters
  }
}

# List block to discover all subaccounts
# Returns only the resource identities by default.
list "btp_subaccount" "all" {
  provider = btp
}

# List block to discover all subaccounts with full resource details
# Setting include_resource = true returns full resource objects (e.g., platform_id, name..)
list "btp_subaccount" "with_resource" {
  provider         = btp
  include_resource = true
}


# List block to discover all subaccounts within a region.
list "btp_subaccount" "with_region_filter" {
  provider         = btp
  config {
  # Optional  
  region = "eu12"

  }
}

list "btp_subaccount" "with_region_filter" {
  provider         = btp
  include_resource = true

  config {
  # Optional
  labels_filter = "my-label=my-value"
  }
}
