# SABA Integration Boundary

This package runs as a private service behind the SABA Go gateway.

Recommended flow:

SABA Core Intelligence
→ SABA Orchestrator/Gateway
→ OpenAI Agents service `/run`
→ specialist agent
→ response
→ SABA Gateway
→ user channel

The Research Agent can use OpenAI-hosted web search. The Coding Coordinator
only creates implementation plans; the separate SABA Coding System/OpenHands
service remains responsible for actual code execution.

Keep port 8090 bound to localhost/private network. Do not expose `/run`
directly to the public internet.

The Go example in `go_bridge_example.go` shows the HTTP boundary.
