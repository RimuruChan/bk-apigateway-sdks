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

package bkapi

import (
	"bytes"
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"gopkg.in/h2non/gentleman.v2/context"
	"gopkg.in/h2non/gentleman.v2/plugin"
)

// Error is returned when the API responds with a non-2xx status.
//
// The fields follow the BlueKing error response, which is like {"error": {"code": ..., "message": ...}}.
type Error struct {
	StatusCode int
	// Code is the category of the error, such as INVALID_ARGUMENT, for the caller to branch on.
	Code string
	// Message describes the error to users.
	Message string
	// System is the system that raises the error, such as bk-apigateway. It is optional.
	System string
	// Details is the JSON array of the sub errors for developers to locate the problem, such as the invalid
	// fields. Each has a code and a message, and other fields that vary with the system. It is nil if empty.
	Details json.RawMessage
	// Data is the JSON object for the caller to handle the error, such as the permissions to apply for
	// IAM_NO_PERMISSION. It is nil if empty.
	Data json.RawMessage
	// RequestID identifies the request in the logs of the gateway.
	RequestID string
}

func (e *Error) Error() string {
	var b strings.Builder
	fmt.Fprintf(&b, "bkapi: status=%d code=%s message=%q", e.StatusCode, e.Code, e.Message)
	if e.System != "" {
		fmt.Fprintf(&b, " system=%s", e.System)
	}
	if len(e.Details) > 0 {
		fmt.Fprintf(&b, " details=%s", e.Details)
	}
	fmt.Fprintf(&b, " request_id=%s", e.RequestID)
	return b.String()
}

// newError creates the Error of a non-2xx response from its body.
//
// The responses that do not follow the BlueKing error response fall back to:
//   - the legacy body like {"code": 1640301, "code_name": "APP_NO_PERMISSION", "message": ...}, such as the
//     requests rejected by the gateway and the APIs not migrated yet
//   - the X-Bkapi-Error-Code and X-Bkapi-Error-Message headers set by the gateway
//   - the text of the status, when the body is not JSON
func newError(res *http.Response, body []byte) *Error {
	var payload struct {
		Error struct {
			Code    string          `json:"code"`
			Message string          `json:"message"`
			System  string          `json:"system"`
			Details json.RawMessage `json:"details"`
			Data    json.RawMessage `json:"data"`
		} `json:"error"`
		Code     json.RawMessage `json:"code"` // a number or a string
		CodeName string          `json:"code_name"`
		Message  string          `json:"message"`
	}
	// Unmarshal keeps the fields it can decode, so the error is ignored
	_ = json.Unmarshal(body, &payload)

	return &Error{
		StatusCode: res.StatusCode,
		Code: cmp.Or(
			payload.Error.Code,
			payload.CodeName,
			strings.Trim(string(payload.Code), `"`),
			res.Header.Get("X-Bkapi-Error-Code"),
		),
		Message: cmp.Or(
			payload.Error.Message,
			payload.Message,
			res.Header.Get("X-Bkapi-Error-Message"),
			http.StatusText(res.StatusCode),
		),
		System:    payload.Error.System,
		Details:   nonEmpty(payload.Error.Details),
		Data:      nonEmpty(payload.Error.Data),
		RequestID: cmp.Or(res.Header.Get("X-Bkapi-Request-Id"), res.Header.Get("X-Request-Id")),
	}
}

// nonEmpty returns nil for the JSON values without content, which servers send for the unset optional fields.
func nonEmpty(raw json.RawMessage) json.RawMessage {
	switch string(raw) {
	case "null", "[]", "{}":
		return nil
	}
	return raw
}

// readBody reads the body of the response, and leaves it to be read again.
func readBody(res *http.Response) []byte {
	if res.Body == nil {
		return nil
	}
	body, _ := io.ReadAll(res.Body)
	_ = res.Body.Close()
	res.Body = io.NopCloser(bytes.NewReader(body))
	return body
}

// errorBodyTransport reads the body of a non-2xx response into memory.
//
// gentleman runs the plugins of a client under a lock shared by all its requests, so checkStatus must not
// read a slow body itself. Reading it here keeps the request timeout and cancellation.
type errorBodyTransport struct {
	base http.RoundTripper
}

func (t *errorBodyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	res, err := t.base.RoundTrip(req)
	if err != nil || isSuccess(res.StatusCode) || res.Body == nil {
		return res, err
	}

	body, err := io.ReadAll(res.Body)
	_ = res.Body.Close()
	if err != nil {
		return nil, err
	}
	res.Body = io.NopCloser(bytes.NewReader(body))
	return res, nil
}

func isSuccess(status int) bool {
	return status >= 200 && status < 300
}

// checkStatus returns the plugin that turns a non-2xx response into an *Error.
func checkStatus() plugin.Plugin {
	return plugin.NewResponsePlugin(func(ctx *context.Context, h context.Handler) {
		res := ctx.Response
		if isSuccess(res.StatusCode) {
			h.Next(ctx)
			return
		}

		h.Error(ctx, newError(res, readBody(res)))
	})
}
