variable "globalaccount" {
  type        = string
  description = "The global account where the project account shall be created in."
}

variable "cli_server" {
  type        = string
  description = "The URL of the CLI server"
}


variable "region" {
  type        = string
  description = "The region where the project account shall be created in."
}

variable "testing_idp" {
  description = "The IDP used for testing. Contains test users and should not be used for other purposes. URL must not contain a protocol prefix. Must not be part of trusted_idps."
  type        = string
}

variable "subaccount_subdomain_extension" {
  type        = string
  description = "Subaccount subdomains are required to be unique across landscapes. Define a custom extension to avoid subdomain collision if operating multiple integration Global Accounts in the same landscape."
}

variable "globalaccount_role_template_app_id" {
  type        = string
  description = "The application identifier (app ID) of the role template at the Global Account level. This value is landscape-specific and is used to assign roles/authorizations within the Global Account."
}

variable "directory_role_template_app_id" {
  type        = string
  description = "The application identifier (app ID) of the role template at the Directory level. This value varies across landscapes and is required for assigning directory-level roles and permissions."
}

variable "subaccount_role_template_app_id" {
  type        = string
  description = "The application identifier (app ID) of the role template at the Subaccount level. This is landscape-dependent and is used to manage role assignments and access within a Subaccount."
}
