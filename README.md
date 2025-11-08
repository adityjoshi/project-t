# Project Synapse - Your Second Brain

A private, intelligent knowledge management system that captures, understands, and retrieves your thoughts using AI.

## Features

- 🧠 **Capture Any Thought**: Save text, URLs, and images instantly
- 🔍 **Semantic Search**: Find items by meaning, not just keywords
- 🤖 **AI-Powered**: Auto-summarization and auto-tagging
- 🔗 **Related Items**: Discover connections between your saved items
- 🌐 **Browser Extension**: One-click capture from any webpage
- 📱 **Beautiful Web Interface**: Clean, modern UI

## Tech Stack

### Backend
- **Go 1.21+** with Gin framework
- **PostgreSQL** for metadata storage
- **ChromaDB** for vector embeddings
- **OpenAI API** for embeddings, summarization, and tagging

### Frontend
- **React** with Vite
- **Tailwind CSS** for styling
- **React Router** for navigation

### Browser Extension
- Chrome Extension Manifest V3
- Content scripts for page extraction

## Prerequisites

1. **Go 1.21+** - [Install Go](https://golang.org/doc/install)
2. **Node.js 18+** - [Install Node.js](https://nodejs.org/)
3. **PostgreSQL** - [Install PostgreSQL](https://www.postgresql.org/download/)
4. **ChromaDB** - Install via pip: `pip install chromadb`
5. **OpenAI API Key** - Get from [OpenAI](https://platform.openai.com/api-keys)

## Quick Start (Docker - Recommended)

The easiest way to run Project Synapse is using Docker Compose:

### 1. Get OpenAI API Key

1. Go to [https://platform.openai.com/api-keys](https://platform.openai.com/api-keys)
2. Sign up/login and create a new API key
3. Copy your key (starts with `sk-`)

### 2. Set Environment Variable

```bash
# Create .env file in project root
echo "OPENAI_API_KEY=sk-your-key-here" > .env
```

### 3. Run with Docker

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f
```

### 4. Access the Application

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080

**See [DOCKER.md](DOCKER.md) for detailed Docker instructions and troubleshooting.**

---

## Manual Setup (Without Docker)

### 1. Database Setup

Create a PostgreSQL database:

```bash
createdb synapse
```

Or using psql:

```sql
CREATE DATABASE synapse;
```

### 2. ChromaDB Setup

Start ChromaDB server:

```bash
chroma run --path ./chroma_db
```

This will start ChromaDB on `http://localhost:8000`

### 3. Backend Setup

```bash
cd backend

# Install dependencies
go mod download

# Create .env file (copy from .env.example)
cp ../.env.example .env

# Edit .env and add your OpenAI API key
# DATABASE_URL=postgres://postgres:postgres@localhost:5432/synapse?sslmode=disable
# CHROMA_URL=http://localhost:8000
# OPENAI_API_KEY=your_key_here

# Run the server
go run cmd/server/main.go
```

The backend will start on `http://localhost:8080`

### 4. Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

The frontend will start on `http://localhost:3000`

### 5. Browser Extension Setup

1. Open Chrome and navigate to `chrome://extensions/`
2. Enable "Developer mode" (toggle in top right)
3. Click "Load unpacked"
4. Select the `extension` folder
5. The extension icon should appear in your toolbar

**Note**: Extension icons are placeholders. You can create custom icons (16x16, 48x48, 128x128 PNG files) and replace the files in `extension/icons/`.

## Usage

### Web Interface

1. Open `http://localhost:3000`
2. Click "+ Capture" to save a new item
3. Use the search bar to find items semantically
4. Click on any item to view details and related items

### Browser Extension

1. Navigate to any webpage
2. Click the Synapse extension icon
3. The popup will show the page title and URL
4. Optionally select text on the page before opening the extension
5. Click "Fill from Page" to extract main content
6. Add notes if needed
7. Click "Save" to capture

## API Endpoints

- `POST /api/items` - Create a new item
- `GET /api/items` - List all items
- `GET /api/items/:id` - Get item details
- `GET /api/items/:id/related` - Get related items
- `DELETE /api/items/:id` - Delete an item
- `GET /api/search?q=query` - Semantic search
- `GET /health` - Health check

## Project Structure

```
project-synapse/
├── backend/
│   ├── cmd/server/        # Main application entry
│   ├── internal/
│   │   ├── handlers/      # HTTP handlers
│   │   ├── services/      # Business logic
│   │   ├── repository/    # Database access
│   │   ├── models/        # Data models
│   │   └── db/            # Database connections
│   └── go.mod
├── frontend/
│   ├── src/
│   │   ├── components/    # React components
│   │   ├── services/      # API client
│   │   └── App.jsx
│   └── package.json
├── extension/
│   ├── manifest.json
│   ├── popup.html/js
│   ├── content.js
│   └── background.js
└── README.md
```

## Development

### Running in Development

1. Start ChromaDB: `chroma run --path ./chroma_db`
2. Start Backend: `cd backend && go run cmd/server/main.go`
3. Start Frontend: `cd frontend && npm run dev`

### Environment Variables

Create a `.env` file in the `backend` directory:

```env
DATABASE_URL=postgres://postgres:postgres@localhost:5432/synapse?sslmode=disable
CHROMA_URL=http://localhost:8000
OPENAI_API_KEY=your_openai_api_key_here
PORT=8080
```

## Features in Detail

### Auto-Summarization
When you save an item, the system automatically generates a 2-3 sentence summary using GPT-4.

### Auto-Tagging
The system extracts 3-5 relevant tags from your content automatically.

### Semantic Search
Search works by understanding the meaning of your query, not just matching keywords. Powered by OpenAI embeddings and vector similarity.

### Related Items
The system discovers connections between your saved items by finding similar content using vector embeddings.

## Troubleshooting

### ChromaDB Connection Issues
- Make sure ChromaDB is running: `chroma run --path ./chroma_db`
- Check that `CHROMA_URL` in `.env` matches the ChromaDB server URL

### Database Connection Issues
- Verify PostgreSQL is running
- Check `DATABASE_URL` in `.env` is correct
- Ensure the database `synapse` exists

### OpenAI API Errors
- Verify your API key is correct in `.env`
- Check you have API credits available
- Ensure you're using a valid OpenAI API key

### Extension Not Working
- Make sure the backend is running on `http://localhost:8080`
- Check browser console for errors
- Verify extension permissions in `chrome://extensions/`

## License

MIT

## Contributing

This is a personal project, but feel free to fork and modify for your own use!

