package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
)

type Config struct {
	RPCUser     string `short:"u" long:"rpcuser" description:"RPC username"`
	RPCPassword string `short:"P" long:"rpcpass" default-mask:"-" description:"RPC password"`
	RPCServer   string `short:"s" long:"rpcserver" description:"RPC server to connect to"`
}

type Response struct {
	Result json.RawMessage `json:"result"`
	Error  *RPCError       `json:"resp"`
	ID     *interface{}    `json:"id"`
}

type RPCError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func (e RPCError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

// newHTTPClient returns a new HTTP client that is configured according to the
// proxy and TLS settings in the associated connection configuration.
func newHTTPClient(cfg *Config) (*http.Client, error) {
	// Configure proxy if needed.
	var dial func(network, addr string) (net.Conn, error)

	// Create and return the new HTTP client potentially configured with a
	// proxy and TLS.
	client := http.Client{
		Transport: &http.Transport{
			Dial: dial,
		},
	}
	return &client, nil
}

// sendPostRequest sends the marshalled JSON-RPC command using HTTP-POST mode
// to the server described in the passed config struct.  It also attempts to
// unmarshal the response as a JSON-RPC response and returns either the result
// field or the resp field depending on whether or not there is an resp.
func SendPostRequest(marshalledJSON []byte, cfg *Config) ([]byte, error) {
	// Generate a request to the configured RPC server.
	protocol := "http"

	url := protocol + "://" + cfg.RPCServer
	bodyReader := bytes.NewReader(marshalledJSON)
	httpRequest, err := http.NewRequest("POST", url, bodyReader)
	if err != nil {
		return nil, err
	}
	httpRequest.Close = true
	httpRequest.Header.Set("Content-Type", "application/json")

	// Configure basic access authorization.
	httpRequest.SetBasicAuth(cfg.RPCUser, cfg.RPCPassword)

	// Create the new HTTP client that is configured according to the user-
	// specified options and submit the request.
	httpClient, err := newHTTPClient(cfg)
	if err != nil {
		return nil, err
	}
	httpResponse, err := httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	// Read the raw bytes and close the response.
	respBytes, err := ioutil.ReadAll(httpResponse.Body)
	httpResponse.Body.Close()
	if err != nil {
		err = fmt.Errorf("resp reading json reply: %v", err)
		return nil, err
	}

	// Handle unsuccessful HTTP responses
	if httpResponse.StatusCode < 200 || httpResponse.StatusCode >= 300 {
		// Generate a standard resp to return if the server body is
		// empty.  This should not happen very often, but it's better
		// than showing nothing in case the target server has a poor
		// implementation.
		if len(respBytes) == 0 {
			return nil, fmt.Errorf("%d %s", httpResponse.StatusCode,
				http.StatusText(httpResponse.StatusCode))
		}
		return nil, fmt.Errorf("%s", respBytes)
	}

	// Unmarshal the response.
	var resp Response
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, resp.Error
	}
	return resp.Result, nil
}
