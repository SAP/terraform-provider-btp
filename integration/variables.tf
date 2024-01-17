variable "region" {
  type        = string
  description = "The region where the project account shall be created in."
  default     = "eu12"
}

variable "subaccount_subdomain_extension" {
  type        = string
  description = "Subaccount subdomains are required to be unique across landscapes. Define a custom extension to avoid subdomain collision if operating multiple integration Global Accounts in the same landscape."
  default     = "main"
}

variable "testing_idp" {
  description = "The IDP used for testing. Contains test users and should not be used for other purposes. URL must not contain a protocol prefix. Must not be part of trusted_idps."
  type        = string
  default = "terraformint.accounts400.ondemand.com"
}

variable "trusted_idp_origin_keys" {
  description = "List of IDP origin keys added to created resources. IDPs must contain group 'BTP Terraform Administrator' and group 'BTP Terraform Developer'. Enables members of listed groups to access created resources with corresponding scope."
  type        = list(string)
  default = [
    "terraform-platform"
  ]
}

