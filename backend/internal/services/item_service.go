package services

import (
	"context"
	"fmt"
	"regexp"
	"strings"
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
	ocrService      *OCRService
	collectionName  string
}

func NewItemService(itemRepo *repository.ItemRepository, aiService *AIService) *ItemService {
	return &ItemService{
		itemRepo:        itemRepo,
		aiService:       aiService,
		metadataService: NewMetadataService(),
		ocrService:      NewOCRService(),
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

	// Generate category, tags, and embedding in parallel (synchronous for initial save)
	type categoryResult struct {
		category string
		err      error
	}
	type tagsResult struct {
		tags []string
		err  error
	}
	type embeddingResult struct {
		embedding []float32
		err       error
	}

	categoryChan := make(chan categoryResult, 1)
	tagsChan := make(chan tagsResult, 1)
	embeddingChan := make(chan embeddingResult, 1)

	// Generate category (AI-powered categorization)
	go func() {
		category, err := s.aiService.CategorizeContent(ctx, req.Title, content, req.Type)
		categoryChan <- categoryResult{category: category, err: err}
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
	categoryRes := <-categoryChan
	tagsRes := <-tagsChan
	embeddingRes := <-embeddingChan

	// Handle errors - make AI features optional if API fails
	if categoryRes.err != nil {
		// If categorization fails, use a default category based on type
		categoryRes.category = s.getDefaultCategory(req.Type, req.SourceURL)
	}
	
	// Override category for YouTube videos - always "Videos & Entertainment"
	if req.SourceURL != "" && (strings.Contains(req.SourceURL, "youtube.com") || strings.Contains(req.SourceURL, "youtu.be")) {
		categoryRes.category = "Videos & Entertainment"
	}
	
	// Override category for video type - always "Videos & Entertainment"
	if req.Type == "video" {
		categoryRes.category = "Videos & Entertainment"
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
		
		// For videos, ALWAYS get embed HTML (required for embedded playback)
		if req.Type == "video" && req.SourceURL != "" {
			embedHTML, imageURL, err = s.metadataService.GetURLMetadata(ctx, req.SourceURL)
			// If GetURLMetadata didn't return embed, try to generate it from URL
			if embedHTML == "" && (strings.Contains(req.SourceURL, "youtube.com") || strings.Contains(req.SourceURL, "youtu.be")) {
				videoID := s.extractYouTubeIDFromURL(req.SourceURL)
				if videoID != "" {
					embedHTML = fmt.Sprintf(`<iframe width="100%%" height="100%%" src="https://www.youtube.com/embed/%s?rel=0" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" allowfullscreen style="position: absolute; top: 0; left: 0; width: 100%%; height: 100%%;"></iframe>`, videoID)
					if imageURL == "" {
						imageURL = fmt.Sprintf("https://img.youtube.com/vi/%s/maxresdefault.jpg", videoID)
					}
				}
			}
		}
		
		// Use pre-extracted image URL if provided (from extension) - but don't override embed
		if req.ImageURL != "" && imageURL == "" {
			imageURL = req.ImageURL
		}
		
		// If URL type (not video), get embed and preview
		if req.Type == "url" && req.SourceURL != "" && embedHTML == "" {
			embedHTML, imageURL, err = s.metadataService.GetURLMetadata(ctx, req.SourceURL)
		}
		
		// Check if URL is a YouTube video even if type is not "video"
		if embedHTML == "" && req.SourceURL != "" && (strings.Contains(req.SourceURL, "youtube.com") || strings.Contains(req.SourceURL, "youtu.be")) {
			videoID := s.extractYouTubeIDFromURL(req.SourceURL)
			if videoID != "" {
				embedHTML = fmt.Sprintf(`<iframe width="100%%" height="100%%" src="https://www.youtube.com/embed/%s?rel=0" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" allowfullscreen style="position: absolute; top: 0; left: 0; width: 100%%; height: 100%%;"></iframe>`, videoID)
				if imageURL == "" {
					imageURL = fmt.Sprintf("https://img.youtube.com/vi/%s/maxresdefault.jpg", videoID)
				}
			}
		}
		
		// For Amazon products, use metadata image if available
		if req.Type == "amazon" && req.Metadata != nil && req.Metadata["image"] != "" {
			imageURL = req.Metadata["image"]
		}
		
		// For blogs, use metadata image if available
		if req.Type == "blog" && req.Metadata != nil && req.Metadata["image"] != "" {
			imageURL = req.Metadata["image"]
		}
		
		// For videos, use thumbnail if available (fallback if GetURLMetadata didn't work)
		if req.Type == "video" && imageURL == "" && req.Metadata != nil && req.Metadata["thumbnail"] != "" {
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
		
		// If still no image, try to fetch a relevant image based on category
		// This should work for all content types (text, blog, etc.)
		if imageURL == "" {
			if categoryRes.category != "" {
				// Use category-based image fetching
				relevantImage, err2 := s.metadataService.FetchRelevantImage(ctx, req.Title, content, req.Type, categoryRes.category)
				if err2 == nil && relevantImage != "" {
					imageURL = relevantImage
				}
			} else if req.Type != "" {
				// Fallback: use type-based default category
				defaultCategory := s.getDefaultCategory(req.Type, req.SourceURL)
				if defaultCategory != "" {
					relevantImage, err2 := s.metadataService.FetchRelevantImage(ctx, req.Title, content, req.Type, defaultCategory)
					if err2 == nil && relevantImage != "" {
						imageURL = relevantImage
					}
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

		// Extract OCR text from images/screenshots asynchronously
		var ocrText string
		if (req.Type == "image" || req.Type == "screenshot") && metadataRes.imageURL != "" {
			// Extract OCR text in background
			go func() {
				extractedText, err := s.ocrService.ExtractTextFromImage(context.Background(), metadataRes.imageURL)
				if err == nil && extractedText != "" {
					// Update item with OCR text
					s.updateOCRText(context.Background(), itemID, extractedText)
				}
			}()
		}

		// Set initial summary (will be replaced by async AI summary)
		initialSummary := ""
		if req.Type == "video" && req.Metadata != nil && req.Metadata["description"] != "" {
			// For videos, use a truncated version of description as initial summary
			// The full description stays in content, summary will be replaced by AI
			desc := req.Metadata["description"]
			if len(desc) > 200 {
				initialSummary = desc[:200] + "..."
			} else {
				initialSummary = desc
			}
		} else if req.Type == "video" && content != "" {
			// Extract description from content if it contains "Description:" marker
			var desc string
			if descIdx := strings.Index(content, "Description:"); descIdx != -1 {
				desc = strings.TrimSpace(content[descIdx+len("Description:"):])
			} else {
				desc = content
			}
			if len(desc) > 200 {
				initialSummary = desc[:200] + "..."
			} else {
				initialSummary = desc
			}
		} else {
			// For non-videos, use truncated content
			if len(content) > 200 {
				initialSummary = content[:200] + "..."
			} else {
				initialSummary = content
			}
		}

		item := &models.Item{
			ID:          itemID,
			Title:       req.Title,
			Content:     content,
			Summary:     initialSummary, // Temporary summary, will be replaced asynchronously
			SourceURL:   req.SourceURL,
			Type:        req.Type,
			Category:    categoryRes.category,
			Tags:        tagsRes.tags,
			EmbeddingID: embeddingID,
			ImageURL:    metadataRes.imageURL,
			EmbedHTML:   metadataRes.embedHTML,
			OcrText:     ocrText, // Will be updated asynchronously for images
			CreatedAt:   time.Now(),
		}

		// Save to database
		if err := s.itemRepo.Create(ctx, item); err != nil {
			return nil, fmt.Errorf("failed to save item: %w", err)
		}

		// Asynchronously generate AI summary (doesn't affect description/content)
		// For videos, extract description and generate a short summary
		if req.Type == "video" && req.SourceURL != "" {
			// Extract description from metadata if available, otherwise use content
			description := ""
			if req.Metadata != nil && req.Metadata["description"] != "" {
				description = req.Metadata["description"]
			} else if content != "" {
				// Try to extract description from content if it contains "Description:" marker
				if descIdx := strings.Index(content, "Description:"); descIdx != -1 {
					description = strings.TrimSpace(content[descIdx+len("Description:"):])
				} else {
					description = content
				}
			}
			
			if description != "" {
				// Generate short AI summary asynchronously (description stays unchanged)
				go s.generateAndUpdateVideoSummaryAsync(context.Background(), itemID, req.SourceURL, req.Title, description)
			}
		} else {
			// For non-videos, generate regular summary
			go s.generateAndUpdateSummaryAsync(context.Background(), itemID, req.Title, content)
		}

	return item, nil
}

// extractYouTubeIDFromURL extracts YouTube video ID from URL
func (s *ItemService) extractYouTubeIDFromURL(url string) string {
	patterns := []string{
		`youtube\.com/watch\?v=([a-zA-Z0-9_-]+)`,
		`youtu\.be/([a-zA-Z0-9_-]+)`,
		`youtube\.com/embed/([a-zA-Z0-9_-]+)`,
	}
	
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(url)
		if len(matches) > 1 {
			return matches[1]
		}
	}
	return ""
}

// generateAndUpdateSummaryAsync generates a semantic summary asynchronously and updates the item
func (s *ItemService) generateAndUpdateSummaryAsync(ctx context.Context, itemID uuid.UUID, title, content string) {
	// Generate semantic summary using Gemini
	summary, err := s.aiService.GenerateSemanticSummary(ctx, title, content)
	if err != nil {
		fmt.Printf("Warning: Failed to generate semantic summary for item %s: %v\n", itemID, err)
		return
	}

	// Update the item's summary in the database
	if err := s.itemRepo.UpdateSummary(ctx, itemID, summary); err != nil {
		fmt.Printf("Warning: Failed to update summary for item %s: %v\n", itemID, err)
		return
	}

	fmt.Printf("Successfully generated and updated semantic summary for item %s\n", itemID)
}

// updateOCRText updates the OCR text for an item
func (s *ItemService) updateOCRText(ctx context.Context, itemID uuid.UUID, ocrText string) {
	if err := s.itemRepo.UpdateOCRText(ctx, itemID, ocrText); err != nil {
		fmt.Printf("Warning: Failed to update OCR text for item %s: %v\n", itemID, err)
		return
	}
	fmt.Printf("Successfully updated OCR text for item %s\n", itemID)
}

// generateAndUpdateVideoSummaryAsync generates a video-specific summary asynchronously
func (s *ItemService) generateAndUpdateVideoSummaryAsync(ctx context.Context, itemID uuid.UUID, videoURL, title, description string) {
	// Log what we're working with
	fmt.Printf("Generating video summary for item %s - Title: %s, Description length: %d\n", itemID, title, len(description))
	
	// Ensure we have a description to work with
	if description == "" {
		fmt.Printf("Warning: No description provided for video summary, item %s\n", itemID)
		// Fallback to regular summary with title
		s.generateAndUpdateSummaryAsync(ctx, itemID, title, title)
		return
	}
	
	// Generate video summary using Gemini
	summary, err := s.aiService.SummarizeYouTubeVideo(ctx, videoURL, title, description)
	if err != nil {
		fmt.Printf("Warning: Failed to generate video summary for item %s: %v\n", itemID, err)
		// Fallback to regular summary
		s.generateAndUpdateSummaryAsync(ctx, itemID, title, description)
		return
	}

	// Ensure we got a valid summary
	if summary == "" {
		fmt.Printf("Warning: Empty summary generated for item %s, using fallback\n", itemID)
		s.generateAndUpdateSummaryAsync(ctx, itemID, title, description)
		return
	}

	// Update the item's summary in the database
	if err := s.itemRepo.UpdateSummary(ctx, itemID, summary); err != nil {
		fmt.Printf("Warning: Failed to update video summary for item %s: %v\n", itemID, err)
		return
	}

	summaryPreview := summary
	if len(summary) > 100 {
		summaryPreview = summary[:100] + "..."
	}
	fmt.Printf("Successfully generated and updated video summary for item %s: %s\n", itemID, summaryPreview)
}

// getDefaultCategory returns a default category based on item type and URL
func (s *ItemService) getDefaultCategory(itemType, sourceURL string) string {
	// Check if it's a YouTube video first
	if sourceURL != "" && (strings.Contains(sourceURL, "youtube.com") || strings.Contains(sourceURL, "youtu.be")) {
		return "Videos & Entertainment"
	}
	
	typeMap := map[string]string{
		"video":   "Videos & Entertainment",
		"book":    "Books & Reading",
		"recipe":  "Food & Recipes",
		"amazon":  "Shopping & Products",
		"blog":    "Articles & News",
		"url":     "Articles & News",
		"text":    "Notes & Ideas",
		"image":   "Design & Inspiration",
		"screenshot": "Notes & Ideas",
	}

	if category, ok := typeMap[itemType]; ok {
		return category
	}
	return "Other"
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

// RefreshImageForItem refreshes the image URL for an existing item
func (s *ItemService) RefreshImageForItem(ctx context.Context, id uuid.UUID) error {
	item, err := s.itemRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Check if image URL is from deprecated Unsplash Source API
	if item.ImageURL != "" && strings.Contains(item.ImageURL, "source.unsplash.com") {
		// Generate new image URL using Picsum Photos
		newImageURL, err := s.metadataService.FetchRelevantImage(ctx, item.Title, item.Content, item.Type, item.Category)
		if err == nil && newImageURL != "" {
			return s.itemRepo.UpdateImageURL(ctx, id, newImageURL)
		}
	}

	// If no image exists, try to fetch one
	if item.ImageURL == "" {
		newImageURL, err := s.metadataService.FetchRelevantImage(ctx, item.Title, item.Content, item.Type, item.Category)
		if err == nil && newImageURL != "" {
			return s.itemRepo.UpdateImageURL(ctx, id, newImageURL)
		}
	}

	return nil
}

// RefreshSummaryForItem regenerates the summary for an existing item
func (s *ItemService) RefreshSummaryForItem(ctx context.Context, id uuid.UUID) error {
	item, err := s.itemRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// For videos, use video-specific summarization
	if item.Type == "video" && item.SourceURL != "" {
		// Extract description from content
		description := ""
		if item.Content != "" {
			// Try to extract description from content if it contains "Description:" marker
			if descIdx := strings.Index(item.Content, "Description:"); descIdx != -1 {
				description = strings.TrimSpace(item.Content[descIdx+len("Description:"):])
			} else {
				// If no "Description:" marker, use the full content
				description = item.Content
			}
		}
		
		if description != "" {
			// Regenerate video summary asynchronously
			go s.generateAndUpdateVideoSummaryAsync(context.Background(), id, item.SourceURL, item.Title, description)
		} else {
			// Fallback to regular summary
			go s.generateAndUpdateSummaryAsync(context.Background(), id, item.Title, item.Content)
		}
	} else {
		// For non-videos, use regular summarization
		go s.generateAndUpdateSummaryAsync(context.Background(), id, item.Title, item.Content)
	}

	return nil
}

