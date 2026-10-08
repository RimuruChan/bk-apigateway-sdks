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
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/define"
)

// VERSION for resource definitions
const VERSION = "1.0.8"

// Client for bkapi bk_apigateway
type Client struct {
	define.BkApiClient
}

// New bk_apigateway client
func New(configProvider define.ClientConfigProvider, opts ...define.BkApiClientOption) (*Client, error) {
	client, err := bkapi.NewBkApiClient("bk-apigateway", configProvider, opts...)
	if err != nil {
		return nil, err
	}

	return &Client{BkApiClient: client}, nil
}

// AddRelatedApps for bkapi resource add_related_apps
// 添加网关关联应用
func (c *Client) AddRelatedApps(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "add_related_apps",
		Method: "POST",
		Path:   "/api/v1/apis/{api_name}/related-apps/",
	}, opts...)
}

// ApplyPermissions for bkapi resource apply_permissions
// 申请网关API访问权限
func (c *Client) ApplyPermissions(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "apply_permissions",
		Method: "POST",
		Path:   "/api/v1/apis/{api_name}/permissions/apply/",
	}, opts...)
}

// CreateResourceVersion for bkapi resource create_resource_version
// 创建资源版本
func (c *Client) CreateResourceVersion(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "create_resource_version",
		Method: "POST",
		Path:   "/api/v1/apis/{api_name}/resource_versions/",
	}, opts...)
}

// GenerateSdk for bkapi resource generate_sdk
// 生成 SDK
func (c *Client) GenerateSdk(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "generate_sdk",
		Method: "POST",
		Path:   "/api/v1/apis/{api_name}/sdk/",
	}, opts...)
}

// GetApigwPublicKey for bkapi resource get_apigw_public_key
// 获取网关公钥
func (c *Client) GetApigwPublicKey(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "get_apigw_public_key",
		Method: "GET",
		Path:   "/api/v1/apis/{api_name}/public_key/",
	}, opts...)
}

// GetApis for bkapi resource get_apis
// 查询网关
func (c *Client) GetApis(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "get_apis",
		Method: "GET",
		Path:   "/api/v1/apis/",
	}, opts...)
}

// GetLatestResourceVersion for bkapi resource get_latest_resource_version
// 获取网关最新版本
func (c *Client) GetLatestResourceVersion(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "get_latest_resource_version",
		Method: "GET",
		Path:   "/api/v1/apis/{api_name}/resource_versions/latest/",
	}, opts...)
}

// GetMicroGatewayAppPermissions for bkapi resource get_micro_gateway_app_permissions
// 获取微网关应用权限
func (c *Client) GetMicroGatewayAppPermissions(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "get_micro_gateway_app_permissions",
		Method: "GET",
		Path:   "/api/v1/edge-controller/micro-gateway/{instance_id}/permissions/",
	}, opts...)
}

// GetMicroGatewayInfo for bkapi resource get_micro_gateway_info
// 获取微网关信息
func (c *Client) GetMicroGatewayInfo(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "get_micro_gateway_info",
		Method: "GET",
		Path:   "/api/v1/edge-controller/micro-gateway/{instance_id}/gateway/",
	}, opts...)
}

// GetMicroGatewayNewestGatewayPermissions for bkapi resource get_micro_gateway_newest_gateway_permissions
// 获取微网关新添加的网关权限
func (c *Client) GetMicroGatewayNewestGatewayPermissions(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "get_micro_gateway_newest_gateway_permissions",
		Method: "GET",
		Path:   "/api/v1/edge-controller/micro-gateway/{instance_id}/permissions/gateway/newest/",
	}, opts...)
}

// GetMicroGatewayNewestResourcePermissions for bkapi resource get_micro_gateway_newest_resource_permissions
// 获取微网关新添加的网关权限
func (c *Client) GetMicroGatewayNewestResourcePermissions(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "get_micro_gateway_newest_resource_permissions",
		Method: "GET",
		Path:   "/api/v1/edge-controller/micro-gateway/{instance_id}/permissions/resource/newest/",
	}, opts...)
}

