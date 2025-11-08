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
			// Generate responsive embed HTML that works with the frontend wrapper
			embedHTML = fmt.Sprintf(`<iframe width="100%%" height="100%%" src="https://www.youtube.com/embed/%s?rel=0" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" allowfullscreen style="position: absolute; top: 0; left: 0; width: 100%%; height: 100%%;"></iframe>`, videoID)
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
	// Extract keywords from recipe title
	keywords := s.extractKeywordsFromTitle(title)
	searchQuery := "recipe"
	if keywords != "" {
		keywordParts := strings.Fields(keywords)
		if len(keywordParts) > 2 {
			keywordParts = keywordParts[:2]
		}
		searchQuery = "recipe," + strings.Join(keywordParts, ",")
	}
	// Use Unsplash Source API with recipe-specific search
	url := fmt.Sprintf("https://source.unsplash.com/400x300/?%s", strings.ReplaceAll(searchQuery, " ", "+"))
	return url, nil
}

// FetchRelevantImage attempts to fetch a relevant image for any content type
func (s *MetadataService) FetchRelevantImage(ctx context.Context, title, content, itemType, category string) (string, error) {
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
		// Try Open Graph image first
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
	case "text", "image", "screenshot":
		// For text notes, images, screenshots - use category-based images
		// Fall through to category-based search
	}
	
	// Category-based image search (works for all types)
	if category != "" {
		return s.getImageByCategory(ctx, title, category)
	}
	
	// Final fallback: generic image
	return s.getImageByCategory(ctx, title, "Other")
}

func (s *MetadataService) getImageByCategory(ctx context.Context, title, category string) (string, error) {
	// Extract keywords from title for more relevant images
	keywords := s.extractKeywordsFromTitle(title)
	
	// Map categories to search terms
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
	
	// Use category as base search term
	searchTerm := categoryMap[category]
	if searchTerm == "" {
		searchTerm = "abstract"
	}
	
	// Combine category with title keywords for better relevance
	// Use Unsplash Source API with search terms (more relevant than Picsum)
	// Format: https://source.unsplash.com/400x300/?keyword1,keyword2
	searchQuery := searchTerm
	if keywords != "" {
		// Limit keywords to avoid too long URLs
		keywordParts := strings.Fields(keywords)
		if len(keywordParts) > 3 {
			keywordParts = keywordParts[:3]
		}
		searchQuery = searchTerm + "," + strings.Join(keywordParts, ",")
	}
	
	// Use Unsplash Source API with search terms
	url := fmt.Sprintf("https://source.unsplash.com/400x300/?%s", strings.ReplaceAll(searchQuery, " ", "+"))
	return url, nil
}

// extractKeywordsFromTitle extracts meaningful keywords from title
func (s *MetadataService) extractKeywordsFromTitle(title string) string {
	// Remove common stop words
	stopWords := map[string]bool{
		"the": true, "a": true, "an": true, "and": true, "or": true, "but": true,
		"in": true, "on": true, "at": true, "to": true, "for": true, "of": true,
		"with": true, "by": true, "from": true, "as": true, "is": true, "are": true,
		"was": true, "were": true, "be": true, "been": true, "have": true, "has": true,
		"had": true, "do": true, "does": true, "did": true, "will": true, "would": true,
		"should": true, "could": true, "may": true, "might": true, "must": true,
		"this": true, "that": true, "these": true, "those": true, "i": true, "you": true,
		"he": true, "she": true, "it": true, "we": true, "they": true,
	}
	
	// Convert to lowercase and split
	words := strings.Fields(strings.ToLower(title))
	
	// Filter out stop words and short words
	keywords := []string{}
	for _, word := range words {
		// Remove punctuation
		word = strings.Trim(word, ".,!?;:()[]{}\"'")
		// Keep words that are meaningful (length > 2 and not a stop word)
		if len(word) > 2 && !stopWords[word] {
			keywords = append(keywords, word)
		}
	}
	
	// Limit to 5 most relevant keywords
	if len(keywords) > 5 {
		keywords = keywords[:5]
	}
	
	return strings.Join(keywords, " ")
}

