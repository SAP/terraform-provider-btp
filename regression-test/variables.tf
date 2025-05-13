variable "globalaccount" {
  type        = string
  description = "The global account where the project account shall be created in."
}

variable "cli_server" {
  type        = string
  description = "The URL of the CLI server"
  default     = "https://cli.btp.cloud.sap"
}


variable "region" {
  type        = string
  description = "The region where the project account shall be created in."
  default     = "us10"
}

variable "testing_idp" {
  description = "The IDP used for testing. Contains test users and should not be used for other purposes. URL must not contain a protocol prefix. Must not be part of trusted_idps."
  type        = string
  default     = "terraformeds2.accounts.ondemand.com"
}

variable "subaccount_subdomain_extension" {
  type        = string
  description = "Subaccount subdomains are required to be unique across landscapes. Define a custom extension to avoid subdomain collision if operating multiple integration Global Accounts in the same landscape."
  default     = "main"
}

