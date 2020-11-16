// Copyright © 2020 The Things Network Foundation, The Things Industries B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package component

import (
	"context"
	"net/http"
	"time"
)

const (
	// defaultTimeout is the default timeout for the HTTP client.
	defaultTimeout = 10 * time.Second
)

// HTTPClientOpt applies an option to the HTTP client.
type HTTPClientOpt func(*http.Client)

// WithTimeout sets the HTTP client timeout.
func WithTimeout(t time.Duration) HTTPClientOpt {
	return HTTPClientOpt(func(c *http.Client) {
		c.Timeout = t
	})
}

// WithTransport sets the HTTP client transport.
func WithTransport(t http.RoundTripper) HTTPClientOpt {
	return HTTPClientOpt(func(c *http.Client) {
		c.Transport = t
	})
}

// HTTPClient returns a new *http.Client with a timeout and a configured
// transport. Uses the http.RoundTripper returned from HTTPTransport(). In case
// of an error, an *http.Client with http.DefaultTransport is returned, along
// with the error. It is up to the caller to decide how to handle the error.
func (c *Component) HTTPClient(ctx context.Context, opts ...HTTPClientOpt) (*http.Client, error) {
	transport, err := c.HTTPTransport(ctx)
	client := &http.Client{
		Timeout:   defaultTimeout,
		Transport: transport,
	}
	for _, opt := range opts {
		opt(client)
	}
	return client, err
}

// HTTPTransport returns a new http.RoundTripper with TLS client configuration.
// If the call fails, http.DefaultTransport is returned along with the error,
// and it is up to the caller to decide how to handle the error.
func (c *Component) HTTPTransport(ctx context.Context) (http.RoundTripper, error) {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	tlsConfig, err := c.GetTLSClientConfig(ctx)
	if err == nil {
		transport.TLSClientConfig = tlsConfig
	}
	return transport, err
}
