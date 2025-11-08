package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"synapse/internal/models"
)

type AIService struct {
	provider      string
	geminiKey     string
	openaiKey     string
	claudeKey     string
	claudeBaseURL string
	client        *http.Client
}

func NewAIService() *AIService {
	provider := os.Getenv("AI_PROVIDER")
	if provider == "" {
		provider = "claude" // Default to Claude
	}

	geminiKey := os.Getenv("GEMINI_API_KEY")
	openaiKey := os.Getenv("OPENAI_API_KEY")
	claudeKey := os.Getenv("ANTHROPIC_AUTH_TOKEN")
	claudeBaseURL := os.Getenv("ANTHROPIC_BASE_URL")
	if claudeBaseURL == "" {
		claudeBaseURL = "https://litellm-339960399182.us-central1.run.app"
	}

	return &AIService{
		provider:      provider,
		geminiKey:     geminiKey,
		openaiKey:     openaiKey,
		claudeKey:     claudeKey,
		claudeBaseURL: claudeBaseURL,
		client:        &http.Client{},
	}
}

func (s *AIService) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	// Use Claude/LiteLLM proxy for embeddings with gemini-embedding-001
	if s.provider == "claude" && s.claudeKey != "" {
		return s.generateEmbeddingClaude(ctx, text)
	}
	if s.provider == "gemini" {
		return s.generateEmbeddingGemini(ctx, text)
	}
	return s.generateEmbeddingOpenAI(ctx, text)
}

// generateEmbeddingClaude uses LiteLLM proxy with gemini-embedding-001 model
func (s *AIService) generateEmbeddingClaude(ctx context.Context, text string) ([]float32, error) {
	url := fmt.Sprintf("%s/v1/embeddings", s.claudeBaseURL)
	
	payload := map[string]interface{}{
		"input": text,
		"model": "gemini-embedding-001",
	}
	
	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.claudeKey)
	
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Claude/LiteLLM API error: %s", string(body))
	}
	
	var result struct {
		Data []struct {
			Embedding []float32 `json:"embedding"`
		} `json:"data"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	if len(result.Data) == 0 {
		return nil, fmt.Errorf("no embedding data returned")
	}
	
	return result.Data[0].Embedding, nil
}

func (s *AIService) generateEmbeddingGemini(ctx context.Context, text string) ([]float32, error) {
	// Gemini doesn't have a direct embeddings API, so we'll use text-embedding-004 model
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/text-embedding-004:embedContent?key=%s", s.geminiKey)
	
	payload := map[string]interface{}{
		"model": "models/text-embedding-004",
		"content": map[string]interface{}{
			"parts": []map[string]string{
				{"text": text},
			},
		},
	}
	
	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Gemini API error: %s", string(body))
	}
	
	var result struct {
		Embedding struct {
			Values []float32 `json:"values"`
		} `json:"embedding"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	if len(result.Embedding.Values) == 0 {
		return nil, fmt.Errorf("no embedding data returned")
	}
	
	return result.Embedding.Values, nil
}

