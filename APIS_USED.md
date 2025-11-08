# External APIs Used in Project Synapse

This document lists all external APIs and services integrated into Project Synapse.

## Total: **8 External APIs**

---

## 1. **Claude API** (Primary AI Provider - via LiteLLM Proxy)

**Base URL**: `https://litellm-339960399182.us-central1.run.app`

**Endpoints Used**:
- **Chat Completions**: `/v1/chat/completions`
  - Models: `claude-sonnet-4-5-20250929`, `claude-opus-4-1-20250805`, `claude-haiku-4-5-20251001`
  - Used for: Summarization, categorization, tag generation, query enhancement, result re-ranking
- **Embeddings**: `/v1/embeddings`
  - Model: `gemini-embedding-001`
  - Used for: Vector embeddings for semantic search

**Purpose**: 
- AI-powered summarization (primary)
- Automatic categorization
- Tag generation
- **Search query enhancement** - converts plain English to searchable terms
- **Search result re-ranking** - improves relevance of search results
- Vector embeddings (via gemini-embedding-001)

**Authentication**: API Key (ANTHROPIC_AUTH_TOKEN)

**Note**: Uses LiteLLM proxy for unified API access

---

## 2. **Google Gemini API** (Fallback/Optional AI Provider)

**Base URL**: `https://generativelanguage.googleapis.com/`

**Endpoints Used**:
- **Text Generation**: `/v1beta/models/{model}:generateContent`
  - Models: `gemini-2.5-pro`, `gemini-2.5-flash`, `gemini-1.5-pro`, `gemini-1.5-flash`
  - Used for: Summarization, categorization, tag generation (fallback)
- **Embeddings**: `/v1beta/models/text-embedding-004:embedContent`
  - Used for: Vector embeddings for semantic search (alternative)
- **Vision API**: `/v1beta/models/gemini-pro-vision:generateContent`
  - Used for: OCR (text extraction from images)

**API Versions**: v1, v1beta

**Purpose**: 
- Fallback for AI features when Claude unavailable
- OCR from images
- Alternative embedding generation

**Authentication**: API Key (GEMINI_API_KEY)

**Rate Limits**: 
- Free tier: 60 requests/minute, 1,500 requests/day
- Paid tier: Higher limits available

**Note**: Optional - only used if `AI_PROVIDER=gemini` or as fallback

---

## 3. **OpenAI API** (Fallback AI Provider)

**Base URL**: `https://api.openai.com/v1/`

**Endpoints Used**:
- **Chat Completions**: `/chat/completions`
  - Model: `gpt-4o-mini`
  - Used for: Summarization (fallback when Gemini quota exceeded)
- **Embeddings**: `/embeddings`
  - Model: `text-embedding-3-small`
  - Used for: Vector embeddings (alternative to Gemini)

**Purpose**: 
- Fallback for summarization when Gemini quota is exceeded
- Alternative embedding generation

**Authentication**: API Key (OPENAI_API_KEY)

**Note**: Optional - only used if configured and Gemini fails

---

## 4. **YouTube API** (Embed & Thumbnails)

**Base URLs**:
- Embed: `https://www.youtube.com/embed/{videoId}`
- Thumbnails: `https://img.youtube.com/vi/{videoId}/maxresdefault.jpg`

**Purpose**: 
- Generate embed HTML for YouTube videos
- Fetch video thumbnails

**Authentication**: None required (public endpoints)

**Usage**: 
- Extracts video ID from YouTube URLs
- Generates responsive iframe embeds
- Fetches high-quality thumbnails

---

## 5. **Open Library API** (Book Covers)

**Base URLs**:
- Cover by ISBN: `https://covers.openlibrary.org/b/isbn/{isbn}-L.jpg`
- Search: `https://openlibrary.org/search.json?title={title}&limit=1`
- Cover by ID: `https://covers.openlibrary.org/b/id/{coverId}-L.jpg`

**Purpose**: 
- Fetch book covers for saved books
- Search books by title
- Get cover images by ISBN or cover ID

**Authentication**: None required (public API)

**Usage**: 
- Detects book content
- Extracts ISBN from content
- Searches by title if ISBN not found
- Fetches cover images

