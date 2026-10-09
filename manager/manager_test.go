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

package manager_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gock "gopkg.in/h2non/gock.v1"

	"github.com/TencentBlueKing/bk-apigateway-sdks/v2/apigateway"
	manager "github.com/TencentBlueKing/bk-apigateway-sdks/v2/manager"
)

var _ = Describe("Manager", func() {
	var (
		ctx      = context.Background()
		endpoint = "http://example.com"
		mgr      *manager.Manager
		bodies   []map[string]interface{}
	)

	captureBody := func(req *http.Request, _ *gock.Request) (bool, error) {
		if req.Body == nil {
			return true, nil
		}
		content, err := io.ReadAll(req.Body)
		if err != nil {
			return false, err
		}
		body := map[string]interface{}{}
		if len(content) > 0 {
			if err := json.Unmarshal(content, &body); err != nil {
				return false, err
			}
		}
		bodies = append(bodies, body)
		return true, nil
	}

	newManager := func(definition map[string]interface{}) *manager.Manager {
		m, err := manager.NewManager(
			"testing",
			apigateway.Config{Endpoint: endpoint, AppCode: "my-app", Transport: gock.NewTransport()},
			manager.NewDefinition(definition),
		)
		Expect(err).To(BeNil())
		return m
	}

	BeforeEach(func() {
		bodies = nil
		mgr = newManager(map[string]interface{}{
			"stages": []interface{}{
				map[string]interface{}{"name": "prod"},
			},
			"grant_permissions": []interface{}{
				map[string]interface{}{"target_app_code": "app1", "grant_dimension": "gateway"},
				map[string]interface{}{
					"target_app_code": "app2",
					"grant_dimension": "resource",
					"resource_names":  []interface{}{"r1"},
				},
			},
			"apply_permissions": []interface{}{
				map[string]interface{}{"gateway_name": "other-gateway"},
			},
			"related_apps": []interface{}{"app1"},
		})
	})

	AfterEach(func() {
		done := gock.IsDone()
		gock.Off()
		Expect(done).To(BeTrue())
	})

	It("should return the data of a v2 response", func() {
		gock.New(endpoint).
			Get("/api/v2/sync/gateways/testing/resource_versions/latest/").
			Reply(200).
			JSON(map[string]interface{}{"data": map[string]interface{}{"version": "1.0.0+prod"}})

		result, err := mgr.GetLatestResourceVersion(ctx)
		Expect(err).To(BeNil())
		Expect(result).To(Equal(map[string]interface{}{"version": "1.0.0+prod"}))
	})

	It("should return the error of a v2 error response", func() {
		gock.New(endpoint).
			Post("/api/v2/sync/gateways/testing/resource_versions/").
			Reply(400).
			JSON(map[string]interface{}{
				"error": map[string]interface{}{"code": "INVALID_ARGUMENT", "message": "invalid version"},
			})

		_, err := mgr.CreateResourceVersion(ctx, "1.0.0+prod+20261008", "")
		Expect(err).To(BeAssignableToTypeOf(&apigateway.Error{}))
		Expect(err.Error()).To(ContainSubstring("INVALID_ARGUMENT"))
		Expect(err.Error()).To(ContainSubstring("invalid version"))
	})

	It("should return an error for a non-2xx response without an error body", func() {
		gock.New(endpoint).
			Post("/api/v2/sync/gateways/testing/resource_versions/release/").
			Reply(500).
			JSON(map[string]interface{}{"data": nil})

		_, err := mgr.Release(ctx, "1.0.0+prod", "")
		Expect(err).To(BeAssignableToTypeOf(&apigateway.Error{}))
		Expect(err.Error()).To(ContainSubstring("500"))
	})

	It("should sync the stage mcp servers whose response data is a list", func() {
		gock.New(endpoint).
			Post("/api/v2/sync/gateways/testing/stages/prod/mcp-servers/").
			Reply(200).
			JSON(map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{"id": 1, "name": "testing-prod-mcp", "action": "created"},
				},
			})

		result, err := mgr.SyncStageMcpConfig(ctx)
		Expect(err).To(BeNil())
		Expect(result["prod"]).To(HaveLen(1))
	})

	It("should check whether the resource version exists", func() {
		gock.New(endpoint).
			Get("/api/v2/sync/gateways/testing/resource_versions/").
			MatchParam("version", `^1\.0\.0\+prod$`).
			Reply(200).
			JSON(map[string]interface{}{"data": map[string]interface{}{"count": 1, "results": []interface{}{}}})
		gock.New(endpoint).
			Get("/api/v2/sync/gateways/testing/resource_versions/").
			MatchParam("version", `^1\.0\.1\+prod$`).
			Reply(200).
			JSON(map[string]interface{}{"data": map[string]interface{}{"count": 0, "results": []interface{}{}}})

		exists, err := mgr.ResourceVersionExists(ctx, "1.0.0+prod")
		Expect(err).To(BeNil())
		Expect(exists).To(BeTrue())

		exists, err = mgr.ResourceVersionExists(ctx, "1.0.1+prod")
		Expect(err).To(BeNil())
		Expect(exists).To(BeFalse())
	})

	It("should grant permissions with the v2 fields", func() {
		gock.New(endpoint).
			Post("/api/v2/sync/gateways/testing/permissions/grant/").
			Times(2).
			AddMatcher(captureBody).
			Reply(201).
			JSON(map[string]interface{}{"data": nil})

		Expect(mgr.GrantPermissions(ctx)).To(Succeed())
		Expect(bodies).To(Equal([]map[string]interface{}{
			{"target_app_code": "app1", "grant_dimension": "gateway"},
			{"target_app_code": "app2", "grant_dimension": "resource", "resource_names": []interface{}{"r1"}},
		}))
	})

	It("should apply permissions to the defined gateway with defaults", func() {
		gock.New(endpoint).
			Post("/api/v2/open/gateways/other-gateway/permissions/apply/").
			AddMatcher(captureBody).
			Reply(200).
			JSON(map[string]interface{}{"data": map[string]interface{}{"record_id": 1}})

		result, err := mgr.ApplyPermissions(ctx)
		Expect(err).To(BeNil())
		Expect(result["result_0"]).To(Equal(map[string]interface{}{"record_id": float64(1)}))
		Expect(bodies).To(Equal([]map[string]interface{}{
			{"target_app_code": "my-app", "applicant": "my-app", "grant_dimension": "gateway"},
		}))
	})

	It("should add related apps", func() {
		gock.New(endpoint).
			Post("/api/v2/sync/gateways/testing/related-apps/").
			AddMatcher(captureBody).
			Reply(201).
			JSON(map[string]interface{}{"data": nil})

		Expect(mgr.AddRelatedApps(ctx)).To(Succeed())
		Expect(bodies).To(Equal([]map[string]interface{}{
			{"related_app_codes": []interface{}{"app1"}},
		}))
	})

	It("should skip adding related apps when there are none", func() {
		// gock fails the request if one is sent
		Expect(newManager(map[string]interface{}{"related_apps": nil}).AddRelatedApps(ctx)).To(Succeed())
	})

	It("should release with the comment", func() {
		gock.New(endpoint).
			Post("/api/v2/sync/gateways/testing/resource_versions/release/").
			AddMatcher(captureBody).
			Reply(200).
			JSON(map[string]interface{}{
				"data": map[string]interface{}{"version": "1.0.0+prod", "stage_names": []interface{}{"prod"}},
			})

		_, err := mgr.Release(ctx, "1.0.0+prod", "release comment")
		Expect(err).To(BeNil())
		Expect(bodies).To(Equal([]map[string]interface{}{
			{"version": "1.0.0+prod", "stage_names": []interface{}{"prod"}, "comment": "release comment"},
		}))
	})
})
