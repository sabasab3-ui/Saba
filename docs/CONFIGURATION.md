# SABA Configuration Reference

## Environment Variables

### Core Server
| Variable | Default | Description |
|----------|---------|-------------|
| `HOST` | `127.0.0.1` | Server bind address |
| `PORT` | `8080` | Main API port |
| `API_PORT` | `8090` | Python agents port |
| `LOG_LEVEL` | `info` | Logging verbosity |
| `LOG_FILE` | `./saba.log` | Log file path |

### OpenAI Configuration
| Variable | Default | Description |
|----------|---------|-------------|
| `OPENAI_API_KEY` | Required | Your OpenAI API key |
| `OPENAI_MODEL` | `gpt-4-turbo` | Model to use |
| `OPENAI_MAX_TURNS` | `12` | Max conversation turns |

### Database
| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_PATH` | `./saba.db` | SQLite database file |
| `DB_BACKUP_INTERVAL` | `3600` | Backup frequency (seconds) |

### Feature Flags
| Variable | Default | Description |
|----------|---------|-------------|
| `ENABLE_ANALYTICS` | `true` | Enable analytics tracking |
| `ENABLE_AUTOMATIONS` | `true` | Enable automation engine |
| `ENABLE_BUSINESS_ANALYZER` | `true` | Enable BI features |
| `ENABLE_CAPABILITY_ENGINE` | `true` | Enable capabilities API |

### Performance
| Variable | Default | Description |
|----------|---------|-------------|
| `MAX_CONCURRENT_TASKS` | `10` | Max concurrent task execution |
| `TASK_TIMEOUT` | `300` | Task timeout (seconds) |
| `CACHE_ENABLED` | `true` | Enable response caching |
| `CACHE_TTL` | `3600` | Cache time-to-live (seconds) |

### Security
| Variable | Default | Description |
|----------|---------|-------------|
| `API_KEY_REQUIRED` | `false` | Require API key for requests |
| `ALLOW_CORS` | `true` | Enable CORS headers |
| `CORS_ORIGINS` | `*` | Allowed CORS origins |

### Integration
| Variable | Default | Description |
|----------|---------|-------------|
| `SABA_GATEWAY_URL` | `http://127.0.0.1:8080` | Gateway URL |
| `SABA_ORCHESTRATOR_URL` | `http://127.0.0.1:9000` | Orchestrator URL |

## Configuration File Format

### .env File Example
```env
# Server
HOST=0.0.0.0
PORT=8080
API_PORT=8090

# OpenAI
OPENAI_API_KEY=sk-...
OPENAI_MODEL=gpt-4-turbo

# Database
DATABASE_PATH=/data/saba.db

# Features
ENABLE_ANALYTICS=true
ENABLE_AUTOMATIONS=true
```

## Runtime Configuration

### Via Environment
```bash
export OPENAI_API_KEY=sk-...
export PORT=9000
./saba
```

### Via Docker
```yaml
environment:
  - OPENAI_API_KEY=${OPENAI_API_KEY}
  - PORT=8080
  - LOG_LEVEL=debug
```

### Via Config File (Future)
```yaml
# config.yaml
server:
  host: 0.0.0.0
  port: 8080
openai:
  api_key: ${OPENAI_API_KEY}
  model: gpt-4-turbo
```

## Validation Rules

- `PORT` must be 1-65535
- `OPENAI_API_KEY` must start with `sk-`
- `TASK_TIMEOUT` must be > 0
- `MAX_CONCURRENT_TASKS` must be 1-1000
- `LOG_LEVEL` must be: debug, info, warn, error

## Defaults Loading Order

1. Hardcoded defaults
2. `.env` file values
3. Environment variables (highest priority)

## Examples

### Development
```env
HOST=127.0.0.1
PORT=8080
LOG_LEVEL=debug
CACHE_ENABLED=false
```

### Production
```env
HOST=0.0.0.0
PORT=8080
LOG_LEVEL=info
API_KEY_REQUIRED=true
CACHE_ENABLED=true
CACHE_TTL=7200
```

### High Performance
```env
MAX_CONCURRENT_TASKS=100
CACHE_ENABLED=true
CACHE_TTL=3600
DB_BACKUP_INTERVAL=7200
```
