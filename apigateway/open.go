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
	"context"
	"net/url"
	"strconv"
)

// The /api/v2/open/ APIs query the gateways and MCP Servers, and apply for permissions.

// OpenListGateways 获取网关列表
func (c *Client) OpenListGateways(ctx context.Context, query url.Values, result any) error {
	return c.send(ctx, c.get("/api/v2/open/gateways/", query), result)
}

// OpenLookupGateways 按 ID 或名称批量查询公开且已发布的网关
func (c *Client) OpenLookupGateways(ctx context.Context, query url.Values, result any) error {
	return c.send(ctx, c.get("/api/v2/open/gateways/-/lookup/", query), result)
}

// OpenGetGateway 获取网关
func (c *Client) OpenGetGateway(ctx context.Context, gatewayName string, query url.Values, result any) error {
	req := c.get("/api/v2/open/gateways/:gateway_name/", query).Param("gateway_name", gatewayName)
	return c.send(ctx, req, result)
}

// OpenGetGatewayPublicKey 获取网关公钥
func (c *Client) OpenGetGatewayPublicKey(ctx context.Context, gatewayName string, query url.Values, result any) error {
	req := c.get("/api/v2/open/gateways/:gateway_name/public_key/", query).Param("gateway_name", gatewayName)
	return c.send(ctx, req, result)
}

// OpenListGatewayResources 获取网关资源列表
func (c *Client) OpenListGatewayResources(ctx context.Context, gatewayName string, query url.Values, result any) error {
	req := c.get("/api/v2/open/gateways/:gateway_name/resources/", query).Param("gateway_name", gatewayName)
	return c.send(ctx, req, result)
}

// OpenRetrieveGatewayAPIDetails 获取网关资源详情
func (c *Client) OpenRetrieveGatewayAPIDetails(
	ctx context.Context, gatewayName, resourceName string, query url.Values, result any,
) error {
	req := c.get("/api/v2/open/gateways/:gateway_name/resources/:resource_name/", query)
	req.Param("gateway_name", gatewayName).Param("resource_name", resourceName)
	return c.send(ctx, req, result)
}

// OpenGetReleasedResources 查询网关指定环境下已发布的资源列表
func (c *Client) OpenGetReleasedResources(
	ctx context.Context, gatewayName, stageName string, query url.Values, result any,
) error {
	req := c.get("/api/v2/open/gateways/:gateway_name/released/stages/:stage_name/resources/", query)
	req.Param("gateway_name", gatewayName).Param("stage_name", stageName)
	return c.send(ctx, req, result)
}

// OpenGetReleasedResource 查询网关指定环境下已发布的资源详情
func (c *Client) OpenGetReleasedResource(
	ctx context.Context, gatewayName, stageName, resourceName string, query url.Values, result any,
) error {
	req := c.get("/api/v2/open/gateways/:gateway_name/released/stages/:stage_name/resources/:resource_name/", query)
	req.Param("gateway_name", gatewayName).Param("stage_name", stageName).Param("resource_name", resourceName)
	return c.send(ctx, req, result)
}

// OpenApplyGatewayPermission 申请网关权限
func (c *Client) OpenApplyGatewayPermission(ctx context.Context, gatewayName string, body, result any) error {
	req := c.post("/api/v2/open/gateways/:gateway_name/permissions/apply/", body).Param("gateway_name", gatewayName)
	return c.send(ctx, req, result)
}

// OpenListMCPServer 获取公开的 MCP Server 列表
func (c *Client) OpenListMCPServer(ctx context.Context, query url.Values, result any) error {
	return c.send(ctx, c.get("/api/v2/open/mcp-servers/", query), result)
}

// OpenListMCPServerCategories 获取 MCP Server 分类列表
func (c *Client) OpenListMCPServerCategories(ctx context.Context, query url.Values, result any) error {
	return c.send(ctx, c.get("/api/v2/open/mcp-servers/categories/", query), result)
}

// OpenBatchQueryMCPServers 按 ID 或名称批量查询 MCP Server 的展示名称、描述和分类
func (c *Client) OpenBatchQueryMCPServers(ctx context.Context, body, result any) error {
	return c.send(ctx, c.post("/api/v2/open/mcp-servers/batch-query/", body), result)
}

