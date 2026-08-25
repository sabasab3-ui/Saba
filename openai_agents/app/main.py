from fastapi import FastAPI, HTTPException
from .config import settings
from .models import RunRequest, RunResponse
from .runner import run_task

app = FastAPI(
    title="SABA OpenAI Agents",
    version="1.1.0",
    description="Multi-agent orchestration layer for SABA.",
)

@app.get("/health")
async def health():
    return {"status": "ok", "service": "saba-openai-agents"}

@app.get("/status")
async def status():
    return {
        "status": "ready",
        "service": "saba-openai-agents",
        "model": settings.model or "sdk-default",
        "saba_gateway_url": settings.saba_gateway_url,
    }

@app.post("/run", response_model=RunResponse)
async def run(request: RunRequest):
    if not settings.openai_api_key:
        raise HTTPException(
            status_code=503,
            detail="OPENAI_API_KEY is not configured.",
        )

    try:
        result = await run_task(
            request.task,
            request.mode,
            request.input,
        )
    except Exception as exc:
        # Do not expose a traceback or secret-bearing internal exception.
        raise HTTPException(
            status_code=500,
            detail="Agent execution failed.",
        ) from exc

    return RunResponse(
        status="completed",
        mode=request.mode,
        agent=result["agent"],
        output=result["output"],
        session_id=request.session_id,
    )

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(
        "app.main:app",
        host=settings.host,
        port=settings.port,
        reload=False,
    )
