# Natural Language Search Features

## Overview

Synapse now supports powerful natural language search that understands your queries and automatically applies filters. The search combines:
- **Semantic Search**: Uses AI embeddings for meaning-based search (when ChromaDB is available)
- **Text Search**: PostgreSQL full-text search as fallback
- **Smart Filters**: Automatically extracts filters from natural language queries

## Supported Query Types

### 1. Date-Based Queries

**Examples:**
- "Show me articles about AI I saved last month"
- "Find that quote from yesterday"
- "My to-do list for yesterday"
- "What did I save last week?"
- "Items from this month"
- "Content from 3 days ago"

**Supported Date Phrases:**
- `last month` - Items from the previous month
- `yesterday` - Items from yesterday
- `last week` - Items from the past 7 days
- `this month` - Items from the current month
- `last year` - Items from the previous year
- `X days ago` - Items from X days ago (e.g., "3 days ago")

### 2. Type-Based Queries

**Examples:**
- "Show me articles about AI"
- "Find videos I saved"
- "My notes about machine learning"
- "Books I've saved"
- "Recipes for dinner"
- "Screenshots from last week"

**Supported Types:**
- `articles` / `article` / `blog` → blog posts
- `notes` / `note` / `handwritten` → text items
- `videos` / `video` / `youtube` → video content
- `products` / `amazon` → Amazon products
- `books` / `book` → books
- `recipes` / `recipe` → recipes
- `images` / `screenshot` → images/screenshots
- `todo` / `to-do` / `list` → todo lists

### 3. Price-Based Queries

**Examples:**
- "Black leather shoes under $300"
- "Products over $100"
- "Items between $50 and $200"
- "Shoes below $250"

**Supported Price Patterns:**
- `under $X` / `below $X` / `less than $X` → Maximum price
- `over $X` / `above $X` / `more than $X` → Minimum price
- `$X to $Y` / `$X-$Y` → Price range

### 4. Author/Source Queries

**Examples:**
- "What did Karpathy say about tokenization?"
- "Find that quote from Karpathy"
- "Articles by John Smith"
- "Papers from that researcher"

**Patterns:**
- `from [Name]` - Searches for content from a specific author
- `by [Name]` - Same as "from"
- `[Name] said` - Searches for quotes/mentions

### 5. Content-Based Queries

**Examples:**
- "Find that quote about new beginnings"
- "Articles about machine learning"
- "Inspiration for a travel blog logo"
- "Content about AI and neural networks"

**How it works:**
- Searches in titles, content, and summaries
- Uses semantic search when available (understands meaning)
- Falls back to text search (keyword matching)

### 6. Combined Queries

You can combine multiple filters in a single query:

**Examples:**
- "Articles about AI I saved last month" (type + content + date)
- "Videos from Karpathy last week" (type + author + date)
- "Black shoes under $300 from Amazon" (content + price + type)
- "My to-do list for yesterday" (type + date)

## How It Works

1. **Query Parsing**: The system parses your natural language query to extract:
   - Search terms (what you're looking for)
   - Date filters (when it was saved)
   - Type filters (what kind of content)
   - Price filters (for products)
   - Author/source filters

2. **Hybrid Search**:
   - **Semantic Search**: Uses AI embeddings to find items with similar meaning (when ChromaDB is working)
   - **Text Search**: Uses PostgreSQL to search titles, content, and summaries
   - **Combined Results**: Merges both result sets, boosting items found in both

3. **Filtering**: Applies all extracted filters to narrow down results

4. **Ranking**: Results are ranked by relevance score

## Technical Details

### Search Service Architecture

```
User Query
    ↓
Natural Language Parser
    ↓
QueryFilters (date, type, price, author, search terms)
    ↓
    ├─→ Semantic Search (ChromaDB + AI embeddings)
    └─→ Text Search (PostgreSQL ILIKE)
    ↓
Combine & Rank Results
    ↓
Apply Post-Filters (price extraction from content)
    ↓
Return Ranked Results
```

### Fallback Behavior

- If ChromaDB is unavailable, search falls back to PostgreSQL text search
- If semantic search fails, text search still works
- Both methods are tried and results are combined for best coverage

## Usage Tips

1. **Be Natural**: Write queries as you would ask a question
   - ✅ "Show me articles about AI from last month"
   - ❌ "type:blog date:last-month tags:ai"

2. **Combine Filters**: Use multiple filters for precise results
   - ✅ "Videos about cooking from last week"
   - ✅ "Books under $20 I saved this month"

3. **Use Specific Terms**: More specific queries get better results
   - ✅ "Articles about neural networks"
   - ❌ "stuff about AI"

4. **Date References**: Use natural date phrases
   - ✅ "last month", "yesterday", "3 days ago"
   - ❌ "2024-10-01" (not yet supported)

## Examples

### Example 1: Finding Articles
**Query:** "Show me articles about AI I saved last month"

**What it does:**
- Extracts: type=blog, search="AI", date=last month
- Searches for blog posts containing "AI" from last month
- Returns ranked results

### Example 2: Finding Quotes
**Query:** "Find that quote about new beginnings from the handwritten note I saved"

**What it does:**
- Extracts: type=text, search="quote about new beginnings"
- Searches text items for content matching the phrase
- Returns matching notes

### Example 3: Finding Specific Content
**Query:** "What did Karpathy say about tokenization in that paper?"

**What it does:**
- Extracts: author="Karpathy", search="tokenization", source="paper"
- Searches for content mentioning Karpathy and tokenization
- Filters for paper-related content

### Example 4: Product Search
**Query:** "Black leather shoes under $300"

**What it does:**
- Extracts: search="black leather shoes", price_max=300
- Searches for items containing "black leather shoes"
- Filters results to only show items with price ≤ $300

### Example 5: Design Inspiration
**Query:** "Inspiration for a travel blog logo"

**What it does:**
- Extracts: search="travel blog logo inspiration"
- Searches all content for matching terms
- Returns relevant images, articles, or design references

### Example 6: Todo Lists
**Query:** "My to-do list for yesterday"

**What it does:**
- Extracts: type=text, search="to-do list", date=yesterday
- Searches text items from yesterday containing "to-do" or "list"
- Returns matching todo items

## API Usage

### Endpoint
```
GET /api/search?q=<query>&limit=<number>
```

### Parameters
- `q` (required): Natural language search query
- `limit` (optional): Maximum number of results (default: 10, max: 50)

### Example Requests

```bash
# Simple search
curl "http://localhost:8080/api/search?q=AI"

# With date filter
curl "http://localhost:8080/api/search?q=articles%20about%20AI%20last%20month"

# With price filter
curl "http://localhost:8080/api/search?q=shoes%20under%20%24300"

# With limit
curl "http://localhost:8080/api/search?q=test&limit=20"
```

### Response Format

```json
[
  {
    "item": {
      "id": "uuid",
      "title": "Item Title",
      "content": "Item content...",
      "summary": "Summary...",
      "source_url": "https://...",
      "type": "blog",
      "tags": ["tag1", "tag2"],
      "image_url": "https://...",
      "embed_html": "...",
      "created_at": "2025-11-08T12:00:00Z"
    },
    "similarity_score": 0.85
  }
]
```

## Future Enhancements

- [ ] Support for exact date ranges ("from January to March")
- [ ] Tag-based filtering ("items tagged with 'important'")
- [ ] Source URL filtering ("from youtube.com")
- [ ] Advanced boolean operators ("AI OR machine learning")
- [ ] Search history and suggestions
- [ ] Saved searches

