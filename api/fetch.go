package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
	"net/http"
	"net/http/httptest"
	"strings"
)

const (
	accessToken = "hidden"
)

var client *github.Client

func handleRepositories(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Set up the GitHub client with authentication
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	tc := oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)

	http.HandleFunc("/repositories", handleRepositories)

	tags := r.URL.Query().Get("tags")
	if tags == "" {
		http.Error(w, "Tags parameter is required", http.StatusBadRequest)
		return
	}

	tagList := strings.Split(tags, ",")

	// Search for repositories with the specified tags
	query := fmt.Sprintf("topic:%s", strings.Join(tagList, " topic:"))

	fmt.Println(context.Background())
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

func main() {
	req, err := http.NewRequest("GET", "/repositories?tags=go", nil)
	if err != nil {
		fmt.Print("wtf")
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleRepositories)

	handler.ServeHTTP(rr, req)

}
