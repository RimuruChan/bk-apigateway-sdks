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

// Package apigateway is a thin client of the bk-apigateway v2 APIs.
//
// The client handles what all the APIs have in common: authentication, response status,
// the {"data": ...} response envelope and logging. Request bodies are any values that
// encoding/json can encode, such as maps or structs. Results are decoded into the value
// that result points to, which can be nil to discard the result.
package apigateway

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	gentleman "gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugins/timeout"
	"gopkg.in/h2non/gentleman.v2/plugins/transport"
)

const defaultTimeout = 60 * time.Second

// Config is the configuration of a Client.
type Config struct {
	// Endpoint is the URL of the bk-apigateway stage, such as https://bkapi.example.com/api/bk-apigateway/prod.
	Endpoint string

	// AppCode, AppSecret and AccessToken are sent in the X-Bkapi-Authorization header.
	AppCode     string
	AppSecret   string
	AccessToken string
	// TenantID is sent in the X-Bk-Tenant-Id header.
	TenantID string

	// Timeout limits the time of each request, including reading the response. Defaults to 60s.
	Timeout time.Duration
	// Transport sends the requests. Defaults to gentleman.DefaultTransport.
	Transport http.RoundTripper
	// Logger logs the requests: successful ones at debug level, failed ones at warn or error level.
	// Defaults to slog.Default().
	Logger *slog.Logger
}

// Client calls the bk-apigateway v2 APIs. It is safe for concurrent use.
type Client struct {
	client *gentleman.Client
	logger *slog.Logger
}

// New creates a Client.
func New(config Config) (*Client, error) {
	endpoint, err := url.Parse(config.Endpoint)
	if err != nil || endpoint.Scheme == "" || endpoint.Host == "" {
		return nil, fmt.Errorf("apigateway: invalid endpoint %q", config.Endpoint)
	}

	authorization, err := json.Marshal(struct {
		AppCode     string `json:"bk_app_code,omitempty"`
		AppSecret   string `json:"bk_app_secret,omitempty"`
		AccessToken string `json:"access_token,omitempty"`
	}{config.AppCode, config.AppSecret, config.AccessToken})
	if err != nil {
		return nil, err
	}

	client := gentleman.New().
		URL(strings.TrimSuffix(config.Endpoint, "/")).
		SetHeader("X-Bkapi-Authorization", string(authorization)).
		Use(timeout.Request(cmp.Or(config.Timeout, defaultTimeout)))
	if config.TenantID != "" {
		client.SetHeader("X-Bk-Tenant-Id", config.TenantID)
	}
	if config.Transport != nil {
		client.Use(transport.Set(config.Transport))
	}

	return &Client{client: client, logger: cmp.Or(config.Logger, slog.Default())}, nil
}

// get, post, put and delete create the request of an API, whose path params are written as :name.
func (c *Client) get(path string, query url.Values) *gentleman.Request {
	return withQuery(c.client.Get().AddPath(path), query)
}

func (c *Client) post(path string, body any) *gentleman.Request {
	return withJSON(c.client.Post().AddPath(path), body)
}

func (c *Client) put(path string, body any) *gentleman.Request {
	return withJSON(c.client.Put().AddPath(path), body)
}

func (c *Client) delete(path string, query url.Values, body any) *gentleman.Request {
	return withJSON(withQuery(c.client.Delete().AddPath(path), query), body)
}

func withQuery(req *gentleman.Request, query url.Values) *gentleman.Request {
	for key, values := range query {
		for _, value := range values {
			req.AddQuery(key, value)
		}
	}
	return req
}

func withJSON(req *gentleman.Request, body any) *gentleman.Request {
	if body != nil {
		req.JSON(body)
	}
	return req
}

// send sends the request, and decodes the response into result.
func (c *Client) send(ctx context.Context, req *gentleman.Request, result any) error {
	req.Context.SetCancelContext(ctx)

	start := time.Now()
	res, err := req.Send()
	var body []byte
	if err == nil {
		body, err = res.Bytes(), res.Error
	}
	if err == nil && (res.StatusCode < 200 || res.StatusCode >= 300) {
		err = newError(res, body)
	}
	c.log(ctx, res, time.Since(start), err)
	if err != nil || len(body) == 0 {
		return err
	}

	// The result is wrapped in data, except for OpenOAuthProtectedResource, whose result is the whole body.
	var envelope struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return fmt.Errorf("apigateway: decode response: %w", err)
	}
	if result == nil {
		return nil
	}
	if envelope.Data != nil {
		body = envelope.Data
	}
	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("apigateway: decode response: %w", err)
	}
	return nil
}

// log logs the request without the headers and bodies, which may contain secrets.
// The status of a request that got no response is 0.
func (c *Client) log(ctx context.Context, res *gentleman.Response, duration time.Duration, err error) {
	level := slog.LevelDebug
	attrs := []slog.Attr{
		slog.String("method", res.RawRequest.Method),
		slog.String("path", res.RawRequest.URL.Path),
		slog.Int("status", res.StatusCode),
		slog.String("request_id", res.Header.Get("X-Bkapi-Request-Id")),
		slog.Duration("duration", duration),
	}
	if err != nil {
		level = slog.LevelError
		if res.ClientError {
			level = slog.LevelWarn
		}
		attrs = append(attrs, slog.Any("error", err))
	}
	c.logger.LogAttrs(ctx, level, "bk-apigateway request", attrs...)
}
