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
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/TencentBlueKing/bk-apigateway-sdks/apigateway"
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/define"
	"github.com/TencentBlueKing/bk-apigateway-sdks/gin_contrib/util"
)

const (
	apiGatewayNamespace   = "apigateway"
	stagesNamespace       = "stages"
	permissionsNamespace  = "grant_permissions"
	resourceDocsNamespace = "resource_docs"
	relatedAppsNamespace  = "related_apps"
)

type apiGatewayV2Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type apiGatewayV2Result struct {
	Data  interface{}        `json:"data"`
	Error *apiGatewayV2Error `json:"error"`
}

// Manager is the manager of apigw, it helps to sync apigw configs and get apigw infomations.
type Manager struct {
	apiName    string
	definition *Definition
	client     *apigateway.Client
	config     *bkapi.ClientConfig
}

func (m *Manager) requestWithBody(
	operation define.Operation,
	body map[string]interface{},
) (map[string]interface{}, error) {
	return m.request(operation.SetBody(body))
}

func (m *Manager) requestWithFile(
	operation define.Operation,
	name string,
	file *os.File,
) (map[string]interface{}, error) {
	return m.request(operation.SetFile(name, file))
}

func (m *Manager) request(operation define.Operation) (map[string]interface{}, error) {
	data, err := m.requestData(operation)
	result, _ := data.(map[string]interface{})
	return result, err
}

// requestData sends a v2 api request and returns the `data` field of the response body.
func (m *Manager) requestData(operation define.Operation) (interface{}, error) {
	var result apiGatewayV2Result
	response, err := operation.
		SetPathParams(map[string]string{
			"gateway_name": m.apiName,
		}).
		SetResult(&result).
		Request()
	if err != nil {
		return nil, errors.Wrapf(err, "request to %v failed", operation)
	}

	if result.Error != nil {
		return nil, errors.Wrapf(
			ErrApigatewayRequest,
			"status: %d, code: %s, message: %s",
			response.StatusCode,
			result.Error.Code,
			result.Error.Message,
		)
	}

	// the operation only reports errors with the X-Bkapi-Error-Code header
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, errors.Wrapf(ErrApigatewayRequest, "status: %d", response.StatusCode)
	}

	return result.Data, nil
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
func (m *Manager) GetPublicKey() (map[string]interface{}, error) {
	return m.request(m.client.V2SyncGetGatewayPublicKey())
}

// GetPublicKey fetch the public key from apigw.
func (m *Manager) GetPublicKeyString() (string, error) {
	info, err := m.GetPublicKey()
	if err != nil {
		return "", err
	}

	value, ok := info["public_key"]
	if !ok {
		return "", errors.Wrap(ErrApiGatewayPublicKeyNotFound, m.apiName)
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
func (m *Manager) GetLatestResourceVersion() (map[string]interface{}, error) {
	return m.request(m.client.V2SyncGetResourceVersionLatest())
}

// ResourceVersionExists check whether the resource version exists in apigw.
func (m *Manager) ResourceVersionExists(version string) (bool, error) {
	data, err := m.requestData(m.client.V2SyncListResourceVersions().SetQueryParams(map[string]string{
		"version": version,
	}))
	if err != nil {
		return false, err
	}

	// apigw < 1.24.0 returns a list instead of the paginated result
	switch result := data.(type) {
	case []interface{}:
		return len(result) != 0, nil
	case map[string]interface{}:
		count, _ := result["count"].(float64)
		return count != 0, nil
	default:
		return false, nil
	}
}

// SyncBasicInfo sync the basic info from definition under the namespace to apigw.
func (m *Manager) SyncBasicInfo() (map[string]interface{}, error) {
	data, err := m.definition.Get(apiGatewayNamespace)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to get %s", apiGatewayNamespace)
	}

	return m.requestWithBody(m.client.V2SyncGateway(), data)
}

// SyncStagesConfig sync the stages config from definition under the namespace to apigw.
func (m *Manager) SyncStagesConfig() (map[string]interface{}, error) {
	stages, err := m.definition.GetArray(stagesNamespace)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to get %s", stagesNamespace)
	}
	resultMap := make(map[string]interface{})
	for _, stage := range stages {
		result, err := m.requestWithBody(m.client.V2SyncStages(), stage)
		if err != nil {
			return nil, errors.WithMessagef(err, "failed to get %s", stagesNamespace)
		}
		resultMap[fmt.Sprintf("%v", stage["name"])] = result
	}
	return resultMap, nil
}

