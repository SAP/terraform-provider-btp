
#subaccount destination with service instance and labels (additional configuration) 
resource "btp_subaccount_destination" "destination" {
  name                = "destination"
  type                = "HTTP"
  proxy_type          = "Internet"
  url                 = "https://myservice.example.com"
  authentication      = "NoAuthentication"
  description         = "resource"
  subaccount_id       = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  service_instance_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  additional_configuration = jsonencode({
    Abc = "good"
  })
}

#subaccount destination creation without service instance and labels (additional configuration) 
resource "btp_subaccount_destination" "destination" {
  name           = "destination"
  type           = "HTTP"
  proxy_type     = "Internet"
  url            = "https://myservice.example.com"
  authentication = "NoAuthentication"
  description    = "resource"
  subaccount_id  = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  additional_configuration = jsonencode({
    Abc = "good"
  })
}

#subaccount destination creation without service instance and without labels (additional configuration) 
resource "btp_subaccount_destination" "destination" {
  name                = "destination"
  type                = "HTTP"
  proxy_type          = "Internet"
  url                 = "https://myservice.example.com"
  authentication      = "NoAuthentication"
  description         = "resource"
  subaccount_id       = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  service_instance_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
}

#subaccount destination creation without service instance and and labels (additional configuration) 
#Note: Auth prpoerties are part of additional configuration
resource "btp_subaccount_destination" "destination" {
  name           = "destination"
  type           = "HTTP"
  proxy_type     = "Internet"
  url            = "https://myservice.example.com"
  authentication = "OAuth2ClientCredentials"
  description    = "resource"
  subaccount_id  = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  additional_configuration = jsonencode({
    ClientID            = "<clientid>"
    ClientSecret        = "<clientsecret>"
    TokenServiceURL     = "https://tokenurl"
    TokenServiceURLType = "dedicated"
  })
}
