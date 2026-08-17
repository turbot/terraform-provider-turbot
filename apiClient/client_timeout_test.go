package apiClient

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/machinebox/graphql"
	"github.com/stretchr/testify/assert"
)

// A request against a server that accepts the connection but never responds must return a timeout
// error within roughly the configured deadline, not hang forever. This is the whole point of
// Defect 3: before the fix doRequest used context.Background() with no deadline.
func TestDoRequestHonoursTimeout(t *testing.T) {
	// A listener that accepts connections and holds them open without ever writing a response.
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	assert.NoError(t, err)
	defer listener.Close()

	accepted := make(chan net.Conn, 16)
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			accepted <- conn // hold the connection open; never respond
		}
	}()
	defer func() {
		close(accepted)
		for conn := range accepted {
			conn.Close()
		}
	}()

	client := &Client{
		Graphql:        graphql.NewClient("http://" + listener.Addr().String() + "/graphql"),
		RequestTimeout: 300 * time.Millisecond,
	}

	done := make(chan error, 1)
	start := time.Now()
	go func() {
		var resp map[string]interface{}
		done <- client.doRequest(`{ __typename }`, nil, &resp)
	}()

	select {
	case err := <-done:
		assert.Error(t, err, "a hung request must return an error, not nil")
		assert.True(t, time.Since(start) < 5*time.Second,
			"the request must abort near the deadline, not hang")
	case <-time.After(5 * time.Second):
		t.Fatal("doRequest hung well past its timeout — the deadline is not being applied")
	}
}

// With no timeout configured (zero), doRequest must NOT impose a deadline — a hand-built client
// keeps the old unbounded behaviour unless it opts in. Asserted by confirming a fast call
// completes normally rather than by waiting forever.
func TestDoRequestZeroTimeoutDoesNotDeadline(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"__typename":"Query"}}`))
	}))
	defer server.Close()

	client := &Client{Graphql: graphql.NewClient(server.URL + "/graphql")} // RequestTimeout == 0
	var resp map[string]interface{}
	assert.NoError(t, client.doRequest(`{ __typename }`, nil, &resp),
		"a zero timeout must not break a normal, fast request")
}

// A slow-but-successful response that completes inside the deadline must succeed — the timeout
// must not clip a legitimately slow Guardrails operation.
func TestDoRequestSlowButWithinDeadlineSucceeds(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"__typename":"Query"}}`))
	}))
	defer server.Close()

	client := &Client{
		Graphql:        graphql.NewClient(server.URL + "/graphql"),
		RequestTimeout: 5 * time.Second,
	}
	var resp map[string]interface{}
	assert.NoError(t, client.doRequest(`{ __typename }`, nil, &resp),
		"a response that arrives well within the deadline must succeed")
}

// CreateClient must install the default timeout when the config leaves it unset, and honour an
// explicit one. This is what guarantees a provider-built client is always bounded.
func TestCreateClientInstallsTimeout(t *testing.T) {
	creds := ClientCredentials{AccessKey: "AK", SecretKey: "SK", Workspace: "https://x.example.com"}

	defaulted, err := CreateClient(ClientConfig{Credentials: creds})
	assert.NoError(t, err)
	assert.Equal(t, DefaultRequestTimeout, defaulted.RequestTimeout,
		"an unset timeout must fall back to DefaultRequestTimeout")

	custom, err := CreateClient(ClientConfig{Credentials: creds, RequestTimeout: 42 * time.Second})
	assert.NoError(t, err)
	assert.Equal(t, 42*time.Second, custom.RequestTimeout,
		"an explicit timeout must be honoured")
}
