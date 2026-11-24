package btpcli

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
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
				assert.Equal(t, fmt.Sprintf("/path-%d", count-1), r.URL.Path)

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
					http.RedirectHandler(fmt.Sprintf("/path-%d", count), http.StatusTemporaryRedirect).ServeHTTP(w, r)
					count = count + 1
					return
				}

				w.WriteHeader(test.returnStatus)
			}))
			defer srv.Close()

			srvURL, err := url.Parse(srv.URL)
			assert.NoError(t, err)

			client := injectBTPCLITransport(srv.Client(), srvURL)

			req, err := http.NewRequest(http.MethodGet, srv.URL+"/path-0", nil)

			assert.NoError(t, err)

			res, err := client.Do(req)

			assert.NoError(t, err)
			assert.Equal(t, test.returnStatus, res.StatusCode)
		})
	}
}

// TestBTPCLITransport_CrossOriginRedirect verifies that sensitive BTP credential headers are NOT
// forwarded when a redirect points to a different origin (CVE / HASI2026375-4790).
func TestBTPCLITransport_CrossOriginRedirect(t *testing.T) {
	sensitiveHeaders := []struct{ key, value string }{
		{HeaderIDToken, "super-secret-id-token"},
		{HeaderCLISubdomain, "my-global-account"},
	}

	for _, hdr := range sensitiveHeaders {
		t.Run("cross-origin redirect must not forward "+hdr.key, func(t *testing.T) {
			// attackerServer records whether it received the sensitive header.
			headerReceived := false
			attackerSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get(hdr.key) != "" {
					headerReceived = true
				}
				w.WriteHeader(http.StatusOK)
			}))
			defer attackerSrv.Close()

			// btpSrv responds with the sensitive header and then redirects to attackerSrv.
			btpSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set(hdr.key, hdr.value)
				http.RedirectHandler(attackerSrv.URL+"/steal", http.StatusTemporaryRedirect).ServeHTTP(w, r)
			}))
			defer btpSrv.Close()

			btpURL, err := url.Parse(btpSrv.URL)
			assert.NoError(t, err)

			// Use btpSrv's client so TLS (if any) is handled; redirects to attackerSrv are cross-origin.
			client := injectBTPCLITransport(btpSrv.Client(), btpURL)

			req, err := http.NewRequest(http.MethodGet, btpSrv.URL+"/start", nil)
			assert.NoError(t, err)

			res, err := client.Do(req)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, res.StatusCode)

			assert.False(t, headerReceived, "header %q must NOT be forwarded to a cross-origin redirect target", hdr.key)
		})
	}
}

// TestBTPCLITransport_NilServerURL verifies that when no serverURL is configured, sensitive headers
// are withheld (safe default).
func TestBTPCLITransport_NilServerURL(t *testing.T) {
	headerReceived := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(HeaderIDToken) != "" || r.Header.Get(HeaderCLISubdomain) != "" {
			headerReceived = true
		}
		// First call: set headers and redirect; second call: respond OK.
		if r.URL.Path == "/start" {
			w.Header().Set(HeaderIDToken, "token-val")
			w.Header().Set(HeaderCLISubdomain, "ga-val")
			http.RedirectHandler("/end", http.StatusTemporaryRedirect).ServeHTTP(w, r)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := injectBTPCLITransport(srv.Client(), nil)

	req, err := http.NewRequest(http.MethodGet, srv.URL+"/start", nil)
	assert.NoError(t, err)

	res, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	assert.False(t, headerReceived, "headers must NOT be forwarded when serverURL is nil")
}
