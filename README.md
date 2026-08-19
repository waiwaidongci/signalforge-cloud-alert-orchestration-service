# 02 SignalForge

SignalForge 是一个纯 Go 实现的云原生告警聚合与降噪编排服务。它接收 Prometheus、Loki、探针和自定义应用告警，统一标准化后执行去重、聚合、静默、抑制、路由、升级和通知编排。

## 架构

项目按领域和分层组织，所有核心依赖都以接口注入：

```text
cmd/
  server/        HTTP API 入口
  scheduler/     后台调度入口
internal/
  alert/         告警接收、去重入口编排、批量确认/关闭
  source/        告警源管理
  dedup/         指纹去重
  silence/       静默与抑制规则
  escalation/    升级策略
  routing/       路由策略
  notification/  通知通道与分发
  incident/      事件组与时间线
  shared/        配置、日志、HTTP、数据库、metrics 等基础设施
  app/           依赖组装与路由
  scheduler/     周期调度 worker
api/             OpenAPI 说明
configs/         本地 YAML 配置
migrations/      SQLite 与 PostgreSQL 迁移
deploy/          Dockerfile 与 Compose 示例
scripts/         本地开发与验证脚本
```

每个核心领域包含 `domain`、`application`、`adapter` 和 `infrastructure` 层。`domain` 定义实体和仓储接口，`application` 编排用例，`adapter` 提供 HTTP 和 SQL 实现，`infrastructure` 提供技术能力。

## 本地运行

默认使用 `modernc.org/sqlite`，无需 Docker 或 PostgreSQL：

```bash
cd /Users/ali/Desktop/moreWork/go-new-projects/02-signalforge
go mod tidy
go run ./cmd/server
```

服务默认监听 `:8080`。若端口被占用：

```bash
SIGNALFORGE_SERVER_ADDR=:18080 go run ./cmd/server
```

也可以使用脚本：

```bash
./scripts/run-dev.sh
```

健康检查：

```bash
curl http://127.0.0.1:8080/healthz
curl http://127.0.0.1:8080/readyz
curl http://127.0.0.1:8080/metrics
```

## 配置

默认配置文件为 `configs/config.yaml`。所有配置均可通过 `SIGNALFORGE_*` 环境变量覆盖，例如：

```bash
export SIGNALFORGE_SERVER_ADDR=:18080
export SIGNALFORGE_DATABASE_DRIVER=sqlite
export SIGNALFORGE_DATABASE_DSN=./data/signalforge.db
export SIGNALFORGE_LOG_LEVEL=debug
export SIGNALFORGE_AUTH_PLACEHOLDER_TOKEN=change-me
```

关键配置项：

- `server.request_timeout`：HTTP 请求超时。
- `server.max_body_bytes`：最大请求体。
- `database.driver`：`sqlite` 或 `postgres`。
- `database.dsn`：数据库连接串。
- `scheduler.enabled`：是否启用后台调度。
- `auth.placeholder_token`：占位认证令牌；为空时关闭认证。

## 数据库迁移

迁移文件位于：

- `migrations/sqlite/0001_init.sql`
- `migrations/postgres/0001_init.sql`

SQLite 迁移会在服务启动时自动执行。PostgreSQL 部署示例见 `deploy/docker-compose.yml`。

## 统一 Webhook 示例

先创建告警源：

```bash
curl -X POST http://127.0.0.1:8080/api/v1/sources \
  -H 'Content-Type: application/json' \
  -d '{"name":"prometheus-main","kind":"prometheus","rate_limit":120}'
```

响应中复制 `api_key`，然后提交统一格式告警：

```bash
curl -X POST http://127.0.0.1:8080/api/v1/webhook \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer sf_<api_key>" \
  -d '{
    "external_id":"evt-001",
    "resource":"api-server-1",
    "severity":"high",
    "title":"API latency high",
    "description":"latency above threshold",
    "labels":{"service":"api","environment":"prod"},
    "annotations":{"summary":"latency high"},
    "timestamp":"2026-08-20T00:00:00Z"
  }'
```

同一 `source_id + external_id` 会幂等合并，相同指纹会聚合到事件组。

## REST API

主要端点：

```text
GET    /api/v1/sources
POST   /api/v1/sources
GET    /api/v1/sources/{id}
PATCH  /api/v1/sources/{id}
DELETE /api/v1/sources/{id}

GET    /api/v1/alerts
GET    /api/v1/alerts/{id}
POST   /api/v1/alerts/acknowledge
POST   /api/v1/alerts/close
POST   /api/v1/webhook

GET    /api/v1/incidents
GET    /api/v1/incidents/{id}
GET    /api/v1/incidents/{id}/timeline
POST   /api/v1/incidents/acknowledge
POST   /api/v1/incidents/close

GET/POST /api/v1/silences
GET/PATCH/DELETE /api/v1/silences/{id}

GET/POST /api/v1/suppressions
GET/PATCH/DELETE /api/v1/suppressions/{id}

GET/POST /api/v1/routing-rules
GET/PATCH/DELETE /api/v1/routing-rules/{id}

GET/POST /api/v1/escalation-policies
GET/PATCH/DELETE /api/v1/escalation-policies/{id}

GET /api/v1/notifications
GET /api/v1/notifications/{id}
```

列表接口支持 `page`、`per_page` 以及领域过滤参数。批量确认和关闭请求体：

```json
{"ids":["alr_...","alr_..."],"actor":"operator"}
```

## Makefile 与脚本

```bash
make tidy
make fmt
make test
make vet
make build
make run-dev
make stats
make verify
```

脚本说明：

- `scripts/run-dev.sh`：启动本地开发服务并等待健康检查通过。
- `scripts/verify.sh`：执行格式化、依赖整理、测试、vet、build 和统计。
- `scripts/stats.sh`：统计非测试 Go 文件数和行数。

