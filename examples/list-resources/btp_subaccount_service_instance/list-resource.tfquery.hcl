# This feature requires Terraform v1.14.0 or later (Stable as of 2026)
# List resources must be defined in .tfquery.hcl files.

# Generic template for a list block
list "btp_subaccount_service_instance" "<label_name>" {
  # (Required) Provider instance to use
  provider = provider_name

   config {
    # Provider specific filters
  }
}

# List block to discover all service instaces for given subaccount
# Returns only the resource identities by default.
list "btp_subaccount_service_instance" "all" {
  provider = btp

  # Required
  config {
    subaccount_id = "<subaccount_id>"
  }
}

# List block to discover all service instaces for given subaccount with full resource details
# Setting include_resource = true returns full resource objects (e.g., platform_id, name..)
list "btp_subaccount_service_instance" "with_resource" {
  provider         = btp
  include_resource = true
  config {
  # Required  
  subaccount_id = "<subaccount_id>"
  }
}

list "btp_subaccount_service_instance" "with_fields_filter" {
  provider         = btp
  include_resource = true

  config {
  # Required  
  subaccount_id = "<subaccount_id>"

  # Optional
  fields_filter = "ready eq 'true'"
  }
}

list "btp_subaccount_service_instance" "with_lables_filter" {
  provider         = btp
  include_resource = true

  config {
  # Required  
  subaccount_id = "<subaccount_id>"

  # Optional
  labels_filter = "cred_revision eq '0'"

  }
}
