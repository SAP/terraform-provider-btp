---
page_title: "btp_globalaccount_role_collection Resource - terraform-provider-btp"
subcategory: ""
description: |-
  Creates a role collection in a global account.
  Tip:
  You must be assigned to the admin role of the global account.
  Further documentation:
  https://help.sap.com/docs/btp/sap-business-technology-platform/role-collections-and-roles-in-global-accounts-directories-and-subaccounts
---

# btp_globalaccount_role_collection (Resource)

Creates a role collection in a global account.

__Tip:__
You must be assigned to the admin role of the global account.

__Further documentation:__
<https://help.sap.com/docs/btp/sap-business-technology-platform/role-collections-and-roles-in-global-accounts-directories-and-subaccounts>

## Example Usage

```terraform
resource "btp_globalaccount_role_collection" "my_collection" {
  name        = "My own role collection"
  description = "A description of what the role collection is supposed to do."

  roles = [
    {
      name                 = "Global Account Admin"
      role_template_app_id = "cis-central!b13"
      role_template_name   = "GlobalAccount_Admin"
    }
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the role collection.
- `roles` (Attributes Set) (see [below for nested schema](#nestedatt--roles))

### Optional

- `description` (String) The description of the role collection.

### Read-Only

- `id` (String, Deprecated) The ID of the role collection.

<a id="nestedatt--roles"></a>
### Nested Schema for `roles`

Required:

- `name` (String) The name of the referenced role.
- `role_template_app_id` (String) The name of the referenced template app id.
- `role_template_name` (String) The name of the referenced role template.

## Import

Import is supported using the following syntax:

```terraform
# terraform import btp_globalaccount_role_collection.<resource_name> '<name>'

terraform import btp_globalaccount_role_collection.globalaccount_viewer 'Global Account Viewer'

#terraform import using id attribute in import block

import {
  to =  btp_globalaccount_role_collection.<resource_name>
  id = "<name>"
}

# this resource supports import using identity attribute from Terraform version 1.12 or higher

import {
to = btp_globalaccount_role_collection.<resource_name>
identity = {
  name = "<name>"
  }
}
```
