# 🚀 Run Commands - Quick Reference

## Prerequisites Check

```bash
# Check if Docker is installed
docker --version

# Check if Docker Compose is installed
docker-compose --version

# Check if Docker is running
docker info
```

## Step-by-Step Run Commands

### 1. Get OpenAI API Key

1. Visit: https://platform.openai.com/api-keys
2. Sign up/login
3. Add payment method at: https://platform.openai.com/account/billing
4. Create API key (starts with `sk-`)
5. Copy the key

### 2. Create Environment File

```bash
# Navigate to project directory
cd /Users/aditya/Desktop/project-t

# Create .env file with your API key
echo "OPENAI_API_KEY=sk-your-actual-api-key-here" > .env

# Verify it was created
cat .env
```

**Important**: Replace `sk-your-actual-api-key-here` with your real OpenAI API key!

### 3. Start All Services

```bash
# Build and start all containers (first time)
docker-compose up -d --build

# Or if already built, just start
docker-compose up -d
```

### 4. Check Service Status

```bash
# See all running containers
docker-compose ps

# Check if backend is healthy
curl http://localhost:8080/health

# Should return: {"status":"ok"}
```

### 5. View Logs

```bash
# View all logs
docker-compose logs -f

# View specific service logs
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f postgres
docker-compose logs -f chromadb
```

### 6. Access the Application

Open in your browser:
- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080

## Common Commands

### Start Services
```bash
docker-compose up -d
```

### Stop Services
```bash
docker-compose down
```

### Stop and Remove All Data
```bash
docker-compose down -v
```

### Restart Services
```bash
docker-compose restart
```

### Restart Specific Service
```bash
docker-compose restart backend
docker-compose restart frontend
```

### Rebuild After Code Changes
```bash
# Rebuild all
docker-compose up -d --build

# Rebuild specific service
docker-compose up -d --build backend
docker-compose up -d --build frontend
```

### View Resource Usage
```bash
docker stats
```

### Access Container Shell
```bash
# Backend container
docker-compose exec backend sh

# PostgreSQL
docker-compose exec postgres psql -U postgres -d synapse

# Frontend (nginx)
docker-compose exec frontend sh
```

## Troubleshooting Commands

### Check if Ports are Available
```bash
# Check what's using the ports
lsof -i :3000
lsof -i :8080
lsof -i :5432
lsof -i :8000
```

### Check Environment Variables
```bash
# Check if API key is set in container
docker-compose exec backend env | grep OPENAI
```

### Check Database Connection
```bash
# Test PostgreSQL connection
docker-compose exec postgres psql -U postgres -d synapse -c "SELECT 1;"
```

### Check ChromaDB
```bash
# Test ChromaDB health
curl http://localhost:8000/api/v1/heartbeat
```

### View Detailed Logs
```bash
# Last 100 lines
docker-compose logs --tail=100

# Follow logs with timestamps
docker-compose logs -f --timestamps
```

### Clean Everything and Start Fresh
```bash
# Stop and remove everything
docker-compose down -v

# Remove all images
docker-compose rm -f

# Remove unused images
docker image prune -a

# Rebuild from scratch
docker-compose up -d --build
```

## Quick Start Script

You can also use the provided script:

```bash
# Make sure it's executable
chmod +x run.sh

# Run it
./run.sh
```

## Development Workflow

### Making Backend Changes
```bash
# After editing Go code
docker-compose up -d --build backend
docker-compose logs -f backend
```

### Making Frontend Changes
```bash
# After editing React code
docker-compose up -d --build frontend
docker-compose logs -f frontend
```

### Running Tests Locally (Optional)
```bash
# Start only databases
docker-compose up -d postgres chromadb

# Run backend locally
cd backend
go run cmd/server/main.go

# Run frontend locally (in another terminal)
cd frontend
npm run dev
```

## Production Deployment

For production, you'll want to:

1. Use proper secrets management
2. Set resource limits in docker-compose.yml
3. Use environment-specific .env files
4. Set up reverse proxy with HTTPS
5. Configure backups

Example production docker-compose override:
```bash
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

