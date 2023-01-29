package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/ytake/kcr/log"
	"github.com/ytake/kcr/payload"
	"io"
	"net/http"
	"net/url"
	"path"
)

// NewConnectors for connectorsURI
func NewConnectors(connectServer string, logger log.Logger) (*Connectors, error) {
	u, err := url.Parse(connectServer)
	if err != nil {
		return nil, err
	}
	return &Connectors{
		RESTClient: RESTClient{
			url:        u,
			HTTPClient: retryClient(&DefaultClient{Logger: logger}),
		},
	}, nil
}

// ErrKafkaConnectServerRequest
var ErrKafkaConnectServerRequest = func(message error) error {
	return fmt.Errorf("request client error: %w", message)
}

func (cs *Connectors) newRequest(ctx context.Context, method string, connect ConnectorsURI, body io.Reader) (*http.Request, error) {
	u := *cs.url
	u.Path = path.Join(string(connect))
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, ErrKafkaConnectServerRequest(err)
	}
	req = req.WithContext(ctx)
	if cs.BasicPassword != "" && cs.BasicUsername != "" {
		req.SetBasicAuth(cs.BasicUsername, cs.BasicPassword)
	}
	req.Header.Set("User-Agent", ua)
	return req, nil
}

// Get a list of active connectors
func (cs *Connectors) Get(ctx context.Context) <-chan payload.ResultConnectors {
	out := make(chan payload.ResultConnectors)
	go func() {
		defer close(out)
		var result payload.ResultConnectors
		req, err := cs.newRequest(ctx, http.MethodGet, connectorsURI, bytes.NewBuffer([]byte{}))
		if err != nil {
			result.Err = err
			out <- result
			return
		}
		res, err := cs.HTTPClient.Do(req)
		defer res.Body.Close()
		if err != nil {
			result.Err = err
			out <- result
			return
		}
		if res.StatusCode == http.StatusOK {
			var definition []string
			if err := decodeBody(res, &definition); err != nil {
				result.Err = err
				out <- result
				return
			}
			result.Connectors = definition
			out <- result
			return
		}
		result.Err = errors.New("kafka connect server error")
		out <- result
	}()
	return out
}