// GetReleasedResources for bkapi resource get_released_resources
// 查询已发布资源列表
func (c *Client) GetReleasedResources(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "get_released_resources",
		Method: "GET",
		Path:   "/api/v1/apis/{api_name}/released/stages/{stage_name}/resources/",
	}, opts...)
}

// GetStages for bkapi resource get_stages
// 查询环境
func (c *Client) GetStages(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "get_stages",
		Method: "GET",
		Path:   "/api/v1/apis/{api_name}/stages/",
	}, opts...)
}

// GetStagesWithResourceVersion for bkapi resource get_stages_with_resource_version
// 查询网关环境资源版本
func (c *Client) GetStagesWithResourceVersion(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "get_stages_with_resource_version",
		Method: "GET",
		Path:   "/api/v1/apis/{api_name}/stages/with-resource-version/",
	}, opts...)
}

// GrantPermissions for bkapi resource grant_permissions
// 网关为应用主动授权
func (c *Client) GrantPermissions(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "grant_permissions",
		Method: "POST",
		Path:   "/api/v1/apis/{api_name}/permissions/grant/",
	}, opts...)
}

// ImportResourceDocsByArchive for bkapi resource import_resource_docs_by_archive
// 通过文档归档文件导入资源文档
func (c *Client) ImportResourceDocsByArchive(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "import_resource_docs_by_archive",
		Method: "POST",
		Path:   "/api/v1/apis/{api_name}/resource-docs/import/by-archive/",
	}, opts...)
}

// ImportResourceDocsBySwagger for bkapi resource import_resource_docs_by_swagger
// 通过 Swagger 格式导入文档
func (c *Client) ImportResourceDocsBySwagger(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "import_resource_docs_by_swagger",
		Method: "POST",
		Path:   "/api/v1/apis/{api_name}/resource-docs/import/by-swagger/",
	}, opts...)
}

// Release for bkapi resource release
// 发布版本
func (c *Client) Release(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "release",
		Method: "POST",
		Path:   "/api/v1/apis/{api_name}/resource_versions/release/",
	}, opts...)
}

// RevokePermissions for bkapi resource revoke_permissions
// 回收应用访问网关 API 的权限
func (c *Client) RevokePermissions(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "revoke_permissions",
		Method: "DELETE",
		Path:   "/api/v1/apis/{api_name}/permissions/revoke/",
	}, opts...)
}

// SyncAccessStrategy for bkapi resource sync_access_strategy
// 同步策略
//
// Deprecated: access strategies have been replaced by plugins, and the api has been removed from apigw.
func (c *Client) SyncAccessStrategy(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "sync_access_strategy",
		Method: "POST",
		Path:   "/api/v1/apis/{api_name}/access_strategies/sync/",
	}, opts...)
}

// SyncAPI for bkapi resource sync_api
// 同步网关
func (c *Client) SyncAPI(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "sync_api",
		Method: "POST",
		Path:   "/api/v1/apis/{api_name}/sync/",
	}, opts...)
}

// SyncResources for bkapi resource sync_resources
// 同步资源
func (c *Client) SyncResources(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "sync_resources",
		Method: "POST",
		Path:   "/api/v1/apis/{api_name}/resources/sync/",
	}, opts...)
}

// SyncStage for bkapi resource sync_stage
// 同步环境
func (c *Client) SyncStage(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "sync_stage",
		Method: "POST",
		Path:   "/api/v1/apis/{api_name}/stages/sync/",
	}, opts...)
}

// SyncStageMcpServers for bkapi resource sync_stage
// 同步环境 MCP Servers
func (c *Client) SyncStageMcpServers(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "v2_sync_stage_mcp_servers",
		Method: "POST",
		Path:   "/api/v2/sync/gateways/{gateway_name}/stages/{stage_name}/mcp-servers/",
	}, opts...)
}

// UpdateMicroGatewayStatus for bkapi resource update_micro_gateway_status
// 更新微网关实例状态
func (c *Client) UpdateMicroGatewayStatus(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "update_micro_gateway_status",
		Method: "PUT",
		Path:   "/api/v1/edge-controller/micro-gateway/{instance_id}/status/",
	}, opts...)
}

