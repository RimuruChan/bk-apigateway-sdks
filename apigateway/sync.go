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
	"io"
	"net/url"
)

// The /api/v2/sync/ APIs sync the gateway definitions, and are called by the gateway's own app or related apps.

// SyncGateway 同步网关，网关不存在时创建
func (c *Client) SyncGateway(ctx context.Context, gatewayName string, body, result any) error {
	req := c.post("/api/v2/sync/gateways/:gateway_name/", body).Param("gateway_name", gatewayName)
	return c.send(ctx, req, result)
}

// SyncGetGatewayPublicKey 获取网关公钥
func (c *Client) SyncGetGatewayPublicKey(ctx context.Context, gatewayName string, query url.Values, result any) error {
	req := c.get("/api/v2/sync/gateways/:gateway_name/public_key/", query).Param("gateway_name", gatewayName)
	return c.send(ctx, req, result)
}

// SyncAddRelatedApps 添加网关的关联应用
func (c *Client) SyncAddRelatedApps(ctx context.Context, gatewayName string, body any) error {
	req := c.post("/api/v2/sync/gateways/:gateway_name/related-apps/", body).Param("gateway_name", gatewayName)
	return c.send(ctx, req, nil)
}

// SyncStages 同步网关环境
func (c *Client) SyncStages(ctx context.Context, gatewayName string, body, result any) error {
	req := c.post("/api/v2/sync/gateways/:gateway_name/stages/", body).Param("gateway_name", gatewayName)
	return c.send(ctx, req, result)
}

// SyncStageMCPServers 同步网关环境的 MCP Server
func (c *Client) SyncStageMCPServers(ctx context.Context, gatewayName, stageName string, body, result any) error {
	req := c.post("/api/v2/sync/gateways/:gateway_name/stages/:stage_name/mcp-servers/", body)
	req.Param("gateway_name", gatewayName).Param("stage_name", stageName)
	return c.send(ctx, req, result)
}

// SyncResources 同步网关资源
func (c *Client) SyncResources(ctx context.Context, gatewayName string, body, result any) error {
	req := c.post("/api/v2/sync/gateways/:gateway_name/resources/", body).Param("gateway_name", gatewayName)
	return c.send(ctx, req, result)
}

// SyncResourceDoc 同步网关资源文档，archive 为 zip 或 tgz 格式的文档归档，是 io.ReadCloser 时读取后会被关闭
func (c *Client) SyncResourceDoc(ctx context.Context, gatewayName string, archive io.Reader) error {
	req := c.post("/api/v2/sync/gateways/:gateway_name/resource-docs/", nil).Param("gateway_name", gatewayName)
	return c.send(ctx, req.File("file", archive), nil)
}

// SyncGenerateSDK 生成网关 SDK
func (c *Client) SyncGenerateSDK(ctx context.Context, gatewayName string, body, result any) error {
	req := c.post("/api/v2/sync/gateways/:gateway_name/sdks/", body).Param("gateway_name", gatewayName)
	return c.send(ctx, req, result)
}

// SyncListPermissions 获取网关权限列表
func (c *Client) SyncListPermissions(ctx context.Context, gatewayName string, query url.Values, result any) error {
	req := c.get("/api/v2/sync/gateways/:gateway_name/permissions/", query).Param("gateway_name", gatewayName)
	return c.send(ctx, req, result)
}

// SyncGrantPermission 为应用主动授权，支持按网关维度或按资源维度
func (c *Client) SyncGrantPermission(ctx context.Context, gatewayName string, body any) error {
	req := c.post("/api/v2/sync/gateways/:gateway_name/permissions/grant/", body).Param("gateway_name", gatewayName)
	return c.send(ctx, req, nil)
}

// SyncRevokePermission 回收应用访问网关 API 的权限，支持按网关维度或按资源维度
func (c *Client) SyncRevokePermission(ctx context.Context, gatewayName string, body any) error {
	req := c.delete("/api/v2/sync/gateways/:gateway_name/permissions/revoke/", nil, body)
	return c.send(ctx, req.Param("gateway_name", gatewayName), nil)
}

// SyncCreateResourceVersion 创建网关资源版本
func (c *Client) SyncCreateResourceVersion(ctx context.Context, gatewayName string, body, result any) error {
	req := c.post("/api/v2/sync/gateways/:gateway_name/resource_versions/", body).Param("gateway_name", gatewayName)
	return c.send(ctx, req, result)
}

// SyncListResourceVersions 查询网关资源版本列表
func (c *Client) SyncListResourceVersions(ctx context.Context, gatewayName string, query url.Values, result any) error {
	req := c.get("/api/v2/sync/gateways/:gateway_name/resource_versions/", query).Param("gateway_name", gatewayName)
	return c.send(ctx, req, result)
}

// SyncLookupResourceVersions 按 ID 或版本号批量查询网关资源版本
func (c *Client) SyncLookupResourceVersions(
	ctx context.Context, gatewayName string, query url.Values, result any,
) error {
	req := c.get("/api/v2/sync/gateways/:gateway_name/resource_versions/-/lookup/", query)
	return c.send(ctx, req.Param("gateway_name", gatewayName), result)
}

// SyncGetResourceVersionLatest 获取网关最新的资源版本
func (c *Client) SyncGetResourceVersionLatest(
	ctx context.Context, gatewayName string, query url.Values, result any,
) error {
	req := c.get("/api/v2/sync/gateways/:gateway_name/resource_versions/latest/", query)
	return c.send(ctx, req.Param("gateway_name", gatewayName), result)
}

// SyncRelease 发布网关资源版本
func (c *Client) SyncRelease(ctx context.Context, gatewayName string, body, result any) error {
	req := c.post("/api/v2/sync/gateways/:gateway_name/resource_versions/release/", body)
	return c.send(ctx, req.Param("gateway_name", gatewayName), result)
}
