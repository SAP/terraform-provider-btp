package cis

const (
	// The CRUD operation or series of operations completed successfully.
	StateOK                          string = "OK"
	StateCanceled                           = "CANCELED"
	StateCreating                           = "CREATING"
	StateCreationFailed                     = "CREATION_FAILED"
	StateDeleting                           = "DELETING"
	StateDeletionFailed                     = "DELETION_FAILED"
	StateMigrating                          = "MIGRATING"
	StateMigrationFailed                    = "MIGRATION_FAILED"
	StateMigrated                           = "MIGRATED"
	StateMoveFailed                         = "MOVE_FAILED"
	StateMoveToOtherGaFailed                = "MOVE_TO_OTHER_GA_FAILED"
	StateMoving                             = "MOVING"
	StateMovingToOtherGa                    = "MOVING_TO_OTHER_GA"
	StatePendingReview                      = "PENDING_REVIEW"
	StateProcessing                         = "PROCESSING"
	StateProcessingFailed                   = "PROCESSING_FAILED"
	StateRollbackMigrationProcessing        = "ROLLBACK_MIGRATION_PROCESSING"
	StateStarted                            = "STARTED"
	StateSuspensionFailed                   = "SUSPENSION_FAILED"
	StateUpdateAccountTypeFailed            = "UPDATE_ACCOUNT_TYPE_FAILED"
	StateUpdateDirectoryTypeFailed          = "UPDATE_DIRECTORY_TYPE_FAILED"
	StateUpdateFailed                       = "UPDATE_FAILED"
	StateUpdating                           = "UPDATING"
)
