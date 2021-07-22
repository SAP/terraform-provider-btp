package btpcli

import (
	"net/http"
)

// injectBTPCLITransport wraps the transport of the given client with the btpcliTransport
func injectBTPCLITransport(client *http.Client) *http.Client {
	parentTransport := http.DefaultTransport

	if client.Transport != nil {
		parentTransport = client.Transport
	}

	client.Transport = &btpcliTransport{transport: parentTransport}
	return client
}

// btpcliTransport implements the http.RoundTripper interface. Its purpose is to copy headers from parent responses (in case of redirects)
// to the actual request.
type btpcliTransport struct {
	transport http.RoundTripper
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
	bt.copyResponseHeaderToRequestHeader(req, HeaderCLIReplacementRefreshToken, HeaderCLIRefreshToken)
	bt.copyResponseHeaderToRequestHeader(req, HeaderCLISubdomain, HeaderCLISubdomain)
	bt.copyResponseHeaderToRequestHeader(req, HeaderIDToken, HeaderIDToken)

	return bt.transport.RoundTrip(req)
}
