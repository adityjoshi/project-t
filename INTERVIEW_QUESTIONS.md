# Project Synapse - Interview Questions & Answers

## Project Description (30-60 Second Pitch)

**Project Synapse** is an AI-powered knowledge management system I built that acts as a "second brain" for capturing, organizing, and retrieving information. 

The system allows users to save any type of content - web pages, YouTube videos, books, recipes, products, notes, and images - through a browser extension or web interface. What makes it intelligent is the AI layer that automatically categorizes content into predefined sections, generates relevant tags, creates semantic summaries, and enables natural language search that understands meaning, not just keywords.

The architecture uses Go for the backend with PostgreSQL for metadata and ChromaDB for vector embeddings, React for the frontend, and Google Gemini for AI features. Key technical highlights include asynchronous AI processing using goroutines, hybrid search combining semantic and text search, specialized views for different content types, and a browser extension that intelligently extracts metadata from various websites.

The project demonstrates full-stack development, AI integration, database design, and creating a seamless user experience for knowledge management.

---

## General Questions

### 1. What is Project Synapse?
**Answer**: Project Synapse is an AI-powered knowledge management system that acts as a "second brain" for capturing, organizing, and retrieving information from various sources like web pages, videos, books, recipes, and notes.

### 2. What problem does it solve?
**Answer**: It solves the problem of information overload by providing a centralized place to save, categorize, and search through all your saved content using AI-powered features like automatic categorization, tagging, and semantic search.

### 3. What is the tech stack?
**Answer**: Backend uses Go with Gin framework, PostgreSQL for metadata, ChromaDB for vector storage, and Google Gemini for AI features. Frontend uses React with Vite and Tailwind CSS. Browser extension uses Chrome Extension Manifest V3.

---

## Architecture Questions

### 4. How is the application structured?
**Answer**: It follows a layered architecture with handlers (HTTP), services (business logic), repositories (data access), and models (data structures). The backend is RESTful, frontend is component-based React, and extension uses content scripts and background workers.

### 5. Why did you choose Go for the backend?
**Answer**: Go provides excellent concurrency with goroutines for async AI operations, fast performance, strong typing, and good standard library support for HTTP and database operations.

### 6. How do you handle asynchronous operations?
**Answer**: Using Go goroutines to run AI operations (summarization, OCR, tagging) in the background without blocking the main request. Items are saved immediately, then updated asynchronously when AI processing completes.

### 7. What databases do you use and why?
**Answer**: PostgreSQL for structured metadata (items, tags, relationships) because it's reliable and supports full-text search. ChromaDB for vector embeddings to enable semantic similarity search.

---

## AI & Machine Learning Questions

### 8. What AI features are implemented?
**Answer**: Automatic categorization into predefined sections, automatic tag generation, semantic summarization, embedding generation for similarity search, and OCR for text extraction from images.

### 9. Why did you choose Google Gemini over OpenAI?
**Answer**: Gemini offers a free tier with good API limits, supports embeddings, vision (OCR), and text generation, making it cost-effective for a personal knowledge management system.

### 10. How does semantic search work?
**Answer**: Content is converted to vector embeddings using AI, stored in ChromaDB, and similarity search finds items with similar meaning rather than just keyword matches. It's combined with text search for hybrid results.

### 11. How do you handle AI API failures?
**Answer**: The system gracefully degrades - if AI categorization fails, it uses type-based defaults. If summarization fails, it keeps the initial truncated summary. The app remains functional even if AI services are unavailable.

### 12. What is the difference between semantic search and text search?
**Answer**: Semantic search understands meaning and context (e.g., "AI" matches "artificial intelligence" and "machine learning"), while text search only matches exact keywords. We combine both for better results.

---

## Content Management Questions

### 13. What content types does the system support?
**Answer**: Text notes, URLs/articles, YouTube videos, images/screenshots, books, recipes, Amazon products, and to-do lists. Each type has specialized display components and extraction logic.

### 14. How does automatic categorization work?
**Answer**: The AI analyzes the title, content, and type to determine the most appropriate category from predefined options like Technology, Food & Recipes, Books, Videos, etc. It uses semantic understanding rather than simple keyword matching.

