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

package manager

import (
	"context"
	"fmt"
	"maps"
	"net/url"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/TencentBlueKing/bk-apigateway-sdks/v2/apigateway"
	"github.com/TencentBlueKing/bk-apigateway-sdks/v2/gin_contrib/util"
)

const (
	apiGatewayNamespace       = "apigateway"
	stagesNamespace           = "stages"
	grantPermissionsNamespace = "grant_permissions"
	applyPermissionsNamespace = "apply_permissions"
	resourceDocsNamespace     = "resource_docs"
	relatedAppsNamespace      = "related_apps"
)

// Manager is the manager of apigw, it helps to sync apigw configs and get apigw infomations.
type Manager struct {
	gatewayName string
	appCode     string
	definition  *Definition
	client      *apigateway.Client
}

// LoadDefinition will load the definition from the file.
func (m *Manager) LoadDefinition(path string) error {
	rendered, err := os.ReadFile(path)
	if err != nil {
		return errors.Wrapf(err, "failed to read %s", path)
	}
	definition, err := NewDefinitionFromYaml(rendered)
	if err != nil {
		return errors.Wrapf(err, "failed to parse %s", path)
	}

	m.definition = definition
	return nil
}

// GetDefinition return the definition.
func (m *Manager) GetDefinition() *Definition {
	return m.definition
}

// GetPublicKey fetch the public key info from apigw.
func (m *Manager) GetPublicKey(ctx context.Context) (map[string]any, error) {
	var result map[string]any
	err := m.client.SyncGetGatewayPublicKey(ctx, m.gatewayName, nil, &result)
	return result, err
}

// GetPublicKeyString fetch the public key from apigw.
func (m *Manager) GetPublicKeyString(ctx context.Context) (string, error) {
	info, err := m.GetPublicKey(ctx)
	if err != nil {
		return "", err
	}

	value, ok := info["public_key"]
	if !ok {
		return "", errors.Wrap(ErrApiGatewayPublicKeyNotFound, m.gatewayName)
	}

	publicKey, ok := value.(string)
	if !ok {
		return "", errors.Wrapf(
			ErrApiGatewayPublicKeyTypeNotSupported,
			"expected %T, got %T", publicKey, value,
		)
	}

	return publicKey, nil
}

// GetLatestResourceVersion get the latest resource version from apigw.
func (m *Manager) GetLatestResourceVersion(ctx context.Context) (map[string]any, error) {
	var result map[string]any
	err := m.client.SyncGetResourceVersionLatest(ctx, m.gatewayName, nil, &result)
	return result, err
}

// ResourceVersionExists check whether the resource version exists in apigw.
func (m *Manager) ResourceVersionExists(ctx context.Context, version string) (bool, error) {
	var result struct {
		Count int `json:"count"`
	}
	err := m.client.SyncListResourceVersions(ctx, m.gatewayName, url.Values{"version": {version}}, &result)
	return result.Count > 0, err
}

// SyncBasicInfo sync the basic info from definition under the namespace to apigw.
func (m *Manager) SyncBasicInfo(ctx context.Context) (map[string]any, error) {
	data, err := m.definition.Get(apiGatewayNamespace)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to get %s", apiGatewayNamespace)
	}

	var result map[string]any
	err = m.client.SyncGateway(ctx, m.gatewayName, data, &result)
	return result, err
}

// SyncStagesConfig sync the stages config from definition under the namespace to apigw.
func (m *Manager) SyncStagesConfig(ctx context.Context) (map[string]any, error) {
	stages, err := m.definition.GetArray(stagesNamespace)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to get %s", stagesNamespace)
	}
	results := make(map[string]any, len(stages))
	for _, stage := range stages {
		var result map[string]any
		if err := m.client.SyncStages(ctx, m.gatewayName, stage, &result); err != nil {
			return nil, errors.WithMessagef(err, "failed to sync stage %v", stage["name"])
		}
		results[fmt.Sprint(stage["name"])] = result
	}
	return results, nil
}

// SyncStageMcpConfig sync the mcp servers of the stages from definition to apigw.
func (m *Manager) SyncStageMcpConfig(ctx context.Context) (map[string]any, error) {
	stages, err := m.definition.GetArray(stagesNamespace)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to get %s", stagesNamespace)
	}
	results := make(map[string]any, len(stages))
	for _, stage := range stages {
		name := fmt.Sprint(stage["name"])
		var result []map[string]any
		if err := m.client.SyncStageMCPServers(ctx, m.gatewayName, name, stage, &result); err != nil {
			return nil, errors.WithMessagef(err, "failed to sync the mcp servers of stage %s", name)
		}
		results[name] = result
	}
	return results, nil
}

// SyncResourcesConfig sync the resources config to apigw.
func (m *Manager) SyncResourcesConfig(ctx context.Context, resources map[string]any) (map[string]any, error) {
	var result map[string]any
	err := m.client.SyncResources(ctx, m.gatewayName, resources, &result)
	return result, err
}

