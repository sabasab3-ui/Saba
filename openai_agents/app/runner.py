from typing import Optional
from agents import Runner
from .agents import build_agents
from .config import settings

AGENTS = build_agents()

async def run_task(
    task: str,
    mode: str = "auto",
    extra_input: Optional[str] = None,
):
    prompt = task.strip()
    if extra_input:
        prompt += f"\n\nAdditional context:\n{extra_input.strip()}"

    result = await Runner.run(
        AGENTS[mode],
        prompt,
        max_turns=settings.max_turns,
    )

    return {
        "agent": result.last_agent.name,
        "output": result.final_output,
    }
