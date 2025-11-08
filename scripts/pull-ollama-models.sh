#!/bin/bash

# Script to pull required Ollama models
# This is run automatically when the backend starts

OLLAMA_URL=${OLLAMA_URL:-http://localhost:11434}

echo "Pulling Ollama models..."
echo "This may take a few minutes on first run..."

# Pull embedding model
echo "Pulling nomic-embed-text (for embeddings)..."
curl -X POST "$OLLAMA_URL/api/pull" -d '{"name": "nomic-embed-text"}'

# Pull chat model
echo "Pulling llama3.2 (for summarization and tagging)..."
curl -X POST "$OLLAMA_URL/api/pull" -d '{"name": "llama3.2"}'

echo "Models pulled successfully!"

