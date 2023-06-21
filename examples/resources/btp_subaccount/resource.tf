# create a subaccount in eu30 region (GCP)
resource "btp_subaccount" "my_project" {
  name      = "My Project"
  subdomain = "my-project"
  region    = "eu30"
}

# create a subaccount in every Azure region which has support for cloud foundry

# Look up all regions via data source
data "btp_regions" "all" {}

# create the subaccounts by iterating over the regions with the defined constraints
resource "btp_subaccount" "my_project_on_azure" {
  for_each = { for dc in data.btp_regions.all.values : dc.region => dc if dc.environment == "cloudfoundry" && dc.iaas_provider == "AZURE" }

  name      = "My CF Project (${each.key})"
  subdomain = "my-project-${each.key}"
  region    = each.key
}
