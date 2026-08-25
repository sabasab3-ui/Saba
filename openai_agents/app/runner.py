import asyncio
from typing import Optional
from agents import Runner

from .agents import build_agents
from .config import settings

AGENTS = build_agents()

async def run_task(task: str, mode: str = "auto", extra_input: Optional[str] = None):
    prompt = task
    if extra_input:
        prompt += f"\n\nAdditional context:\n{extra_input}"

    agent = AGENTS[mode]
    result = await Runner.run(
        agent,
        prompt,
        max_turns=settings.max_turns,
    )
    return {
        "agent": result.last_agent.name,
        "output": result.final_output,
    }

def run_task_sync(task: str, mode: str = "auto", extra_input: Optional[str] = None):
    return asyncio.run(run_task(task, mode, extra_input))
