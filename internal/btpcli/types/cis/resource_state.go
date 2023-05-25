package cis

const (
	// The CRUD operation or series of operations completed successfully.
	StateOK                          string = "OK"
	StateCanceled                    string = "CANCELED"
	StateCreating                    string = "CREATING"
	StateCreationFailed              string = "CREATION_FAILED"
	StateDeleting                    string = "DELETING"
	StateDeletionFailed              string = "DELETION_FAILED"
	StateMigrating                   string = "MIGRATING"
	StateMigrationFailed             string = "MIGRATION_FAILED"
	StateMigrated                    string = "MIGRATED"
	StateMoveFailed                  string = "MOVE_FAILED"
	StateMoveToOtherGaFailed         string = "MOVE_TO_OTHER_GA_FAILED"
	StateMoving                      string = "MOVING"
	StateMovingToOtherGa             string = "MOVING_TO_OTHER_GA"
	StatePendingReview               string = "PENDING_REVIEW"
	StateProcessing                  string = "PROCESSING"
	StateProcessingFailed            string = "PROCESSING_FAILED"
	StateRollbackMigrationProcessing string = "ROLLBACK_MIGRATION_PROCESSING"
	StateStarted                     string = "STARTED"
	StateSuspensionFailed            string = "SUSPENSION_FAILED"
	StateUpdateAccountTypeFailed     string = "UPDATE_ACCOUNT_TYPE_FAILED"
	StateUpdateDirectoryTypeFailed   string = "UPDATE_DIRECTORY_TYPE_FAILED"
	StateUpdateFailed                string = "UPDATE_FAILED"
	StateUpdating                    string = "UPDATING"
)
