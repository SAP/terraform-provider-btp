# About
This directory contains Terraform definitions for setting up a Global Account for integration testing and fixture recording.

# Prerequisite
- Global Account with:
  - Use For: Testing
  - Services:
    - SAP HANA Cloud | hana-cloud
    - Data Privacy Integration | data-privacy-integration-service
    - Alert Notification | alert-notification
    - Audit Log Service | auditlog
    - Malware Scanning Service | malware-scanner
  - Entitlements:
    - HANA Cloud: hana (Canary | Quota: 3)
    - Data Privacy Integration | data-privacy-integration-service: standard (Quota: 3)
    - Alert Notification | alert-notification: lite (Quota: 1), free (Quota: 1)
    - Audit Log Service | auditlog: standard (Quota: 1)
    - Malware Scanning Service | malware-scanner: clamav (Quota: 1)
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

