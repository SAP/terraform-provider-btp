# Drift Detection

## Overview

In general, Terraform enables you to provision and manage your Infrastructure as Code. The management part also comprises the ability to detect and reconcile configuration drifts in your infrastructure.

The mechanism to detect drifts in your infrastructure is provided by the `terraform plan` command that compares the current state of the infrastructure with the Terraform state. You find the details of the `terraform plan` command in the [official Terraform documentation](https://developer.hashicorp.com/terraform/cli/commands/plan). Technically you check for a configuration drift by executing the `terraform plan` with the option `-detailed-exitcode`. This will return an exit code of `2` if there are changes to be applied, `1` if there is an error, and `0` if the infrastructure matches the Terraform state.

In this document we discuss the drift detection for the Terraform provider for SAP BTP and what needs to be considered.

## Prerequisites

A drift can only be detected if the resources on SAP BTP have either been *created* by Terraform or if they have been *imported* into the Terraform state. Any resources that have been created manually and are not reflected in the Terraform state are "invisible" to Terraform and a drift can consequently not be detected.

Consequently a drift will only show up for changes in the resource configuration as defined by the corresponding Terraform resource or deletion of resources that have been created by Terraform.

## Resource Overview

From a technical perspective the drift detection requires the ability to compare the current state of the resources on SAP BTP with the Terraform state. This is achieved by the Terraform provider for SAP BTP by querying the platform APIs for the current state of the resources. Unfortunately, not all resources on SAP BTP support this i.e., the query of the current state of the resource on the platform is either not supported by the platform APIs at all or it does not return the full set of parameters.

The following overview list des resources and their support for drift detection (state: 09.04.2024):

| Resource                                     | Drift Detection Support | Comments                                                                                                                                  |
|---                                           |---                      |---                                                                                                                                        |
| btp_directory                                | Yes                     | -                                                                                                                                         |
| btp_directory_api_credential                 | No                      | -                                                                                                                                         |
| btp_directory_entitlement                    | Yes                     | -                                                                                                                                         |
| btp_directory_role                           | Yes                     | -                                                                                                                                         |
| btp_directory_role_collection                | Yes                     | -                                                                                                                                         |
| btp_directory_role_collection_assignment     | No                      | READ capability of resource not available. Improvement planned for H2/2024 see [issue](https://github.com/SAP/terraform-provider-btp/issues/748) |
| btp_globalaccount_api_credential             | No                      | -                                                                                                                                         |
| btp_globalaccount_resource_provider          | Yes                     | -                                                                                                                                         |
| btp_globalaccount_role                       | Yes                     | -                                                                                                                                         |
| btp_globalaccount_role_collection            | Yes                     | -                                                                                                                                         |
| btp_globalaccount_role_collection_assignment | No                      | READ capability of resource not available. Improvement planned for H2/2024 see [issue](https://github.com/SAP/terraform-provider-btp/issues/748) |
| btp_globalaccount_security_settings          | Yes                     | -                                                                                                                                         |
| btp_globalaccount_trust_configuration        | Yes                     | -                                                                                                                                         |
| btp_subaccount                               | Yes                     | -                                                                                                                                         |
| btp_subaccount_api_credential                | No                      | -                                                                                                                                         |
| btp_subaccount_entitlement                   | Yes                     | -                                                                                                                                         |
| btp_subaccount_environment_instance          | Yes                     | -                                                                                                                                         |
| btp_subaccount_role                          | Yes                     | -                                                                                                                                         |
| btp_subaccount_role_collection               | Yes                     | -                                                                                                                                         |
| btp_subaccount_role_collection_assignment    | No                      | READ capability of resource not available. Improvement planned for H2/2024 see [issue](https://github.com/SAP/terraform-provider-btp/issues/748) |
| btp_subaccount_security_settings             | Yes                     | -                                                                                                                                         |
| btp_subaccount_service_binding               | Yes                     | -                                                                                                                                         |
| btp_subaccount_service_broker                | Yes                      | -                                                                                                                                         |
| btp_subaccount_service_instance              | Yes with restrictions   | The parameters defined via `parameters` are not tracked due to missing READ functionality depending on the service offering configuration |
| btp_subaccount_subscription                  | Yes                     | -                                                                                                                                         |
| btp_subaccount_trust_configuration           | Yes                     | -                                                                                                                                         |

## Further options

Besides the `terraform plan` command there are further options to detect drifts in your infrastructure. You can also create custom checks by leveraging the data sources of the Terraform provider and combine the results with custom logic e.g., in a CI/CD pipeline. The concrete setup depends on your requirements and no generic solution can be provided.

## Next Steps

After a configuration drift has been detected you must analyze the changes and decide how to proceed. In general you have two options:

- You can either reconcile the *infrastructure setup* by applying the changes via the `terraform apply` command. This will apply the change to the platform so that the infrastructure matches your Terraform state again.
- You can adjust the *Terraform state* without applying changes to the infrastructure. The process to sync the state with the infrastructure on the platform is described in the [official Terraform documentation](https://developer.hashicorp.com/terraform/tutorials/state/refresh) leveraging the `-refresh-only` mode for `terraform plan` and `terraform apply`.
