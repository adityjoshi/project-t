package services

import (
	"context"
	"fmt"
	"regexp"
	"synapse/internal/db"
	"synapse/internal/models"
	"synapse/internal/repository"

	"github.com/google/uuid"
)

type SearchService struct {
	aiService      *AIService
	itemRepo       *repository.ItemRepository
	collectionName string
}

func NewSearchService(aiService *AIService, itemRepo *repository.ItemRepository) *SearchService {
	return &SearchService{
		aiService:      aiService,
		itemRepo:       itemRepo,
		collectionName: "synapse_items",
	}
}

// Search performs hybrid search: semantic (ChromaDB) + text (PostgreSQL) with natural language parsing
func (s *SearchService) Search(ctx context.Context, query string, limit int) ([]models.SearchResult, error) {
	// Parse natural language query
	filters := ParseNaturalLanguageQuery(query)

	// Try semantic search first (if ChromaDB is available)
	semanticResults, semanticErr := s.semanticSearch(ctx, filters.SearchTerms, limit*2)
	
	// Always do text search as fallback/combination
	textResults, textErr := s.itemRepo.SearchItems(ctx, filters, limit*2)
	
	if semanticErr != nil && textErr != nil {
		// Both failed, return empty
		return []models.SearchResult{}, fmt.Errorf("search failed: semantic=%v, text=%v", semanticErr, textErr)
	}

	// Combine results
	results := s.combineResults(semanticResults, textResults, limit)

	// Apply post-filters (price, etc. that aren't in SQL)
	results = s.applyPostFilters(results, filters)

	return results, nil
}

func (s *SearchService) semanticSearch(ctx context.Context, query string, limit int) ([]models.SearchResult, error) {
	// Generate embedding for query
	queryEmbedding, err := s.aiService.GenerateEmbedding(ctx, query)
	if err != nil {
		return nil, err
	}

	// Query ChromaDB
	ids, distances, err := db.Chroma.Query(s.collectionName, queryEmbedding, limit)
	if err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return []models.SearchResult{}, nil
	}

	// Convert string IDs to UUIDs
	var itemIDs []uuid.UUID
	for _, id := range ids {
		itemID, err := uuid.Parse(id)
		if err != nil {
			continue
		}
		itemIDs = append(itemIDs, itemID)
	}

	// Get items from database
	items, err := s.itemRepo.GetByIDs(ctx, itemIDs)
	if err != nil {
		return nil, err
	}

	// Create map for quick lookup
	itemMap := make(map[uuid.UUID]models.Item)
	for _, item := range items {
		itemMap[item.ID] = item
	}

	// Build results with similarity scores
	var results []models.SearchResult
	for i, id := range ids {
		itemID, err := uuid.Parse(id)
		if err != nil {
			continue
		}

		item, exists := itemMap[itemID]
		if !exists {
			continue
		}

		// Convert distance to similarity score (1 - distance)
		similarity := 1.0 - distances[i]
		if similarity < 0 {
			similarity = 0
		}

		results = append(results, models.SearchResult{
			Item:            item,
			SimilarityScore: similarity,
		})
	}

	return results, nil
}

func (s *SearchService) combineResults(semanticResults []models.SearchResult, textResults []models.Item, limit int) []models.SearchResult {
	// Create a map to deduplicate and combine scores
	resultMap := make(map[uuid.UUID]models.SearchResult)

	// Add semantic results with their scores
	for _, result := range semanticResults {
		resultMap[result.Item.ID] = result
	}

	// Add text results, combining scores if they exist
	for _, item := range textResults {
		if existing, exists := resultMap[item.ID]; exists {
			// Item found in both - boost the score
			existing.SimilarityScore = existing.SimilarityScore*0.7 + 0.3
			resultMap[item.ID] = existing
		} else {
			// New item from text search - give it a base score
			resultMap[item.ID] = models.SearchResult{
				Item:            item,
				SimilarityScore: 0.5, // Base score for text matches
			}
		}
	}

	// Convert map to slice and sort by score
	results := make([]models.SearchResult, 0, len(resultMap))
	for _, result := range resultMap {
		results = append(results, result)
	}

	// Simple sort by similarity score (descending)
	for i := 0; i < len(results)-1; i++ {
		for j := i + 1; j < len(results); j++ {
			if results[i].SimilarityScore < results[j].SimilarityScore {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	// Limit results
	if len(results) > limit {
		results = results[:limit]
	}

	return results
}

func (s *SearchService) applyPostFilters(results []models.SearchResult, filters *models.QueryFilters) []models.SearchResult {
	if filters.PriceMax == nil && filters.PriceMin == nil {
		return results
	}

	filtered := []models.SearchResult{}
	for _, result := range results {
		// Extract price from content (for Amazon products)
		price := extractPriceFromContent(result.Item.Content)
		if price == 0 {
			// No price found, include it anyway
			filtered = append(filtered, result)
			continue
		}

		// Apply price filters
		if filters.PriceMax != nil && price > *filters.PriceMax {
			continue
		}
		if filters.PriceMin != nil && price < *filters.PriceMin {
			continue
		}

		filtered = append(filtered, result)
	}

	return filtered
}

func extractPriceFromContent(content string) float64 {
	// Try to extract price from content (e.g., "Price: $299.99")
	priceRe := regexp.MustCompile(`(?i)price[:\s]+\$?(\d+(?:\.\d+)?)`)
	if match := priceRe.FindStringSubmatch(content); match != nil {
		var price float64
		fmt.Sscanf(match[1], "%f", &price)
		return price
	}
	return 0
}
