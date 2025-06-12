resource "btp_subaccount_role" "xsuaa_auditor" {
  subaccount_id      = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  name               = "XSUAA Auditor"
  role_template_name = "xsuaa_auditor"
  app_id             = "xsuaa!t1"
}

#Example with custom attribute
resource "btp_subaccount_role" "custom" {
  subaccount_id      = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  name               = "Custom Role with Attributes"
  role_template_name = "Application_Frontend_Developer"
  app_id             = "eu12-appfront!b390135"
  attribute_list = [
    {
      attribute_name         = "namespace"
      attribute_value_origin = "saml"
      attribute_values       = ["custom value"]
    },
  ]
}
