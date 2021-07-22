package btpcli

// TODO generate

const (
	ActionAdd         Action = "add"
	ActionAssign             = "assign"
	ActionCreate             = "create"
	ActionDelete             = "delete"
	ActionDisable            = "disable"
	ActionEnable             = "enable"
	ActionGet                = "get"
	ActionList               = "list"
	ActionRegister           = "register"
	ActionRemove             = "remove"
	ActionShare              = "share"
	ActionSubscribe          = "subscribe"
	ActionUnassign           = "unassign"
	ActionUnregister         = "unregister"
	ActionUnshare            = "unshare"
	ActionUnsubscribe        = "unsubscribe"
	ActionUpdate             = "update"
)

// NewAddRequest creates a new add request
func NewAddRequest(command string, args any) *CommandRequest {
	return NewCommandRequest(ActionAdd, command, args)
}

// NewAssignRequest creates a new assign request
func NewAssignRequest(command string, args any) *CommandRequest {
	return NewCommandRequest(ActionAssign, command, args)
}

// NewCreateRequest creates a new create request
func NewCreateRequest(command string, args any) *CommandRequest {
	return NewCommandRequest(ActionCreate, command, args)
}

// NewDeleteRequest creates a new delete request
func NewDeleteRequest(command string, args any) *CommandRequest {
	return NewCommandRequest(ActionDelete, command, args)
}

// NewDisableRequest creates a new disable request
func NewDisableRequest(command string, args any) *CommandRequest {
	return NewCommandRequest(ActionDisable, command, args)
}

// NewEnableRequest creates a new enable request
func NewEnableRequest(command string, args any) *CommandRequest {
	return NewCommandRequest(ActionEnable, command, args)
}

// NewGetRequest creates a new get request
func NewGetRequest(command string, args any) *CommandRequest {
	return NewCommandRequest(ActionGet, command, args)
}

// NewListRequest creates a new list request
func NewListRequest(command string, args any) *CommandRequest {
	return NewCommandRequest(ActionList, command, args)
}

// NewRegisterRequest creates a new register request
func NewRegisterRequest(command string, args any) *CommandRequest {
	return NewCommandRequest(ActionRegister, command, args)
}

// NewRemoveRequest creates a new remove request
func NewRemoveRequest(command string, args any) *CommandRequest {
	return NewCommandRequest(ActionRemove, command, args)
}

// NewShareRequest creates a new share request
func NewShareRequest(command string, args any) *CommandRequest {
	return NewCommandRequest(ActionShare, command, args)
}

// NewSubscribeRequest creates a new subscribe request
func NewSubscribeRequest(command string, args any) *CommandRequest {
	return NewCommandRequest(ActionSubscribe, command, args)
}

// NewUnassignRequest creates a new unassign request
func NewUnassignRequest(command string, args any) *CommandRequest {
	return NewCommandRequest(ActionUnassign, command, args)
}

// NewUnregisterRequest creates a new unregister request
func NewUnregisterRequest(command string, args any) *CommandRequest {
	return NewCommandRequest(ActionUnregister, command, args)
}

// NewUnshareRequest creates a new unshare request
func NewUnshareRequest(command string, args any) *CommandRequest {
	return NewCommandRequest(ActionUnshare, command, args)
}

// NewUnsubscribeRequest creates a new unsubscribe request
func NewUnsubscribeRequest(command string, args any) *CommandRequest {
	return NewCommandRequest(ActionUnsubscribe, command, args)
}

// NewUpdateRequest creates a new update request
func NewUpdateRequest(command string, args any) *CommandRequest {
	return NewCommandRequest(ActionUpdate, command, args)
}
