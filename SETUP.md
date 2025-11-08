# Quick Setup Guide

## Step 1: Install Prerequisites

### PostgreSQL
```bash
# macOS
brew install postgresql
brew services start postgresql

# Linux (Ubuntu/Debian)
sudo apt-get install postgresql postgresql-contrib
sudo systemctl start postgresql

# Create database
createdb synapse
```

### ChromaDB
```bash
pip install chromadb
```

### Go & Node.js
- Install Go: https://golang.org/doc/install
- Install Node.js: https://nodejs.org/

## Step 2: Start ChromaDB

```bash
chroma run --path ./chroma_db
```

Keep this terminal open. ChromaDB will run on `http://localhost:8000`

## Step 3: Configure Backend

```bash
cd backend

# Copy environment template
cp ../.env.example .env

# Edit .env and add your OpenAI API key
# Required: OPENAI_API_KEY=sk-...
```

## Step 4: Start Backend

```bash
cd backend
go mod download
go run cmd/server/main.go
```

Backend runs on `http://localhost:8080`

## Step 5: Start Frontend

```bash
cd frontend
npm install
npm run dev
```

Frontend runs on `http://localhost:3000`

## Step 6: Install Browser Extension

1. Open Chrome: `chrome://extensions/`
2. Enable "Developer mode"
3. Click "Load unpacked"
4. Select the `extension` folder

**Note**: Extension icons are placeholders. The extension will work without custom icons, but you can create 16x16, 48x48, and 128x128 PNG icons and place them in `extension/icons/` for a better experience.

## You're Ready!

1. Open `http://localhost:3000`
2. Click "+ Capture" to save your first item
3. Use the browser extension to capture web pages
4. Search your knowledge base semantically

## Troubleshooting

### ChromaDB not starting?
- Make sure port 8000 is available
- Try: `chroma run --path ./chroma_db --port 8000`

### Database connection error?
- Verify PostgreSQL is running: `pg_isready`
- Check database exists: `psql -l | grep synapse`
- Verify DATABASE_URL in `.env`

### OpenAI API errors?
- Check your API key is correct
- Verify you have API credits
- Check rate limits

### Extension not working?
- Make sure backend is running on port 8080
- Check browser console for errors
- Verify extension permissions

