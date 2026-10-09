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
	"os"
	"strings"
)

// ConfigFromEnv returns the Config of an app deployed on BlueKing PaaS, which calls the prod stage of bk-apigateway.
//
//   - Endpoint: BK_API_URL_TMPL, such as http://bkapi.example.com/api/{api_name}/, where {api_name} or
//     {gateway_name} is the gateway name
//   - AppCode: BK_APP_CODE, BKPAAS_APP_ID or APP_CODE
//   - AppSecret: BK_APP_SECRET, BKPAAS_APP_SECRET or SECRET_KEY
//   - TenantID: BKPAAS_APP_TENANT_ID, which is empty for global tenant apps that belong to the system tenant,
//     or BK_APP_TENANT_ID for apps not deployed on PaaS
func ConfigFromEnv() Config {
	var endpoint string
	if tmpl := os.Getenv("BK_API_URL_TMPL"); tmpl != "" {
		endpoint = strings.NewReplacer("{api_name}", "bk-apigateway", "{gateway_name}", "bk-apigateway").Replace(tmpl)
		endpoint = strings.TrimSuffix(endpoint, "/") + "/prod"
	}

	tenantID, ok := os.LookupEnv("BKPAAS_APP_TENANT_ID")
	if ok {
		tenantID = cmp.Or(tenantID, "system")
	} else {
		tenantID = os.Getenv("BK_APP_TENANT_ID")
	}

	return Config{
		Endpoint:  endpoint,
		AppCode:   cmp.Or(os.Getenv("BK_APP_CODE"), os.Getenv("BKPAAS_APP_ID"), os.Getenv("APP_CODE")),
		AppSecret: cmp.Or(os.Getenv("BK_APP_SECRET"), os.Getenv("BKPAAS_APP_SECRET"), os.Getenv("SECRET_KEY")),
		TenantID:  tenantID,
	}
}
