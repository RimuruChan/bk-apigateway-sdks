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
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/TencentBlueKing/bk-apigateway-sdks/v2/apigateway"
)

// TestAPIs checks the request of every API, which is formatted as "METHOD PATH?QUERY BODY".
func TestAPIs(t *testing.T) {
	var got string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		got = strings.TrimSpace(r.Method + " " + r.URL.RequestURI() + " " + strings.TrimSpace(string(body)))
	}))
	defer server.Close()

	c := newClient(t, apigateway.Config{Endpoint: server.URL})
	ctx := context.Background()
	q := url.Values{"q": {"1"}}
	b := map[string]int{"b": 1}
	check := func(err error, want string) {
		t.Helper()
		if err != nil || got != want {
			t.Errorf("got %q, error %v, want %q", got, err, want)
		}
	}

	// sync
	check(c.SyncGateway(ctx, "g", b, nil), `POST /api/v2/sync/gateways/g/ {"b":1}`)
	check(c.SyncGetGatewayPublicKey(ctx, "g", q, nil), `GET /api/v2/sync/gateways/g/public_key/?q=1`)
	check(c.SyncAddRelatedApps(ctx, "g", b), `POST /api/v2/sync/gateways/g/related-apps/ {"b":1}`)
	check(c.SyncStages(ctx, "g", b, nil), `POST /api/v2/sync/gateways/g/stages/ {"b":1}`)
	check(c.SyncStageMCPServers(ctx, "g", "s", b, nil), `POST /api/v2/sync/gateways/g/stages/s/mcp-servers/ {"b":1}`)
	check(c.SyncResources(ctx, "g", b, nil), `POST /api/v2/sync/gateways/g/resources/ {"b":1}`)
	check(c.SyncGenerateSDK(ctx, "g", b, nil), `POST /api/v2/sync/gateways/g/sdks/ {"b":1}`)
	check(c.SyncListPermissions(ctx, "g", q, nil), `GET /api/v2/sync/gateways/g/permissions/?q=1`)
	check(c.SyncGrantPermission(ctx, "g", b), `POST /api/v2/sync/gateways/g/permissions/grant/ {"b":1}`)
	check(c.SyncRevokePermission(ctx, "g", b), `DELETE /api/v2/sync/gateways/g/permissions/revoke/ {"b":1}`)
	check(c.SyncCreateResourceVersion(ctx, "g", b, nil), `POST /api/v2/sync/gateways/g/resource_versions/ {"b":1}`)
	check(c.SyncListResourceVersions(ctx, "g", q, nil), `GET /api/v2/sync/gateways/g/resource_versions/?q=1`)
	check(c.SyncLookupResourceVersions(ctx, "g", q, nil), `GET /api/v2/sync/gateways/g/resource_versions/-/lookup/?q=1`)
	check(c.SyncGetResourceVersionLatest(ctx, "g", q, nil), `GET /api/v2/sync/gateways/g/resource_versions/latest/?q=1`)
	check(c.SyncRelease(ctx, "g", b, nil), `POST /api/v2/sync/gateways/g/resource_versions/release/ {"b":1}`)
	// SyncResourceDoc is checked by TestSyncResourceDoc

	// open
	check(c.OpenListGateways(ctx, q, nil), `GET /api/v2/open/gateways/?q=1`)
	check(c.OpenLookupGateways(ctx, q, nil), `GET /api/v2/open/gateways/-/lookup/?q=1`)
	check(c.OpenGetGateway(ctx, "g", q, nil), `GET /api/v2/open/gateways/g/?q=1`)
	check(c.OpenGetGatewayPublicKey(ctx, "g", q, nil), `GET /api/v2/open/gateways/g/public_key/?q=1`)
	check(c.OpenListGatewayResources(ctx, "g", q, nil), `GET /api/v2/open/gateways/g/resources/?q=1`)
	check(c.OpenRetrieveGatewayAPIDetails(ctx, "g", "r", q, nil), `GET /api/v2/open/gateways/g/resources/r/?q=1`)
	check(c.OpenGetReleasedResources(ctx, "g", "s", q, nil),
		`GET /api/v2/open/gateways/g/released/stages/s/resources/?q=1`)
	check(c.OpenGetReleasedResource(ctx, "g", "s", "r", q, nil),
		`GET /api/v2/open/gateways/g/released/stages/s/resources/r/?q=1`)
	check(c.OpenApplyGatewayPermission(ctx, "g", b, nil), `POST /api/v2/open/gateways/g/permissions/apply/ {"b":1}`)
	check(c.OpenListMCPServer(ctx, q, nil), `GET /api/v2/open/mcp-servers/?q=1`)
	check(c.OpenListMCPServerCategories(ctx, q, nil), `GET /api/v2/open/mcp-servers/categories/?q=1`)
	check(c.OpenBatchQueryMCPServers(ctx, b, nil), `POST /api/v2/open/mcp-servers/batch-query/ {"b":1}`)
	check(c.OpenRetrieveMCPServer(ctx, 1, q, nil), `GET /api/v2/open/mcp-servers/1/?q=1`)
	check(c.OpenListMCPServerPermissions(ctx, 1, q, nil), `GET /api/v2/open/mcp-servers/1/permissions/?q=1`)
	check(c.OpenListMCPServerAppPermissions(ctx, q, nil), `GET /api/v2/open/mcp-servers/permissions/?q=1`)
	check(c.OpenMCPServerAppPermissionsApply(ctx, b, nil), `POST /api/v2/open/mcp-servers/permissions/apply/ {"b":1}`)
	check(c.OpenListMCPServerAppPermissionsApplyRecords(ctx, q, nil),
		`GET /api/v2/open/mcp-servers/permissions/apply-records/?q=1`)
	check(c.OpenLookupMCPServerAppPermissionApplyRecords(ctx, q, nil),
		`GET /api/v2/open/mcp-servers/permissions/apply-records/-/lookup/?q=1`)
	check(c.OpenListUserMCPServer(ctx, q, nil), `GET /api/v2/open/user/mcp-servers/?q=1`)
	check(c.OpenToolsTimeGetDatetime(ctx, q, nil), `GET /api/v2/open/tools/time/get_datetime/?q=1`)
	check(c.OpenToolsTimeGetCurrentUnixTimestamp(ctx, q, nil),
		`GET /api/v2/open/tools/time/get_current_unix_timestamp/?q=1`)
	check(c.OpenToolsTimeParseDatetimeStrToTimestamp(ctx, b, nil),
		`POST /api/v2/open/tools/time/parse_datetime_str_to_timestamp/ {"b":1}`)
	check(c.OpenToolsLogQueryByRequestID(ctx, q, nil), `GET /api/v2/open/tools/log/query_by_request_id/?q=1`)
	check(c.OpenOAuthProtectedResource(ctx, q, nil), `GET /api/v2/open/.well-known/oauth-protected-resource?q=1`)

	// inner
	check(c.InnerITSMCallback(ctx, b), `POST /api/v2/inner/itsm/callback/ {"b":1}`)
	check(c.InnerMonitorAlarmCallback(ctx, "a", q, b),
		`POST /api/v2/inner/monitor/alarm-types/a/callback/?q=1 {"b":1}`)
	check(c.InnerListAppAlarmRecords(ctx, "a", q, nil), `GET /api/v2/inner/apps/a/monitor/alarm-records/?q=1`)
	check(c.InnerListAppRequestLogs(ctx, "a", q, nil), `GET /api/v2/inner/apps/a/monitor/request-logs/?q=1`)
	check(c.InnerListGateways(ctx, q, nil), `GET /api/v2/inner/gateways/?q=1`)
	check(c.InnerLookupGateways(ctx, q, nil), `GET /api/v2/inner/gateways/-/lookup/?q=1`)
	check(c.InnerGetGateway(ctx, "g", q, nil), `GET /api/v2/inner/gateways/g/?q=1`)
	check(c.InnerDeleteGateway(ctx, "g"), `DELETE /api/v2/inner/gateways/g/`)
	check(c.InnerUpdateGatewayStatus(ctx, "g", b), `PUT /api/v2/inner/gateways/g/status/ {"b":1}`)
	check(c.InnerListGatewayReleasedResources(ctx, "g", q, nil),
		`GET /api/v2/inner/gateways/g/released-resources/?q=1`)
	check(c.InnerLookupGatewayReleasedResources(ctx, "g", q, nil),
		`GET /api/v2/inner/gateways/g/released-resources/-/lookup/?q=1`)
	check(c.InnerListGatewayPermissionResources(ctx, "g", q, nil),
		`GET /api/v2/inner/gateways/g/permissions/resources/?q=1`)
	check(c.InnerCheckIsAllowedApplyByGateway(ctx, "g", q, nil),
		`GET /api/v2/inner/gateways/g/permissions/app-permissions/allow-apply-by-gateway/?q=1`)
	check(c.InnerApplyGatewayResourcePermission(ctx, "g", b, nil),
		`POST /api/v2/inner/gateways/g/permissions/app-permissions/apply/ {"b":1}`)
	check(c.InnerRenewResourcePermission(ctx, b), `POST /api/v2/inner/gateways/permissions/renew/ {"b":1}`)
	check(c.InnerListAppResourcePermissions(ctx, q, nil), `GET /api/v2/inner/gateways/permissions/app-permissions/?q=1`)
	check(c.InnerListResourcePermissionApplyRecords(ctx, q, nil),
		`GET /api/v2/inner/gateways/permissions/apply-records/?q=1`)
	check(c.InnerRetrieveResourcePermissionApplyRecord(ctx, 1, q, nil),
		`GET /api/v2/inner/gateways/permissions/apply-records/1/?q=1`)
	check(c.InnerDeleteResourcePermissionApplyRecord(ctx, 1, q),
		`DELETE /api/v2/inner/gateways/permissions/apply-records/1/?q=1`)
	check(c.InnerCancelResourcePermissionApplyRecord(ctx, 1, b),
		`POST /api/v2/inner/gateways/permissions/apply-records/1/cancel/ {"b":1}`)
	check(c.InnerListESBSystems(ctx, q, nil), `GET /api/v2/inner/esb/systems/?q=1`)
	check(c.InnerListESBSystemPermissionComponents(ctx, 1, q, nil),
		`GET /api/v2/inner/esb/systems/1/permissions/components/?q=1`)
	check(c.InnerApplyESBSystemComponentPermissions(ctx, 1, b, nil),
		`POST /api/v2/inner/esb/systems/1/permissions/apply/ {"b":1}`)
	check(c.InnerRenewESBComponentPermissions(ctx, b), `POST /api/v2/inner/esb/systems/permissions/renew/ {"b":1}`)
	check(c.InnerListAppESBComponentPermissions(ctx, q, nil),
		`GET /api/v2/inner/esb/systems/permissions/app-permissions/?q=1`)
	check(c.InnerListAppESBComponentPermissionApplyRecords(ctx, q, nil),
		`GET /api/v2/inner/esb/systems/permissions/apply-records/?q=1`)
	check(c.InnerGetAppESBComponentPermissionApplyRecord(ctx, 1, q, nil),
		`GET /api/v2/inner/esb/systems/permissions/apply-records/1/?q=1`)
	check(c.InnerListMCPServerPermissions(ctx, q, nil), `GET /api/v2/inner/mcp-server/permissions/?q=1`)
	check(c.InnerApplyMCPServerPermission(ctx, b, nil), `POST /api/v2/inner/mcp-server/permissions/apply/ {"b":1}`)
	check(c.InnerListMCPServerAppPermissions(ctx, q, nil),
		`GET /api/v2/inner/mcp-server/permissions/app-permissions/?q=1`)
	check(c.InnerListMCPServerPermissionApplyRecords(ctx, q, nil),
		`GET /api/v2/inner/mcp-server/permissions/apply-records/?q=1`)
	check(c.InnerRetrieveMCPServerPermissionApplyRecord(ctx, 1, q, nil),
		`GET /api/v2/inner/mcp-server/permissions/apply-records/1/?q=1`)
	check(c.InnerDeleteMCPServerPermissionApplyRecord(ctx, 1, q),
		`DELETE /api/v2/inner/mcp-server/permissions/apply-records/1/?q=1`)
	check(c.InnerCancelMCPServerPermissionApplyRecord(ctx, 1, b),
		`POST /api/v2/inner/mcp-server/permissions/apply-records/1/cancel/ {"b":1}`)
	check(c.InnerListMCPServer(ctx, q, nil), `GET /api/v2/inner/mcp-servers/?q=1`)
	check(c.InnerLookupMCPServers(ctx, q, nil), `GET /api/v2/inner/mcp-servers/-/lookup/?q=1`)
	check(c.InnerListOAuth2MCPServerScopes(ctx, q, nil), `GET /api/v2/inner/oauth2/client-scopes/mcp-servers/?q=1`)
	check(c.InnerListOAuth2ResourceScopes(ctx, q, nil), `GET /api/v2/inner/oauth2/client-scopes/resources/?q=1`)
}
