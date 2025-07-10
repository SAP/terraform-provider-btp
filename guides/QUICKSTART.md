# Quick Start Guide

## Introduction

The Terraform provider for SAP BTP enables you to automate the provisioning, management, and configuration of resources on [SAP Business Technology Platform](https://account.hana.ondemand.com/). By leveraging this provider, you can simplify and streamline the deployment and maintenance of BTP services and applications.

## Prerequisites

To follow along with this tutorial, ensure you have access to a [free BTP trial account](https://developers.sap.com/tutorials/hcp-create-trial-account.html) and Terraform installed on your machine. You can download it from the official [Terraform website](https://developer.hashicorp.com/terraform/downloads).

## Authentication

Before running your Terraform script, one have to choose the Authentication mechanism. Terraform Provider for BTP supports username/password based authentication and Single Sign-On (only for desktop environment)

### Username/Password

The SAP BTP provider offers the authentication via username and password. Be aware that this authentication is not compatible with the SAP Universal ID. For details on how to resolve this please see [SAP Note 3085908 - Getting an error (e.g. invalid credentials) in certain applications (e.g. SAP Download Manager) when using S-user ID or SAP Universal ID](https://me.sap.com/notes/3085908).

#### Windows

For Windows you have two options to export the environment variables:

If you use Windows CMD, do the export via the following commands:

```Shell
set BTP_USERNAME=<your_username>
set BTP_PASSWORD=<your_password>
```

If you use Powershell, do the export via the following commands:

```Shell
$Env:BTP_USERNAME = '<your_username>'
$Env:BTP_PASSWORD = '<your_password>'
```

#### Mac

For Mac OS export the environment variables via:

```Shell
export BTP_USERNAME=<your_username>
export BTP_PASSWORD=<your_password>
```

#### Linux

For Linux export the environment variables via:

```Shell
export BTP_USERNAME=<your_username>
export BTP_PASSWORD=<your_password>
```

Replace `<your_username>` and `<your_password>` with your actual BTP username and password.

### SSO

The provider supports login via Single Sign-On (SSO). To enable this you need to set the environment variable BTP_ENABLE_SSO to true. Additionally, ensure that you run your scripts in a desktop environment. It's important to note that the SSO login feature is not intended for use in containerized environments or CI/CD pipelines.

#### Windows

For Windows you have two options to export the environment variables:

If you use Windows CMD, do the export via the following commands:

```Shell
set BTP_ENABLE_SSO=true
```

If you use Powershell, do the export via the following commands:

```Shell
$Env:BTP_ENABLE_SSO = "true"
```

#### Mac

For Mac OS export the environment variables via:

```Shell
export BTP_ENABLE_SSO=true
```

#### Linux

For Linux export the environment variables via:

```Shell
export BTP_ENABLE_SSO=true
```

## Samples

- [Get Started with the Terraform Provider for BTP](https://developers.sap.com/tutorials/btp-terraform-get-started.html)
- More [use case samples](https://github.com/SAP-samples/btp-terraform-samples/tree/main/released/usecases)

## Documentation

Terraform Provider for SAP BTP [Documentation](https://registry.terraform.io/providers/SAP/btp/latest/docs)
