package btpcli

import (
	"net/http"
	"net/url"
)

// injectBTPCLITransport wraps the transport of the given client with the btpcliTransport.
// serverURL is the original BTP CLI server URL; sensitive headers are only forwarded to the same origin.
func injectBTPCLITransport(client *http.Client, serverURL *url.URL) *http.Client {
	parentTransport := http.DefaultTransport

	if client.Transport != nil {
		parentTransport = client.Transport
	}

	client.Transport = &btpcliTransport{transport: parentTransport, serverURL: serverURL}
	return client
}

// btpcliTransport implements the http.RoundTripper interface. Its purpose is to copy headers from parent
// responses (in case of redirects) to the actual request, but only when the redirect target shares the
// same origin (scheme + host) as the configured BTP CLI server to prevent credential leakage.
type btpcliTransport struct {
	transport http.RoundTripper
	serverURL *url.URL
}

// sameOrigin returns true when a and b share the same scheme and host (including port).
func sameOrigin(a, b *url.URL) bool {
	return a.Scheme == b.Scheme && a.Host == b.Host
}

// copyResponseHeaderToRequestHeader copies a header from the parent response (source) to the request (target).
func (bt *btpcliTransport) copyResponseHeaderToRequestHeader(req *http.Request, source string, target string) {
	if req.Response == nil {
		return
	}

	if value := req.Response.Header.Get(source); len(value) > 0 {
		req.Header.Set(target, value)
	}
}

func (bt *btpcliTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Only forward BTP credential headers to the original server origin; a cross-origin redirect
	// must not receive session or identity material (CWE-200, CWE-601).
	if bt.serverURL != nil && sameOrigin(req.URL, bt.serverURL) {
		bt.copyResponseHeaderToRequestHeader(req, HeaderCLISubdomain, HeaderCLISubdomain)
		bt.copyResponseHeaderToRequestHeader(req, HeaderIDToken, HeaderIDToken)
	}

	return bt.transport.RoundTrip(req)
}
