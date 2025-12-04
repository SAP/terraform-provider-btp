# Subaccount specific Destination Certificates
data "btp_subaccount_destination_certificates" "all" {
    subaccountd_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"            
}

# Both Subaccount and Service Instance specific Destination Certificates
data "btp_subaccount_destination_certificates" "all" {
    subaccountd_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"    
    service_instance_id = "bc8a216f-1184-49dc-b4b4-17cfe2828965"        
}