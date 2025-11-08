package services

import (
	"context"
	"fmt"
	"synapse/internal/db"
	"synapse/internal/models"
	"synapse/internal/repository"
	"time"

	"github.com/google/uuid"
)

type ItemService struct {
	itemRepo        *repository.ItemRepository
	aiService       *AIService
	metadataService *MetadataService
	collectionName  string
}

func NewItemService(itemRepo *repository.ItemRepository, aiService *AIService) *ItemService {
	return &ItemService{
		itemRepo:        itemRepo,
		aiService:       aiService,
		metadataService: NewMetadataService(),
		collectionName:  "synapse_items",
	}
}

func (s *ItemService) CreateItem(ctx context.Context, req *models.CreateItemRequest) (*models.Item, error) {
	// Generate ID
	itemID := uuid.New()
	embeddingID := itemID.String()

	// Prepare content for processing
	content := req.Content
	if content == "" {
		content = req.Title
	}

	// Generate summary and tags in parallel using goroutines
	type summaryResult struct {
		summary string
		err     error
	}
	type tagsResult struct {
		tags []string
		err  error
	}
	type embeddingResult struct {
		embedding []float32
		err       error
	}

	summaryChan := make(chan summaryResult, 1)
	tagsChan := make(chan tagsResult, 1)
	embeddingChan := make(chan embeddingResult, 1)

	// Generate summary
	go func() {
		summary, err := s.aiService.SummarizeContent(ctx, content)
		summaryChan <- summaryResult{summary: summary, err: err}
	}()

	// Generate tags
	go func() {
		tags, err := s.aiService.GenerateTags(ctx, content)
		tagsChan <- tagsResult{tags: tags, err: err}
	}()

	// Generate embedding
	go func() {
		embedding, err := s.aiService.GenerateEmbedding(ctx, content)
		embeddingChan <- embeddingResult{embedding: embedding, err: err}
	}()

	// Wait for all results
	summaryRes := <-summaryChan
	tagsRes := <-tagsChan
	embeddingRes := <-embeddingChan

	// Handle errors - make AI features optional if API fails
	if summaryRes.err != nil {
		// If summary fails, use a truncated version of content
		if len(content) > 200 {
			summaryRes.summary = content[:200] + "..."
		} else {
			summaryRes.summary = content
		}
	}
	if tagsRes.err != nil {
		// Tags are optional, continue with empty tags
		tagsRes.tags = []string{}
	}
	if embeddingRes.err != nil {
		// If embedding fails, we can't proceed - return error
		return nil, fmt.Errorf("failed to generate embedding (check AI API key): %w", embeddingRes.err)
	}

	// Get metadata (embeds, covers, images) in parallel
	type metadataResult struct {
		embedHTML string
		imageURL  string
		err       error
	}
	metadataChan := make(chan metadataResult, 1)
	
	go func() {
		var embedHTML, imageURL string
		var err error
		
		// Use pre-extracted image URL if provided (from extension)
		if req.ImageURL != "" {
			imageURL = req.ImageURL
		}
		
		// If URL type, get embed and preview
		if req.Type == "url" && req.SourceURL != "" && imageURL == "" {
			embedHTML, imageURL, err = s.metadataService.GetURLMetadata(ctx, req.SourceURL)
		}
		
		// For Amazon products, use metadata image if available
		if req.Type == "amazon" && req.Metadata != nil && req.Metadata["image"] != "" {
			imageURL = req.Metadata["image"]
		}
		
		// For blogs, use metadata image if available
		if req.Type == "blog" && req.Metadata != nil && req.Metadata["image"] != "" {
			imageURL = req.Metadata["image"]
		}
		
		// For videos, use thumbnail if available
		if req.Type == "video" && req.Metadata != nil && req.Metadata["thumbnail"] != "" {
			imageURL = req.Metadata["thumbnail"]
		}
		
		// Detect and get book cover
		if imageURL == "" {
			bookCover, err2 := s.metadataService.DetectBookAndGetCover(ctx, req.Title, content)
			if err2 == nil && bookCover != "" {
				imageURL = bookCover
				if req.Type == "" {
					req.Type = "book"
				}
			}
		}
		
		// Detect and get recipe image
		if imageURL == "" {
			recipeImage, err2 := s.metadataService.DetectRecipeAndGetImage(ctx, req.Title, content)
			if err2 == nil && recipeImage != "" {
				imageURL = recipeImage
				if req.Type == "" {
					req.Type = "recipe"
				}
			}
		}
		
		metadataChan <- metadataResult{embedHTML: embedHTML, imageURL: imageURL, err: err}
	}()
	
	metadataRes := <-metadataChan

	// Store embedding in ChromaDB (optional - if it fails, continue without vector search)
	metadata := map[string]interface{}{
		"title": req.Title,
		"type":  req.Type,
	}
	if err := db.Chroma.AddEmbedding(s.collectionName, embeddingID, embeddingRes.embedding, metadata); err != nil {
		// Log error but continue - item will be saved without embedding
		fmt.Printf("Warning: Failed to store embedding in ChromaDB: %v\n", err)
		fmt.Println("Item will be saved but semantic search may not work until ChromaDB is fixed")
		// Continue without embedding - item can still be saved
	}

	// Create item
	item := &models.Item{
		ID:          itemID,
		Title:       req.Title,
		Content:     content,
		Summary:     summaryRes.summary,
		SourceURL:   req.SourceURL,
		Type:        req.Type,
		Tags:        tagsRes.tags,
		EmbeddingID: embeddingID,
		ImageURL:    metadataRes.imageURL,
		EmbedHTML:   metadataRes.embedHTML,
		CreatedAt:   time.Now(),
	}

	// Save to database
	if err := s.itemRepo.Create(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to save item: %w", err)
	}

	return item, nil
}

func (s *ItemService) GetItem(ctx context.Context, id uuid.UUID) (*models.Item, error) {
	return s.itemRepo.GetByID(ctx, id)
}

func (s *ItemService) GetAllItems(ctx context.Context) ([]models.Item, error) {
	return s.itemRepo.GetAll(ctx)
}

func (s *ItemService) DeleteItem(ctx context.Context, id uuid.UUID) error {
	return s.itemRepo.Delete(ctx, id)
}

