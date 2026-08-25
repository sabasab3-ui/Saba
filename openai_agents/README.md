# SABA OpenAI Agents — Multi-Agent Workflow Package

This package adds an OpenAI Agents SDK orchestration layer for SABA.

Architecture:

SABA Core Intelligence -> SABA Orchestrator -> OpenAI Agents -> specialist agents

Specialists included:
- Research Agent
- Analysis Agent
- Reasoning Agent
- Business Agent
- Coding Coordinator

The Coding Coordinator does NOT execute code itself. It prepares coding work for the existing SABA Coding System/OpenHands layer.

## Requirements

- Python 3.10+
- An OpenAI API key
- `pip install -r requirements.txt`

The official OpenAI Agents SDK uses `Agent`, `Runner`, handoffs, tools, guardrails, sessions and tracing for multi-agent workflows.

## Quick start

```bash
cd SABA_OpenAI_Agents
python -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
cp .env.example .env
# edit .env and set OPENAI_API_KEY
./run.sh
```

The service listens on `127.0.0.1:8090` by default.

Endpoints:
- GET `/health`
- GET `/status`
- POST `/run`

Example:

```bash
curl -X POST http://127.0.0.1:8090/run \
  -H 'Content-Type: application/json' \
  -d '{"task":"Research how SABA could help small businesses in Uganda.","mode":"auto"}'
```

## Modes

`auto` — triage to the most suitable specialist.
`research` — research-focused agent.
`analysis` — analysis-focused agent.
`reasoning` — reasoning/planning-focused agent.
`business` — business-focused agent.
`coding` — coding-task coordinator; sends a structured coding plan rather than executing code.

## SABA integration

Set:

```env
SABA_GATEWAY_URL=http://127.0.0.1:8080
```

The adapter can call a future HTTP gateway endpoint without changing the core agent workflow. The current SABA Go gateway can be wired to this service as the next integration step.

Security:
- Keep the service on localhost/private network until authentication is added.
- Never commit `.env` or API keys.
- Do not expose `/run` directly to the public internet.
