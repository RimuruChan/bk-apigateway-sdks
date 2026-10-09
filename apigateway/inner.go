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

// The /api/v2/inner/ APIs are called by BlueKing platforms, such as PaaS, which are granted the permissions.

// InnerITSMCallback ITSM 工单审批结果回调
func (c *Client) InnerITSMCallback(ctx context.Context, body any) error {
	return c.send(ctx, c.post("/api/v2/inner/itsm/callback/", body), nil)
}

// InnerMonitorAlarmCallback 监控告警回调，query 中的 token 用于校验回调
func (c *Client) InnerMonitorAlarmCallback(ctx context.Context, alarmType string, query url.Values, body any) error {
	req := withQuery(c.post("/api/v2/inner/monitor/alarm-types/:alarm_type/callback/", body), query)
	return c.send(ctx, req.Param("alarm_type", alarmType), nil)
}

// InnerListAppAlarmRecords 获取应用的告警记录列表
func (c *Client) InnerListAppAlarmRecords(ctx context.Context, appCode string, query url.Values, result any) error {
	req := c.get("/api/v2/inner/apps/:app_code/monitor/alarm-records/", query).Param("app_code", appCode)
	return c.send(ctx, req, result)
}

// InnerListAppRequestLogs 获取应用的调用流水日志列表
func (c *Client) InnerListAppRequestLogs(ctx context.Context, appCode string, query url.Values, result any) error {
	req := c.get("/api/v2/inner/apps/:app_code/monitor/request-logs/", query).Param("app_code", appCode)
	return c.send(ctx, req, result)
}

// InnerListGateways 获取网关列表
func (c *Client) InnerListGateways(ctx context.Context, query url.Values, result any) error {
	return c.send(ctx, c.get("/api/v2/inner/gateways/", query), result)
}

// InnerLookupGateways 按名称批量查询网关
func (c *Client) InnerLookupGateways(ctx context.Context, query url.Values, result any) error {
	return c.send(ctx, c.get("/api/v2/inner/gateways/-/lookup/", query), result)
}

// InnerGetGateway 获取网关
func (c *Client) InnerGetGateway(ctx context.Context, gatewayName string, query url.Values, result any) error {
	req := c.get("/api/v2/inner/gateways/:gateway_name/", query).Param("gateway_name", gatewayName)
	return c.send(ctx, req, result)
}

// InnerDeleteGateway 删除网关，仅支持已停用的 bp- 开头的网关
func (c *Client) InnerDeleteGateway(ctx context.Context, gatewayName string) error {
	req := c.delete("/api/v2/inner/gateways/:gateway_name/", nil, nil).Param("gateway_name", gatewayName)
	return c.send(ctx, req, nil)
}

// InnerUpdateGatewayStatus 启用或停用网关，仅支持 bp- 开头的网关
func (c *Client) InnerUpdateGatewayStatus(ctx context.Context, gatewayName string, body any) error {
	req := c.put("/api/v2/inner/gateways/:gateway_name/status/", body).Param("gateway_name", gatewayName)
	return c.send(ctx, req, nil)
}

// InnerListGatewayReleasedResources 获取网关已发布的资源，仅包含启用的环境
func (c *Client) InnerListGatewayReleasedResources(
	ctx context.Context, gatewayName string, query url.Values, result any,
) error {
	req := c.get("/api/v2/inner/gateways/:gateway_name/released-resources/", query)
	return c.send(ctx, req.Param("gateway_name", gatewayName), result)
}

// InnerLookupGatewayReleasedResources 按名称查询网关已发布的资源，仅包含启用的环境
func (c *Client) InnerLookupGatewayReleasedResources(
	ctx context.Context, gatewayName string, query url.Values, result any,
) error {
	req := c.get("/api/v2/inner/gateways/:gateway_name/released-resources/-/lookup/", query)
	return c.send(ctx, req.Param("gateway_name", gatewayName), result)
}

// InnerListGatewayPermissionResources 获取网关可申请权限的资源
func (c *Client) InnerListGatewayPermissionResources(
	ctx context.Context, gatewayName string, query url.Values, result any,
) error {
	req := c.get("/api/v2/inner/gateways/:gateway_name/permissions/resources/", query)
	return c.send(ctx, req.Param("gateway_name", gatewayName), result)
}

