from app.models import RunRequest

def test_default_mode():
    assert RunRequest(task="test").mode == "auto"

def test_all_modes():
    modes = ["auto", "research", "analysis", "reasoning", "business", "coding"]
    for mode in modes:
        assert RunRequest(task="test", mode=mode).mode == mode

def test_limits():
    request = RunRequest(task="hello", input="context")
    assert request.task == "hello"
    assert request.input == "context"
