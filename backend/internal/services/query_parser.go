package services

import (
	"fmt"
	"regexp"
	"strings"
	"synapse/internal/models"
	"time"
)

func ParseNaturalLanguageQuery(query string) *models.QueryFilters {
	filters := &models.QueryFilters{
		SearchTerms: query,
	}

	lowerQuery := strings.ToLower(query)

	// Extract date filters
	filters.DateFrom, filters.DateTo = extractDateRange(lowerQuery)

	// Extract type filters
	filters.Type = extractType(lowerQuery)

	// Extract price filters
	filters.PriceMin, filters.PriceMax = extractPriceRange(lowerQuery)

	// Extract author/source mentions
	filters.Author = extractAuthor(lowerQuery)
	filters.Source = extractSource(lowerQuery)

	// Extract tags (common patterns)
	filters.Tags = extractTags(lowerQuery)

	// Clean search terms (remove filter phrases)
	filters.SearchTerms = cleanSearchTerms(query, filters)

	return filters
}

func extractDateRange(query string) (*time.Time, *time.Time) {
	now := time.Now()
	var from, to *time.Time

	// "last month"
	if strings.Contains(query, "last month") {
		lastMonth := time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, now.Location())
		from = &lastMonth
		thisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		to = &thisMonth
	}

	// "yesterday"
	if strings.Contains(query, "yesterday") {
		yesterday := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, now.Location())
		from = &yesterday
		today := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
		to = &today
	}

	// "last week"
	if strings.Contains(query, "last week") {
		lastWeek := now.AddDate(0, 0, -7)
		from = &lastWeek
		to = &now
	}

	// "this month"
	if strings.Contains(query, "this month") {
		thisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		from = &thisMonth
		to = &now
	}

	// "last year"
	if strings.Contains(query, "last year") {
		lastYear := time.Date(now.Year()-1, 1, 1, 0, 0, 0, 0, now.Location())
		from = &lastYear
		thisYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		to = &thisYear
	}

	// Relative days: "3 days ago", "5 days ago"
	daysAgoRe := regexp.MustCompile(`(\d+)\s*days?\s*ago`)
	if match := daysAgoRe.FindStringSubmatch(query); match != nil {
		days, _ := time.ParseDuration(match[1] + "h")
		daysAgo := now.Add(-days * 24)
		from = &daysAgo
		to = &now
	}

	return from, to
}

func extractType(query string) string {
	// Only extract type if there are contextual words (like "show me", "my", "I saved")
	// This prevents single-word searches like "video" from being treated as type filters
	hasContext := strings.Contains(query, "show me") ||
		strings.Contains(query, "my ") ||
		strings.Contains(query, "i saved") ||
		strings.Contains(query, "find ") ||
		strings.Contains(query, "get ") ||
		strings.Contains(query, "list of") ||
		strings.Contains(query, "all ")

	// If it's a single word or very short query without context, don't extract type
	words := strings.Fields(query)
	if len(words) <= 2 && !hasContext {
		return ""
	}

	typeMap := map[string]string{
		"article":     "blog",
		"articles":    "blog",
		"blog":        "blog",
		"blog post":   "blog",
		"note":        "text",
		"notes":       "text",
		"handwritten": "text",
		"video":       "video",
		"videos":      "video",
		"youtube":     "video",
		"product":     "amazon",
		"amazon":      "amazon",
		"book":        "book",
		"books":       "book",
		"recipe":      "recipe",
		"recipes":     "recipe",
		"image":       "image",
		"screenshot":  "image",
		"todo":        "text",
		"to-do":       "text",
		"to do":       "text",
		"list":        "text",
	}

	for keyword, itemType := range typeMap {
		if strings.Contains(query, keyword) {
			return itemType
		}
	}

	return ""
}

