package saas_manager_service

const (
	StateInProcess              string = "IN_PROCESS"
	StateNotSubscribed          string = "NOT_SUBSCRIBED"
	StateSubscribed             string = "SUBSCRIBED"
	StateSubscribeFailed        string = "SUBSCRIBE_FAILED"
	StateUnsubscribeFailed      string = "UNSUBSCRIBE_FAILED"
	StateUpdateFailed           string = "UPDATE_FAILED"
	StateUpdateParametersFailed string = "UPDATE_PARAMETERS_FAILED"
)
