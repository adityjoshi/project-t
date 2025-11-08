package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type MetadataService struct {
	client *http.Client
}

func NewMetadataService() *MetadataService {
	return &MetadataService{
		client: &http.Client{},
	}
}

// GetURLMetadata extracts metadata from a URL including embed HTML and images
func (s *MetadataService) GetURLMetadata(ctx context.Context, url string) (embedHTML string, imageURL string, err error) {
	// For YouTube URLs, generate embed
	if strings.Contains(url, "youtube.com") || strings.Contains(url, "youtu.be") {
		videoID := s.extractYouTubeID(url)
		if videoID != "" {
			embedHTML = fmt.Sprintf(`<iframe width="560" height="315" src="https://www.youtube.com/embed/%s" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>`, videoID)
			imageURL = fmt.Sprintf("https://img.youtube.com/vi/%s/maxresdefault.jpg", videoID)
			return embedHTML, imageURL, nil
		}
	}

	// For other URLs, try to get Open Graph image
	imageURL, _ = s.getOpenGraphImage(ctx, url)
	
	// Generate simple embed for other URLs
	if imageURL != "" {
		embedHTML = fmt.Sprintf(`<div class="url-preview"><img src="%s" alt="Preview" style="max-width: 100%%; border-radius: 8px;" /></div>`, imageURL)
	}

	return embedHTML, imageURL, nil
}

// DetectBookAndGetCover detects if content is about a book and fetches cover
func (s *MetadataService) DetectBookAndGetCover(ctx context.Context, title, content string) (string, error) {
	// Simple detection: check if title/content mentions "book" or common book patterns
	bookKeywords := []string{"book", "author", "published", "isbn", "chapter", "novel", "read"}
	lowerTitle := strings.ToLower(title)
	lowerContent := strings.ToLower(content)
	
	isBook := false
	for _, keyword := range bookKeywords {
		if strings.Contains(lowerTitle, keyword) || strings.Contains(lowerContent, keyword) {
			isBook = true
			break
		}
	}
	
	if !isBook {
		return "", nil
	}

	// Try to extract ISBN
	isbn := s.extractISBN(content)
	if isbn != "" {
		return s.getBookCoverByISBN(ctx, isbn)
	}

	// Try Open Library API with title
	return s.getBookCoverByTitle(ctx, title)
}

// DetectRecipeAndGetImage detects if content is a recipe and fetches image
func (s *MetadataService) DetectRecipeAndGetImage(ctx context.Context, title, content string) (string, error) {
	// Simple detection: check for recipe keywords
	recipeKeywords := []string{"recipe", "ingredients", "cook", "bake", "prep time", "servings", "cups", "tablespoons", "tsp", "tbsp"}
	lowerTitle := strings.ToLower(title)
	lowerContent := strings.ToLower(content)
	
	isRecipe := false
	for _, keyword := range recipeKeywords {
		if strings.Contains(lowerTitle, keyword) || strings.Contains(lowerContent, keyword) {
			isRecipe = true
			break
		}
	}
	
	if !isRecipe {
		return "", nil
	}

	// Try to get recipe image from content or use a placeholder service
	// For now, we'll use a recipe image API or extract from content
	return s.getRecipeImage(ctx, title)
}

func (s *MetadataService) extractYouTubeID(url string) string {
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

func (s *MetadataService) getOpenGraphImage(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; SynapseBot/1.0)")
	
	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	
	// Extract og:image
	re := regexp.MustCompile(`<meta\s+property=["']og:image["']\s+content=["']([^"']+)["']`)
	matches := re.FindStringSubmatch(string(body))
	if len(matches) > 1 {
		return matches[1], nil
	}
	
	// Try twitter:image
	re = regexp.MustCompile(`<meta\s+name=["']twitter:image["']\s+content=["']([^"']+)["']`)
	matches = re.FindStringSubmatch(string(body))
	if len(matches) > 1 {
		return matches[1], nil
	}
	
	return "", nil
}

