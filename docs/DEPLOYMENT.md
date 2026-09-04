# SABA - Complete Deployment Guide

## Quick Start

### Prerequisites
- Docker & Docker Compose
- Python 3.10+
- Go 1.20+
- Node.js 16+ (optional, for frontend)
- OpenAI API Key

### 1. Local Development Setup

```bash
# Clone the repository
git clone https://github.com/sabasab3-ui/Saba.git
cd Saba

# Create environment file
cp .env.example .env

# Edit .env with your configuration
nano .env
```

### 2. Docker Deployment (Recommended)

```bash
# Build and start services
docker-compose up -d

# Check service status
docker-compose ps

# View logs
docker-compose logs -f openai_agents
```

### 3. Manual Setup

#### Go Backend
```bash
# Install dependencies
go mod download

# Run tests
go test ./...

# Build binary
go build -o saba ./cmd/saba

# Run server
./saba
```

#### Python Agents Service
```bash
# Create virtual environment
cd openai_agents
python3 -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate

# Install dependencies
pip install -r requirements.txt

# Set environment variables
export OPENAI_API_KEY=your_key_here

# Run service
python -m app.main
```

## Architecture

```
┌─────────────────┐
│   Web UI        │
│ (Dashboard)     │
└────────┬────────┘
         │
┌────────▼─────────────────┐
│   SABA Gateway (Go)      │
│ - Request routing        │
│ - Authentication         │
│ - Task orchestration     │
└────────┬─────────────────┘
         │
    ┌────┴────┬──────────────┬────────────┐
    │         │              │            │
┌───▼──┐ ┌───▼──┐      ┌────▼────┐ ┌───▼──┐
│OpenAI│ │Agent │      │Database │ │Intel │
│Agents│ │Layer │      │(SQLite) │ │(BI)  │
└──────┘ └──────┘      └─────────┘ └──────┘
```

## Configuration

### Environment Variables

```env
# Core
HOST=0.0.0.0
PORT=8080
API_PORT=8090

# OpenAI
OPENAI_API_KEY=sk-...
OPENAI_MODEL=gpt-4-turbo

# Database
DATABASE_PATH=./saba.db

# Features
ENABLE_ANALYTICS=true
ENABLE_AUTOMATIONS=true
ENABLE_BUSINESS_ANALYZER=true
```

## API Endpoints

### Health Check
```bash
GET /health
→ {"status": "ok", "service": "saba"}
```

### Run Agent
```bash
POST /api/agents/run
{
  "task": "Analyze market in Uganda",
  "mode": "auto",
  "country": "Uganda",
  "industry": "Retail"
}
```

### Get Platform Stats
```bash
GET /api/stats
→ {"health": {...}, "usage_by_agent": {...}}
```

## Database Initialization

```bash
# Connect to SQLite
sqlite3 saba.db

# Tables are auto-created on first run
# Tables: customers, orders, automations, agents, analytics
```

## Performance Tuning

### For High Load
```env
MAX_CONCURRENT_TASKS=50
CACHE_ENABLED=true
CACHE_TTL=3600
DB_BACKUP_INTERVAL=1800
```

### For Development
```env
MAX_CONCURRENT_TASKS=5
LOG_LEVEL=debug
```

## Monitoring

### Logs
```bash
# Go service
docker-compose logs -f saba

# Python agents
docker-compose logs -f openai_agents

# Web UI
tail -f logs/dashboard.log
```

### Metrics
- Access `/api/stats` for platform metrics
- Analytics data stored in database
- Performance monitoring via logs

## Troubleshooting

### OpenAI API Key Error
```bash
# Verify key is set
echo $OPENAI_API_KEY

# Test API access
curl -H "Authorization: Bearer $OPENAI_API_KEY" \
  https://api.openai.com/v1/models
```

### Database Lock
```bash
# Close all connections
docker-compose restart

# Or reset database
rm saba.db
```

### Port Already in Use
```bash
# Change port in .env
PORT=8081
API_PORT=8091

# Restart services
docker-compose restart
```

## Production Deployment

### AWS ECS
```bash
# Build and push image
aws ecr get-login-password | docker login --username AWS --password-stdin <account>.dkr.ecr.<region>.amazonaws.com
docker tag saba:latest <account>.dkr.ecr.<region>.amazonaws.com/saba:latest
docker push <account>.dkr.ecr.<region>.amazonaws.com/saba:latest
```

### Kubernetes
```yaml
apiVersion: v1
kind: Deployment
metadata:
  name: saba
spec:
  replicas: 3
  template:
    spec:
      containers:
      - name: saba
        image: saba:latest
        env:
        - name: OPENAI_API_KEY
          valueFrom:
            secretKeyRef:
              name: saba-secrets
              key: openai-api-key
```

## Security

- Always use HTTPS in production
- Rotate API keys regularly
- Use environment variables for secrets
- Enable authentication for production
- Set up firewall rules
- Regular security audits

## Support

For issues and questions:
- GitHub Issues: https://github.com/sabasab3-ui/Saba/issues
- Documentation: See `docs/` directory
- Email: support@saba.ai
