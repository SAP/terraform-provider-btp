# Import

## Overview

In general Terraform supports the *import* of resources into the Terraform state. You find the official documentation on how to achieve this [here](https://developer.hashicorp.com/terraform/cli/import).

The Terraform provider for SAP BTP supports the import of resources as well. [The documentation](https://registry.terraform.io/providers/SAP/btp/latest/docs) of the Terraform provider for SAP BTP provides the necessary information on how to import a resource and which keys to use on the level of each resource.

To get a quick overview of the resources and if they support the import functionality, you can refer to the [Resource Overview](#resource-overview) section in this document.

## Resource Overview

The following list provides an overview of the resources and their support for the import functionality (state: 01.01.2025)

| Resource                                     | Import Support
|---                                           |---
| btp_directory                                | Yes
| btp_directory_api_credential                 | No
| btp_directory_entitlement                    | Yes
| btp_directory_role                           | Yes
| btp_directory_role_collection                | Yes
| btp_directory_role_collection_assignment     | No
| btp_globalaccount_api_credential             | No
| btp_globalaccount_resource_provider          | Yes
| btp_globalaccount_role                       | Yes
| btp_globalaccount_role_collection            | Yes
| btp_globalaccount_role_collection_assignment | No
| btp_globalaccount_security_settings          | Yes
| btp_globalaccount_trust_configuration        | Yes
| btp_subaccount                               | Yes with restrictions (see [documentation](https://registry.terraform.io/providers/SAP/btp/latest/docs/resources/subaccount#restriction))
| btp_subaccount_api_credential                | No
| btp_subaccount_entitlement                   | Yes
| btp_subaccount_environment_instance          | Yes
| btp_subaccount_role                          | Yes
| btp_subaccount_role_collection               | Yes
| btp_subaccount_role_collection_assignment    | No
| btp_subaccount_security_settings             | Yes
| btp_subaccount_service_binding               | Yes
| btp_subaccount_service_broker                | Yes
| btp_subaccount_service_instance              | Yes with restrictions (see [documentation](https://registry.terraform.io/providers/SAP/btp/latest/docs/resources/subaccount_service_instance#restriction))
| btp_subaccount_subscription                  | Yes
| btp_subaccount_trust_configuration           | Yes
