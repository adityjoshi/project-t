# Project Synapse - Complete Features & Services Documentation

## Table of Contents
1. [Overview](#overview)
2. [Core Features](#core-features)
3. [AI-Powered Services](#ai-powered-services)
4. [Content Types & Specialized Views](#content-types--specialized-views)
5. [Browser Extension](#browser-extension)
6. [Search & Discovery](#search--discovery)
7. [API Endpoints](#api-endpoints)
8. [Architecture & Services](#architecture--services)
9. [Technology Stack](#technology-stack)

---

## Overview

**Project Synapse** is an intelligent, AI-powered knowledge management system designed to be your "second brain." It captures, understands, organizes, and retrieves information from various sources, making it easy to build and search your personal knowledge base.

### Key Capabilities
- **Universal Content Capture**: Save text, URLs, images, videos, books, recipes, products, and more
- **AI-Powered Intelligence**: Automatic categorization, tagging, summarization, and image fetching
- **Semantic Search**: Find content by meaning, not just keywords
- **Smart Organization**: Automatic categorization into predefined sections
- **Related Items Discovery**: Discover connections between saved content
- **Beautiful UI**: Specialized views for different content types

---

## Core Features

### 1. Content Capture

#### Manual Entry
- Save text notes, ideas, and thoughts
- Add URLs and web links
- Upload images and screenshots
- Create to-do lists

#### Browser Extension
- One-click capture from any webpage
- Context menu integration (right-click to save)
- Automatic content extraction
- Support for multiple content types

#### Quick Save Options
- **Save Page**: Capture entire webpage
- **Save Selected Text**: Capture highlighted text
- **Save Screenshot**: Capture visible area
- **Fill from Page**: Auto-extract content from current page

### 2. Content Types Supported

The system intelligently detects and handles:

- **Text/Notes**: Plain text, notes, ideas
- **URLs/Articles**: Web pages, blog posts, articles
- **Videos**: YouTube, Vimeo, and other video platforms
- **Images**: Screenshots, photos, images
- **Books**: Book titles with automatic cover fetching
- **Recipes**: Cooking recipes with ingredient parsing
- **Products**: Amazon products with price and rating
- **To-Do Lists**: Structured task lists with progress tracking

### 3. Content Organization

#### Automatic Categorization
Content is automatically categorized into:
- Technology
- Food & Recipes
- Books & Reading
- Videos & Entertainment
- Shopping & Products
- Articles & News
- Notes & Ideas
- Design & Inspiration
- Travel
- Health & Fitness
- Education & Learning
- Other

#### Tags
- Automatic AI-generated tags
- Keyword extraction from content
- Searchable tag system

---

## AI-Powered Services

### 1. Automatic Categorization Service

**Service**: `AIService.CategorizeContent()`

**What it does**:
- Analyzes title, content, and type
- Categorizes into predefined sections
- Uses Claude AI (via LiteLLM proxy) for intelligent classification
- Falls back to type-based defaults if AI fails

**How it works**:
1. Sends content to Claude AI
2. AI analyzes semantic meaning
3. Returns appropriate category
4. Category displayed as purple badge in UI

### 2. Automatic Tagging Service

**Service**: `AIService.GenerateTags()`

**What it does**:
- Generates relevant tags automatically
- Extracts key concepts and topics
- Creates searchable metadata

**How it works**:
1. Analyzes content with Claude AI
2. Claude extracts important keywords and concepts
3. Generates 3-5 relevant tags
4. Stores as searchable array in database

### 3. Semantic Summarization Service

**Service**: `AIService.GenerateSemanticSummary()` / `AIService.SummarizeYouTubeVideo()`

**What it does**:
- Generates concise summaries (2-3 sentences) using Claude AI
- Optimized for search with key concepts
- Video-specific summarization for YouTube content
- Asynchronous processing (non-blocking)

**How it works**:
1. Item saved immediately with temporary summary
2. Background goroutine processes summary
3. AI generates semantic summary
4. Summary updated in database asynchronously

**For Videos**:
- Extracts full description from YouTube
- Generates short, focused summary (2-3 sentences)
- Preserves full description in content field
- Summary appears after a few seconds

### 4. Embedding Generation Service

**Service**: `AIService.GenerateEmbedding()`

**What it does**:
- Creates vector embeddings for semantic search
- Uses gemini-embedding-001 model (via Claude/LiteLLM proxy) or Gemini text-embedding-004
- Stores in ChromaDB for similarity search

**How it works**:
1. Content sent to embedding API (via LiteLLM proxy when using Claude)
2. Returns vector embeddings
3. Stored in ChromaDB with item ID
4. Used for semantic similarity search

### 5. OCR (Optical Character Recognition) Service

**Service**: `OCRService.ExtractTextFromImage()`

**What it does**:
- Extracts text from images and screenshots
- Uses Gemini Vision API (optional fallback)
- Makes image content searchable
- Asynchronous processing

**How it works**:
1. Image uploaded or screenshot taken
2. Background process extracts text
3. Text stored in `ocr_text` field
4. Included in search queries

### 6. Metadata Extraction Service

**Service**: `MetadataService.GetURLMetadata()`

**What it does**:
- Extracts metadata from URLs
- Fetches Open Graph images
- Generates embed HTML for videos
- Detects books and recipes

**Features**:
- **YouTube Videos**: Extracts description, thumbnail, generates embed HTML
- **Books**: Detects ISBN, fetches cover from Open Library
- **Recipes**: Identifies recipe content, fetches images
- **Generic URLs**: Extracts title, description, images

### 7. Image Fetching Service

**Service**: `MetadataService.FetchRelevantImage()`

**What it does**:
- Automatically fetches relevant images when none exist
- Uses category and title keywords
- Multiple image sources

**Image Sources** (in priority order):
1. Pre-extracted images (from extension)
2. YouTube thumbnails (for videos)
3. Book covers (Open Library API)
4. Recipe images (Unsplash)
5. Open Graph images (for URLs)
6. Category-based images (Unsplash with keywords)

**How it works**:
1. Checks for existing images
2. Tries type-specific sources
3. Extracts keywords from title
4. Searches Unsplash by category + keywords
5. Falls back to generic category images

---

## Content Types & Specialized Views

### 1. Video Cards & Views

**Component**: `VideoCard.jsx`, `ItemDetail.jsx` (video section)

**Features**:
- Thumbnail display with play button overlay
- Embedded video player in detail view
- Full description extraction from YouTube
- Channel and platform information
- Automatic categorization as "Videos & Entertainment"
- Responsive embed HTML generation

**Special Handling**:
- YouTube videos always categorized as "Videos & Entertainment"
- Full description saved in content
- Short AI summary generated asynchronously
- Video ID verification to prevent stale data

### 2. Book Cards & Views

**Component**: `BookCard.jsx`

**Features**:
- Book cover display
- Automatic cover fetching from Open Library
- ISBN detection
- Title and author information
- Category: "Books & Reading"

### 3. Recipe Cards & Views

**Component**: `RecipeCard.jsx`, `RecipeView.jsx`

**Features**:
- Recipe image display
- Ingredient list parsing
- Instruction formatting
- Automatic image fetching from Unsplash
- Category: "Food & Recipes"

**Recipe Parser**:
- Detects recipe keywords
- Extracts ingredients and instructions
- Formats for easy reading
- Fetches relevant recipe images

### 4. Product Cards & Views

**Component**: `ProductCard.jsx`, `ProductView.jsx`

**Features**:
- Product image display
- Price extraction and display
- Rating display
- ASIN tracking (Amazon)
- Category: "Shopping & Products"

**Amazon Integration**:
- Automatic product detection
- Price and rating extraction
- Product image capture
- Metadata storage

### 5. Article Cards & Views

**Component**: `ArticleCard.jsx`, `ReaderMode.jsx`

**Features**:
- Clean article display
- Reader mode for focused reading
- Open Graph image support
- Content extraction from web pages

**Reader Mode**:
- Distraction-free reading
- Clean typography
- Focused content display
- Toggle on/off

### 6. To-Do List Cards & Views

**Component**: `TodoCard.jsx`, `TodoView.jsx`

**Features**:
- Task list formatting
- Progress tracking
- Checkbox support
- Automatic detection from text patterns

**Detection Patterns**:
- Lines starting with `-`, `*`, `•`
- Numbered lists
- `[ ]` or `[x]` checkbox patterns
- Keywords: "todo", "to-do", "task", "item"

### 7. Default Item Cards

**Component**: `DefaultItemCard.jsx` (via `ItemCard.jsx`)

**Features**:
- Generic content display
- Image preview
- Summary display
- Category badges
- Tag display

---

## Browser Extension

### Extension Features

#### 1. Content Type Detection
- **Amazon Products**: Detects product pages, extracts price, rating, ASIN
- **YouTube Videos**: Extracts title, channel, description, thumbnail
- **Blog Posts**: Detects articles, extracts content, author, date
- **Generic Pages**: Extracts main content, title, images

#### 2. Content Extraction

**YouTube Videos**:
- Full description extraction (with "Show more" expansion)
- Video ID verification to prevent stale data
- Channel and platform information
- Thumbnail capture
- Multiple extraction strategies (internal data, DOM, meta tags)

**Amazon Products**:
- Product title and description
- Price extraction
- Rating extraction
- ASIN identification
- Product image

**Blog Posts**:
- Article content extraction
- Author detection
- Publication date
- Featured image

**Generic Pages**:
- Main content extraction
- Title extraction
- Image extraction
- Clean content (removes ads, navigation)

#### 3. Context Menu Integration

**Options**:
- **Save to Synapse**: Save current page
- **Save selected text to Synapse**: Save highlighted text
- **Save this page to Synapse**: Quick save current page

**How it works**:
1. Right-click on page or selected text
2. Choose context menu option
3. Extension popup opens with pre-filled data
4. User can edit and save

#### 4. Screenshot Capture

**Features**:
- Capture visible area
- Automatic image upload
- OCR text extraction
- Image storage

#### 5. Quick Save

**Features**:
- One-click save from context menu
- Pre-filled form data
- Fast capture workflow
- Badge notification

### Extension Architecture

**Files**:
- `manifest.json`: Extension configuration
- `popup.html/js`: Extension popup UI
- `content.js`: Content extraction scripts
- `background.js`: Service worker for context menus

**Permissions**:
- `activeTab`: Access current tab
- `storage`: Store quick save data
- `tabs`: Tab information
- `desktopCapture`: Screenshot capture
- `contextMenus`: Right-click menu
- `scripting`: Inject content scripts

---

## Search & Discovery

### 1. Natural Language Search

**Service**: `SearchService.Search()`, `QueryParser.ParseNaturalLanguageQuery()`

**Capabilities**:
- Understands natural language queries
- Extracts filters automatically
- Combines semantic and text search
- Intelligent query parsing

**Query Types Supported**:

#### Date-Based Queries
- "Show me articles about AI I saved last month"
- "Find that quote from yesterday"
- "My to-do list for yesterday"
- "Items from this month"
- "Content from 3 days ago"

**Supported Date Phrases**:
- `last month`, `yesterday`, `last week`
- `this month`, `last year`
- `X days ago` (e.g., "3 days ago")

#### Type-Based Queries
- "Show me articles about AI"
- "Find videos I saved"
- "My notes about machine learning"
- "Books I've saved"
- "Recipes for dinner"

**Supported Types**:
- `articles` / `blog` → blog posts
- `notes` / `handwritten` → text items
- `videos` / `youtube` → video content
- `products` / `amazon` → Amazon products
- `books` → books
- `recipes` → recipes
- `images` / `screenshot` → images
- `todo` / `to-do` → todo lists

#### Price-Based Queries
- "Black leather shoes under $300"
- "Products over $100"
- "Items between $50 and $200"

**Supported Patterns**:
- `under $X` / `below $X` → Maximum price
- `over $X` / `above $X` → Minimum price
- `$X to $Y` → Price range

#### Author/Source Queries
- "What did Karpathy say about tokenization?"
- "Find that quote from Karpathy"
- "Articles by John Smith"

**Patterns**:
- `from [Name]` - Searches for content from author
- `by [Name]` - Same as "from"
- `[Name] said` - Searches for quotes/mentions

#### Quote Searches
- "Find that quote about new beginnings"
- "That quote from the handwritten note"
- "Quote about new beginnings"

**How it works**:
- Detects quote-related keywords
- Searches in content, summaries, and OCR text
- Prioritizes exact phrase matches
- Uses semantic search for meaning

#### Combined Queries
- "Articles about AI I saved last month" (type + content + date)
- "Videos from Karpathy last week" (type + author + date)
- "Black shoes under $300 from Amazon" (content + price + type)

### 2. Hybrid Search System (AI-Optimized)

**Service**: `SearchService.Search()`

**Components**:

#### Query Enhancement (Claude AI)
- **Plain English Understanding**: Converts natural language queries into searchable terms
- **Query Expansion**: Adds synonyms and related terms automatically
- **Example**: "things about AI" → "artificial intelligence machine learning neural networks AI"
- **Intent Understanding**: Understands what users really want to find

#### Semantic Search
- Uses AI embeddings (gemini-embedding-001 via Claude/LiteLLM, or Gemini text-embedding-004)
- Stored in ChromaDB
- Finds items by meaning, not just keywords
- Vector similarity search
- Uses enhanced queries from Claude for better semantic matching

#### Text Search (Enhanced)
- PostgreSQL ILIKE queries with multi-term matching
- Searches titles, content, summaries, OCR text
- **Enhanced**: Matches individual terms from Claude-expanded queries
- Finds content even when exact phrase doesn't match
- Fallback when ChromaDB unavailable

#### Result Re-ranking (Claude AI)
- Claude re-ranks search results by relevance to original query
- Considers user intent and context
- Improves result ordering for better user experience

#### Result Combination
- Merges semantic and text results
- Boosts items found in both
- Removes duplicates
- Ranks by relevance score (enhanced by Claude)

### 3. Related Items Discovery

**Service**: `RelationService.FindRelatedItems()`

**What it does**:
- Finds similar content using vector embeddings
- Discovers connections between items
- Shows related items on detail pages

**How it works**:
1. Gets embedding for current item
2. Searches ChromaDB for similar embeddings
3. Calculates similarity scores
4. Returns top 5 related items
5. Displays on item detail page

### 4. Search Features

#### Query Enhancement
- **Passage Search**: Preserves context for quote searches
- **Exact Match Boosting**: Increases relevance for exact phrases
- **OCR Text Inclusion**: Searches text extracted from images

#### Filtering
- Date range filtering
- Type filtering
- Category filtering (dropdown in UI)
- Price range filtering
- Author/source filtering
- Tag filtering

#### Ranking
- Relevance score calculation
- Exact match boosting
- Semantic similarity scoring
- Combined score ranking

---

## API Endpoints

### Items API

#### Create Item
```
POST /api/items
Content-Type: application/json

{
  "title": "Item Title",
  "content": "Item content",
  "source_url": "https://example.com",
  "type": "url",
  "image_url": "https://example.com/image.jpg",
  "metadata": {
    "description": "Video description",
    "channel": "Channel Name",
    "platform": "YouTube"
  }
}
```

**Response**: Created item object

#### Get All Items
```
GET /api/items
```

**Response**: Array of all items (sorted by created_at DESC)

#### Get Item by ID
```
GET /api/items/:id
```

**Response**: Single item object

#### Delete Item
```
DELETE /api/items/:id
```

**Response**: `{"message": "item deleted"}`

#### Get Related Items
```
GET /api/items/:id/related
```

**Response**: Array of related items with similarity scores

#### Refresh Image
```
POST /api/items/:id/refresh-image
```

**Response**: Updated item object

**What it does**: Fetches a new relevant image for the item

#### Refresh Summary
```
POST /api/items/:id/refresh-summary
```

**Response**: `{"message": "Summary regeneration started", "item": {...}}`

**What it does**: Regenerates AI summary for the item (asynchronous)

### Search API

#### Natural Language Search
```
GET /api/search?q=your query&limit=10
```

**Query Parameters**:
- `q` (required): Search query (natural language)
- `limit` (optional): Maximum results (default: 10, max: 50)

**Response**: Array of search results with similarity scores

**Examples**:
- `/api/search?q=articles about AI I saved last month`
- `/api/search?q=black shoes under $300`
- `/api/search?q=that quote about new beginnings`

### Health Check

```
GET /health
```

**Response**: `{"status": "ok"}`

---

## Architecture & Services

### Backend Architecture

#### Database Layer

**PostgreSQL** (`internal/db/postgres.go`):
- Stores item metadata
- Full-text search capabilities
- Tag arrays (GIN index)
- Timestamps and relationships

**Schema**:
```sql
items (
  id UUID PRIMARY KEY,
  title TEXT,
  content TEXT,
  summary TEXT,
  source_url TEXT,
  type TEXT,
  category TEXT,
  tags TEXT[],
  embedding_id TEXT,
  image_url TEXT,
  embed_html TEXT,
  ocr_text TEXT,
  created_at TIMESTAMP
)

item_relations (
  item_id UUID,
  related_item_id UUID,
  similarity_score FLOAT,
  created_at TIMESTAMP
)
```

**ChromaDB** (`internal/db/chroma.go`):
- Vector embeddings storage
- Semantic similarity search
- Collection: "synapse_items"
- Embedding dimension: 768

#### Service Layer

**ItemService** (`internal/services/item_service.go`):
- Item creation and management
- Orchestrates AI services
- Metadata extraction
- Image fetching
- Summary generation

**AIService** (`internal/services/ai_service.go`):
- Embedding generation
- Summarization
- Tag generation
- Categorization
- Claude API integration (via LiteLLM proxy) - Primary
- Gemini API integration (optional fallback)

**SearchService** (`internal/services/search_service.go`):
- Hybrid search (semantic + text)
- Query parsing
- Result ranking
- Filter application

**MetadataService** (`internal/services/metadata_service.go`):
- URL metadata extraction
- Image fetching
- Book cover detection
- Recipe detection
- Embed HTML generation

**OCRService** (`internal/services/ocr_service.go`):
- Text extraction from images
- Gemini Vision API integration (optional fallback)
- Asynchronous processing

**RelationService** (`internal/services/relation_service.go`):
- Related items discovery
- Similarity calculation
- Relationship storage

**QueryParser** (`internal/services/query_parser.go`):
- Natural language parsing
- Filter extraction
- Date parsing
- Type detection
- Price range extraction

#### Handler Layer

**ItemHandler** (`internal/handlers/items.go`):
- HTTP request handling
- Request validation
- Response formatting
- Error handling

**SearchHandler** (`internal/handlers/search.go`):
- Search query handling
- Parameter validation
- Result formatting

#### Repository Layer

**ItemRepository** (`internal/repository/item_repo.go`):
- Database operations
- CRUD operations
- Search queries
- Update operations

**RelationRepository** (`internal/repository/relation_repo.go`):
- Relationship storage
- Similarity score storage
- Related items queries

### Frontend Architecture

#### Component Structure

**Layout Components**:
- `App.jsx`: Main application component
- `Dashboard.jsx`: Main dashboard with grid/swipe view
- `ItemDetail.jsx`: Item detail page

**Card Components**:
- `ItemCard.jsx`: Router to specialized cards
- `VideoCard.jsx`: Video display
- `BookCard.jsx`: Book display
- `RecipeCard.jsx`: Recipe display
- `ArticleCard.jsx`: Article display
- `ProductCard.jsx`: Product display
- `TodoCard.jsx`: To-do list display

**View Components**:
- `RecipeView.jsx`: Detailed recipe view
- `ProductView.jsx`: Detailed product view
- `TodoView.jsx`: Interactive to-do view
- `ReaderMode.jsx`: Clean reading mode

**Feature Components**:
- `CaptureForm.jsx`: Content capture form
- `SearchBar.jsx`: Search interface
- `RelatedItems.jsx`: Related items display
- `SwipeableView.jsx`: Swipeable card navigation

#### State Management

- React hooks (`useState`, `useEffect`)
- Context API (if needed)
- Local component state
- API service layer

#### Routing

- React Router v6
- Routes:
  - `/`: Dashboard
  - `/items/:id`: Item detail
  - `/capture`: Capture form

### Extension Architecture

#### Content Scripts

**content.js**:
- Runs on all web pages
- Extracts content based on page type
- Detects content type
- Extracts metadata
- Handles YouTube description extraction

#### Background Service Worker

**background.js**:
- Context menu creation
- Quick save handling
- Message routing
- Storage management

#### Popup

**popup.html/js**:
- Extension UI
- Form for saving content
- Quick save handling
- Fill from page functionality
- Screenshot capture

---

## Technology Stack

### Backend

- **Language**: Go 1.21+
- **Framework**: Gin (HTTP web framework)
- **Database**: PostgreSQL 15
- **Vector DB**: ChromaDB
- **AI Provider**: Claude (via LiteLLM proxy) - Primary
  - Models: claude-sonnet-4-5-20250929, claude-opus-4-1-20250805, claude-haiku-4-5-20251001
  - Embeddings: gemini-embedding-001 (via LiteLLM)
  - Used for: Summarization, categorization, tagging, search query enhancement, result re-ranking
- **AI Provider**: Google Gemini (optional fallback)
  - Models: gemini-2.5-flash, gemini-2.5-pro
  - Embeddings: text-embedding-004
  - Vision: gemini-pro-vision (OCR)
- **Image APIs**: 
  - Unsplash (recipe and category images)
  - Open Library (book covers)
  - Open Graph (web page images)

### Frontend

- **Framework**: React 18
- **Build Tool**: Vite
- **Styling**: Tailwind CSS
- **Routing**: React Router v6
- **HTTP Client**: Axios

### Extension

- **Manifest**: Chrome Extension Manifest V3
- **Languages**: JavaScript (ES6+)
- **APIs**: Chrome Extension APIs

### Infrastructure

- **Containerization**: Docker & Docker Compose
- **Web Server**: Nginx (for frontend)
- **Reverse Proxy**: Nginx (API proxying)

### Development Tools

- **Package Manager**: npm (frontend), go mod (backend)
- **Environment**: .env files
- **Logging**: Standard Go logging

---

## Data Flow

### Item Creation Flow

```
1. User captures content (extension/manual)
   ↓
2. Frontend sends POST /api/items
   ↓
3. ItemService.CreateItem()
   ├─→ Extract metadata (parallel)
   ├─→ Generate category (AI, parallel)
   ├─→ Generate tags (AI, parallel)
   ├─→ Generate embedding (AI, parallel)
   ├─→ Fetch image (if needed)
   └─→ Save to PostgreSQL
   ↓
4. Asynchronous processing (goroutines)
   ├─→ Generate AI summary
   ├─→ Extract OCR text (if image)
   └─→ Update database
   ↓
5. Return item to frontend
```

### Search Flow

```
1. User enters natural language query
   ↓
2. Frontend sends GET /api/search?q=query
   ↓
3. SearchService.Search()
   ├─→ QueryParser.ParseNaturalLanguageQuery()
   │   └─→ Extract filters (date, type, price, author)
   ├─→ Semantic Search (ChromaDB)
   ├─→ Text Search (PostgreSQL)
   └─→ Combine & Rank Results
   ↓
4. Apply filters
   ↓
5. Return ranked results
```

### Related Items Flow

```
1. User views item detail
   ↓
2. Frontend requests GET /api/items/:id/related
   ↓
3. RelationService.FindRelatedItems()
   ├─→ Get item embedding
   ├─→ Search ChromaDB for similar embeddings
   ├─→ Calculate similarity scores
   └─→ Return top 5 related items
   ↓
4. Display on detail page
```

---

## Performance Optimizations

### Asynchronous Processing
- AI operations run in background goroutines
- Non-blocking item creation
- Summary generation doesn't delay save
- OCR processing is asynchronous

### Parallel Operations
- Category, tags, and embedding generated in parallel
- Multiple goroutines for concurrent AI calls
- Faster item creation

### Caching & Optimization
- Embedding reuse (if possible)
- Image URL caching
- Database connection pooling

### Search Optimization
- GIN indexes on tags
- Indexed timestamps
- Efficient query patterns
- Result limiting

---

## Security Features

### Input Validation
- Request validation in handlers
- SQL injection prevention (parameterized queries)
- XSS prevention (React escaping)

### API Security
- CORS configuration
- Input sanitization
- Error message sanitization

### Extension Security
- Content Security Policy
- Minimal permissions
- Secure message passing

---

## Future Enhancements (Potential)

- User authentication and multi-user support
- Sharing and collaboration features
- Advanced filtering and sorting
- Export functionality (PDF, Markdown)
- Mobile app
- Offline support
- Advanced analytics
- Custom categories
- Folder/collection organization
- Full-text search improvements
- Real-time updates (WebSockets)

---

## Support & Documentation

- **Setup Guide**: See `SETUP.md`
- **Docker Guide**: See `DOCKER.md`
- **Extension Setup**: See `EXTENSION_SETUP.md`
- **Search Features**: See `SEARCH_FEATURES.md`
- **AI Features**: See `AI_FEATURES.md`
- **Quick Start**: See `QUICKSTART.md`

---

## License

[Add your license information here]

---

**Project Synapse** - Your intelligent second brain for capturing, understanding, and retrieving knowledge.

