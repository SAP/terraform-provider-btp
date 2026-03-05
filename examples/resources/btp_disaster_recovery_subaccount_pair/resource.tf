# Create a subaccount pair for SAP BTP Disaster Recovery
resource "btp_disaster_recovery_subaccount_pair" "dr_pair" {
  subaccount_id        = "dd005d8b-1fee-4e6b-b6ff-cb9a197b7fe0"
  paired_subaccount_id = "cc116e9c-3xdd-5f6b-c7gg-db9a197b7ff1"
}
