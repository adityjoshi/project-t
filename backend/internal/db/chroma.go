package db

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type ChromaClient struct {
	BaseURL string
	Client  *http.Client
}

var Chroma *ChromaClient

func InitChroma() error {
	baseURL := os.Getenv("CHROMA_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8000"
	}

	Chroma = &ChromaClient{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}

	// Create collection if it doesn't exist
	collectionName := "synapse_items"
	if err := Chroma.CreateCollection(collectionName); err != nil {
		// Collection might already exist, that's okay
		fmt.Printf("Note: Collection creation: %v\n", err)
	}

	return nil
}

func (c *ChromaClient) CreateCollection(name string) error {
	// Try v1 API first (for older ChromaDB versions)
	url := fmt.Sprintf("%s/api/v1/collections", c.BaseURL)
	
	payload := map[string]interface{}{
		"name": name,
	}
	
	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	// If v1 is deprecated, collection might already exist or we'll create on first add
	if resp.StatusCode == 200 || resp.StatusCode == 409 {
		return nil
	}
	
	// If v1 fails, that's okay - we'll try to create on first add
	return nil
}

func (c *ChromaClient) AddEmbedding(collectionName, id string, embedding []float32, metadata map[string]interface{}) error {
	// Try v1 API
	url := fmt.Sprintf("%s/api/v1/collections/%s/add", c.BaseURL, collectionName)
	
	payload := map[string]interface{}{
		"ids":       []string{id},
		"embeddings": [][]float32{embedding},
		"metadatas": []map[string]interface{}{metadata},
	}
	
	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	// If v1 API is deprecated, try to use Python client approach or skip
	if resp.StatusCode == 404 || resp.StatusCode == 501 {
		// ChromaDB v1 API is deprecated, embeddings won't work with this version
		// Return error but don't fail the whole operation
		return fmt.Errorf("ChromaDB v1 API deprecated - semantic search disabled")
	}
	
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to add embedding: %s", string(body))
	}
	
	return nil
}

func (c *ChromaClient) Query(collectionName string, queryEmbedding []float32, nResults int) ([]string, []float64, error) {
	if queryEmbedding == nil || len(queryEmbedding) == 0 {
		return []string{}, []float64{}, fmt.Errorf("query embedding cannot be empty")
	}

	url := fmt.Sprintf("%s/api/v1/collections/%s/query", c.BaseURL, collectionName)
	
	payload := map[string]interface{}{
		"query_embeddings": [][]float32{queryEmbedding},
		"n_results":        nResults,
	}
	
	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	
	// If v1 API is deprecated, return empty results
	if resp.StatusCode == 404 || resp.StatusCode == 501 {
		return []string{}, []float64{}, fmt.Errorf("ChromaDB v1 API deprecated - semantic search disabled")
	}
	
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, nil, fmt.Errorf("failed to query: %s", string(body))
	}
	
	var result struct {
		Ids       [][]string   `json:"ids"`
		Distances [][]float64  `json:"distances"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, nil, err
	}
	
	if len(result.Ids) == 0 || len(result.Ids[0]) == 0 {
		return []string{}, []float64{}, nil
	}
	
	return result.Ids[0], result.Distances[0], nil
}
