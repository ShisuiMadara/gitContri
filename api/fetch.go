package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-github/v39/github"
)

const (
	accessToken = "sd"
)

var client *github.Client

func handleRepositories(w http.ResponseWriter, r *http.Request) {
	tags := r.URL.Query().Get("tags")
	if tags == "" {
		http.Error(w, "Tags parameter is required", http.StatusBadRequest)
		return
	}

	tagList := strings.Split(tags, ",")

	// Search for repositories with the specified tags
	query := fmt.Sprintf("topic:%s", strings.Join(tagList, " topic:"))
	repos, _, err := client.Search.Repositories(context.Background(), query, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error searching repositories: %s", err), http.StatusInternalServerError)
		return
	}

	var repoNames []string
	for _, repo := range repos.Repositories {
		repoNames = append(repoNames, *repo.Name)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"repositories": %s}`, strings.Join(repoNames, ", "))
}