// OpenRetrieveMCPServer 获取 MCP Server 详情，需要用户认证
func (c *Client) OpenRetrieveMCPServer(ctx context.Context, mcpServerID int, query url.Values, result any) error {
	req := c.get("/api/v2/open/mcp-servers/:mcp_server_id/", query).Param("mcp_server_id", strconv.Itoa(mcpServerID))
	return c.send(ctx, req, result)
}

// OpenListMCPServerPermissions 获取 MCP Server 的权限列表
func (c *Client) OpenListMCPServerPermissions(
	ctx context.Context, mcpServerID int, query url.Values, result any,
) error {
	req := c.get("/api/v2/open/mcp-servers/:mcp_server_id/permissions/", query)
	return c.send(ctx, req.Param("mcp_server_id", strconv.Itoa(mcpServerID)), result)
}

// OpenListMCPServerAppPermissions 获取应用的 MCP Server 权限列表
func (c *Client) OpenListMCPServerAppPermissions(ctx context.Context, query url.Values, result any) error {
	return c.send(ctx, c.get("/api/v2/open/mcp-servers/permissions/", query), result)
}

// OpenMCPServerAppPermissionsApply 为应用申请 MCP Server 权限
func (c *Client) OpenMCPServerAppPermissionsApply(ctx context.Context, body, result any) error {
	return c.send(ctx, c.post("/api/v2/open/mcp-servers/permissions/apply/", body), result)
}

// OpenListMCPServerAppPermissionsApplyRecords 获取应用的 MCP Server 权限申请记录列表
func (c *Client) OpenListMCPServerAppPermissionsApplyRecords(ctx context.Context, query url.Values, result any) error {
	return c.send(ctx, c.get("/api/v2/open/mcp-servers/permissions/apply-records/", query), result)
}

// OpenLookupMCPServerAppPermissionApplyRecords 按 ID 批量查询应用的 MCP Server 权限申请记录
func (c *Client) OpenLookupMCPServerAppPermissionApplyRecords(ctx context.Context, query url.Values, result any) error {
	return c.send(ctx, c.get("/api/v2/open/mcp-servers/permissions/apply-records/-/lookup/", query), result)
}

// OpenListUserMCPServer 获取当前用户的 MCP Server 列表，需要用户认证
func (c *Client) OpenListUserMCPServer(ctx context.Context, query url.Values, result any) error {
	return c.send(ctx, c.get("/api/v2/open/user/mcp-servers/", query), result)
}

// OpenToolsTimeGetDatetime 获取当前时间
func (c *Client) OpenToolsTimeGetDatetime(ctx context.Context, query url.Values, result any) error {
	return c.send(ctx, c.get("/api/v2/open/tools/time/get_datetime/", query), result)
}

// OpenToolsTimeGetCurrentUnixTimestamp 获取当前时间戳
func (c *Client) OpenToolsTimeGetCurrentUnixTimestamp(ctx context.Context, query url.Values, result any) error {
	return c.send(ctx, c.get("/api/v2/open/tools/time/get_current_unix_timestamp/", query), result)
}

// OpenToolsTimeParseDatetimeStrToTimestamp 将时间字符串转换为时间戳
func (c *Client) OpenToolsTimeParseDatetimeStrToTimestamp(ctx context.Context, body, result any) error {
	return c.send(ctx, c.post("/api/v2/open/tools/time/parse_datetime_str_to_timestamp/", body), result)
}

// OpenToolsLogQueryByRequestID 根据 request_id 查询日志
func (c *Client) OpenToolsLogQueryByRequestID(ctx context.Context, query url.Values, result any) error {
	return c.send(ctx, c.get("/api/v2/open/tools/log/query_by_request_id/", query), result)
}

// OpenOAuthProtectedResource 获取 OAuth 保护资源元数据，result 接收完整的响应
func (c *Client) OpenOAuthProtectedResource(ctx context.Context, query url.Values, result any) error {
	return c.send(ctx, c.get("/api/v2/open/.well-known/oauth-protected-resource", query), result)
}
