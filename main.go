package main

import (
	"context"
	"github.com/google/go-github/github"
	"github.com/paynejacob/release-note-manager/pkg/configuration"
	"github.com/paynejacob/release-note-manager/pkg/readme"
	"golang.org/x/oauth2"
	"net/http"
	"github.com/gorilla/mux"
)

var githubClient *github.Client

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/oauth2/callback/", oauthCallback).Methods("GET")
	r.HandleFunc("/{owner}/{repo}/milestones/{milestone}/", registerRelease).Methods("PUT")
	r.HandleFunc("/webhook/", webhook)
}

func registerRelease(w http.ResponseWriter, r *http.Request) {
	owner := mux.Vars(r)["owner"]
	repo := mux.Vars(r)["repo"]
	milestone := mux.Vars(r)["milestone"]

	issueChan := make(chan *github.Issue, 0)
	ctx := context.Background()

	// TODO: create release empty and store id, else lookup id, CAHCHCHHCHCHCHCHCHCHCHOla here

	release := github.RepositoryRelease{}
	release.Draft = getBoolPtr(true)
	release.Name = &milestone

	// TODO: get the configuration file from the repo
	sections := []configuration.Section{}


	var err error
	go func() {
		err = listMilestoneIssues(ctx, githubClient, issueChan, owner, repo, milestone)
	}()
	body := readme.GenerateMarkdown(readme.ReadmeFromIssue(issueChan, sections), owner, repo, sections)
	release.Body = &body

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	githubClient.Repositories.CreateRelease(context.TODO(), owner, repo, &release)
}

func oauthCallback(w http.ResponseWriter, r *http.Request) {
	// create user session
	// get a token for the app
	// store it in memory
}

func webhook(w http.ResponseWriter, r *http.Request) {
	// watch for issue an milestone changes
	// update existing notes with changes
}

func getBoolPtr(v bool) *bool {
	t := true
	f := false

	if v {
		return &t
	}

	return &f
}

func listMilestoneIssues(ctx context.Context, ghc *github.Client, issueChan chan *github.Issue, owner, repo, milestone string) error {
	opts := github.IssueListByRepoOptions{Milestone: milestone}
	for {
		issues, resp, err := ghc.Issues.ListByRepo(ctx, owner, repo, &opts)
		if err != nil {
			close(issueChan)
			return err
		}

		for i := 0; i < len(issues); i++ {
			issueChan <- issues[i]
		}

		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	close(issueChan)

	return nil
}