// SyncResourceDocByArchive sync the resource doc from archive to apigw.
func (m *Manager) SyncResourceDocByArchive(ctx context.Context) error {
	data, err := m.definition.Get(resourceDocsNamespace)
	if err != nil {
		return errors.WithMessagef(err, "failed to get %s", resourceDocsNamespace)
	}
	baseDir := data["basedir"].(string)
	zipPath := filepath.Join(baseDir, "resources_docs.zip")
	err = util.ZipDirectory(baseDir, zipPath, ".md")
	if err != nil {
		return errors.WithMessagef(err, "failed to zip %s", baseDir)
	}
	// 上传资源文档
	archive, err := os.Open(zipPath)
	if err != nil {
		return errors.WithMessagef(err, "failed to read %s", zipPath)
	}
	defer archive.Close()

	return m.client.SyncResourceDoc(ctx, m.gatewayName, archive)
}

// targetGateway removes the gateway_name from the permission, which is the gateway to apply for or grant,
// and defaults to the managed gateway.
func (m *Manager) targetGateway(permission map[string]any) string {
	name, ok := permission["gateway_name"]
	if !ok {
		return m.gatewayName
	}
	delete(permission, "gateway_name")
	return fmt.Sprint(name)
}

// ApplyPermissions apply for the permissions defined in apply_permissions, which can be of other gateways.
func (m *Manager) ApplyPermissions(ctx context.Context) (map[string]any, error) {
	permissions, err := m.definition.GetArray(applyPermissionsNamespace)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to get %s", applyPermissionsNamespace)
	}

	results := make(map[string]any, len(permissions))
	for i, definedPermission := range permissions {
		permission := maps.Clone(definedPermission)
		if _, ok := permission["target_app_code"]; !ok {
			permission["target_app_code"] = m.appCode
		}
		if _, ok := permission["applicant"]; !ok {
			permission["applicant"] = permission["target_app_code"]
		}
		if permission["grant_dimension"] == nil {
			permission["grant_dimension"] = "gateway"
		}

		gatewayName := m.targetGateway(permission)
		var result map[string]any
		if err := m.client.OpenApplyGatewayPermission(ctx, gatewayName, permission, &result); err != nil {
			return nil, err
		}
		results[fmt.Sprintf("result_%d", i)] = result
	}
	return results, nil
}

// GrantPermissions grant the permissions defined in grant_permissions to apps.
func (m *Manager) GrantPermissions(ctx context.Context) error {
	permissions, err := m.definition.GetArray(grantPermissionsNamespace)
	if err != nil {
		return errors.WithMessagef(err, "failed to get %s", grantPermissionsNamespace)
	}

	for _, definedPermission := range permissions {
		permission := maps.Clone(definedPermission)
		if permission["grant_dimension"] == nil {
			permission["grant_dimension"] = "gateway"
		}
		if err := m.client.SyncGrantPermission(ctx, m.targetGateway(permission), permission); err != nil {
			return err
		}
	}
	return nil
}

// AddRelatedApps add the apps defined in related_apps as the related apps of the gateway.
func (m *Manager) AddRelatedApps(ctx context.Context) error {
	relatedApps, ok := m.definition.definition[relatedAppsNamespace].([]any)
	if !ok || len(relatedApps) == 0 {
		return nil
	}

	return m.client.SyncAddRelatedApps(ctx, m.gatewayName, map[string]any{"related_app_codes": relatedApps})
}

// CreateResourceVersion create a resource version.
func (m *Manager) CreateResourceVersion(ctx context.Context, version, comment string) (map[string]any, error) {
	data := map[string]any{
		"version": version,
		"comment": comment,
	}
	var result map[string]any
	err := m.client.SyncCreateResourceVersion(ctx, m.gatewayName, data, &result)
	return result, err
}

// Release release the resource version to the stages defined in the definition.
func (m *Manager) Release(ctx context.Context, version, comment string) (map[string]any, error) {
	stages, err := m.definition.GetArray(stagesNamespace)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to get %s", stagesNamespace)
	}
	stageNames := make([]string, 0, len(stages))
	for _, stage := range stages {
		stageNames = append(stageNames, fmt.Sprint(stage["name"]))
	}
	data := map[string]any{
		"stage_names": stageNames,
		"version":     version,
		"comment":     comment,
	}
	var result map[string]any
	err = m.client.SyncRelease(ctx, m.gatewayName, data, &result)
	return result, err
}

// NewManager create a new manager.
func NewManager(gatewayName string, config apigateway.Config, definition *Definition) (*Manager, error) {
	client, err := apigateway.New(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create apigateway client")
	}

	return &Manager{
		gatewayName: gatewayName,
		appCode:     config.AppCode,
		definition:  definition,
		client:      client,
	}, nil
}

// NewDefaultManager create a new default manager.
func NewDefaultManager(gatewayName string, config apigateway.Config) (*Manager, error) {
	return NewManager(gatewayName, config, nil)
}

// NewManagerFrom file will create a new manager from the file.
func NewManagerFrom(gatewayName string, config apigateway.Config, path string) (*Manager, error) {
	manager, err := NewDefaultManager(gatewayName, config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create manager")
	}

	return manager, manager.LoadDefinition(path)
}
