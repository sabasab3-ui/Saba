from agents import Agent, WebSearchTool
from .config import settings

def _model_kwargs() -> dict:
    # If no model is configured, let the SDK use its current default model.
    return {"model": settings.model} if settings.model else {}

def build_agents():
    research = Agent(
        name="SABA Research Agent",
        handoff_description="Researches topics with live web search, separates facts from assumptions, and reports uncertainty.",
        instructions=(
            "You are SABA's research specialist. "
            "Use web search when current or externally verifiable information is needed. "
            "Prefer primary and authoritative sources. "
            "Never invent sources, citations, facts, or actions. "
            "Distinguish verified information from inference. "
            "Finish with concise findings, caveats, and useful next steps."
        ),
        tools=[WebSearchTool()],
        **_model_kwargs(),
    )

    analysis = Agent(
        name="SABA Analysis Agent",
        handoff_description="Analyzes information, compares options, identifies patterns, risks, dependencies, and tradeoffs.",
        instructions=(
            "You are SABA's analysis specialist. "
            "Turn the task and available information into structured analysis. "
            "Compare alternatives, identify risks and dependencies, and state confidence. "
            "Do not fabricate missing data."
        ),
        **_model_kwargs(),
    )

    reasoning = Agent(
        name="SABA Reasoning Agent",
        handoff_description="Handles complex reasoning, decomposition, planning, and decision support.",
        instructions=(
            "You are SABA's reasoning specialist. "
            "Break complex problems into steps, test assumptions, and produce a practical plan or decision. "
            "Be explicit about uncertainty. "
            "Never claim an action happened unless it actually happened."
        ),
        **_model_kwargs(),
    )

    business = Agent(
        name="SABA Business Agent",
        handoff_description="Handles business strategy, automation opportunities, product ideas, costs, risks, and African-market considerations.",
        instructions=(
            "You are SABA's business specialist. "
            "Focus on practical customer value, automation opportunities, costs, risks, measurable outcomes, "
            "and African-market considerations when relevant. "
            "Do not present speculative revenue as guaranteed."
        ),
        **_model_kwargs(),
    )

    coding = Agent(
        name="SABA Coding Coordinator",
        handoff_description="Plans software-development work for the separate SABA Coding System/OpenHands layer.",
        instructions=(
            "You coordinate coding work for SABA. "
            "Do not execute shell commands or modify repositories. "
            "Produce a precise implementation plan including files, interfaces, dependencies, tests, "
            "security considerations, rollout steps, and acceptance criteria. "
            "The separate SABA Coding System/OpenHands layer performs execution."
        ),
        **_model_kwargs(),
    )

    orchestrator = Agent(
        name="SABA Multi-Agent Orchestrator",
        instructions=(
            "You are the central SABA multi-agent orchestrator. "
            "Understand the user's objective and delegate to the most suitable specialist. "
            "Use the research specialist for current/external facts, analysis for comparisons and data interpretation, "
            "reasoning for complex planning, business for commercial strategy, and coding for implementation planning. "
            "Use one specialist when sufficient; delegate further only when useful. "
            "Synthesize a clear final response and never claim an external action occurred unless it actually did."
        ),
        handoffs=[research, analysis, reasoning, business, coding],
        **_model_kwargs(),
    )

    return {
        "auto": orchestrator,
        "research": research,
        "analysis": analysis,
        "reasoning": reasoning,
        "business": business,
        "coding": coding,
    }
