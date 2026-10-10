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

	"gopkg.in/h2non/gentleman.v2/context"
	"gopkg.in/h2non/gentleman.v2/plugin"
)

// Error is returned when the API responds with a non-2xx status.
type Error struct {
	StatusCode int
	// Code and Message describe the error, such as INVALID_ARGUMENT. They come from the error in the
	// response body, or from the X-Bkapi-Error-Code and X-Bkapi-Error-Message headers when the gateway
	// rejects the request, such as when the app fails to authenticate.
	Code    string
	Message string
	// RequestID identifies the request in the logs of the gateway.
	RequestID string
}

func (e *Error) Error() string {
	return fmt.Sprintf("bkapi: status=%d code=%s message=%q request_id=%s",
		e.StatusCode, e.Code, e.Message, e.RequestID)
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

		var body struct {
			Error struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
		}
		data, _ := io.ReadAll(res.Body)
		_ = res.Body.Close()
		_ = json.Unmarshal(data, &body) // the body may not be json, then the headers are used

		h.Error(ctx, &Error{
			StatusCode: res.StatusCode,
			Code:       cmp.Or(body.Error.Code, res.Header.Get("X-Bkapi-Error-Code")),
			Message: cmp.Or(
				body.Error.Message,
				res.Header.Get("X-Bkapi-Error-Message"),
				http.StatusText(res.StatusCode),
			),
			RequestID: res.Header.Get("X-Bkapi-Request-Id"),
		})
	})
}
