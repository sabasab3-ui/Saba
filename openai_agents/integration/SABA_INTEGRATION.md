# Adding the Agents Layer to SABA

Your current SABA repository already contains an `internal/gateway` package and an
`internal/intelligence` package. The gateway currently accepts an `AgentRequest`
and returns an `AgentResponse`.

Recommended integration:

1. Run this Python service privately on the same VPS as SABA.
2. Add an HTTP client in the Go gateway for `/run`.
3. Route requests whose agent/mode is `openai_agents` to this service.
4. Keep SABA's existing intelligence analyzer as a native fallback.
5. Keep OpenHands behind the separate Coding System boundary.

Suggested request:

```json
{
  "task": "Analyze this business problem",
  "mode": "auto",
  "input": "..."
}
```

Suggested response:

```json
{
  "status": "completed",
  "mode": "auto",
  "agent": "SABA Analysis Agent",
  "output": "..."
}
```

Do not expose this service directly to the public internet. Put it behind the
SABA gateway, authentication, and a private network/VPS firewall.
