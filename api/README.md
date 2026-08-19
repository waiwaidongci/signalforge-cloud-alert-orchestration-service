# SignalForge REST API

The OpenAPI description in `openapi.yaml` covers the REST surface exposed by the service. The actual route
registration lives in `internal/app/router.go`.

Public base path: `/api/v1`

Operational endpoints:

- `GET /healthz`
- `GET /readyz`
- `GET /metrics`

