# apigateway

蓝鲸 API 网关（bk-apigateway）v2 接口的轻量客户端，基于 [gentleman](https://github.com/h2non/gentleman)。

客户端只处理所有接口共有的部分：认证、响应状态、`{"data": ...}` 响应外壳和日志。请求体和结果的类型由调用方决定，不需要定义额外的请求或响应结构。

## 创建客户端

```go
import "github.com/TencentBlueKing/bk-apigateway-sdks/v2/apigateway"

// 部署在蓝鲸 PaaS 上的应用，从环境变量读取配置
client, err := apigateway.New(apigateway.ConfigFromEnv())

// 或者直接指定配置
client, err := apigateway.New(apigateway.Config{
	Endpoint:  "https://bkapi.example.com/api/bk-apigateway/prod",
	AppCode:   "my-app",
	AppSecret: "my-secret",
})
```

| 字段 | 说明 | 环境变量 |
| --- | --- | --- |
| `Endpoint` | bk-apigateway 的环境地址，必填 | `BK_API_URL_TMPL`，其中的 `{api_name}` 或 `{gateway_name}` 替换为 `bk-apigateway`，环境为 `prod` |
| `AppCode` | 应用 ID | `BK_APP_CODE`、`BKPAAS_APP_ID`、`APP_CODE` |
| `AppSecret` | 应用密钥 | `BK_APP_SECRET`、`BKPAAS_APP_SECRET`、`SECRET_KEY` |
| `AccessToken` | 可选，用户或应用的 access_token | |
| `TenantID` | 租户 ID，通过 `X-Bk-Tenant-Id` 请求头传递 | `BKPAAS_APP_TENANT_ID`（为空表示全租户应用，使用 `system`），未部署在 PaaS 时为 `BK_APP_TENANT_ID` |
| `Timeout` | 单次请求超时，默认 60 秒 | |
| `Transport` | 可选，自定义 `http.RoundTripper`，如代理、链路追踪、测试 mock | |
| `Logger` | 可选，默认 `slog.Default()` | |

`ConfigFromEnv` 只返回配置，可以修改后再创建客户端。客户端可以并发使用。

## 调用接口

每个接口对应一个方法，方法名为网关资源的 operationId 去掉 `v2_` 前缀，如 `v2_sync_gateway` 对应 `SyncGateway`。参数依次为：

1. `ctx`：控制取消和超时；
2. 路径参数，如 `gatewayName`；
3. 请求参数：查询参数为 `url.Values`；请求体可以是 map 或结构体，`nil` 表示没有请求体。以方法签名为准，如 `GET` 接口只有查询参数；
4. `result`：结果的指针，`nil` 表示不需要结果。成功时不返回数据的接口没有这个参数，如授权、回收权限、上传文档。

```go
ctx := context.Background()

// 同步网关
var gateway struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
err := client.SyncGateway(ctx, "my-gateway", map[string]any{
	"description": "my gateway",
	"maintainers": []string{"admin"},
	"is_public":   true,
}, &gateway)

// 查询资源版本，结果也可以是 map
var versions map[string]any
err = client.SyncListResourceVersions(ctx, "my-gateway", url.Values{"version": {"1.0.0"}}, &versions)

// 上传资源文档归档
archive, err := os.Open("docs.zip")
err = client.SyncResourceDoc(ctx, "my-gateway", archive)
```

`result` 接收响应的 `data` 字段。`OpenOAuthProtectedResource` 的响应没有 `data` 外壳，`result` 接收完整的响应。

## 错误处理

响应状态不是 2xx 时返回 `*apigateway.Error`：

```go
var apiErr *apigateway.Error
if errors.As(err, &apiErr) {
	log.Println(apiErr.StatusCode, apiErr.Code, apiErr.Message, apiErr.RequestID)
}
```

- `Code`、`Message` 来自响应体中的 `error`，如 `INVALID_ARGUMENT`；请求被网关拦截时（如应用认证失败），来自 `X-Bkapi-Error-Code`、`X-Bkapi-Error-Message` 响应头；
- `RequestID` 来自 `X-Bkapi-Request-Id` 响应头，用于在网关中查询请求日志。

网络错误、超时和取消原样返回，可以用 `errors.Is(err, context.DeadlineExceeded)` 判断。成功的响应不是 JSON 时返回解析错误。

## 日志

每个请求记录一条日志，包含方法、路径、状态码、request_id、耗时和错误，不记录请求头和请求/响应体。成功的请求为 DEBUG 级别，4xx 为 WARN，其他错误为 ERROR。不需要日志时设置 `Logger: slog.New(slog.DiscardHandler)`。

## 接口范围

包含 bk-apigateway 的全部 v2 接口，与[网关资源定义](https://github.com/TencentBlueKing/blueking-apigateway/blob/master/src/dashboard/apigateway/apigateway/data/apigw-definitions/bk-apigateway-resources.yaml)一致：

- [sync.go](sync.go)：`/api/v2/sync/`，同步网关、环境、资源、文档、权限、版本，发布；
- [open.go](open.go)：`/api/v2/open/`，查询网关、资源、MCP Server，申请权限等；
- [inner.go](inner.go)：`/api/v2/inner/`，供蓝鲸平台调用，需要网关主动授权。

不包含 MCP Server 的 SSE 和 Streamable HTTP 代理接口（请使用 MCP 客户端），以及已废弃的 `/api/v2/open/gateway/{gateway_name}/public_key/`。
