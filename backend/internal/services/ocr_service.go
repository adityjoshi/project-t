package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type OCRService struct {
	client    *http.Client
	geminiKey string
}

func NewOCRService() *OCRService {
	return &OCRService{
		client:    &http.Client{},
		geminiKey: os.Getenv("GEMINI_API_KEY"),
	}
}

// ExtractTextFromImage uses Gemini Vision API to extract text from images
func (s *OCRService) ExtractTextFromImage(ctx context.Context, imageURL string) (string, error) {
	if s.geminiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY not set")
	}

	// Download image
	resp, err := http.Get(imageURL)
	if err != nil {
		return "", fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download image: status %d", resp.StatusCode)
	}

	// Read image data
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read image data: %w", err)
	}

	// Encode to base64
	base64Image := base64.StdEncoding.EncodeToString(imageData)

	// Use Gemini Vision API to extract text
	return s.callGeminiVision(ctx, base64Image, resp.Header.Get("Content-Type"))
}

// ExtractTextFromImageData extracts text from image data (for direct uploads)
func (s *OCRService) ExtractTextFromImageData(ctx context.Context, imageData []byte, mimeType string) (string, error) {
	if s.geminiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY not set")
	}

	base64Image := base64.StdEncoding.EncodeToString(imageData)
	return s.callGeminiVision(ctx, base64Image, mimeType)
}

func (s *OCRService) callGeminiVision(ctx context.Context, base64Image, mimeType string) (string, error) {
	// Default to jpeg if mime type not provided
	if mimeType == "" {
		mimeType = "image/jpeg"
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-pro-vision:generateContent?key=%s", s.geminiKey)

	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{
						"text": "Extract all text from this image. Return only the text content, no explanations or formatting.",
					},
					{
						"inline_data": map[string]interface{}{
							"mime_type": mimeType,
							"data":      base64Image,
						},
					},
				},
			},
		},
	}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call Gemini Vision API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Gemini Vision API error: %s", string(body))
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
		return "", fmt.Errorf("no text extracted from image")
	}

	text := strings.TrimSpace(result.Candidates[0].Content.Parts[0].Text)
	return text, nil
}

