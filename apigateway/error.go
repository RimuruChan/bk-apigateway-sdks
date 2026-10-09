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

package apigateway

import (
	"cmp"
	"encoding/json"
	"fmt"
	"net/http"

	gentleman "gopkg.in/h2non/gentleman.v2"
)

// Error is returned when bk-apigateway responds with a non-2xx status.
type Error struct {
	StatusCode int
	// Code and Message describe the error, such as INVALID_ARGUMENT. They come from the error in the
	// response body, or from the X-Bkapi-Error-Code and X-Bkapi-Error-Message headers when the gateway
	// rejects the request, such as when the app fails to authenticate.
	Code    string
	Message string
	// RequestID identifies the request in the logs of bk-apigateway.
	RequestID string
}

func (e *Error) Error() string {
	return fmt.Sprintf("apigateway: status=%d code=%s message=%q request_id=%s",
		e.StatusCode, e.Code, e.Message, e.RequestID)
}

func newError(res *gentleman.Response, body []byte) *Error {
	var payload struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	_ = json.Unmarshal(body, &payload) // the body may not be json, then the headers are used

	return &Error{
		StatusCode: res.StatusCode,
		Code:       cmp.Or(payload.Error.Code, res.Header.Get("X-Bkapi-Error-Code")),
		Message: cmp.Or(
			payload.Error.Message,
			res.Header.Get("X-Bkapi-Error-Message"),
			http.StatusText(res.StatusCode),
		),
		RequestID: res.Header.Get("X-Bkapi-Request-Id"),
	}
}
