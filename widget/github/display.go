package github

import (
	"fmt"

	"github.com/google/go-github/github"
)

func (widget *Widget) display() {
	repo := widget.currentGithubRepo()
	if repo == nil {
		widget.TextView.SetText(" GitHub repo data is unavailable ")
		return
	}

	widget.TextView.SetTitle(widget.ContextualTitle(fmt.Sprintf("%s - %s", widget.title, widget.title(repo))))

	str := widget.formatter.SigilStr(len(widget.githubRepos), widget.idx, widget.TextView) + "\n"
	str = str + " [red]Stats[white]\n"
	str = str + widget.displayStats(repo)
	str = str + "\n"
	str = str + " [red]Open Review Requests[white]\n"
	str = str + widget.displayMyReviewRequests(repo, widget.config.Username)
	str = str + "\n"
	str = str + " [red]My Pull Requests[white]\n"
	str = str + widget.displayMyPullRequests(repo, widget.config.Username)

	widget.TextView.SetText(str)
}

func (widget *Widget) displayMyPullRequests(repo *Repo, username string) string {
	prs := repo.myPullRequests(username)

	if len(prs) == 0 {
		return " [grey]none[white]\n"
	}

	str := ""
	for _, pr := range prs {
		str = str + fmt.Sprintf(" %s[green]%4d[white] %s\n", widget.prMergeStatus(pr), *pr.Number, *pr.Title)
	}

	return str
}

func (widget *Widget) displayMyReviewRequests(repo *Repo, username string) string {
	prs := repo.myReviewRequests(username)

	if len(prs) == 0 {
		return " [grey]none[white]\n"
	}

	str := ""
	for _, pr := range prs {
		str = str + fmt.Sprintf(" [green]%4d[white] %s\n", *pr.Number, *pr.Title)
	}

	return str
}

func (widget *Widget) displayStats(repo *Repo) string {
	str := fmt.Sprintf(
		" PRs: %d  Issues: %d  Stars: %d\n",
		repo.PullRequestCount(),
		repo.IssueCount(),
		repo.StarCount(),
	)

	return str
}

func (widget *Widget) title(repo *Repo) string {
	return fmt.Sprintf("[green]%s - %s[white]", repo.Owner, repo.Name)
}

var mergeIcons = map[string]string{
	"dirty":    "[red]![white] ",
	"clean":    "[green]✔[white] ",
	"unstable": "[red]✖[white] ",
	"blocked":  "[red]✖[white] ",
}

func (widget *Widget) prMergeStatus(pr *github.PullRequest) string {
	if widget.config.ShowStatus == false {
		return ""
	}
	if str, ok := mergeIcons[pr.GetMergeableState()]; ok {
		return str
	}
	return "? "
}
