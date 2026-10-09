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
	"os"
	"testing"

	"github.com/TencentBlueKing/bk-apigateway-sdks/v2/apigateway"
)

func setenv(t *testing.T, env map[string]string) {
	t.Helper()
	for _, key := range []string{
		"BK_API_URL_TMPL", "BK_APP_CODE", "BKPAAS_APP_ID", "APP_CODE", "BK_APP_SECRET", "BKPAAS_APP_SECRET",
		"SECRET_KEY", "BKPAAS_APP_TENANT_ID", "BK_APP_TENANT_ID",
	} {
		t.Setenv(key, "") // restores the variable after the test
		if err := os.Unsetenv(key); err != nil {
			t.Fatal(err)
		}
	}
	for key, value := range env {
		t.Setenv(key, value)
	}
}

func TestConfigFromEnv(t *testing.T) {
	for _, tc := range []struct {
		name string
		env  map[string]string
		want apigateway.Config
	}{
		{
			name: "paas app",
			env: map[string]string{
				"BK_API_URL_TMPL":      "http://bkapi.example.com/api/{api_name}/",
				"BKPAAS_APP_ID":        "paas-app",
				"BKPAAS_APP_SECRET":    "paas-secret",
				"BKPAAS_APP_TENANT_ID": "paas-tenant",
				"BK_APP_TENANT_ID":     "app-tenant",
			},
			want: apigateway.Config{
				Endpoint:  "http://bkapi.example.com/api/bk-apigateway/prod",
				AppCode:   "paas-app",
				AppSecret: "paas-secret",
				TenantID:  "paas-tenant",
			},
		},
		{
			name: "global tenant paas app",
			env:  map[string]string{"BKPAAS_APP_TENANT_ID": "", "BK_APP_TENANT_ID": "app-tenant"},
			want: apigateway.Config{TenantID: "system"},
		},
		{
			name: "gateway name in url template",
			env:  map[string]string{"BK_API_URL_TMPL": "http://{gateway_name}.apigw.example.com"},
			want: apigateway.Config{Endpoint: "http://bk-apigateway.apigw.example.com/prod"},
		},
		{
			name: "app not deployed on paas",
			env: map[string]string{
				"BK_APP_CODE":      "app",
				"BKPAAS_APP_ID":    "paas-app",
				"BK_APP_SECRET":    "secret",
				"SECRET_KEY":       "legacy-secret",
				"BK_APP_TENANT_ID": "app-tenant",
			},
			want: apigateway.Config{AppCode: "app", AppSecret: "secret", TenantID: "app-tenant"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			setenv(t, tc.env)
			if got := apigateway.ConfigFromEnv(); got != tc.want {
				t.Errorf("config = %+v, want %+v", got, tc.want)
			}
		})
	}
}
