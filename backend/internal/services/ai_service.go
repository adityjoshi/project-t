package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type AIService struct {
	provider   string
	geminiKey  string
	openaiKey  string
	client     *http.Client
}

func NewAIService() *AIService {
	provider := os.Getenv("AI_PROVIDER")
	if provider == "" {
		provider = "gemini" // Default to Gemini (free)
	}

	geminiKey := os.Getenv("GEMINI_API_KEY")
	openaiKey := os.Getenv("OPENAI_API_KEY")

	return &AIService{
		provider:  provider,
		geminiKey: geminiKey,
		openaiKey: openaiKey,
		client:    &http.Client{},
	}
}

func (s *AIService) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	if s.provider == "gemini" {
		return s.generateEmbeddingGemini(ctx, text)
	}
	return s.generateEmbeddingOpenAI(ctx, text)
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
	
	if s.provider == "gemini" {
		return s.callGemini(ctx, prompt, 150)
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
	
	if s.provider == "gemini" {
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
	
	if s.provider == "gemini" {
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
	
	if s.provider == "gemini" {
		return s.callGemini(ctx, prompt, 200)
	}
	return s.callChatGPT(ctx, prompt, 200)
}

// SummarizeYouTubeVideo generates a summary for a YouTube video using Gemini
func (s *AIService) SummarizeYouTubeVideo(ctx context.Context, videoURL, title, description string) (string, error) {
	// Truncate description if too long (keep it reasonable for the API)
	truncatedDesc := description
	if len(description) > 5000 {
		truncatedDesc = description[:5000] + "..."
	}
	
	prompt := fmt.Sprintf(
		`You are summarizing a YouTube video. Create a concise, informative summary (3-4 sentences) that captures the main topics, key points, and important information discussed in the video. Focus on what the video is actually about, not just promotional content or links. Make it useful for someone who wants to understand the video's content without watching it.

Video Title: %s
Video Description: %s
Video URL: %s

Based on the title and description above, provide a clear summary of what this video covers:`,
		title, truncatedDesc, videoURL,
	)
	
	if s.provider == "gemini" {
		return s.callGemini(ctx, prompt, 300)
	}
	return s.callChatGPT(ctx, prompt, 300)
}

func (s *AIService) callGemini(ctx context.Context, prompt string, maxTokens int) (string, error) {
	// Use gemini-1.5-flash model (faster and more available than gemini-pro)
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=%s", s.geminiKey)
	
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
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call Gemini API: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Gemini API error: %s", string(body))
	}
	
	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}
	
	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response from Gemini")
	}
	
	return strings.TrimSpace(result.Candidates[0].Content.Parts[0].Text), nil
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
