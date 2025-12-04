package btpcli

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	command      = "connectivity/destination-certificate"
	subaccountId = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	instanceId   = "bc8a216f-1184-49dc-b4b4-17cfe2828965"
	fileName     = "test.p12"
)

func TestConnectivityDestinationCertificateFacade_Get(t *testing.T) {
	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionGet, map[string]string{
				"subaccount":      subaccountId,
				"serviceInstance": instanceId,
				"certName":        fileName,
			})

		}))
		defer srv.Close()

		_, res, err := uut.Connectivity.DestinationCertificate.Get(context.TODO(), &DestinationCertificateGetInput{
			SubaccountId:      subaccountId,
			ServiceInstanceId: instanceId,
			CertificateName:   fileName,
		})

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestConnectivityDestinationCertificateFacade_Delete(t *testing.T) {
	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionDelete, map[string]string{
				"subaccount":      subaccountId,
				"serviceInstance": instanceId,
				"certName":        fileName,
			})

		}))
		defer srv.Close()

		res, err := uut.Connectivity.DestinationCertificate.Delete(context.TODO(), &DestinationCertificateGetInput{
			SubaccountId:      subaccountId,
			ServiceInstanceId: instanceId,
			CertificateName:   fileName,
		})

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestConnectivityDestinationCertificateFacade_List(t *testing.T) {
	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		var payload struct {
			ParamValues map[string]string `json:"paramValues"`
		}

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			bodyBytes, _ := io.ReadAll(r.Body)

			// Decode the payload
			_ = json.Unmarshal(bodyBytes, &payload)

			// Restore the body for assertCall
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			if len(payload.ParamValues) == 2 {

				assertCall(t, r, command, ActionList, map[string]string{
					"subaccount": subaccountId,
					"namesOnly":  "false",
				})
			} else if len(payload.ParamValues) == 3 {

				assertCall(t, r, command, ActionList, map[string]string{
					"subaccount":      subaccountId,
					"serviceInstance": instanceId,
					"namesOnly":       "false",
				})
			}

		}))
		defer srv.Close()

		_, res, err := uut.Connectivity.DestinationCertificate.List(context.TODO(), subaccountId, instanceId)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 0, res.StatusCode)
		}
	})
}
