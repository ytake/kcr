package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/ytake/kcr/log"
	"io"
	"net/http"
	"net/url"
	"runtime"
)

const (
	version            = 1.0
	// connectorsURI Get a list of active connectors
	// GET or PUT
	connectorsURI ConnectorsURI = "/connectors"
	// connectorConfigURI the configuration for the connector.
	// GET or PUT
	connectorConfigURI ConnectorsURI = "/connectors/%s/config"
)

type (
	ConnectorsURI string
	RESTClient struct {
		url        *url.URL
		HTTPClient *http.Client
	}
	// Connectors for Get a list of active connectors or new connector
	Connectors struct {
		RESTClient
		BasicAuth
	}
	// Requester リクエスト仕様インターフェース
	Requester interface {
		newRequest(ctx context.Context, method string, connect ConnectorsURI, body io.Reader) (*http.Request, error)
	}
	BasicAuth struct {
		BasicUsername string
		BasicPassword string
	}
)

// user agent
var ua = fmt.Sprintf("kcr/%.1f (%s)", version, runtime.Version())

// retryClient internal
func retryClient(logger log.Logger) *http.Client {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5
	retryClient.Logger = logger
	return retryClient.StandardClient()
}

// decodeBody internal
func decodeBody(res *http.Response, out interface{}) error {
	b, _ := io.ReadAll(res.Body)
	err := json.Unmarshal(b, out)
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}
