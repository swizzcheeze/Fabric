package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LMStudio Structures
type LMStudioRequest struct {
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type LMStudioResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

func main() {
	http.HandleFunc("/", handleRequest)
	fmt.Println("Starting server on :3333")
	http.ListenAndServe(":3333", nil)
}

// Input represents the JSON payload from the user
type Input struct {
	Input string `json:"input"`
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/" {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		listPatterns(w, r)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	patternName := strings.TrimPrefix(path, "/")
	executePattern(w, r, patternName)
}

func listPatterns(w http.ResponseWriter, r *http.Request) {
	patternsDir := "./patterns"
	files, err := os.ReadDir(patternsDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var patterns []string
	for _, file := range files {
		if file.IsDir() {
			patterns = append(patterns, file.Name())
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(patterns)
}

func executePattern(w http.ResponseWriter, r *http.Request, patternName string) {
	systemFile := filepath.Join("./patterns", patternName, "system.md")
	if _, err := os.Stat(systemFile); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	var userInput Input
	err := json.NewDecoder(r.Body).Decode(&userInput)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	systemContent, err := os.ReadFile(systemFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lmStudioResponse, err := callLMStudio(string(systemContent), userInput.Input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lmStudioResponse)
}

func callLMStudio(systemPrompt, userPrompt string) (*LMStudioResponse, error) {
	lmStudioEndpoint := os.Getenv("LM_STUDIO_ENDPOINT")
	if lmStudioEndpoint == "" {
		lmStudioEndpoint = "http://localhost:1234/v1/chat/completions"
	}

	reqBody := LMStudioRequest{
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", lmStudioEndpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var lmStudioResponse LMStudioResponse
	err = json.NewDecoder(resp.Body).Decode(&lmStudioResponse)
	if err != nil {
		return nil, err
	}

	return &lmStudioResponse, nil
}