func (s *AIService) generateEmbeddingOpenAI(ctx context.Context, text string) ([]float32, error) {
	url := "https://api.openai.com/v1/embeddings"
	
	payload := map[string]interface{}{
		"input": text,
		"model": "text-embedding-3-small",
	}
	
	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.openaiKey)
	
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		var apiError struct {
			Error struct {
				Message string `json:"message"`
				Type    string `json:"type"`
				Code    string `json:"code"`
			} `json:"error"`
		}
		if err := json.Unmarshal(body, &apiError); err == nil && apiError.Error.Message != "" {
			return nil, fmt.Errorf("OpenAI API error: %s (code: %s)", apiError.Error.Message, apiError.Error.Code)
		}
		return nil, fmt.Errorf("OpenAI API error: %s", string(body))
	}
	
	var result struct {
		Data []struct {
			Embedding []float32 `json:"embedding"`
		} `json:"data"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	if len(result.Data) == 0 {
		return nil, fmt.Errorf("no embedding data returned")
	}
	
	return result.Data[0].Embedding, nil
}

func (s *AIService) SummarizeContent(ctx context.Context, content string) (string, error) {
	prompt := fmt.Sprintf(
		"Summarize the following content in 2-3 concise sentences. Focus on the key points:\n\n%s",
		content,
	)
	
	if s.provider == "claude" && s.claudeKey != "" {
		return s.callClaude(ctx, prompt, 150)
	}
	if s.provider == "gemini" {
		return s.callGeminiPro(ctx, prompt, 150)
	}
	return s.callChatGPT(ctx, prompt, 150)
}

func (s *AIService) GenerateTags(ctx context.Context, content string) ([]string, error) {
	// Truncate content if too long
	truncated := content
	if len(content) > 2000 {
		truncated = content[:2000]
	}
	
	prompt := fmt.Sprintf(
		"Extract 3-5 relevant tags for this content. Return only comma-separated tags, no explanations, no numbering, just tags separated by commas:\n\n%s",
		truncated,
	)
	
	var response string
	var err error
	
	if s.provider == "claude" && s.claudeKey != "" {
		response, err = s.callClaude(ctx, prompt, 50)
	} else if s.provider == "gemini" {
		response, err = s.callGemini(ctx, prompt, 50)
	} else {
		response, err = s.callChatGPT(ctx, prompt, 50)
	}
	
	if err != nil {
		return nil, err
	}
	
	// Parse comma-separated tags
	tagsStr := strings.TrimSpace(response)
	tags := strings.Split(tagsStr, ",")
	
	var cleanedTags []string
	for _, tag := range tags {
		cleaned := strings.TrimSpace(tag)
		if cleaned != "" {
			cleanedTags = append(cleanedTags, cleaned)
		}
	}
	
	return cleanedTags, nil
}

// EnhanceSearchQuery uses Claude to understand and enhance search queries
// Converts plain English into searchable terms with synonyms and related concepts
func (s *AIService) EnhanceSearchQuery(ctx context.Context, query string) (string, error) {
	prompt := fmt.Sprintf(`You are a search query enhancement assistant. Your goal is to help users find content even when they use plain English that doesn't match exact words in the content.

Analyze the following search query and return an improved search query that will find relevant content using semantic understanding.

Examples:
- "things about AI" → "artificial intelligence machine learning neural networks AI"
- "cooking ideas" → "recipes cooking food preparation ingredients"
- "workout tips" → "exercise fitness training health workout"
- "money saving" → "budget savings finance frugal economical"

Your task:
1. Understand the user's intent and what they're really looking for
2. Expand with relevant synonyms, related terms, and alternative phrasings
3. Include both formal and informal terms
4. Keep the original meaning but add searchable keywords
5. Return ONLY the enhanced query with expanded terms, nothing else

Original query: "%s"

Enhanced query:`, query)
	
	if s.provider == "claude" && s.claudeKey != "" {
		enhanced, err := s.callClaude(ctx, prompt, 150)
		if err == nil && enhanced != "" {
			return enhanced, nil
		}
	}
	// Fallback: return original query if Claude not available
	return query, nil
}

// ReRankSearchResults uses Claude to re-rank search results by relevance
func (s *AIService) ReRankSearchResults(ctx context.Context, query string, results []models.SearchResult, topK int) ([]models.SearchResult, error) {
	if len(results) == 0 {
		return results, nil
	}
	
	// Build context for Claude
	var itemsContext strings.Builder
	itemsContext.WriteString(fmt.Sprintf("Search query: %s\n\n", query))
	itemsContext.WriteString("Search results to rank:\n")
	
	for i, result := range results {
		if i >= 10 { // Limit to top 10 for Claude context
			break
		}
		itemsContext.WriteString(fmt.Sprintf("%d. Title: %s\n   Summary: %s\n   Type: %s\n\n", 
			i+1, result.Item.Title, result.Item.Summary, result.Item.Type))
	}
	
	prompt := fmt.Sprintf(`You are a search result ranking assistant. Given a search query and a list of search results, rank them by relevance to the query.

%s

Return ONLY a comma-separated list of numbers (1, 2, 3, etc.) representing the order of relevance, with the most relevant first. For example: "3,1,5,2,4"

Ranked order:`, itemsContext.String())
	
	if s.provider == "claude" && s.claudeKey != "" {
		rankedOrder, err := s.callClaude(ctx, prompt, 50)
		if err != nil {
			// If Claude fails, return original order
			return results, nil
		}
		
		// Parse the ranked order
		indices := parseRankedIndices(rankedOrder, len(results))
		if len(indices) > 0 {
			// Reorder results based on Claude's ranking
			reordered := make([]models.SearchResult, 0, len(indices))
			for _, idx := range indices {
				if idx >= 0 && idx < len(results) {
					reordered = append(reordered, results[idx])
				}
			}
			// Add any remaining results that weren't ranked
			rankedSet := make(map[int]bool)
			for _, idx := range indices {
				rankedSet[idx] = true
			}
			for i, result := range results {
				if !rankedSet[i] {
					reordered = append(reordered, result)
				}
			}
			return reordered, nil
		}
	}
	
	return results, nil
}

// parseRankedIndices parses Claude's ranked order response
func parseRankedIndices(rankedOrder string, maxLen int) []int {
	// Clean the response
	rankedOrder = strings.TrimSpace(rankedOrder)
	rankedOrder = strings.Trim(rankedOrder, "\"")
	rankedOrder = strings.Trim(rankedOrder, "'")
	
	// Extract numbers
	re := regexp.MustCompile(`\d+`)
	matches := re.FindAllString(rankedOrder, -1)
	
	var indices []int
	for _, match := range matches {
		idx := 0
		if _, err := fmt.Sscanf(match, "%d", &idx); err == nil {
			// Convert from 1-based to 0-based index
			idx--
			if idx >= 0 && idx < maxLen {
				indices = append(indices, idx)
			}
		}
	}
	
	return indices
}

// CategorizeContent uses AI to automatically categorize content into sections
func (s *AIService) CategorizeContent(ctx context.Context, title, content, itemType string) (string, error) {
	// Truncate content if too long
	truncated := content
	if len(content) > 1500 {
		truncated = content[:1500]
	}
	
	prompt := fmt.Sprintf(
		`Categorize this content into ONE of these specific sections:
- Technology
- Food & Recipes
- Books & Reading
- Videos & Entertainment
- Shopping & Products
- Articles & News
- Notes & Ideas
- Design & Inspiration
- Travel
- Health & Fitness
- Education & Learning
- Other

Title: %s
Type: %s
Content: %s

Return ONLY the category name, nothing else.`,
		title, itemType, truncated,
	)
	
	var response string
	var err error
	
	if s.provider == "claude" && s.claudeKey != "" {
		response, err = s.callClaude(ctx, prompt, 20)
	} else if s.provider == "gemini" {
		response, err = s.callGemini(ctx, prompt, 20)
	} else {
		response, err = s.callChatGPT(ctx, prompt, 20)
	}
	
	if err != nil {
		return "", err
	}
	
	category := strings.TrimSpace(response)
	// Clean up any extra text
	if strings.Contains(category, "\n") {
		category = strings.Split(category, "\n")[0]
	}
	
	return category, nil
}

// GenerateSemanticSummary creates a concise semantic summary optimized for search
// Uses Claude via LiteLLM proxy, falls back to Gemini/OpenAI if needed
func (s *AIService) GenerateSemanticSummary(ctx context.Context, title, content string) (string, error) {
	// Truncate content if too long
	truncated := content
	if len(content) > 3000 {
		truncated = content[:3000]
	}
	
	prompt := fmt.Sprintf(
		`Create a concise semantic summary (2-3 sentences) of this content that captures key concepts, topics, and ideas. This summary will be used for search, so include important keywords and concepts:
    
    Title: %s
    Content: %s
    
    Summary:`,
		title, truncated,
	)
	
	// Use Claude if available
	if s.provider == "claude" && s.claudeKey != "" {
		return s.callClaude(ctx, prompt, 200)
	}
	
	// Try Gemini first (if provider is gemini)
	if s.provider == "gemini" {
		summary, err := s.callGeminiPro(ctx, prompt, 200)
		// If Gemini fails due to quota/rate limit and OpenAI is available, fallback to OpenAI
		if err != nil && s.openaiKey != "" {
			if strings.Contains(err.Error(), "quota") || strings.Contains(err.Error(), "429") || strings.Contains(err.Error(), "rate limit") || strings.Contains(err.Error(), "503") {
				fmt.Printf("Gemini quota exceeded, falling back to OpenAI for summary generation\n")
				return s.callChatGPT(ctx, prompt, 200)
			}
		}
		return summary, err
	}
	return s.callChatGPT(ctx, prompt, 200)
}

// SummarizeYouTubeVideo generates a short summary for a YouTube video
// Uses Claude via LiteLLM proxy, falls back to Gemini/OpenAI if needed
func (s *AIService) SummarizeYouTubeVideo(ctx context.Context, videoURL, title, description string) (string, error) {
	// Truncate description if too long (keep it reasonable for the API)
	truncatedDesc := description
	if len(description) > 5000 {
		truncatedDesc = description[:5000] + "..."
	}
	
	prompt := fmt.Sprintf(
		`Create a SHORT, concise summary (2-3 sentences maximum) of this YouTube video. Focus only on the main topic and key points. Be brief and informative.

Video Title: %s
Video Description: %s

Provide a brief summary:`,
		title, truncatedDesc,
	)
	
	// Use Claude if available
	if s.provider == "claude" && s.claudeKey != "" {
		return s.callClaude(ctx, prompt, 150)
	}
	
	// Try Gemini first (if provider is gemini)
	if s.provider == "gemini" {
		summary, err := s.callGeminiPro(ctx, prompt, 150)
		// If Gemini fails due to quota/rate limit and OpenAI is available, fallback to OpenAI
		if err != nil && s.openaiKey != "" {
			if strings.Contains(err.Error(), "quota") || strings.Contains(err.Error(), "429") || strings.Contains(err.Error(), "rate limit") || strings.Contains(err.Error(), "503") {
				fmt.Printf("Gemini quota exceeded, falling back to OpenAI for summary generation\n")
				return s.callChatGPT(ctx, prompt, 150)
			}
		}
		return summary, err
	}
	return s.callChatGPT(ctx, prompt, 150)
}

// callGeminiPro specifically uses Gemini 2.5 Pro for better quality summaries
func (s *AIService) callGeminiPro(ctx context.Context, prompt string, maxTokens int) (string, error) {
	// Prioritize Gemini 2.5 Pro for summaries, with fallbacks
	// Try v1 API first, then v1beta, with multiple model options
	models := []struct {
		apiVersion string
		modelName  string
	}{
		{"v1", "gemini-2.5-flash"},              // Try v1 API with Flash
		{"v1", "gemini-2.5-pro"},                // Try v1 API with Pro
		{"v1beta", "gemini-2.5-flash"},          // Try v1beta Flash
		{"v1beta", "gemini-2.5-pro"},            // Try v1beta Pro
		{"v1beta", "gemini-2.5-flash-preview-05-20"},
		{"v1beta", "gemini-2.5-pro-preview-06-05"},
	}
	
	return s.callGeminiWithModels(ctx, prompt, maxTokens, models)
}

func (s *AIService) callGemini(ctx context.Context, prompt string, maxTokens int) (string, error) {
	// Try multiple model names and API versions as fallback
	// Updated to use Gemini 2.5 models which are currently available
	models := []struct {
		apiVersion string
		modelName  string
	}{
		{"v1beta", "gemini-2.5-flash"},
		{"v1beta", "gemini-2.5-pro"},
		{"v1beta", "gemini-2.5-flash-preview-05-20"},
		{"v1beta", "gemini-2.5-pro-preview-06-05"},
		{"v1beta", "gemini-1.5-flash-latest"},
		{"v1beta", "gemini-1.5-pro-latest"},
	}
	
	return s.callGeminiWithModels(ctx, prompt, maxTokens, models)
}

func (s *AIService) callGeminiWithModels(ctx context.Context, prompt string, maxTokens int, models []struct {
	apiVersion string
	modelName  string
}) (string, error) {
	
	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"maxOutputTokens": maxTokens,
			"temperature":     0.7,
		},
	}
	
	jsonData, _ := json.Marshal(payload)
	
	var lastErr error
	for _, model := range models {
		url := fmt.Sprintf("https://generativelanguage.googleapis.com/%s/models/%s:generateContent?key=%s", 
			model.apiVersion, model.modelName, s.geminiKey)
		
		req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := s.client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("failed to call Gemini API: %w", err)
			continue
		}
		
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		
		if resp.StatusCode == http.StatusOK {
			var result struct {
				Candidates []struct {
					Content struct {
						Parts []struct {
							Text string `json:"text"`
						} `json:"parts"`
					} `json:"content"`
				} `json:"candidates"`
				Error *struct {
					Code    int    `json:"code"`
					Message string `json:"message"`
					Status  string `json:"status"`
				} `json:"error"`
			}
			
			if err := json.Unmarshal(body, &result); err != nil {
				lastErr = fmt.Errorf("failed to decode response: %w", err)
				continue
			}
			
			// Check for API errors in response
			if result.Error != nil {
				lastErr = fmt.Errorf("Gemini API error (model: %s, code: %d): %s", model.modelName, result.Error.Code, result.Error.Message)
				// If it's a temporary error (503, 429), continue to next model
				if result.Error.Code == 503 || result.Error.Code == 429 {
					continue
				}
				// For other errors, also continue to try next model
				continue
			}
			
			if len(result.Candidates) == 0 {
				lastErr = fmt.Errorf("no candidates in response from model %s", model.modelName)
				continue
			}
			
			// Check if we have parts with text
			if len(result.Candidates[0].Content.Parts) > 0 {
				text := result.Candidates[0].Content.Parts[0].Text
				if text != "" {
					return strings.TrimSpace(text), nil
				}
			}
			
			// If no text, return error
			if len(result.Candidates) > 0 {
				lastErr = fmt.Errorf("no text content in response from model %s", model.modelName)
			} else {
				lastErr = fmt.Errorf("empty response from model %s", model.modelName)
			}
			continue
		}
		
		// Handle non-200 status codes
		var apiError struct {
			Error struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
				Status  string `json:"status"`
			} `json:"error"`
		}
		if err := json.Unmarshal(body, &apiError); err == nil && apiError.Error.Message != "" {
			lastErr = fmt.Errorf("Gemini API error (model: %s, code: %d): %s", model.modelName, apiError.Error.Code, apiError.Error.Message)
			// For temporary errors, continue to next model
			if apiError.Error.Code == 503 || apiError.Error.Code == 429 {
				continue
			}
		} else {
			lastErr = fmt.Errorf("Gemini API error (model: %s, status: %d): %s", model.modelName, resp.StatusCode, string(body))
		}
	}
	
	return "", fmt.Errorf("all Gemini models failed, last error: %w", lastErr)
}

// callClaude uses Claude API via LiteLLM proxy for text generation
func (s *AIService) callClaude(ctx context.Context, prompt string, maxTokens int) (string, error) {
	url := fmt.Sprintf("%s/v1/chat/completions", s.claudeBaseURL)
	
	// Try different Claude model names available via LiteLLM proxy
	models := []string{
		"claude-sonnet-4-5-20250929",  // Claude Sonnet 4.5 (best quality)
		"claude-opus-4-1-20250805",    // Claude Opus 4.1
		"claude-haiku-4-5-20251001",   // Claude Haiku 4.5 (fastest)
	}
	
	var lastErr error
	for _, model := range models {
		payload := map[string]interface{}{
			"model": model,
			"messages": []map[string]interface{}{
				{
					"role":    "user",
					"content": prompt,
				},
			},
			"max_tokens": maxTokens,
			"temperature": 0.7,
		}
		
		jsonData, _ := json.Marshal(payload)
		req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.claudeKey)
		
		resp, err := s.client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("failed to call Claude API: %w", err)
			continue
		}
		defer resp.Body.Close()
		
		if resp.StatusCode == http.StatusOK {
			var result struct {
				Choices []struct {
					Message struct {
						Content string `json:"content"`
					} `json:"message"`
				} `json:"choices"`
			}
			
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				lastErr = fmt.Errorf("failed to decode response: %w", err)
				continue
			}
			
			if len(result.Choices) == 0 {
				lastErr = fmt.Errorf("no response from Claude")
				continue
			}
			
			return strings.TrimSpace(result.Choices[0].Message.Content), nil
		}
		
		body, _ := io.ReadAll(resp.Body)
		var apiError struct {
			Error struct {
				Message string `json:"message"`
				Type    string `json:"type"`
				Code    string `json:"code"`
			} `json:"error"`
		}
		if err := json.Unmarshal(body, &apiError); err == nil && apiError.Error.Message != "" {
			lastErr = fmt.Errorf("Claude API error (model: %s): %s (code: %s)", model, apiError.Error.Message, apiError.Error.Code)
			// Continue to next model if this one fails
			continue
		}
		lastErr = fmt.Errorf("Claude API error (model: %s): %s", model, string(body))
	}
	
	return "", fmt.Errorf("all Claude models failed, last error: %w", lastErr)
}

func (s *AIService) callChatGPT(ctx context.Context, prompt string, maxTokens int) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"
	
	payload := map[string]interface{}{
		"model": "gpt-4o-mini",
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"max_tokens": maxTokens,
		"temperature": 0.7,
	}
	
	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.openaiKey)
	
	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call OpenAI API: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		var apiError struct {
			Error struct {
				Message string `json:"message"`
				Type    string `json:"type"`
				Code    string `json:"code"`
			} `json:"error"`
		}
		if err := json.Unmarshal(body, &apiError); err == nil && apiError.Error.Message != "" {
			return "", fmt.Errorf("OpenAI API error: %s (code: %s)", apiError.Error.Message, apiError.Error.Code)
		}
		return "", fmt.Errorf("OpenAI API error: %s", string(body))
	}
	
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}
	
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}
	
	return strings.TrimSpace(result.Choices[0].Message.Content), nil
}
