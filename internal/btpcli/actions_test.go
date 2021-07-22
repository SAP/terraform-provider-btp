package btpcli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertAction(t *testing.T, expectAction Action, initFn func(command string, args any) *CommandRequest) {
	t.Helper()

	assert.Equal(t, expectAction, initFn("", map[string]string{}).Action)
}

func TestNewAddRequest(t *testing.T) {
	assertAction(t, ActionAdd, NewAddRequest)
}

func TestNewAssignRequest(t *testing.T) {
	assertAction(t, ActionAssign, NewAssignRequest)
}

func TestNewCreateRequest(t *testing.T) {
	assertAction(t, ActionCreate, NewCreateRequest)
}

func TestNewDeleteRequest(t *testing.T) {
	assertAction(t, ActionDelete, NewDeleteRequest)
}

func TestNewDisableRequest(t *testing.T) {
	assertAction(t, ActionDisable, NewDisableRequest)
}

func TestNewEnableRequest(t *testing.T) {
	assertAction(t, ActionEnable, NewEnableRequest)
}

func TestNewGetRequest(t *testing.T) {
	assertAction(t, ActionGet, NewGetRequest)
}

func TestNewListRequest(t *testing.T) {
	assertAction(t, ActionList, NewListRequest)
}

func TestNewRegisterRequest(t *testing.T) {
	assertAction(t, ActionRegister, NewRegisterRequest)
}

func TestNewRemoveRequest(t *testing.T) {
	assertAction(t, ActionRemove, NewRemoveRequest)
}

func TestNewShareRequest(t *testing.T) {
	assertAction(t, ActionShare, NewShareRequest)
}

func TestNewSubscribeRequest(t *testing.T) {
	assertAction(t, ActionSubscribe, NewSubscribeRequest)
}

func TestNewUnassignRequest(t *testing.T) {
	assertAction(t, ActionUnassign, NewUnassignRequest)
}

func TestNewUnregisterRequest(t *testing.T) {
	assertAction(t, ActionUnregister, NewUnregisterRequest)
}

func TestNewUnshareRequest(t *testing.T) {
	assertAction(t, ActionUnshare, NewUnshareRequest)
}

func TestNewUnsubscribeRequest(t *testing.T) {
	assertAction(t, ActionUnsubscribe, NewUnsubscribeRequest)
}

func TestNewUpdateRequest(t *testing.T) {
	assertAction(t, ActionUpdate, NewUpdateRequest)
}
