from agents import Agent
from .config import settings

def build_agents():
    research = Agent(
        name="SABA Research Agent",
        model=settings.model,
        handoff_description="Researches a topic, gathers relevant facts, identifies uncertainty, and presents useful findings.",
        instructions=(
            "You are SABA's research specialist. "
            "Investigate the user's task carefully. Separate known facts from assumptions. "
            "Do not invent sources or claim you browsed the web unless a web tool is actually available. "
            "Return concise findings, important caveats, and recommended next steps."
        ),
    )

    analysis = Agent(
        name="SABA Analysis Agent",
        model=settings.model,
        handoff_description="Analyzes information, compares options, identifies patterns, risks, and tradeoffs.",
        instructions=(
            "You are SABA's analysis specialist. "
            "Turn the supplied task and information into structured analysis. "
            "Compare alternatives, identify risks and dependencies, and state your confidence. "
            "Do not fabricate missing data."
        ),
    )

    reasoning = Agent(
        name="SABA Reasoning Agent",
        model=settings.model,
        handoff_description="Handles complex reasoning, planning, decomposition, and decision support.",
        instructions=(
            "You are SABA's reasoning specialist. "
            "Break complex problems into manageable steps, test assumptions, and produce a practical decision or plan. "
            "Be explicit about uncertainty and avoid pretending to have performed actions you did not perform."
        ),
    )

    business = Agent(
        name="SABA Business Agent",
        model=settings.model,
        handoff_description="Handles business strategy, automation opportunities, product ideas, and African-market considerations.",
        instructions=(
            "You are SABA's business specialist. "
            "Focus on practical business value, automation opportunities, customers, costs, risks, and measurable outcomes. "
            "Do not present speculative revenue as guaranteed."
        ),
    )

    coding = Agent(
        name="SABA Coding Coordinator",
        model=settings.model,
        handoff_description="Plans software-development work and prepares tasks for the separate SABA Coding System/OpenHands layer.",
        instructions=(
            "You coordinate coding work for SABA. "
            "Do not execute shell commands or modify repositories in this agent. "
            "Produce a precise implementation plan: files, interfaces, dependencies, tests, security considerations, and acceptance criteria. "
            "The separate SABA Coding System/OpenHands layer is responsible for execution."
        ),
    )

    triage = Agent(
        name="SABA Multi-Agent Orchestrator",
        model=settings.model,
        instructions=(
            "You are the central SABA multi-agent orchestrator. "
            "Understand the user's objective and delegate to the best specialist. "
            "Use one specialist when possible; use additional delegation only when it materially improves the answer. "
            "Synthesize the final answer clearly and do not claim that an external action occurred unless it actually did."
        ),
        handoffs=[research, analysis, reasoning, business, coding],
    )

    return {
        "auto": triage,
        "research": research,
        "analysis": analysis,
        "reasoning": reasoning,
        "business": business,
        "coding": coding,
    }
