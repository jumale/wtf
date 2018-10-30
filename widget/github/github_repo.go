package github

import (
	"context"
	ghb "github.com/google/go-github/github"
	"github.com/senorprogrammer/wtf/wtf"
	"golang.org/x/oauth2"
	"net/http"
)

type Repo struct {
	Name         string
	Owner        string
	PullRequests []*ghb.PullRequest
	RemoteRepo   *ghb.Repository
	config       *Config
	logger       wtf.Logger
}

// Refresh reloads the github data via the Github API
func (repo *Repo) Refresh() {
	repo.logger.Debugf("Github: refresh %s repo", repo.Name)
	repo.PullRequests, _ = repo.loadPullRequests()
	repo.RemoteRepo, _ = repo.loadRemoteRepository()
}

/* -------------------- Counts -------------------- */

func (repo *Repo) IssueCount() int {
	if repo.RemoteRepo == nil {
		return 0
	}

	return *repo.RemoteRepo.OpenIssuesCount
}

func (repo *Repo) PullRequestCount() int {
	return len(repo.PullRequests)
}

func (repo *Repo) StarCount() int {
	if repo.RemoteRepo == nil {
		return 0
	}

	return *repo.RemoteRepo.StargazersCount
}

/* -------------------- Unexported Functions -------------------- */

func (repo *Repo) isGitHubEnterprise() bool {
	if len(repo.config.BaseURL) > 0 {
		if len(repo.config.UploadURL) == 0 {
			repo.config.UploadURL = repo.config.BaseURL
		}
		return true
	}
	return false
}

func (repo *Repo) oauthClient() *http.Client {
	tokenService := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: repo.config.ApiKey},
	)

	return oauth2.NewClient(context.Background(), tokenService)
}

func (repo *Repo) githubClient() (*ghb.Client, error) {
	oauthClient := repo.oauthClient()

	if repo.isGitHubEnterprise() {
		return ghb.NewEnterpriseClient(repo.config.BaseURL, repo.config.UploadURL, oauthClient)
	}

	return ghb.NewClient(oauthClient), nil
}

// myPullRequests returns a list of pull requests created by username on this repo
func (repo *Repo) myPullRequests(username string) []*ghb.PullRequest {
	var prs []*ghb.PullRequest

	for _, pr := range repo.PullRequests {
		user := *pr.User

		if *user.Login == username {
			prs = append(prs, pr)
		}
	}

	if repo.config.ShowStatus {
		prs = repo.individualPRs(prs)
	}

	return prs
}

// individualPRs takes a list of pull requests (presumably returned from
// github.PullRequests.List) and fetches them individually to get more detailed
// status info on each. see: https://developer.github.com/v3/git/#checking-mergeability-of-pull-requests
func (repo *Repo) individualPRs(prs []*ghb.PullRequest) []*ghb.PullRequest {
	github, err := repo.githubClient()
	if err != nil {
		return prs
	}

	var ret []*ghb.PullRequest
	for i := range prs {
		pr, _, err := github.PullRequests.Get(context.Background(), repo.Owner, repo.Name, prs[i].GetNumber())
		if err != nil {
			// worst case, just keep the original one
			ret = append(ret, prs[i])
		} else {
			ret = append(ret, pr)
		}
	}
	return ret
}

// myReviewRequests returns a list of pull requests for which username has been
// requested to do a code review
func (repo *Repo) myReviewRequests(username string) []*ghb.PullRequest {
	var prs []*ghb.PullRequest

	for _, pr := range repo.PullRequests {
		for _, reviewer := range pr.RequestedReviewers {
			if *reviewer.Login == username {
				prs = append(prs, pr)
			}
		}
	}

	return prs
}

func (repo *Repo) loadPullRequests() ([]*ghb.PullRequest, error) {
	github, err := repo.githubClient()

	if err != nil {
		return nil, err
	}

	opts := &ghb.PullRequestListOptions{}

	prs, _, err := github.PullRequests.List(context.Background(), repo.Owner, repo.Name, opts)

	if err != nil {
		return nil, err
	}

	return prs, nil
}

func (repo *Repo) loadRemoteRepository() (*ghb.Repository, error) {
	github, err := repo.githubClient()

	if err != nil {
		return nil, err
	}

	repository, _, err := github.Repositories.Get(context.Background(), repo.Owner, repo.Name)

	if err != nil {
		return nil, err
	}

	return repository, nil
}
