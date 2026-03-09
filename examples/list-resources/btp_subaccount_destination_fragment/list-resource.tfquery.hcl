# This feature requires Terraform v1.14.0 or later (Stable as of 2026)
# List resources must be defined in .tfquery.hcl files.

# Generic template for a list block
list "btp_subaccount_destination_fragment" "<label_name>" {
  # (Required) Provider instance to use
  provider = provider_name

   config {
    # Provider specific filters
  }
}

# List block to discover all destination fragments for given subaccount
# Returns only the resource identities by default.
list "btp_subaccount_destination_fragment" "all" {
  provider = btp

  # Required
  config {
    subaccount_id = "<subaccount_id>"
  }
}

# List block to discover all destination fragments for given subaccount with full resource details
# Setting include_resource = true returns full resource objects (e.g., destination_configuration, name..)
list "btp_subaccount_destination_fragment" "with_resource" {
  provider         = btp
  include_resource = true
  config {
  # Required  
  subaccount_id = "<subaccount_id>"
  }
}


# List block to discover all destination fragments for given subaccount and service instance id
# Returns only the resource identities by default.
list "btp_subaccount_destination_fragment" "by_service_instance" {
  provider = btp

  config {
    # Required
    subaccount_id       = "<subaccount_id>"

    # Optional
    service_instance_id = "<service_instance_id>"
  }
}
