package client

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type ts struct {
	name                string
	mockResponseBody    string
	expectedMethod      string
	expectedRequestPath string
	expectedErrMessage  string
}

const connectorsResponseOK = `["my-jdbc-source", "my-hdfs-sink"]`

func TestConnectors_Get(t *testing.T) {
	var tt []ts
	tt = append(tt, ts{
		name:                "get connectors",
		expectedMethod:      http.MethodGet,
		expectedRequestPath: "/connectors",
		mockResponseBody:    connectorsResponseOK,
	})
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if req.Method != tc.expectedMethod {
					t.Fatalf("request method wrong. want=%s, got=%s", tc.expectedMethod, req.Method)
				}
				if req.URL.Path != tc.expectedRequestPath {
					t.Fatalf("request path wrong. want=%s, got=%s", tc.expectedRequestPath, req.URL.Path)
				}
				w.WriteHeader(http.StatusOK)
				bodyBytes, _ := io.ReadAll(strings.NewReader(tc.mockResponseBody))
				w.Write(bodyBytes)
			}))
			defer server.Close()
			serverURL, _ := url.Parse(server.URL)
			cs := &Connectors{
				RESTClient: RESTClient{
					url:        serverURL,
					HTTPClient: server.Client(),
				},
			}
			r := cs.Get(context.TODO())
			rs := <-r
			if len(rs.Connectors) != 2 {
				t.Fatalf("wrong number of elements returned. want=%d, got=%d", 2, len(rs.Connectors))
			}
		})
	}
}
