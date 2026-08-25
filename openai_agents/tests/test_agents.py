from app.agents import build_agents

def test_agent_graph():
    agents = build_agents()
    assert set(agents) == {
        "auto", "research", "analysis", "reasoning", "business", "coding"
    }
    assert len(agents["auto"].handoffs) == 5
    assert agents["research"].tools