func (s *MetadataService) extractISBN(content string) string {
	// Extract ISBN-13 or ISBN-10
	patterns := []string{
		`ISBN[-\s]*(?:13)?[:\s]*([0-9]{13})`,
		`ISBN[-\s]*(?:10)?[:\s]*([0-9X]{10})`,
		`([0-9]{3}[- ]?[0-9]{10})`,
	}
	
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(content)
		if len(matches) > 1 {
			return strings.ReplaceAll(strings.ReplaceAll(matches[1], "-", ""), " ", "")
		}
	}
	return ""
}

func (s *MetadataService) getBookCoverByISBN(ctx context.Context, isbn string) (string, error) {
	// Use Open Library Covers API
	url := fmt.Sprintf("https://covers.openlibrary.org/b/isbn/%s-L.jpg", isbn)
	
	req, _ := http.NewRequestWithContext(ctx, "HEAD", url, nil)
	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == 200 {
		return url, nil
	}
	return "", nil
}

func (s *MetadataService) getBookCoverByTitle(ctx context.Context, title string) (string, error) {
	// Use Open Library Search API
	searchURL := fmt.Sprintf("https://openlibrary.org/search.json?title=%s&limit=1", strings.ReplaceAll(title, " ", "+"))
	
	req, _ := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	var result struct {
		Docs []struct {
			CoverI int `json:"cover_i"`
		} `json:"docs"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	
	if len(result.Docs) > 0 && result.Docs[0].CoverI > 0 {
		return fmt.Sprintf("https://covers.openlibrary.org/b/id/%d-L.jpg", result.Docs[0].CoverI), nil
	}
	
	return "", nil
}

func (s *MetadataService) getRecipeImage(ctx context.Context, title string) (string, error) {
	// Use Unsplash API for recipe images (free, no key needed for basic usage)
	// Or use a recipe API
	searchQuery := strings.ReplaceAll(title, " ", "+")
	url := fmt.Sprintf("https://source.unsplash.com/400x300/?recipe,%s", searchQuery)
	return url, nil
}

// FetchRelevantImage attempts to fetch a relevant image for any content type
func (s *MetadataService) FetchRelevantImage(ctx context.Context, title, content, itemType, category string) (string, error) {
	// If we already have an image, return it
	// This should be checked before calling this function
	
	// Try different strategies based on type and category
	switch itemType {
	case "video":
		// Already handled in GetURLMetadata
		return "", nil
	case "book":
		return s.DetectBookAndGetCover(ctx, title, content)
	case "recipe":
		return s.DetectRecipeAndGetImage(ctx, title, content)
	case "amazon":
		// Amazon products should have images from metadata
		return "", nil
	case "blog", "url":
		// Try Open Graph image
		if content != "" {
			// Extract URL from content if possible
			urlRe := regexp.MustCompile(`https?://[^\s]+`)
			matches := urlRe.FindStringSubmatch(content)
			if len(matches) > 0 {
				imageURL, _ := s.getOpenGraphImage(ctx, matches[0])
				if imageURL != "" {
					return imageURL, nil
				}
			}
		}
		// Fall through to category-based search
	}
	
	// Category-based image search using Unsplash
	return s.getImageByCategory(ctx, title, category)
}

func (s *MetadataService) getImageByCategory(ctx context.Context, title, category string) (string, error) {
	// Map categories to Unsplash search terms
	categoryMap := map[string]string{
		"Technology":        "technology",
		"Food & Recipes":    "food",
		"Books & Reading":   "books",
		"Videos & Entertainment": "entertainment",
		"Shopping & Products": "product",
		"Articles & News":   "news",
		"Notes & Ideas":     "notebook",
		"Design & Inspiration": "design",
		"Travel":            "travel",
		"Health & Fitness":  "fitness",
		"Education & Learning": "education",
	}
	
	searchTerm := categoryMap[category]
	if searchTerm == "" {
		searchTerm = "abstract"
	}
	
	// Use Unsplash Source API (free, no key needed)
	searchQuery := strings.ReplaceAll(title, " ", "+")
	url := fmt.Sprintf("https://source.unsplash.com/400x300/?%s,%s", searchTerm, searchQuery)
	return url, nil
}

