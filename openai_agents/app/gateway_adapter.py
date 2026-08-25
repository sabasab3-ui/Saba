import httpx
from typing import Any

class SabaGatewayClient:
    """Optional adapter for the SABA Go gateway.

    The current SABA gateway is an in-process Go component. This client is
    intentionally isolated so an HTTP endpoint can be added later without
    coupling the OpenAI Agents layer to Go internals.
    """

    def __init__(self, base_url: str, timeout: float = 15.0):
        self.base_url = base_url.rstrip("/")
        self.timeout = timeout

    async def health(self) -> dict[str, Any]:
        async with httpx.AsyncClient(timeout=self.timeout) as client:
            response = await client.get(f"{self.base_url}/health")
            response.raise_for_status()
            return response.json()
