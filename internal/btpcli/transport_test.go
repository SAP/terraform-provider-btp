package btpcli

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBTPCLITransport(t *testing.T) {
	tests := []struct {
		description string

		redirectXtimes int

		returnStatus int

		returnAdditionalHeaders       map[string]string
		returnAdditionalHeadersInCall int
	}{
		{
			description:    "happy path",
			redirectXtimes: 3,
			returnStatus:   http.StatusOK,
		},
		{
			description:                   "happy path - forwards idtoken in case of redirect",
			redirectXtimes:                3,
			returnAdditionalHeadersInCall: 3,
			returnAdditionalHeaders: map[string]string{
				HeaderIDToken: "id-token-val",
			},
			returnStatus: http.StatusOK,
		},
		{
			description:                   "happy path - forwards subdomain in case of redirect",
			redirectXtimes:                3,
			returnAdditionalHeadersInCall: 3,
			returnAdditionalHeaders: map[string]string{
				HeaderCLISubdomain: "another-ga",
			},
			returnStatus: http.StatusOK,
		},
		{
			description:                   "happy path - forwards idtoken and subdomain",
			redirectXtimes:                2,
			returnAdditionalHeadersInCall: 2,
			returnAdditionalHeaders: map[string]string{
				HeaderIDToken:      "id-token-val",
				HeaderCLISubdomain: "another-ga",
			},
			returnStatus: http.StatusOK,
		},
		{
			description:  "error path",
			returnStatus: http.StatusForbidden,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			count := 1

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, fmt.Sprintf("token-%d", count-1), r.Header.Get(HeaderCLIRefreshToken))

				w.Header().Set(HeaderCLIReplacementRefreshToken, fmt.Sprintf("token-%d", count))

				if test.returnAdditionalHeadersInCall == count {
					for k, v := range test.returnAdditionalHeaders {
						w.Header().Set(k, v)
					}
				}

				if test.returnAdditionalHeadersInCall+1 == count {
					for k, v := range test.returnAdditionalHeaders {
						assert.Equal(t, v, r.Header.Get(k), fmt.Sprintf("value of header '%s' is expected to be '%s', but is '%s'", k, v, r.Header.Get(k)))
					}
				}

				if count < test.redirectXtimes+1 {
					count = count + 1
					http.RedirectHandler("", http.StatusTemporaryRedirect).ServeHTTP(w, r)
					return
				}

				w.WriteHeader(test.returnStatus)
			}))
			defer srv.Close()

			client := injectBTPCLITransport(srv.Client())

			req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
			req.Header.Set(HeaderCLIRefreshToken, "token-0")

			assert.NoError(t, err)

			res, err := client.Do(req)

			assert.NoError(t, err)
			assert.Equal(t, test.returnStatus, res.StatusCode)
			assert.Equal(t, fmt.Sprintf("token-%d", test.redirectXtimes+1), res.Header.Get(HeaderCLIReplacementRefreshToken))
		})
	}
}
