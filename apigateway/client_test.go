/**
 * TencentBlueKing is pleased to support the open source community by
 * making 蓝鲸智云-蓝鲸 PaaS 平台(BlueKing-PaaS) available.
 * Copyright (C) 2025 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

package apigateway_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/TencentBlueKing/bk-apigateway-sdks/v2/apigateway"
)

func newClient(t *testing.T, config apigateway.Config) *apigateway.Client {
	t.Helper()
	if config.Logger == nil {
		config.Logger = slog.New(slog.DiscardHandler)
	}
	c, err := apigateway.New(config)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// serve starts a server that replies with the status, headers and body.
func serve(t *testing.T, status int, header http.Header, body string) string {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for key, values := range header {
			w.Header()[key] = values
		}
		w.WriteHeader(status)
		fmt.Fprint(w, body)
	}))
	t.Cleanup(server.Close)
	return server.URL
}

func TestNewWithInvalidEndpoint(t *testing.T) {
	for _, endpoint := range []string{"", "example.com", "/api/bk-apigateway/prod", "http://"} {
		if _, err := apigateway.New(apigateway.Config{Endpoint: endpoint}); err == nil {
			t.Errorf("New(%q) should fail", endpoint)
		}
	}
}

func TestRequest(t *testing.T) {
	var got *http.Request
	var gotBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		got, gotBody = r, string(body)
		fmt.Fprint(w, `{"data": {"id": 1, "name": "prod"}}`)
	}))
	defer server.Close()

	c := newClient(t, apigateway.Config{
		Endpoint:  server.URL + "/api/bk-apigateway/prod/",
		AppCode:   "my-app",
		AppSecret: "my-secret",
		TenantID:  "my-tenant",
	})
	var result struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	err := c.SyncStages(context.Background(), "my-gateway", map[string]string{"name": "prod"}, &result)
	if err != nil {
		t.Fatal(err)
	}

	if result.ID != 1 || result.Name != "prod" {
		t.Errorf("result = %+v", result)
	}
	const wantRequest = "POST /api/bk-apigateway/prod/api/v2/sync/gateways/my-gateway/stages/"
	if request := got.Method + " " + got.URL.Path; request != wantRequest {
		t.Errorf("request = %s, want %s", request, wantRequest)
	}
	for key, want := range map[string]string{
		"X-Bkapi-Authorization": `{"bk_app_code":"my-app","bk_app_secret":"my-secret"}`,
		"X-Bk-Tenant-Id":        "my-tenant",
		"Content-Type":          "application/json",
	} {
		if value := got.Header.Get(key); value != want {
			t.Errorf("header %s = %q, want %q", key, value, want)
		}
	}
	if gotBody != `{"name":"prod"}`+"\n" {
		t.Errorf("body = %q", gotBody)
	}
}

func TestAccessToken(t *testing.T) {
	var authorization string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization = r.Header.Get("X-Bkapi-Authorization")
	}))
	defer server.Close()

	c := newClient(t, apigateway.Config{Endpoint: server.URL, AppCode: "my-app", AccessToken: "my-token"})
	if err := c.OpenListGateways(context.Background(), nil, nil); err != nil {
		t.Fatal(err)
	}
	if authorization != `{"bk_app_code":"my-app","access_token":"my-token"}` {
		t.Errorf("authorization = %s", authorization)
	}
}

func TestQuery(t *testing.T) {
	var query url.Values
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query = r.URL.Query()
		fmt.Fprint(w, `{"data": {"count": 0, "results": []}}`)
	}))
	defer server.Close()

	want := url.Values{"version": {"1.0.0+prod"}, "id": {"1", "2"}}
	c := newClient(t, apigateway.Config{Endpoint: server.URL})
	if err := c.SyncListResourceVersions(context.Background(), "my-gateway", want, nil); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(query, want) {
		t.Errorf("query = %v, want %v", query, want)
	}
}

func TestResult(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		body   string
		want   any
	}{
		{"object", 200, `{"data": {"id": 1}}`, map[string]any{"id": float64(1)}},
		{"list", 201, `{"data": [{"id": 1}]}`, []any{map[string]any{"id": float64(1)}}},
		{"null", 201, `{"data": null}`, nil},
		{"no content", 204, ``, nil},
		{"not wrapped in data", 200, `{"resource": "https://example.com"}`, map[string]any{
			"resource": "https://example.com",
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := newClient(t, apigateway.Config{Endpoint: serve(t, tc.status, nil, tc.body)})
			var result any
			if err := c.OpenListGateways(context.Background(), nil, &result); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(result, tc.want) {
				t.Errorf("result = %#v, want %#v", result, tc.want)
			}
		})
	}
}

func TestInvalidResponse(t *testing.T) {
	c := newClient(t, apigateway.Config{Endpoint: serve(t, 200, nil, `<html>login</html>`)})
	// the response is checked even if the result is discarded
	if err := c.OpenListGateways(context.Background(), nil, nil); err == nil {
		t.Error("expected a decode error")
	}
}

func TestInvalidBody(t *testing.T) {
	var called bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true }))
	defer server.Close()

	c := newClient(t, apigateway.Config{Endpoint: server.URL})
	if err := c.SyncGateway(context.Background(), "my-gateway", make(chan int), nil); err == nil || called {
		t.Errorf("error = %v, request sent = %v", err, called)
	}
}

func TestError(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		header http.Header
		body   string
		want   apigateway.Error
	}{
		{
			name:   "api error",
			status: 400,
			header: http.Header{"X-Bkapi-Request-Id": {"request-1"}},
			body:   `{"error": {"code": "INVALID_ARGUMENT", "message": "invalid version", "details": [], "data": {}}}`,
			want:   apigateway.Error{StatusCode: 400, Code: "INVALID_ARGUMENT", Message: "invalid version", RequestID: "request-1"},
		},
		{
			name:   "gateway error",
			status: 403,
			header: http.Header{
				"X-Bkapi-Request-Id":    {"request-2"},
				"X-Bkapi-Error-Code":    {"1640301"},
				"X-Bkapi-Error-Message": {"App has no permission to the resource"},
			},
			body: `{"code": 1640301, "code_name": "APP_NO_PERMISSION", "message": "App has no permission to the resource"}`,
			want: apigateway.Error{
				StatusCode: 403, Code: "1640301", Message: "App has no permission to the resource", RequestID: "request-2",
			},
		},
		{
			name:   "unknown error",
			status: 502,
			body:   `<html>502 Bad Gateway</html>`,
			want:   apigateway.Error{StatusCode: 502, Message: "Bad Gateway"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := newClient(t, apigateway.Config{Endpoint: serve(t, tc.status, tc.header, tc.body)})
			err := c.OpenListGateways(context.Background(), nil, nil)

			var apiErr *apigateway.Error
			if !errors.As(err, &apiErr) || *apiErr != tc.want {
				t.Errorf("error = %v, want %v", err, &tc.want)
			}
		})
	}
}

func TestSyncResourceDoc(t *testing.T) {
	var path, content string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		file, _, err := r.FormFile("file")
		if err == nil {
			data, _ := io.ReadAll(file)
			content = string(data)
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"data": null}`)
	}))
	defer server.Close()

	c := newClient(t, apigateway.Config{Endpoint: server.URL})
	err := c.SyncResourceDoc(context.Background(), "my-gateway", strings.NewReader("zip content"))
	if err != nil {
		t.Fatal(err)
	}
	if path != "/api/v2/sync/gateways/my-gateway/resource-docs/" || content != "zip content" {
		t.Errorf("uploaded %q to %s", content, path)
	}
}

func TestContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	c := newClient(t, apigateway.Config{Endpoint: server.URL})
	if err := c.OpenListGateways(ctx, nil, nil); !errors.Is(err, context.Canceled) {
		t.Errorf("error = %v, want context.Canceled", err)
	}

	c = newClient(t, apigateway.Config{Endpoint: server.URL, Timeout: 50 * time.Millisecond})
	if err := c.OpenListGateways(context.Background(), nil, nil); !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("error = %v, want context.DeadlineExceeded", err)
	}
}

func TestTransport(t *testing.T) {
	var called bool
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		called = true
		return http.DefaultTransport.RoundTrip(r)
	})
	c := newClient(t, apigateway.Config{Endpoint: serve(t, 200, nil, `{"data": []}`), Transport: transport})
	if err := c.OpenListGateways(context.Background(), nil, nil); err != nil || !called {
		t.Errorf("error = %v, transport called = %v", err, called)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestLog(t *testing.T) {
	var logs bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logs, &slog.HandlerOptions{Level: slog.LevelDebug}))
	header := http.Header{"X-Bkapi-Request-Id": {"request-1"}}
	c := newClient(t, apigateway.Config{
		Endpoint:  serve(t, 400, header, `{"error": {"code": "INVALID_ARGUMENT", "message": "invalid"}}`),
		AppSecret: "my-secret",
		Logger:    logger,
	})
	_ = c.SyncGateway(context.Background(), "my-gateway", map[string]string{"description": "my-description"}, nil)

	for _, want := range []string{
		"level=WARN", "method=POST", "path=/api/v2/sync/gateways/my-gateway/", "status=400", "request_id=request-1",
	} {
		if !strings.Contains(logs.String(), want) {
			t.Errorf("log %q does not contain %q", logs.String(), want)
		}
	}
	for _, secret := range []string{"my-secret", "my-description"} {
		if strings.Contains(logs.String(), secret) {
			t.Errorf("log %q contains %q", logs.String(), secret)
		}
	}
}

func TestConcurrentRequests(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"data": %q}`, r.URL.Path)
	}))
	defer server.Close()

	c := newClient(t, apigateway.Config{Endpoint: server.URL})
	var wg sync.WaitGroup
	for i := range 20 {
		wg.Go(func() {
			name := fmt.Sprint("gateway-", i)
			var path string
			if err := c.OpenGetGateway(context.Background(), name, nil, &path); err != nil ||
				path != "/api/v2/open/gateways/"+name+"/" {
				t.Errorf("path = %s, error = %v", path, err)
			}
		})
	}
	wg.Wait()
}
