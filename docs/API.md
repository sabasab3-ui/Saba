# SABA API Documentation

## Overview
SABA provides a comprehensive REST API for accessing AI agents, automations, and business intelligence capabilities.

## Base URL
```
http://localhost:8080/api
```

## Authentication
No authentication required for local deployment. For production, configure API keys in `.env`.

## Endpoints

### Health & Status

#### GET /health
Check system health
```bash
curl http://localhost:8080/health
```

Response:
```json
{
  "status": "ok",
  "service": "saba",
  "version": "2.0.0"
}
```

### Agents

#### POST /api/agents/run
Execute an agent task
```bash
curl -X POST http://localhost:8080/api/agents/run \
  -H 'Content-Type: application/json' \
  -d '{
    "task": "Analyze market opportunities in Uganda",
    "mode": "auto",
    "user_id": "user123"
  }'
```

#### GET /api/agents
List available agents
```bash
curl http://localhost:8080/api/agents
```

#### GET /api/agents/:id
Get agent details
```bash
curl http://localhost:8080/api/agents/research-001
```

### Automations

#### POST /api/automations
Create an automation
```bash
curl -X POST http://localhost:8080/api/automations \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "Daily Inventory Check",
    "type": "schedule",
    "config": {"cron": "0 9 * * *"}
  }'
```

#### GET /api/automations
List automations
```bash
curl http://localhost:8080/api/automations
```

#### POST /api/automations/:id/run
Execute an automation
```bash
curl -X POST http://localhost:8080/api/automations/auto-001/run
```

### Customers

#### POST /api/customers
Add a customer
```bash
curl -X POST http://localhost:8080/api/customers \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "Example Corp",
    "email": "contact@example.com",
    "country": "Uganda",
    "industry": "Retail"
  }'
```

#### GET /api/customers
List customers
```bash
curl http://localhost:8080/api/customers
```

### Analytics

#### GET /api/analytics
Get analytics data
```bash
curl http://localhost:8080/api/analytics?metric_type=agents_used&days=30
```

#### POST /api/analytics
Record an analytics event
```bash
curl -X POST http://localhost:8080/api/analytics \
  -H 'Content-Type: application/json' \
  -d '{
    "metric_type": "agents_used",
    "metric_value": 1,
    "dimensions": {"agent_type": "research", "country": "Uganda"}
  }'
```

## Error Handling

All errors return a consistent format:
```json
{
  "error": "Error message",
  "status": "error",
  "code": 400
}
```

## Rate Limiting
No rate limiting in development mode.
Production: 1000 requests/hour per API key

## Pagination
List endpoints support pagination:
```bash
curl http://localhost:8080/api/automations?page=1&limit=20
```