// V2OpenApplyGatewayPermission for bkapi resource v2_open_apply_gateway_permission
// 申请网关API访问权限
func (c *Client) V2OpenApplyGatewayPermission(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "v2_open_apply_gateway_permission",
		Method: "POST",
		Path:   "/api/v2/open/gateways/{gateway_name}/permissions/apply/",
	}, opts...)
}

// V2SyncAddRelatedApps for bkapi resource v2_sync_add_related_apps
// 添加网关关联应用
func (c *Client) V2SyncAddRelatedApps(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "v2_sync_add_related_apps",
		Method: "POST",
		Path:   "/api/v2/sync/gateways/{gateway_name}/related-apps/",
	}, opts...)
}

// V2SyncCreateResourceVersion for bkapi resource v2_sync_create_resource_version
// 创建资源版本
func (c *Client) V2SyncCreateResourceVersion(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "v2_sync_create_resource_version",
		Method: "POST",
		Path:   "/api/v2/sync/gateways/{gateway_name}/resource_versions/",
	}, opts...)
}

// V2SyncGateway for bkapi resource v2_sync_gateway
// 同步网关
func (c *Client) V2SyncGateway(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "v2_sync_gateway",
		Method: "POST",
		Path:   "/api/v2/sync/gateways/{gateway_name}/",
	}, opts...)
}

// V2SyncGetGatewayPublicKey for bkapi resource v2_sync_get_gateway_public_key_new
// 获取网关公钥
func (c *Client) V2SyncGetGatewayPublicKey(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "v2_sync_get_gateway_public_key_new",
		Method: "GET",
		Path:   "/api/v2/sync/gateways/{gateway_name}/public_key/",
	}, opts...)
}

// V2SyncGetResourceVersionLatest for bkapi resource v2_sync_get_resource_version_latest
// 获取网关最新版本
func (c *Client) V2SyncGetResourceVersionLatest(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "v2_sync_get_resource_version_latest",
		Method: "GET",
		Path:   "/api/v2/sync/gateways/{gateway_name}/resource_versions/latest/",
	}, opts...)
}

// V2SyncGrantPermission for bkapi resource v2_sync_grant_permission
// 网关为应用主动授权
func (c *Client) V2SyncGrantPermission(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "v2_sync_grant_permission",
		Method: "POST",
		Path:   "/api/v2/sync/gateways/{gateway_name}/permissions/grant/",
	}, opts...)
}

// V2SyncListResourceVersions for bkapi resource v2_sync_list_resource_versions
// 查询资源版本列表
func (c *Client) V2SyncListResourceVersions(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "v2_sync_list_resource_versions",
		Method: "GET",
		Path:   "/api/v2/sync/gateways/{gateway_name}/resource_versions/",
	}, opts...)
}

// V2SyncRelease for bkapi resource v2_sync_release
// 发布版本
func (c *Client) V2SyncRelease(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "v2_sync_release",
		Method: "POST",
		Path:   "/api/v2/sync/gateways/{gateway_name}/resource_versions/release/",
	}, opts...)
}

// V2SyncResourceDoc for bkapi resource v2_sync_resource_doc
// 通过文档归档文件导入资源文档
func (c *Client) V2SyncResourceDoc(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "v2_sync_resource_doc",
		Method: "POST",
		Path:   "/api/v2/sync/gateways/{gateway_name}/resource-docs/",
	}, opts...)
}

// V2SyncResources for bkapi resource v2_sync_resources
// 同步资源
func (c *Client) V2SyncResources(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "v2_sync_resources",
		Method: "POST",
		Path:   "/api/v2/sync/gateways/{gateway_name}/resources/",
	}, opts...)
}

// V2SyncStages for bkapi resource v2_sync_stages
// 同步环境
func (c *Client) V2SyncStages(opts ...define.OperationOption) define.Operation {
	return c.BkApiClient.NewOperation(bkapi.OperationConfig{
		Name:   "v2_sync_stages",
		Method: "POST",
		Path:   "/api/v2/sync/gateways/{gateway_name}/stages/",
	}, opts...)
}
