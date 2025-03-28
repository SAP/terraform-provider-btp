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

variable "subaccount_subdomain_extension" {
  type        = string
  description = "Subaccount subdomains are required to be unique across landscapes. Define a custom extension to avoid subdomain collision if operating multiple integration Global Accounts in the same landscape."
  default     = "main"
}

