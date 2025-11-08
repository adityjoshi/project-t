# Docker Setup Guide

This guide will help you run Project Synapse using Docker Compose.

## Prerequisites

1. **Docker** - [Install Docker](https://docs.docker.com/get-docker/)
2. **Docker Compose** - Usually included with Docker Desktop
3. **OpenAI API Key** - See instructions below

## Getting Your OpenAI API Key

### Step 1: Create an OpenAI Account
1. Go to [https://platform.openai.com](https://platform.openai.com)
2. Click "Sign up" and create an account
3. Verify your email address

### Step 2: Add Payment Method
1. Go to [https://platform.openai.com/account/billing](https://platform.openai.com/account/billing)
2. Click "Add payment method"
3. Add a credit card (OpenAI charges per API usage)

### Step 3: Get Your API Key
1. Go to [https://platform.openai.com/api-keys](https://platform.openai.com/api-keys)
2. Click "Create new secret key"
3. Give it a name (e.g., "Synapse Project")
4. **Copy the key immediately** - you won't be able to see it again!
5. Save it securely

**Important**: 
- Never share your API key publicly
- The key starts with `sk-`
- You'll be charged based on API usage (embeddings are cheap, GPT-4 is more expensive)

## Running with Docker

### Step 1: Set Your OpenAI API Key

Create a `.env` file in the project root:

```bash
# In the project root directory
echo "OPENAI_API_KEY=sk-your-actual-api-key-here" > .env
```

Or manually create `.env` file:
```env
OPENAI_API_KEY=sk-your-actual-api-key-here
```

**Replace `sk-your-actual-api-key-here` with your actual OpenAI API key!**

### Step 2: Build and Start All Services

```bash
# Build and start all containers
docker-compose up -d

# View logs
docker-compose logs -f
```

This will:
- Start PostgreSQL database
- Start ChromaDB vector database
- Build and start the Go backend
- Build and start the React frontend

### Step 3: Wait for Services to be Ready

Wait about 30-60 seconds for all services to start. Check status:

```bash
# Check if all services are running
docker-compose ps

# Check backend health
curl http://localhost:8080/health
```

### Step 4: Access the Application

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **ChromaDB**: http://localhost:8000
- **PostgreSQL**: localhost:5432

## Useful Docker Commands

### View Logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f postgres
docker-compose logs -f chromadb
```

### Stop Services
```bash
# Stop all services
docker-compose down

# Stop and remove volumes (deletes all data!)
docker-compose down -v
```

### Restart Services
```bash
# Restart all services
docker-compose restart

# Restart specific service
docker-compose restart backend
```

### Rebuild After Code Changes
```bash
# Rebuild and restart
docker-compose up -d --build

# Rebuild specific service
docker-compose up -d --build backend
```

### Check Service Status
```bash
# See running containers
docker-compose ps

# Check resource usage
docker stats
```

### Access Container Shell
```bash
# Backend container
docker-compose exec backend sh

# PostgreSQL
docker-compose exec postgres psql -U postgres -d synapse

# ChromaDB (if needed)
docker-compose exec chromadb sh
```

## Troubleshooting

### Services Won't Start
```bash
# Check logs for errors
docker-compose logs

# Verify .env file exists and has OPENAI_API_KEY
cat .env

# Check if ports are already in use
lsof -i :3000
lsof -i :8080
lsof -i :5432
lsof -i :8000
```

### Backend Can't Connect to Database
```bash
# Check if postgres is healthy
docker-compose ps postgres

# Check postgres logs
docker-compose logs postgres

# Wait a bit longer - postgres needs time to initialize
```

### ChromaDB Connection Issues
```bash
# Check chromadb logs
docker-compose logs chromadb

# Verify chromadb is accessible
curl http://localhost:8000/api/v1/heartbeat
```

### OpenAI API Errors
```bash
# Verify your API key is set
docker-compose exec backend env | grep OPENAI

# Check backend logs for API errors
docker-compose logs backend | grep -i openai
```

### Rebuild Everything from Scratch
```bash
# Stop and remove everything
docker-compose down -v

# Remove all images
docker-compose rm -f

# Rebuild and start
docker-compose up -d --build
```

## Development Workflow

### Making Code Changes

1. **Backend changes**: 
   ```bash
   docker-compose up -d --build backend
   ```

2. **Frontend changes**:
   ```bash
   docker-compose up -d --build frontend
   ```

3. **View logs**:
   ```bash
   docker-compose logs -f backend
   ```

### Running Without Docker (Local Development)

If you want to develop locally without Docker:

1. Start only databases:
   ```bash
   docker-compose up -d postgres chromadb
   ```

2. Run backend locally:
   ```bash
   cd backend
   go run cmd/server/main.go
   ```

3. Run frontend locally:
   ```bash
   cd frontend
   npm run dev
   ```

## Production Deployment

For production, you should:

1. Use environment-specific `.env` files
2. Set up proper secrets management
3. Use a reverse proxy (nginx/traefik)
4. Enable HTTPS
5. Set up database backups
6. Configure proper resource limits in docker-compose.yml

## Data Persistence

Data is stored in Docker volumes:
- `postgres_data` - PostgreSQL database
- `chromadb_data` - ChromaDB vector store

To backup:
```bash
docker-compose exec postgres pg_dump -U postgres synapse > backup.sql
```

To restore:
```bash
docker-compose exec -T postgres psql -U postgres synapse < backup.sql
```