// InnerCheckIsAllowedApplyByGateway 检查是否允许按网关申请资源权限
func (c *Client) InnerCheckIsAllowedApplyByGateway(
	ctx context.Context, gatewayName string, query url.Values, result any,
) error {
	req := c.get("/api/v2/inner/gateways/:gateway_name/permissions/app-permissions/allow-apply-by-gateway/", query)
	return c.send(ctx, req.Param("gateway_name", gatewayName), result)
}

// InnerApplyGatewayResourcePermission 申请网关资源权限
func (c *Client) InnerApplyGatewayResourcePermission(
	ctx context.Context, gatewayName string, body, result any,
) error {
	req := c.post("/api/v2/inner/gateways/:gateway_name/permissions/app-permissions/apply/", body)
	return c.send(ctx, req.Param("gateway_name", gatewayName), result)
}

// InnerRenewResourcePermission 网关资源权限续期
func (c *Client) InnerRenewResourcePermission(ctx context.Context, body any) error {
	return c.send(ctx, c.post("/api/v2/inner/gateways/permissions/renew/", body), nil)
}

// InnerListAppResourcePermissions 获取应用已申请的网关资源权限列表
func (c *Client) InnerListAppResourcePermissions(ctx context.Context, query url.Values, result any) error {
	return c.send(ctx, c.get("/api/v2/inner/gateways/permissions/app-permissions/", query), result)
}

// InnerListResourcePermissionApplyRecords 获取网关资源权限申请记录列表
func (c *Client) InnerListResourcePermissionApplyRecords(ctx context.Context, query url.Values, result any) error {
	return c.send(ctx, c.get("/api/v2/inner/gateways/permissions/apply-records/", query), result)
}

// InnerRetrieveResourcePermissionApplyRecord 获取网关资源权限申请记录详情
func (c *Client) InnerRetrieveResourcePermissionApplyRecord(
	ctx context.Context, recordID int, query url.Values, result any,
) error {
	req := c.get("/api/v2/inner/gateways/permissions/apply-records/:record_id/", query)
	return c.send(ctx, req.Param("record_id", strconv.Itoa(recordID)), result)
}

// InnerDeleteResourcePermissionApplyRecord 删除已取消的网关资源权限申请记录
func (c *Client) InnerDeleteResourcePermissionApplyRecord(ctx context.Context, recordID int, query url.Values) error {
	req := c.delete("/api/v2/inner/gateways/permissions/apply-records/:record_id/", query, nil)
	return c.send(ctx, req.Param("record_id", strconv.Itoa(recordID)), nil)
}

// InnerCancelResourcePermissionApplyRecord 取消待审批的网关资源权限申请
func (c *Client) InnerCancelResourcePermissionApplyRecord(ctx context.Context, recordID int, body any) error {
	req := c.post("/api/v2/inner/gateways/permissions/apply-records/:record_id/cancel/", body)
	return c.send(ctx, req.Param("record_id", strconv.Itoa(recordID)), nil)
}

// InnerListESBSystems 获取组件系统列表
func (c *Client) InnerListESBSystems(ctx context.Context, query url.Values, result any) error {
	return c.send(ctx, c.get("/api/v2/inner/esb/systems/", query), result)
}

// InnerListESBSystemPermissionComponents 获取组件系统可申请权限的组件
func (c *Client) InnerListESBSystemPermissionComponents(
	ctx context.Context, systemID int, query url.Values, result any,
) error {
	req := c.get("/api/v2/inner/esb/systems/:system_id/permissions/components/", query)
	return c.send(ctx, req.Param("system_id", strconv.Itoa(systemID)), result)
}

// InnerApplyESBSystemComponentPermissions 申请组件系统的组件权限
func (c *Client) InnerApplyESBSystemComponentPermissions(
	ctx context.Context, systemID int, body, result any,
) error {
	req := c.post("/api/v2/inner/esb/systems/:system_id/permissions/apply/", body)
	return c.send(ctx, req.Param("system_id", strconv.Itoa(systemID)), result)
}

// InnerRenewESBComponentPermissions 组件权限续期
func (c *Client) InnerRenewESBComponentPermissions(ctx context.Context, body any) error {
	return c.send(ctx, c.post("/api/v2/inner/esb/systems/permissions/renew/", body), nil)
}

// InnerListAppESBComponentPermissions 获取应用已申请的组件权限列表
func (c *Client) InnerListAppESBComponentPermissions(ctx context.Context, query url.Values, result any) error {
	return c.send(ctx, c.get("/api/v2/inner/esb/systems/permissions/app-permissions/", query), result)
}