func extractPriceRange(query string) (*float64, *float64) {
	var min, max *float64

	// "under $300", "below $300", "less than $300"
	underRe := regexp.MustCompile(`(under|below|less than)\s*\$?(\d+(?:\.\d+)?)`)
	if match := underRe.FindStringSubmatch(query); match != nil {
		price := parsePrice(match[2])
		if price > 0 {
			max = &price
		}
	}

	// "over $100", "above $100", "more than $100"
	overRe := regexp.MustCompile(`(over|above|more than)\s*\$?(\d+(?:\.\d+)?)`)
	if match := overRe.FindStringSubmatch(query); match != nil {
		price := parsePrice(match[2])
		if price > 0 {
			min = &price
		}
	}

	// "$100 to $300", "$100-$300"
	rangeRe := regexp.MustCompile(`\$?(\d+(?:\.\d+)?)\s*(?:to|-)\s*\$?(\d+(?:\.\d+)?)`)
	if match := rangeRe.FindStringSubmatch(query); match != nil {
		price1 := parsePrice(match[1])
		price2 := parsePrice(match[2])
		if price1 > 0 && price2 > 0 {
			if price1 < price2 {
				min = &price1
				max = &price2
			} else {
				min = &price2
				max = &price1
			}
		}
	}

	return min, max
}

func parsePrice(s string) float64 {
	var val float64
	_, err := fmt.Sscanf(s, "%f", &val)
	if err == nil {
		return val
	}
	return 0
}

func extractAuthor(query string) string {
	// "from Karpathy", "by Karpathy", "Karpathy said"
	patterns := []string{
		`from\s+([A-Z][a-z]+(?:\s+[A-Z][a-z]+)?)`,
		`by\s+([A-Z][a-z]+(?:\s+[A-Z][a-z]+)?)`,
		`([A-Z][a-z]+(?:\s+[A-Z][a-z]+)?)\s+said`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if match := re.FindStringSubmatch(query); match != nil {
			return match[1]
		}
	}

	return ""
}

func extractSource(query string) string {
	// "from that paper", "in that article", "from the blog"
	if strings.Contains(query, "paper") {
		return "paper"
	}
	if strings.Contains(query, "article") {
		return "article"
	}
	return ""
}

func extractTags(query string) []string {
	// Don't extract tags from simple search queries
	// Tags should only be extracted from explicit tag mentions like #ai
	// This prevents removing important search terms
	var tags []string
	re := regexp.MustCompile(`#(\w+)`) // Finds words prefixed with #
	matches := re.FindAllStringSubmatch(query, -1)
	for _, m := range matches {
		tags = append(tags, m[1])
	}
	return tags
}

func cleanSearchTerms(originalQuery string, filters *models.QueryFilters) string {
	query := originalQuery

	// Remove date phrases
	datePhrases := []string{
		"last month", "yesterday", "last week", "this month", "last year",
		"days ago", "weeks ago", "months ago",
	}
	for _, phrase := range datePhrases {
		query = strings.ReplaceAll(strings.ToLower(query), phrase, "")
	}

	// Only remove type phrases if a type filter was actually set
	// This prevents removing search terms when type wasn't meant to be a filter
	if filters.Type != "" {
		typePhrases := []string{
			"articles", "article", "notes", "note", "videos", "video",
			"products", "product", "books", "book", "recipes", "recipe",
			"images", "image", "screenshots", "screenshot", "todo", "to-do", "list",
		}
		for _, phrase := range typePhrases {
			// Only remove if it matches the detected type
			expectedType := ""
			switch phrase {
			case "articles", "article", "blog", "blog post":
				expectedType = "blog"
			case "notes", "note", "handwritten", "todo", "to-do", "to do", "list":
				expectedType = "text"
			case "videos", "video", "youtube":
				expectedType = "video"
			case "products", "product", "amazon":
				expectedType = "amazon"
			case "books", "book":
				expectedType = "book"
			case "recipes", "recipe":
				expectedType = "recipe"
			case "images", "image", "screenshots", "screenshot":
				expectedType = "image"
			}
			if expectedType == filters.Type {
				query = strings.ReplaceAll(strings.ToLower(query), phrase, "")
			}
		}
	}

	// Remove price phrases
	priceRe := regexp.MustCompile(`(under|below|over|above|less than|more than)\s*\$?\d+`)
	query = priceRe.ReplaceAllString(query, "")

	// Remove author phrases
	authorRe := regexp.MustCompile(`(from|by)\s+[A-Z][a-z]+`)
	query = authorRe.ReplaceAllString(query, "")

	// Clean up extra spaces
	query = regexp.MustCompile(`\s+`).ReplaceAllString(query, " ")
	query = strings.TrimSpace(query)

	// If query is too short after cleaning, use original
	if len(query) < 2 {
		return originalQuery
	}

	return query
}

