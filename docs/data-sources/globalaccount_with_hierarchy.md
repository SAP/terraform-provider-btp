---
page_title: "btp_globalaccount_with_hierarchy Data Source - terraform-provider-btp"
subcategory: ""
description:   Gets details about a global account's hierarchy structure
  Tip:
  You must be assigned to the admin or viewer role of the global account.
  Further documentation:
---

# btp_globalaccount_with_hierarchy (Data Source)

Gets details about a global account's hierarchy structure

__Tip:__
You must be assigned to the admin or viewer role of the global account.

__Further documentation:__
<https://help.sap.com/docs/btp/sap-business-technology-platform/account-model>

<br>

## Example Usage

```terraform
data "btp_globalaccount_with_hierarchy" "this"{}
```
<br>

## Schema

### Read Only

- `created_date` (String) The date and time when the resource was created in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.
- `id` (String) The ID of the global account.
- `last_modified` (String) The date and time when the resource was last modified in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.
- `name` (String) The display name of the global account.
- `state` (String) The current state of the global account. See [State Values](#state-values) below for the possible values. 
- `subdomain` (String) The subdomain is part of the path used to access the authorization tenant of the global account.
- `type` (String) The type of the resource, in this case this will have value 'Global Account'.

### Optional
- `directories` (Attributes List) The directories contained in the global account (see [Directories](#directories) below for the schema) 
- `subaccounts` (Attributes List) The subaccounts contained in the globalaccount (see [Subaccounts](#subaccounts) below for the schema)

<br>

<h2><a id="directories">Directories</a></h2>

> **Note:** The directory object is recursive as directories can hold other directories within itself. Currently upto 5 levels of directories are supported.

### Read Only

- `created_by` (String) The details of the user that created the directory.
- `created_date` (String) The date and time when the directory was created in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.
- `features` (Set of String) The features that are enabled for the directory. Possible values are: 

  | value | description | 
  | --- | --- | 
  | `DEFAULT`  | All directories have the following basic feature enabled:<br> 1. Group and filter subaccounts for reports and filters <br> 2. Monitor usage and costs on a directory level (costs only available for contracts that use the consumption-based commercial model)<br> 3. Set custom properties and tags to the directory for identification and reporting purposes. | 
  | `ENTITLEMENTS` | Allows the assignment of a quota for services and applications to the directory from the global account quota for distribution to the subaccounts under this directory. | 
  | `AUTHORIZATIONS` | Allows the assignment of users as administrators or viewers of this directory. You must apply this feature in combination with the `ENTITLEMENTS` feature. |
- `id` (String) The ID of the directory.
- `last_modified` (String) The date and time when the directory was last modified in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.
- `name` (String) The display name of the directory.
- `parent_id` (String) The ID of the directory's parent entity. Typically this is the global account.
- `parent_name` (String) The name of the directory's parent entity. Typically this is the global account.
- `parent_type` (String) The type of the directory's parent entity. Typically this will have the value 'Global Account'.
- `state` (String) The current state of the directory. See [State Values](#state-values) below for the possible values. 
- `type` (String) The type of resource, in this case will have the value 'Directory'.

### Optional

- `directories` (Attributes List) The list of directories contained in this directory. (see [Directories](#directories) above for the schema)
- `subaccounts` (Attributes List) The subaccounts contained in this directory. (see [Subaccounts](#subaccounts) below for the schema)
- `subdomain` (String) This applies only to directories that have the user authorization management feature enabled. The subdomain is part of the path used to access the authorization tenant of the directory.

<br>

<h2><a id="subaccounts">Subaccounts</a></h2>

### Read Only

- `created_by` (String) The details of the user that created the subaccount.
- `created_date` (String) The date and time when the subaccount was created in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.
- `id` (String) The ID of the subaccount.
- `last_modified` (String) The date and time when the subaccount was last modified in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.
- `name` (String) A descriptive name of the subaccount for customer-facing UIs.
- `parent_id` (String) The ID of the subaccount's parent entity. If the subaccount is located directly in the global account (not in a directory), then this is the ID of the global account.
- `parent_name` (String) The name of the subaccount's parent entity. If the subaccount is located directly in the global account (not in a directory), then this is the name of the global account.
- `parent_type` (String) The type of the subaccount's parent entity. If the subaccount is located directly in the global account (not in a directory), then this will have the value 'Global Account' .
- `region` (String) The region in which the subaccount was created.
- `state` (String) The current state of the subaccount. See [State Values](#state-values) below for the possible values. 
- `subdomain` (String) The subdomain that becomes part of the path used to access the authorization tenant of the subaccount. Must be unique within the defined region. Use only letters (a-z), digits (0-9), and hyphens (not at the start or end). Maximum length is 63 characters. Cannot be changed after the subaccount has been created.
- `type` (String) The type of resource, in this case this will be 'Subaccount'.

<br>

<h2><a id="state-values">State Values </a></h2>

Below are the possible values for the state Attribute:
  | state | description | 
  | --- | --- | 
  | `OK` | The CRUD operation or series of operations completed successfully. | 
  | `STARTED` | CRUD operation on an entity has started. | 
  | `CANCELLED` | The operation or processing was canceled by the operator. | 
  | `PROCESSING` | A series of operations related to the entity is in progress. | 
  | `PROCESSING_FAILED` | The processing operations failed. | 
  | `CREATING` | Creating entity operation is in progress. | 
  | `CREATION_FAILED` | The creation operation failed, and the entity was not created or was created but cannot be used. | 
  | `UPDATING` | Updating entity operation is in progress. | 
  | `UPDATE_FAILED` | The update operation failed, and the entity was not updated. | 
  | `DELETING` | Deleting entity operation is in progress. | 
  | `DELETION_FAILED` | The delete operation failed, and the entity was not deleted. | 
  | `MOVING` | Moving entity operation is in progress. | 
  | `MOVE_FAILED` | Entity could not be moved to a different location. | 
  | `PENDING REVIEW` | The processing operation has been stopped for reviewing and can be restarted by the operator. | 
  | `MIGRATING` | Migrating entity from Neo to Cloud Foundry. |