---

## 6. **Unsplash Source API** (Images)

**Base URL**: `https://source.unsplash.com/{width}x{height}/?{keywords}`

**Purpose**: 
- Fetch relevant images for recipes
- Fetch category-based images
- Generic image fetching when no image exists

**Authentication**: None required (public API)

**Usage**: 
- Recipe images: Searches with recipe keywords
- Category images: Searches with category + title keywords
- Format: `400x300/?keyword1,keyword2`

**Note**: Uses Unsplash Source API (simpler, no API key needed)

---

## 7. **Open Graph Protocol** (Web Page Metadata)

**Base URL**: Any public URL

**Purpose**: 
- Extract metadata from web pages
- Fetch Open Graph images
- Get page titles and descriptions

**Authentication**: None required (HTTP scraping)

**Usage**: 
- Fetches `<meta property="og:image">` tags
- Extracts page metadata
- Used for URL previews

**Implementation**: Custom HTTP client with HTML parsing

---

## 8. **ChromaDB API** (Vector Database)

**Base URL**: `http://chromadb:8000/api/v1/` (internal Docker network)

**Endpoints Used**:
- Collections: `/collections`
- Add embeddings: `/collections/{name}/add`
- Query: `/collections/{name}/query`

**Purpose**: 
- Store vector embeddings
- Semantic similarity search
- Related items discovery

**Authentication**: None (internal service)

**Note**: 
- Currently using deprecated v1 API
- Semantic search is disabled due to API deprecation
- Can be upgraded to newer ChromaDB API version

---

## Summary Table

| # | API/Service | Purpose | Authentication | Required |
|---|-------------|---------|----------------|----------|
| 1 | Claude (LiteLLM) | AI (summarization, categorization, search optimization, embeddings) | API Key | ✅ Yes |
| 2 | Google Gemini | AI fallback (summarization, embeddings, OCR) | API Key | ❌ Optional |
| 3 | OpenAI | AI fallback (summarization, embeddings) | API Key | ❌ Optional |
| 4 | YouTube | Video embeds & thumbnails | None | ✅ Yes |
| 5 | Open Library | Book covers | None | ✅ Yes |
| 6 | Unsplash Source | Recipe & category images | None | ✅ Yes |
| 7 | Open Graph | Web page metadata | None | ✅ Yes |
| 8 | ChromaDB | Vector storage & search | None | ✅ Yes |

---

## API Key Requirements

**Required**:
- `ANTHROPIC_AUTH_TOKEN` - For Claude AI features (summarization, search optimization)
- `ANTHROPIC_BASE_URL` - LiteLLM proxy URL (default: https://litellm-339960399182.us-central1.run.app)

**Optional**:
- `GEMINI_API_KEY` - For fallback AI features and OCR
- `OPENAI_API_KEY` - For additional fallback support

**No API Keys Needed**:
- YouTube (public endpoints)
- Open Library (public API)
- Unsplash Source (public API)
- Open Graph (HTTP scraping)
- ChromaDB (internal service)

---

## Rate Limits & Quotas

1. **Claude API (via LiteLLM)**: 
   - Depends on your API provider/plan
   - Check with your provider for specific limits

2. **Gemini API**: 
   - Free: 60 req/min, 1,500 req/day
   - Paid: Higher limits

3. **OpenAI API**: 
   - Varies by plan
   - Pay-per-use pricing

4. **Other APIs**: 
   - No rate limits (or very high limits)
   - Public/free services

---

## Error Handling

- **Claude API Failures**: Automatically falls back to Gemini or OpenAI if configured
- **Gemini Quota Exceeded**: Automatically falls back to OpenAI if configured
- **API Failures**: Graceful degradation - core functionality remains unaffected
- **Missing Images**: Falls back through multiple image sources
- **ChromaDB Unavailable**: Falls back to PostgreSQL text search only
- **Search Optimization**: If Claude query enhancement fails, uses original query

---

## Future API Integrations (Potential)

- **Pexels API**: Alternative image source
- **Google Books API**: Enhanced book metadata
- **Spotify API**: Music/audio content
- **Twitter/X API**: Social media content
- **Reddit API**: Forum discussions

