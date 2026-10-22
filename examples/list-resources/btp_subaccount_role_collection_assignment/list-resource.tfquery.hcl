# This feature requires Terraform v1.14.0 or later.
# List resources must be defined in .tfquery.hcl files.

# Discover assignments from all role collections and origins in a subaccount.
list "btp_subaccount_role_collection_assignment" "all" {
  provider = btp

  config {
    subaccount_id = "<subaccount_id>"
  }
}

# Discover assignments from one role collection and origin.
list "btp_subaccount_role_collection_assignment" "by_role_collection" {
  provider = btp

  config {
    subaccount_id        = "<subaccount_id>"
    role_collection_name = "<role_collection_name>"
    origin               = "<origin>" # Optional; omit or set to an empty string for all origins.
  }
}

# Return complete managed resource values in addition to resource identities.
list "btp_subaccount_role_collection_assignment" "with_resource" {
  provider         = btp
  include_resource = true

  config {
    subaccount_id        = "<subaccount_id>"
    origin               = "<origin>"
    role_collection_name = "<role_collection_name>"
  }
}