### 15. How are images automatically fetched?
**Answer**: The system tries multiple sources in priority: pre-extracted images from extension, YouTube thumbnails, book covers from Open Library, recipe images from Unsplash, Open Graph images, and finally category-based images using keywords extracted from titles.

### 16. How do you extract YouTube video descriptions?
**Answer**: Multiple strategies are used: YouTube's internal data structures (ytInitialPlayerResponse, ytInitialData), DOM extraction with "Show more" button expansion, meta tags, and script tag parsing. Video ID verification ensures we get the current video's description, not cached data.

### 17. What is OCR and how is it used?
**Answer**: OCR (Optical Character Recognition) extracts text from images and screenshots using Gemini Vision API. The extracted text is stored separately and included in search queries, making image content searchable.

---

## Search & Discovery Questions

### 18. How does natural language search work?
**Answer**: The query parser extracts filters (dates, types, prices, authors) from natural language, then performs hybrid search combining semantic search (meaning-based) and text search (keyword-based), applies filters, and ranks results by relevance.

### 19. What types of queries are supported?
**Answer**: Date-based ("last month", "yesterday"), type-based ("videos", "articles"), price-based ("under $300"), author-based ("from Karpathy"), quote searches ("that quote about..."), and combined queries mixing multiple filters.

### 20. How do you find related items?
**Answer**: Using vector embeddings stored in ChromaDB, we calculate similarity scores between items. When viewing an item, we find other items with similar embeddings and display the top 5 most related items.

### 21. What is hybrid search?
**Answer**: Combining semantic search (AI embeddings for meaning) with text search (PostgreSQL for keywords). Results from both are merged, with items found in both methods getting boosted relevance scores.

---

## Browser Extension Questions

### 22. How does the browser extension work?
**Answer**: Content scripts extract page data, background service worker handles context menus and quick saves, and popup provides the UI. It detects content types (Amazon, YouTube, blogs) and extracts relevant metadata automatically.

### 23. What content can the extension capture?
**Answer**: Amazon products (price, rating, images), YouTube videos (description, channel, thumbnail), blog posts (content, author, date), generic web pages (main content, images), selected text, and screenshots.

### 24. How do you prevent stale data when extracting YouTube descriptions?
**Answer**: We verify the video ID in YouTube's internal data matches the current URL's video ID before extracting. We also wait for the page to load and poll until the data matches, preventing extraction of old video data.

---

## Performance & Scalability Questions

### 25. How do you optimize performance?
**Answer**: Asynchronous AI operations using goroutines, parallel generation of category/tags/embeddings, database connection pooling, indexed database queries, result limiting, and efficient vector search in ChromaDB.

### 26. How would you scale this for multiple users?
**Answer**: Add user authentication, partition data by user ID, implement connection pooling, add caching layers (Redis), use load balancers, scale ChromaDB horizontally, and implement rate limiting for AI API calls.

### 27. What happens if ChromaDB is unavailable?
**Answer**: The system falls back to PostgreSQL text search only. Semantic search is disabled, but text search continues to work. The application remains functional, just without semantic similarity features.

---

## Data Flow Questions

### 28. Walk me through what happens when a user saves a YouTube video.
**Answer**: Extension extracts video metadata and description, sends to backend. Backend saves item immediately with initial summary, then asynchronously generates AI summary, creates embedding, categorizes content, generates tags, and fetches thumbnail. All AI operations happen in background without blocking.

### 29. How does the swipeable card view work?
**Answer**: Uses touch and mouse drag events to detect swipe gestures, tracks offset for visual feedback, differentiates between clicks and swipes using timing and distance thresholds, and allows keyboard navigation with arrow keys. Navigation buttons are positioned on sides of content.

### 30. How do you ensure data consistency?
**Answer**: Database transactions for critical operations, UUIDs for unique identification, foreign key constraints for relationships, validation at handler level, and error handling with rollback capabilities. Async operations update existing records rather than creating new ones.

---

## Key Takeaways

- **Architecture**: Clean layered architecture with separation of concerns
- **AI Integration**: Multiple AI services working together (categorization, tagging, summarization, embeddings, OCR)
- **User Experience**: Specialized views for different content types, natural language search, related items discovery
- **Performance**: Asynchronous processing, parallel operations, hybrid search
- **Reliability**: Graceful degradation, fallback mechanisms, error handling
- **Extensibility**: Modular services, easy to add new content types or AI features

