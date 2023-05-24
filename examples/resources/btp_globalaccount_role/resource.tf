resource "btp_globalaccount_role" "xsuaa_admin" {
  name               = "My Role"
  role_template_name = "xsuaa_admin"
  app_id             = "xsuaa!t1"
}
