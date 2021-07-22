resource "btp_directory_role" "dirrole" {
  directory_id       = "ddfc2206-5f11-48ed-a1ec-29010af70050"
  name               = "DirUsageRepViewTest"
  role_template_name = "Directory_Usage_Reporting_Viewer"
  app_id             = "uas!b36585"
}
