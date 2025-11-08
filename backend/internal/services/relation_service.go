package services

import (
	"context"
	"synapse/internal/db"
	"synapse/internal/models"
	"synapse/internal/repository"

	"github.com/google/uuid"
)

type RelationService struct {
	itemRepo       *repository.ItemRepository
	relationRepo   *repository.RelationRepository
	aiService      *AIService
	collectionName string
}

func NewRelationService(itemRepo *repository.ItemRepository, relationRepo *repository.RelationRepository, aiService *AIService) *RelationService {
	return &RelationService{
		itemRepo:       itemRepo,
		relationRepo:   relationRepo,
		aiService:      aiService,
		collectionName: "synapse_items",
	}
}

func (s *RelationService) FindRelatedItems(ctx context.Context, itemID uuid.UUID, limit int) ([]models.RelatedItem, error) {
	// Check cache first
	cached, err := s.relationRepo.GetRelated(ctx, itemID, limit)
	if err == nil && len(cached) > 0 {
		return cached, nil
	}

	// Get the item
	item, err := s.itemRepo.GetByID(ctx, itemID)
	if err != nil {
		return nil, err
	}

	// Generate embedding for the item's content to use for similarity search
	// (We could store this, but for MVP we'll regenerate)
	searchText := item.Title + " " + item.Content
	if len(searchText) > 1000 {
		searchText = searchText[:1000]
	}
	
	embedding, err := s.aiService.GenerateEmbedding(ctx, searchText)
	if err != nil {
		return nil, err
	}

	// Query for similar items (limit+1 to potentially exclude the item itself)
	ids, distances, err := db.Chroma.Query(s.collectionName, embedding, limit+10)
	if err != nil {
		return nil, err
	}

	// Filter out the current item and convert to UUIDs
	var relatedIDs []uuid.UUID
	var relatedDistances []float64
	for i, id := range ids {
		relatedID, err := uuid.Parse(id)
		if err != nil || relatedID == itemID {
			continue
		}
		relatedIDs = append(relatedIDs, relatedID)
		relatedDistances = append(relatedDistances, distances[i])
	}

	if len(relatedIDs) == 0 {
		return []models.RelatedItem{}, nil
	}

	// Limit to requested number
	if len(relatedIDs) > limit {
		relatedIDs = relatedIDs[:limit]
		relatedDistances = relatedDistances[:limit]
	}

	// Get items
	items, err := s.itemRepo.GetByIDs(ctx, relatedIDs)
	if err != nil {
		return nil, err
	}

	// Create map for quick lookup
	itemMap := make(map[uuid.UUID]models.Item)
	for _, it := range items {
		itemMap[it.ID] = it
	}

	// Build results and cache them
	var results []models.RelatedItem
	for i, id := range relatedIDs {
		it, exists := itemMap[id]
		if !exists {
			continue
		}

		similarity := 1.0 - relatedDistances[i]
		if similarity < 0 {
			similarity = 0
		}

		// Cache the relation
		s.relationRepo.Create(ctx, itemID, id, similarity)

		results = append(results, models.RelatedItem{
			Item:            it,
			SimilarityScore: similarity,
		})
	}

	return results, nil
}

