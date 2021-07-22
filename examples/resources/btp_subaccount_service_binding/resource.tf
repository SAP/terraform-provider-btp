# create a service binding in a subaccount
resource "btp_subaccount_service_binding" "my_binding" {
  subaccount_id       = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  service_instance_id = "8911491d-0e1d-425d-a233-785512602d6f"
  name                = "my binding"
}

# create a parameterized service binding in a subaccount
resource "btp_subaccount_service_binding" "my_parameterized_binding" {
  subaccount_id       = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  service_instance_id = "8911491d-0e1d-425d-a233-785512602d6f"
  name                = "my parameterized binding"
  parameters = jsonencode({
    param_a = ""
    param_b = ""
  })
}
