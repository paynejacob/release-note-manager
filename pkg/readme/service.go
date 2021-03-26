package readme

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/gorilla/mux"
	"github.com/paynejacob/release-note-manager/pkg/configuration"
	"net/http"
)

type ReleaseManager struct {
	ghc * github.Client
	releases map[string]int64
}

func (r *ReleaseManager) GetReleaseID(owner, repo, milestone string) (int, error) {

}

func (r *ReleaseManager) findRelease(owner, repo, milestone string) (bool, error) {
	opts := github.ListOptions{}

	for {
		releases, resp, err := r.ghc.Repositories.ListReleases(context.TODO(), owner, repo, &opts)

		if err != nil {
			return false, err
		}

		for i := 0; i < len(releases); i++ {
			if releases[i].GetName() == milestone {
				r.releases[fmt.Sprintf("%s/%s", owner, repo)] = releases[i].GetID()
				return true, nil
			}
		}

		if resp.NextPage == 0 {
			return false, nil
		}

		opts.Page = resp.NextPage
	}
}

func (r *ReleaseManager) createReleaseDraft(owner, repo, milestone string) (int64, error) {
	release := github.RepositoryRelease{}
	release.Name = &milestone
	release.Draft = getBoolPtr(true)

	newRelease, _, err := r.ghc.Repositories.CreateRelease(context.TODO(), owner, repo, &release)



	return newRelease.GetID(), err
}

type Service struct {

}

func (s *Service) registerRelease(w http.ResponseWriter, r *http.Request) {
	owner := mux.Vars(r)["owner"]
	repo := mux.Vars(r)["repo"]
	milestone := mux.Vars(r)["milestone"]

	issueChan := make(chan *github.Issue, 0)
	ctx := context.Background()

	// TODO: create release empty and store id, else lookup id, cache: here

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

func (s *Service) oauthCallback(w http.ResponseWriter, r *http.Request) {
	// create user session
	// get a token for the app
	// store it in memory
}

func (s *Service) webhook(w http.ResponseWriter, r *http.Request) {
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
