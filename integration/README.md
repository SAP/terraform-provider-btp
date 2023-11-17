# About
This directory contains Terraform definitions for setting up a Global Account for integration testing and fixture recording.

# Prerequisite
- Global Account with:
  - Services:
    - SAP HANA Cloud | hana-cloud
  - Entitlements:
    - HANA Cloud: hana (Canary | Quota: 3)
- IDP with Technical User
  - Groups need to be properly configured

# Setup
To setup a global account set the following environment variables:
```sh
BTP_USERNAME=<value>
BTP_PASSWORD=<value>
BTP_IDP=<value> # IDP of technical user
```

Afterwards use Terraform to setup the global account:
```sh
terraform apply
```

