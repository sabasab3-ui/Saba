# SABA Unified Tool Package

This package turns the existing SABA intelligence service into a single
discoverable platform for business and finance tools.

## Included now

- `inventory` — manage and query inventory levels
- `orders` — create, list, and update orders
- `payments` — provider interface for MTN Mobile Money and Airtel Money
- `/tools` — API endpoint exposing the registered tool catalog
- Existing intelligence pipeline remains available through `/analyze`

## Architecture

Question -> Research -> Analysis -> Decision -> Report

Business tools are registered separately:

SABA -> ToolKit -> Inventory / Orders / Payments

This keeps the intelligence engine independent from individual integrations
and gives us a clean place to add future tools such as:

- web search/fetch
- memory
- finance and market data
- notifications
- automation/scheduling
- Telegram/WhatsApp
- VPS/server monitoring
- analytics
- customer/business management

## Applying the package

Copy:

- `internal/agent/registry.go`
- `cmd/saba/main.go`

into the SABA repository, replacing the existing `cmd/saba/main.go`.

Then run:

```bash
go test ./...
go run ./cmd/saba
```

The service should expose:

```text
GET  /health
GET  /tools
POST /analyze
```
