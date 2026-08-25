# SABA OpenAI Agents — Production-Ready Multi-Agent Layer

This package is the OpenAI Agents orchestration layer for SABA.

## Architecture

SABA Core Intelligence
        ↓
SABA Orchestrator / Gateway
        ↓
OpenAI Agents
   ├── Research Agent (web search)
   ├── Analysis Agent
   ├── Reasoning Agent
   ├── Business Agent
   └── Coding Coordinator
        ↓
SABA Coding System / OpenHands

The Coding Coordinator prepares implementation plans. It does not execute shell
commands or modify repositories. Execution remains inside the separate SABA
Coding System/OpenHands boundary.

## Included

- OpenAI Agents SDK
- Multi-agent handoffs
- Real hosted web search for the research specialist
- Optional explicit model configuration
- FastAPI service
- `/health`
- `/status`
- `/run`
- Input validation
- Safe API-key handling through environment variables
- SABA gateway adapter
- Go bridge example
- Local syntax/test scripts
- Docker support
- Termux installation script
- No secrets committed

## Requirements

Python 3.10+ and an OpenAI API key.

The SDK can use its current default model when `OPENAI_MODEL` is empty. Set
`OPENAI_MODEL` only when you explicitly want a different supported model.

## Termux / Linux

```bash
bash install_termux.sh
cp .env.example .env
nano .env
./run.sh
```

Set:

```env
OPENAI_API_KEY=your_key_here
```

Do not commit `.env`.

## API

Health:

```bash
curl http://127.0.0.1:8090/health
```

Status:

```bash
curl http://127.0.0.1:8090/status
```

Run an automatic workflow:

```bash
curl -X POST http://127.0.0.1:8090/run \
  -H 'Content-Type: application/json' \
  -d '{"task":"Analyze how SABA could help small businesses in Uganda.","mode":"auto"}'
```

Available modes:

- `auto`
- `research`
- `analysis`
- `reasoning`
- `business`
- `coding`

## SABA integration

The package is deliberately isolated under `openai_agents/`. The included
`integration/` directory documents the Go HTTP boundary for connecting the
existing SABA gateway to this service.

Keep this service private. Put authentication and network access control at
the SABA gateway/VPS layer before exposing it beyond localhost/private network.

## Testing

```bash
bash test.sh
```

This performs syntax compilation and package tests. An OpenAI API key is not
required for the local syntax/unit tests.
