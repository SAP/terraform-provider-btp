data "btp_directory_role" "user_role_auditor" {
  directory_id       = "dd005d8b-1fee-4e6b-b6ff-cb9a197b7fe0"
  name               = "User and Role Auditor"
  role_template_name = "xsuaa_auditor"
  app_id             = "xsuaa!t1"
}
