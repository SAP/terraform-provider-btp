---
page_title: "btp_subaccount_roles Data Source - terraform-provider-btp"
subcategory: ""
description: |-
  Gets all roles.
  Tip:
  You must be assigned to the admin or viewer role of the subaccount.
  Further documentation:
  https://help.sap.com/docs/btp/sap-business-technology-platform/role-collections-and-roles-in-global-accounts-directories-and-subaccounts
---

# btp_subaccount_roles (Data Source)

Gets all roles.

__Tip:__
You must be assigned to the admin or viewer role of the subaccount.

__Further documentation:__
<https://help.sap.com/docs/btp/sap-business-technology-platform/role-collections-and-roles-in-global-accounts-directories-and-subaccounts>

## Example Usage

```terraform
data "btp_subaccount_roles" "all" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `subaccount_id` (String) The ID of the subaccount.

### Read-Only

- `id` (String, Deprecated) The ID of the subaccount.
- `values` (Attributes List) (see [below for nested schema](#nestedatt--values))

<a id="nestedatt--values"></a>
### Nested Schema for `values`

Read-Only:

- `app_id` (String) The id of the application that provides the role template and the role.
- `app_name` (String) The name of the application that provides the role template and the role.
- `attribute_list` (Attributes List) The attributes assigned to this role. (see [below for nested schema](#nestedatt--values--attribute_list))
- `description` (String) The description of the role.
- `name` (String) The name of the role.
- `read_only` (Boolean) Shows whether the role can be modified or not.
- `role_template_name` (String) The name of the role template.
- `scopes` (Attributes List) The scopes available with this role. (see [below for nested schema](#nestedatt--values--scopes))

<a id="nestedatt--values--attribute_list"></a>
### Nested Schema for `values.attribute_list`

Read-Only:

- `attribute_name` (String) The name of the role attribute.
- `attribute_value_origin` (String) The origin of the attribute value.
- `attribute_values` (Set of String)
- `value_required` (Boolean) Shows whether the value is required.


<a id="nestedatt--values--scopes"></a>
### Nested Schema for `values.scopes`

Read-Only:

- `custom_grant_as_authority_to_apps` (Set of String)
- `custom_granted_apps` (Set of String)
- `description` (String) The description of the scope.
- `grant_as_authority_to_apps` (Set of String)
- `granted_apps` (Set of String)
- `name` (String) The name of the scope.