# SABA Development Guide

## Getting Started

### Fork & Clone
```bash
git clone https://github.com/sabasab3-ui/Saba.git
cd Saba
git checkout -b feature/my-feature
```

### Setup Development Environment

#### Go Backend
```bash
# Install Go 1.20+
go version

# Download dependencies
go mod download
go mod tidy

# Run tests
go test ./... -v

# Build
go build -o saba ./cmd/saba

# Run
./saba
```

#### Python Agents
```bash
cd openai_agents
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
pip install -e .  # Development mode

# Run tests
pytest tests/ -v

# Run service
python -m app.main
```

## Code Structure

```
Saba/
├── cmd/
│   └── saba/           # CLI entry point
├── internal/
│   ├── agent/          # Agent implementations
│   ├── database/       # Database layer
│   ├── gateway/        # API gateway
│   └── intelligence/   # BI & analysis
├── openai_agents/
│   ├── app/            # FastAPI application
│   ├── tests/          # Tests
│   └── integration/    # External connectors
├── web/                # Frontend assets
│   ├── css/
│   ├── js/
│   └── dashboard.html
├── docs/               # Documentation
└── docker-compose.yml  # Container setup
```

## Development Workflow

### 1. Create Feature Branch
```bash
git checkout -b feature/add-new-agent
```

### 2. Make Changes
```bash
# Edit files
nano internal/agent/new_agent.go
```

### 3. Test Locally
```bash
# Run tests
go test ./internal/agent

# Manual testing
curl -X POST http://localhost:8080/api/agents/run \
  -H 'Content-Type: application/json' \
  -d '{"task": "Test task"}'
```

### 4. Commit Changes
```bash
git add .
git commit -m "feat: add new agent capability"
```

### 5. Push & Create PR
```bash
git push origin feature/add-new-agent
# Create PR on GitHub
```

## Adding New Features

### Adding an Agent

1. Create agent file in `internal/agent/`:
```go
package agent

type MyAgent struct {
    name string
    tools []Tool
}

func (a *MyAgent) Process(ctx context.Context, task string) (string, error) {
    // Implementation
    return result, nil
}
```

2. Register in agent registry:
```go
registry.Register("my_agent", myAgent)
```

### Adding a Database Schema

1. Create migration file:
```go
// internal/database/my_table.go
func (db *DB) CreateMyTable(ctx context.Context) error {
    query := `CREATE TABLE IF NOT EXISTS my_table (...)`
    _, err := db.conn.ExecContext(ctx, query)
    return err
}
```

2. Call in initialization:
```go
db.CreateMyTable(ctx)
```

### Adding Python Module

1. Create module:
```python
# openai_agents/app/my_module.py
class MyFeature:
    def __init__(self):
        pass
    
    async def perform_action(self):
        pass
```

2. Import in main.py:
```python
from .my_module import MyFeature
```

## Testing

### Unit Tests
```go
// internal/agent/agent_test.go
func TestAgent(t *testing.T) {
    agent := NewAgent()
    result, err := agent.Process(context.Background(), "test")
    assert.NoError(t, err)
    assert.NotEmpty(t, result)
}
```

```python
# openai_agents/tests/test_agents.py
def test_agent_initialization():
    agent = build_agents()
    assert "research" in agent
    assert "analysis" in agent
```

### Integration Tests
```bash
# Run full test suite
./test.sh

# Run specific tests
go test ./internal/database -v
pytest openai_agents/tests/test_config.py -v
```

## Code Standards

### Go
- Follow `go fmt` standards
- Use meaningful variable names
- Add comments for exported functions
- Handle errors explicitly
- Write tests for new code

### Python
- Follow PEP 8
- Use type hints
- Write docstrings
- Use async/await
- Add pytest tests

## Debugging

### Go
```bash
# Run with debug output
GO_DEBUG=1 ./saba

# Use delve debugger
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug ./cmd/saba
```

### Python
```python
# Add debugging
import pdb; pdb.set_trace()

# Or use logging
import logging
logging.basicConfig(level=logging.DEBUG)
logger = logging.getLogger(__name__)
logger.debug("Debug message")
```

## Performance Profiling

### Go
```bash
# CPU profiling
go test -cpuprofile=cpu.prof
go tool pprof cpu.prof
```

### Python
```python
import cProfile
cProfile.run('main_function()')
```

## PR Checklist

- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] No breaking changes
- [ ] Code follows style guide
- [ ] All tests pass
- [ ] Clear commit messages
- [ ] No debug code left

## Questions?

Open an issue on GitHub or check documentation in `docs/`
