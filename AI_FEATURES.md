# AI-Powered Features Documentation

## Overview

Project Synapse now includes a comprehensive AI-powered system that automatically categorizes content, fetches relevant images, and generates semantic summaries for enhanced searchability.

## Features Implemented

### 1. **Automatic AI Categorization**

The system automatically categorizes incoming content into specific sections using Google Gemini:

**Available Categories:**
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

**How it works:**
- When an item is created, the AI analyzes the title, content, and type
- It categorizes the content into one of the predefined sections
- If categorization fails, a default category is assigned based on item type
- Categories are displayed as purple badges in the frontend

### 2. **Automatic Image Fetching**

The system automatically fetches relevant images when none exist:

**Image Sources:**
- **Pre-extracted images**: From browser extension (Amazon products, blog posts, videos)
- **YouTube thumbnails**: For video content
- **Book covers**: From Open Library API (for books)
- **Recipe images**: From Unsplash (for recipes)
- **Open Graph images**: For URLs and blog posts
- **Category-based images**: From Unsplash based on the AI-determined category

**How it works:**
1. First checks for pre-extracted images from the extension
2. Tries type-specific image sources (YouTube, books, recipes)
3. Falls back to category-based image search using Unsplash
4. Images are automatically associated with items

### 3. **Asynchronous Semantic Summarization**

The system generates concise semantic summaries asynchronously using Gemini:

**How it works:**
1. Item is saved immediately with a temporary summary (first 200 chars)
2. A background goroutine generates a semantic summary using Gemini
3. The summary is optimized for search with key concepts and keywords
4. The database is updated asynchronously with the new summary
5. Users can search across all content using natural language

**Benefits:**
- Fast item creation (doesn't wait for summary generation)
- Better search results (semantic summaries include important keywords)
- Non-blocking user experience

### 4. **Enhanced Natural Language Search**

The search system now supports:
- **Category filtering**: "Show me technology articles"
- **Type filtering**: "Find videos about AI"
- **Date filtering**: "Articles from last month"
- **Price filtering**: "Products under $300"
- **Author filtering**: "What did Karpathy say?"
- **Tag filtering**: Using #hashtags
- **Semantic search**: Using embeddings for meaning-based search

## Technical Implementation

### Database Schema

Added `category` column to `items` table:
```sql
ALTER TABLE items ADD COLUMN category TEXT;
```

### API Changes

**Item Model:**
- Added `category` field (string)
- Category is automatically populated by AI

**Item Creation Flow:**
1. Generate category (parallel with tags and embedding)
2. Fetch images (if not provided)
3. Save item with temporary summary
4. Generate semantic summary asynchronously
5. Update item summary in background

### Frontend Changes

- Category badges displayed on item cards (purple badges)
- Categories shown in item detail view
- Search supports category filtering

## Usage Examples

### Creating Items

Items are automatically categorized and processed:

```bash
POST /api/items
{
  "title": "Introduction to Neural Networks",
  "content": "Neural networks are computing systems inspired by biological neural networks...",
  "type": "article"
}
```

**Result:**
- Automatically categorized as "Technology" or "Education & Learning"
- Relevant image fetched automatically
- Semantic summary generated asynchronously
- Ready for natural language search

### Searching by Category

```bash
GET /api/search?q=technology articles
```

Returns items in the "Technology" category containing "articles" in the content.

### Natural Language Queries

- "Show me articles about AI I saved last month"
- "Find technology content"
- "Books about machine learning"
- "Videos from last week"

## Configuration

### Environment Variables

- `GEMINI_API_KEY`: Your Google Gemini API key (required)
- `AI_PROVIDER`: Set to "gemini" (default) or "openai"

### API Keys

Get your Gemini API key from: https://makersuite.google.com/app/apikey

## Performance

- **Item Creation**: ~500-800ms (includes categorization, tagging, embedding)
- **Image Fetching**: ~200-500ms (varies by source)
- **Semantic Summary**: Generated asynchronously (doesn't block creation)
- **Search**: ~400-600ms (hybrid semantic + text search)

## Future Enhancements

- Category-based filtering in UI
- Category statistics and insights
- Custom category definitions
- Batch processing for bulk imports
- Image caching and optimization

