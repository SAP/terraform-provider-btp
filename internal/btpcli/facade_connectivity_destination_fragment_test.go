package btpcli

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectivityDestinationFragmentFacade_GetBySubaccount(t *testing.T) {
	command := "connectivity/destination-fragment"
	subaccountId := "sub-123"
	frag := "my-frag"

	var srvCalled bool
	uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srvCalled = true

		assertCall(t, r, command, ActionGet, map[string]string{
			"name":       frag,
			"subaccount": subaccountId,
		})
	}))
	defer srv.Close()

	_, res, err := uut.Connectivity.DestinationFragment.GetBySubaccount(context.TODO(), subaccountId, frag)

	if assert.True(t, srvCalled) && assert.NoError(t, err) {
		assert.Equal(t, 200, res.StatusCode)
	}
}

func TestConnectivityDestinationFragmentFacade_GetByServiceInstance(t *testing.T) {
	command := "connectivity/destination-fragment"
	subaccountId := "sub-123"
	serviceInstanceId := "inst-001"
	frag := "frag1"

	var srvCalled bool
	uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srvCalled = true

		assertCall(t, r, command, ActionGet, map[string]string{
			"name":            frag,
			"subaccount":      subaccountId,
			"serviceInstance": serviceInstanceId,
		})
	}))
	defer srv.Close()

	_, res, err := uut.Connectivity.DestinationFragment.GetByServiceInstance(context.TODO(), subaccountId, frag, serviceInstanceId)

	if assert.True(t, srvCalled) && assert.NoError(t, err) {
		assert.Equal(t, 200, res.StatusCode)
	}
}

func TestConnectivityDestinationFragmentFacade_ListBySubaccount(t *testing.T) {
	command := "connectivity/destination-fragment"
	subaccountId := "sub-123"

	var srvCalled bool
	uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srvCalled = true

		assertCall(t, r, command, ActionList, map[string]string{
			"subaccount": subaccountId,
		})
	}))
	defer srv.Close()

	_, res, err := uut.Connectivity.DestinationFragment.ListBySubaccount(context.TODO(), subaccountId)

	if assert.True(t, srvCalled) && assert.NoError(t, err) {
		assert.Equal(t, 200, res.StatusCode)
	}
}

func TestConnectivityDestinationFragmentFacade_ListByServiceInstance(t *testing.T) {
	command := "connectivity/destination-fragment"
	subaccountId := "sub-123"
	serviceInstanceId := "inst-001"

	var srvCalled bool
	uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srvCalled = true

		assertCall(t, r, command, ActionList, map[string]string{
			"subaccount":      subaccountId,
			"serviceInstance": serviceInstanceId,
		})
	}))
	defer srv.Close()

	_, res, err := uut.Connectivity.DestinationFragment.ListByServiceInstance(context.TODO(), subaccountId, serviceInstanceId)

	if assert.True(t, srvCalled) && assert.NoError(t, err) {
		assert.Equal(t, 200, res.StatusCode)
	}
}

func TestConnectivityDestinationFragmentFacade_CreateBySubaccount(t *testing.T) {
	command := "connectivity/destination-fragment"
	subId := "sub-123"
	content := map[string]string{"k": "v"}
	expectedJSON := `{"k":"v"}`

	var srvCalled bool
	uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srvCalled = true

		assertCall(t, r, command, ActionCreate, map[string]string{
			"subaccount": subId,
			"content":    expectedJSON,
		})
	}))
	defer srv.Close()

	_, res, err := uut.Connectivity.DestinationFragment.CreateBySubaccount(context.TODO(), subId, content)

	if assert.True(t, srvCalled) && assert.NoError(t, err) {
		assert.Equal(t, 200, res.StatusCode)
	}
}

func TestConnectivityDestinationFragmentFacade_CreateByServiceInstance(t *testing.T) {
	command := "connectivity/destination-fragment"
	subId := "sub-123"
	instId := "inst-001"
	content := map[string]string{"foo": "bar"}
	expectedJSON := `{"foo":"bar"}`

	var srvCalled bool
	uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srvCalled = true

		assertCall(t, r, command, ActionCreate, map[string]string{
			"subaccount":      subId,
			"serviceInstance": instId,
			"content":         expectedJSON,
		})
	}))
	defer srv.Close()

	_, res, err := uut.Connectivity.DestinationFragment.CreateByServiceInstance(context.TODO(), subId, instId, content)

	if assert.True(t, srvCalled) && assert.NoError(t, err) {
		assert.Equal(t, 200, res.StatusCode)
	}
}

func TestConnectivityDestinationFragmentFacade_UpdateBySubaccount(t *testing.T) {
	command := "connectivity/destination-fragment"
	subId := "sub-123"
	content := map[string]string{"x": "y"}
	expectedJSON := `{"x":"y"}`

	var srvCalled bool
	uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srvCalled = true

		assertCall(t, r, command, ActionUpdate, map[string]string{
			"subaccount": subId,
			"content":    expectedJSON,
		})
	}))
	defer srv.Close()

	_, res, err := uut.Connectivity.DestinationFragment.UpdateBySubaccount(context.TODO(), subId, content)

	if assert.True(t, srvCalled) && assert.NoError(t, err) {
		assert.Equal(t, 200, res.StatusCode)
	}
}

func TestConnectivityDestinationFragmentFacade_UpdateByServiceInstance(t *testing.T) {
	command := "connectivity/destination-fragment"
	subId := "sub-123"
	instId := "inst-001"
	content := map[string]string{"a": "b"}
	expectedJSON := `{"a":"b"}`

	var srvCalled bool
	uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srvCalled = true

		assertCall(t, r, command, ActionUpdate, map[string]string{
			"subaccount":      subId,
			"serviceInstance": instId,
			"content":         expectedJSON,
		})
	}))
	defer srv.Close()

	_, res, err := uut.Connectivity.DestinationFragment.UpdateByServiceInstance(context.TODO(), subId, instId, content)

	if assert.True(t, srvCalled) && assert.NoError(t, err) {
		assert.Equal(t, 200, res.StatusCode)
	}
}

func TestConnectivityDestinationFragmentFacade_DeleteBySubaccount(t *testing.T) {
	command := "connectivity/destination-fragment"
	subId := "sub-123"
	name := "frag1"

	var srvCalled bool
	uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srvCalled = true

		assertCall(t, r, command, ActionDelete, map[string]string{
			"subaccount": subId,
			"name":       name,
		})
	}))
	defer srv.Close()

	_, res, err := uut.Connectivity.DestinationFragment.DeleteBySubaccount(context.TODO(), subId, name)

	if assert.True(t, srvCalled) && assert.NoError(t, err) {
		assert.Equal(t, 200, res.StatusCode)
	}
}

func TestConnectivityDestinationFragmentFacade_DeleteByServiceInstance(t *testing.T) {
	command := "connectivity/destination-fragment"
	subId := "sub-123"
	instId := "inst-001"
	name := "frag1"

	var srvCalled bool
	uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srvCalled = true

		assertCall(t, r, command, ActionDelete, map[string]string{
			"subaccount":      subId,
			"name":            name,
			"serviceInstance": instId,
		})
	}))
	defer srv.Close()

	_, res, err := uut.Connectivity.DestinationFragment.DeleteByServiceInstance(context.TODO(), subId, name, instId)

	if assert.True(t, srvCalled) && assert.NoError(t, err) {
		assert.Equal(t, 200, res.StatusCode)
	}
}
