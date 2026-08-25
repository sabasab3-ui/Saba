from app.models import RunRequest

def test_request_defaults():
    req = RunRequest(task="test")
    assert req.mode == "auto"

def test_modes():
    for mode in ["auto", "research", "analysis", "reasoning", "business", "coding"]:
        assert RunRequest(task="test", mode=mode).mode == mode