// SyncStageMcpConfig sync mcp config from definition under the namespace to apigw.
func (m *Manager) SyncStageMcpConfig() (map[string]interface{}, error) {
	stages, err := m.definition.GetArray(stagesNamespace)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to get %s", stagesNamespace)
	}
	resultMap := make(map[string]interface{})
	for _, stage := range stages {
		result, err := m.requestData(m.client.SyncStageMcpServers().SetPathParams(
			map[string]string{"stage_name": fmt.Sprintf("%v", stage["name"])}).SetBody(stage))
		if err != nil {
			return nil, errors.WithMessagef(err, "failed to get %s", stagesNamespace)
		}
		resultMap[fmt.Sprintf("%v", stage["name"])] = result
	}
	return resultMap, nil
}

// SyncPluginConfig sync the plugin config from definition under the namespace to apigw.
//
// Deprecated: access strategies are not supported by the apigw v2 apis, it does nothing now.
func (m *Manager) SyncPluginConfig(namespace string) (map[string]interface{}, error) {
	log.Printf("SyncPluginConfig is deprecated, and now it does nothing")
	return map[string]interface{}{}, nil
}

// SyncResourcesConfig sync the resources config from definition under the namespace to apigw.
func (m *Manager) SyncResourcesConfig(resources map[string]interface{}) (map[string]interface{}, error) {
	return m.requestWithBody(m.client.V2SyncResources(), resources)
}

// SyncResourceDocByArchive sync the resource doc from archive to apigw.
func (m *Manager) SyncResourceDocByArchive() (map[string]interface{}, error) {
	data, err := m.definition.Get(resourceDocsNamespace)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to get %s", resourceDocsNamespace)
	}
	baseDir := data["basedir"].(string)
	zipPath := filepath.Join(baseDir, "resources_docs.zip")
	err = util.ZipDirectory(baseDir, zipPath, ".md")
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to zip %s", baseDir)
	}
	// 上传资源文档
	resourceDocsFile, err := os.Open(zipPath)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to read %s", zipPath)
	}
	defer resourceDocsFile.Close()
	return m.requestWithFile(m.client.V2SyncResourceDoc(), "file", resourceDocsFile)
}

func copyMap(data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{}, len(data))
	for k, v := range data {
		result[k] = v
	}
	return result
}

// withGatewayName moves the gateway name in the permission to the path params,
// so that the permission can be applied to or granted for other gateways.
func withGatewayName(operation define.Operation, permission map[string]interface{}) define.Operation {
	if apiName, ok := permission["api_name"]; ok {
		permission["gateway_name"] = apiName
		delete(permission, "api_name")
	}

	gatewayName, ok := permission["gateway_name"]
	if !ok {
		return operation
	}
	delete(permission, "gateway_name")

	return operation.SetPathParams(map[string]string{"gateway_name": fmt.Sprintf("%v", gatewayName)})
}

// ApplyPermissions apply the permissions under the namespace to apigw.
func (m *Manager) ApplyPermissions(namespace string) (map[string]interface{}, error) {
	permissions, err := m.definition.GetArray(namespace)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to get %s", namespace)
	}

	appCode := m.config.ProvideConfig(m.client.Name()).(*bkapi.ClientConfig).AppCode
	resultMap := make(map[string]interface{})
	for i, definedPermission := range permissions {
		permission := copyMap(definedPermission)
		if _, ok := permission["target_app_code"]; !ok {
			permission["target_app_code"] = appCode
		}
		if _, ok := permission["applicant"]; !ok {
			permission["applicant"] = permission["target_app_code"]
		}
		if permission["grant_dimension"] == nil {
			permission["grant_dimension"] = "gateway"
		}

		operation := withGatewayName(m.client.V2OpenApplyGatewayPermission(), permission)
		result, err := m.requestWithBody(operation, permission)
		if err != nil {
			return nil, err
		}
		resultMap[fmt.Sprintf("result_%d", i)] = result
	}
	return resultMap, nil
}

