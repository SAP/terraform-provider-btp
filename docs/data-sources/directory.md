---
page_title: "btp_directory Data Source - terraform-provider-btp"
subcategory: ""
description: |-
  Gets the details about a directory.
  Tip:
  You must be assigned to the global account admin role, or the directory admin if the directory is configured to manage its authorizations.
  Further documentation:
  https://help.sap.com/docs/btp/sap-business-technology-platform/account-model
---

# btp_directory (Data Source)

Gets the details about a directory.

__Tip:__
You must be assigned to the global account admin role, or the directory admin if the directory is configured to manage its authorizations.

__Further documentation:__
<https://help.sap.com/docs/btp/sap-business-technology-platform/account-model>

## Example Usage

```terraform
data "btp_directory" "by_id" {
  id = "dd005d8b-1fee-4e6b-b6ff-cb9a197b7fe0"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) The ID of the directory.

### Read-Only

- `created_by` (String) The details of the user that created the directory.
- `created_date` (String) The date and time when the resource was created in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.
- `description` (String) The description of the directory.
- `features` (Set of String) The features that are enabled for the directory. Possible values are: 

  | value | description | 
  | --- | --- | 
  | `DEFAULT`  | All directories have the following basic feature enabled:<br> 1. Group and filter subaccounts for reports and filters <br> 2. Monitor usage and costs on a directory level (costs only available for contracts that use the consumption-based commercial model)<br> 3. Set custom properties and tags to the directory for identification and reporting purposes. | 
  | `ENTITLEMENTS` | Allows the assignment of a quota for services and applications to the directory from the global account quota for distribution to the subaccounts under this directory. | 
  | `AUTHORIZATIONS` | Allows the assignment of users as administrators or viewers of this directory. You must apply this feature in combination with the `ENTITLEMENTS` feature. |
- `labels` (Map of Set of String) The set of words or phrases assigned to the directory.
- `last_modified` (String) The date and time when the resource was last modified in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.
- `name` (String) The display name of the directory.
- `parent_id` (String) The ID of the directory's parent entity. Typically this is the global account.
- `state` (String) The current state of the directory. Possible values are: 

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
- `subdomain` (String) This applies only to directories that have the user authorization management feature enabled. The subdomain is part of the path used to access the authorization tenant of the directory.