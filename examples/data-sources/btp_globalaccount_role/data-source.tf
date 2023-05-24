data "btp_globalaccount_role" "user_role_auditor" {
  name               = "User and Role Auditor"
  role_template_name = "xsuaa_auditor"
  app_id             = "xsuaa!t1"
}
