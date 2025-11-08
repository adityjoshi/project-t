package models

import (
	"time"

	"github.com/google/uuid"
)

type QueryFilters struct {
	SearchTerms   string
	Type          string
	DateFrom      *time.Time
	DateTo        *time.Time
	Tags          []string
	PriceMax      *float64
	PriceMin      *float64
	Author        string
	Source        string
}

type Item struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Summary     string    `json:"summary"`
	SourceURL   string    `json:"source_url"`
	Type        string    `json:"type"` // "text", "url", "image", "book", "recipe"
	Category    string    `json:"category"` // AI-categorized section: "Technology", "Food & Recipes", "Books", "Videos", "Shopping", "Articles", "Notes", etc.
	Tags        []string  `json:"tags"`
	EmbeddingID string    `json:"embedding_id"`
	ImageURL    string    `json:"image_url"`    // For book covers, recipe images, or page previews
	EmbedHTML   string    `json:"embed_html"`   // For URL embeds/previews
	CreatedAt   time.Time `json:"created_at"`
}

type CreateItemRequest struct {
	Title     string            `json:"title"`
	Content   string            `json:"content"`
	SourceURL string            `json:"source_url"`
	Type      string            `json:"type"` // "text", "url", "image", "amazon", "blog", "video"
	ImageURL  string            `json:"image_url"` // For pre-extracted images
	Metadata  map[string]string `json:"metadata"` // Additional metadata (price, rating, etc.)
}

type RelatedItem struct {
	Item           Item    `json:"item"`
	SimilarityScore float64 `json:"similarity_score"`
}

type SearchResult struct {
	Item           Item    `json:"item"`
	SimilarityScore float64 `json:"similarity_score"`
}

