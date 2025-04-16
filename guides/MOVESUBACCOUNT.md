# Limitation - Move of Subaccounts

## Status Quo

In general, SAP BTP offers the option to move subaccounts between directories or from within a directory to a global account. This operation is possible in the [SAP BTP cockpit](https://help.sap.com/docs/btp/sap-business-technology-platform/change-subaccount-details) as well as via the [BTP CLI](https://help.sap.com/docs/btp/btp-cli-command-reference/btp-move-accounts-subaccount). Be aware that the Terraform provider for SAP BTP currently does **not** support this functionality.

In case you change the attribute [`parent_id`](https://registry.terraform.io/providers/SAP/btp/latest/docs/resources/subaccount#parent_id-1) in your resource configuration of [`btp_subaccount`](https://registry.terraform.io/providers/SAP/btp/latest/docs/resources/subaccount) this will trigger a *deletion* of the subaccount and a *recreation* of the subaccount under the new parent. This is not the same as moving a subaccount, as it will cause a deletion and recreation of all resources in the subaccount in accordance to your Terraform configuration.

> [!CAUTION]
> Make sure that you always validate the outcome of the Terraform plan phase to avoid an accidental deletion of a subaccount and the contained resources. In case of CI/CD pipelines you can ensure this by implementing automatic checks on the plan e.g., by leveraging tools like the [Open Policy Agent](https://www.openpolicyagent.org/)

A feature request exists for this functionality, and you can vote for it in the [corresponding issue](https://github.com/SAP/terraform-provider-btp/issues/1020) on GitHub.

## Workaround

If you want to restructure your subaccount landscape and must move existing subaccounts, you can use the following workaround:

1. Move the subaccount to the new parent using either the SAP BTP Cockpit or the BTP CLI.
2. Backup your Terraform state file and your Terraform configuration.
3. Remove the subaccount from the Terraform state file. For details see the Terraform documentation on [removing resources](https://developer.hashicorp.com/terraform/language/resources/syntax#removing-resources) from the state.
4. Adjust the `parent_id` attribute in your resource configuration of `btp_subaccount` to point to the new parent.
5 Import the subaccount back into the Terraform state. For details see the Terraform documentation on [importing resources](https://developer.hashicorp.com/terraform/language/import) into the state.

> [!IMPORTANT]
> Make sure that step 2 is executed successfully and a backup of your Terraform state and your configuration is available.