// GrantPermissions grant the permissions under the namespace to apigw.
func (m *Manager) GrantPermissions() (map[string]interface{}, error) {
	permissions, err := m.definition.GetArray(permissionsNamespace)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to get %s", permissionsNamespace)
	}
	resultMap := make(map[string]interface{})
	for i, definedPermission := range permissions {
		permission := copyMap(definedPermission)
		if _, ok := permission["target_app_code"]; !ok {
			permission["target_app_code"] = permission["bk_app_code"]
		}
		delete(permission, "bk_app_code")
		if dimension := permission["grant_dimension"]; dimension == nil || dimension == "api" {
			permission["grant_dimension"] = "gateway"
		}

		operation := withGatewayName(m.client.V2SyncGrantPermission(), permission)
		result, err := m.requestWithBody(operation, permission)
		if err != nil {
			return nil, err
		}
		resultMap[fmt.Sprintf("result_%d", i)] = result
	}
	return resultMap, nil
}

// AddRelatedApps add the related apps under the namespace to apigw.
func (m *Manager) AddRelatedApps() (map[string]interface{}, error) {
	relatedApps, ok := m.definition.definition[relatedAppsNamespace].([]interface{})
	if !ok || len(relatedApps) == 0 {
		return map[string]interface{}{}, nil
	}

	return m.requestWithBody(m.client.V2SyncAddRelatedApps(), map[string]interface{}{
		"related_app_codes": relatedApps,
	})
}

// CreateResourceVersion create a resource version defined in the namespace.
func (m *Manager) CreateResourceVersion(version string, comment string) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"version": version,
		"comment": comment,
	}
	return m.requestWithBody(m.client.V2SyncCreateResourceVersion(), data)
}

// Release release the resource version defined in the namespace.
func (m *Manager) Release(version string) (map[string]interface{}, error) {
	return m.ReleaseWithComment(version, "")
}

// ReleaseWithComment release the resource version defined in the namespace with the comment.
func (m *Manager) ReleaseWithComment(version string, comment string) (map[string]interface{}, error) {
	stages, err := m.definition.GetArray(stagesNamespace)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to get %s", stagesNamespace)
	}
	var stageNames []string
	for _, stage := range stages {
		stageNames = append(stageNames, stage["name"].(string))
	}
	data := map[string]interface{}{
		"stage_names": stageNames,
		"version":     version,
		"comment":     comment,
	}
	return m.requestWithBody(m.client.V2SyncRelease(), data)
}

// NewManager create a new manager.
func NewManager(
	apiName string,
	config bkapi.ClientConfig,
	definition *Definition,
	clientFactory func(
		configProvider define.ClientConfigProvider, opts ...define.BkApiClientOption,
	) (*apigateway.Client, error),
) (*Manager, error) {
	client, err := clientFactory(config, bkapi.OptJsonBodyProvider(), bkapi.JsonResultProvider())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create apigateway client")
	}

	return &Manager{
		apiName:    apiName,
		config:     &config,
		client:     client,
		definition: definition,
	}, nil
}

// NewDefaultManager create a new default manager.
func NewDefaultManager(apiName string, config bkapi.ClientConfig) (*Manager, error) {
	return NewManager(apiName, config, nil, apigateway.New)
}

// NewManagerFrom file will create a new manager from the file.
func NewManagerFrom(
	apiName string,
	config bkapi.ClientConfig,
	path string,
) (*Manager, error) {
	manager, err := NewDefaultManager(apiName, config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create manager")
	}

	return manager, manager.LoadDefinition(path)
}
