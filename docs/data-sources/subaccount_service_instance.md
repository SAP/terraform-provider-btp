---
page_title: "btp_subaccount_service_instance Data Source - terraform-provider-btp"
subcategory: ""
description: |-
  Gets details about a specific provisioned service instance, such as its name, id,  platform to which it belongs, and the last operation performed.
  Tip:
  You must be assigned to the admin or viewer role of the subaccount.
---

# btp_subaccount_service_instance (Data Source)

Gets details about a specific provisioned service instance, such as its name, id,  platform to which it belongs, and the last operation performed.

__Tip:__
You must be assigned to the admin or viewer role of the subaccount.

## Example Usage

```terraform
# look up a service instance by its ID and subaccount ID
data "btp_subaccount_service_instance" "by_id" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  id            = "bc8a216f-1184-49dc-b4b4-17cfe2828965"
}

# look up a service instance by its name and subaccount ID
data "btp_subaccount_service_instance" "by_name" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  name          = "my-xsuaa-application"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `subaccount_id` (String) The ID of the subaccount.

### Optional

- `id` (String) The ID of the service instance.
- `name` (String) The name of the service instance.

### Read-Only

- `context` (String) Contextual data for the resource.
- `created_date` (String) The date and time when the resource was created in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.
- `dashboard_url` (String) The URL of the web-based management UI for the service instance.
- `labels` (Map of Set of String) The set of words or phrases assigned to the service instance.
- `last_modified` (String) The date and time when the resource was last modified in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.
- `parameters` (String) The configuration parameters for the service instance.
- `platform_id` (String) The platform ID.
- `ready` (Boolean)
- `referenced_instance_id` (String) The ID of the instance to which the service instance refers.
- `serviceplan_id` (String) The ID of the service plan.
- `shared` (Boolean) Shows whether the service instance is shared.
- `state` (String) The current state of the service instance.
- `usable` (Boolean) Shows whether the resource can be used.