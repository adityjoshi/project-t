# 🚀 Quick Start Guide

## Get Your OpenAI API Key

1. **Go to OpenAI Platform**: https://platform.openai.com/api-keys
2. **Sign up/Login** to your account
3. **Add Payment Method** (required - you pay per API usage):
   - Go to: https://platform.openai.com/account/billing
   - Click "Add payment method"
   - Add your credit card
4. **Create API Key**:
   - Go back to: https://platform.openai.com/api-keys
   - Click "Create new secret key"
   - Name it (e.g., "Synapse")
   - **Copy the key immediately** (starts with `sk-`)
   - Save it somewhere safe

**Note**: OpenAI charges based on usage:
- Embeddings: ~$0.0001 per 1K tokens (very cheap)
- GPT-4: ~$0.01-0.03 per 1K tokens (more expensive)
- Typical usage: $1-5/month for personal use

## Run the Project

### Option 1: Using the Quick Start Script (Easiest)

```bash
# 1. Create .env file with your API key
echo "OPENAI_API_KEY=sk-your-actual-key-here" > .env

# 2. Run the script
./run.sh
```

### Option 2: Using Docker Compose Directly

```bash
# 1. Create .env file
echo "OPENAI_API_KEY=sk-your-actual-key-here" > .env

# 2. Start all services
docker-compose up -d

# 3. View logs
docker-compose logs -f
```

### Option 3: Manual Setup (Without Docker)

See [SETUP.md](SETUP.md) for manual installation instructions.

## Access the Application

Once services are running:

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080

Wait 30-60 seconds for all services to fully start.

## Verify Everything Works

```bash
# Check if services are running
docker-compose ps

# Check backend health
curl http://localhost:8080/health

# View logs
docker-compose logs -f
```

## Common Commands

```bash
# Start services
docker-compose up -d

# Stop services
docker-compose down

# View logs
docker-compose logs -f

# Restart a service
docker-compose restart backend

# Rebuild after code changes
docker-compose up -d --build
```

## Troubleshooting

### Services won't start?
```bash
# Check logs
docker-compose logs

# Verify .env file exists
cat .env
```

### Backend errors?
```bash
# Check backend logs
docker-compose logs backend

# Verify API key is set
docker-compose exec backend env | grep OPENAI
```

### Port already in use?
```bash
# Check what's using the ports
lsof -i :3000
lsof -i :8080
lsof -i :5432
lsof -i :8000
```

## Next Steps

1. Open http://localhost:3000 in your browser
2. Click "+ Capture" to save your first item
3. Try the semantic search
4. Install the browser extension (see README.md)

For detailed information, see:
- [DOCKER.md](DOCKER.md) - Complete Docker guide
- [README.md](README.md) - Full documentation
- [SETUP.md](SETUP.md) - Manual setup guide

