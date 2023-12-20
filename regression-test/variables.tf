variable "globalaccount" {
  type        = string
  description = "The global account where the project account shall be created in."
}


variable "idp" {
  type        = string
  description = "The IDP to use for authentication."  
}

variable "cli_server" {
  type        = string
  description = "The URL of the CLI server"
  default     = "https://cpcli.cf.eu10.hana.ondemand.com"
}


variable "region" {
  type        = string
  description = "The region where the project account shall be created in."
  default     = "us10"
}