// InnerListAppESBComponentPermissionApplyRecords 获取应用的组件权限申请记录列表
func (c *Client) InnerListAppESBComponentPermissionApplyRecords(
	ctx context.Context, query url.Values, result any,
) error {
	return c.send(ctx, c.get("/api/v2/inner/esb/systems/permissions/apply-records/", query), result)
}

// InnerGetAppESBComponentPermissionApplyRecord 获取应用的组件权限申请记录详情
func (c *Client) InnerGetAppESBComponentPermissionApplyRecord(
	ctx context.Context, recordID int, query url.Values, result any,
) error {
	req := c.get("/api/v2/inner/esb/systems/permissions/apply-records/:record_id/", query)
	return c.send(ctx, req.Param("record_id", strconv.Itoa(recordID)), result)
}

// InnerListMCPServerPermissions 获取可申请权限的 MCP Server 列表
func (c *Client) InnerListMCPServerPermissions(ctx context.Context, query url.Values, result any) error {
	return c.send(ctx, c.get("/api/v2/inner/mcp-server/permissions/", query), result)
}

// InnerApplyMCPServerPermission 申请 MCP Server 权限，支持批量申请
func (c *Client) InnerApplyMCPServerPermission(ctx context.Context, body, result any) error {
	return c.send(ctx, c.post("/api/v2/inner/mcp-server/permissions/apply/", body), result)
}

// InnerListMCPServerAppPermissions 获取应用已申请的 MCP Server 权限列表
func (c *Client) InnerListMCPServerAppPermissions(ctx context.Context, query url.Values, result any) error {
	return c.send(ctx, c.get("/api/v2/inner/mcp-server/permissions/app-permissions/", query), result)
}

// InnerListMCPServerPermissionApplyRecords 获取 MCP Server 权限申请记录列表
func (c *Client) InnerListMCPServerPermissionApplyRecords(ctx context.Context, query url.Values, result any) error {
	return c.send(ctx, c.get("/api/v2/inner/mcp-server/permissions/apply-records/", query), result)
}

// InnerRetrieveMCPServerPermissionApplyRecord 获取 MCP Server 权限申请记录详情
func (c *Client) InnerRetrieveMCPServerPermissionApplyRecord(
	ctx context.Context, recordID int, query url.Values, result any,
) error {
	req := c.get("/api/v2/inner/mcp-server/permissions/apply-records/:record_id/", query)
	return c.send(ctx, req.Param("record_id", strconv.Itoa(recordID)), result)
}

// InnerDeleteMCPServerPermissionApplyRecord 删除已取消的 MCP Server 权限申请记录
func (c *Client) InnerDeleteMCPServerPermissionApplyRecord(ctx context.Context, recordID int, query url.Values) error {
	req := c.delete("/api/v2/inner/mcp-server/permissions/apply-records/:record_id/", query, nil)
	return c.send(ctx, req.Param("record_id", strconv.Itoa(recordID)), nil)
}

// InnerCancelMCPServerPermissionApplyRecord 取消待审批的 MCP Server 权限申请
func (c *Client) InnerCancelMCPServerPermissionApplyRecord(ctx context.Context, recordID int, body any) error {
	req := c.post("/api/v2/inner/mcp-server/permissions/apply-records/:record_id/cancel/", body)
	return c.send(ctx, req.Param("record_id", strconv.Itoa(recordID)), nil)
}

// InnerListMCPServer 获取全量的 MCP Server 列表
func (c *Client) InnerListMCPServer(ctx context.Context, query url.Values, result any) error {
	return c.send(ctx, c.get("/api/v2/inner/mcp-servers/", query), result)
}

// InnerLookupMCPServers 按 ID 或名称批量查询 MCP Server
func (c *Client) InnerLookupMCPServers(ctx context.Context, query url.Values, result any) error {
	return c.send(ctx, c.get("/api/v2/inner/mcp-servers/-/lookup/", query), result)
}

// InnerListOAuth2MCPServerScopes 获取 OAuth2 客户端可选的 MCP Server 范围
func (c *Client) InnerListOAuth2MCPServerScopes(ctx context.Context, query url.Values, result any) error {
	return c.send(ctx, c.get("/api/v2/inner/oauth2/client-scopes/mcp-servers/", query), result)
}

// InnerListOAuth2ResourceScopes 获取 OAuth2 客户端可选的 API 资源范围
func (c *Client) InnerListOAuth2ResourceScopes(ctx context.Context, query url.Values, result any) error {
	return c.send(ctx, c.get("/api/v2/inner/oauth2/client-scopes/resources/", query), result)
}
