package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	http.HandleFunc("/", handleRequest)
	fmt.Println("Starting server on :3333")
	http.ListenAndServe(":3333", nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/" {
		listPatterns(w, r)
		return
	}

	patternName := strings.TrimPrefix(path, "/")
	executePattern(w, r, patternName)
}

func listPatterns(w http.ResponseWriter, r *http.Request) {
	patternsDir := "data/patterns"
	files, err := ioutil.ReadDir(patternsDir)
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
	systemFile := filepath.Join("data", "patterns", patternName, "system.md")
	if _, err := os.Stat(systemFile); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	content, err := ioutil.ReadFile(systemFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(content)
}
