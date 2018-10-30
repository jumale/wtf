package github

import (
	"github.com/gdamore/tcell"
	"github.com/senorprogrammer/wtf/wtf"
)

const HelpText = `
  Keyboard commands for GitHub:

    /: Show/hide this help window
    h: Previous git repository
    l: Next git repository
    r: Refresh the data

    arrow left:  Previous git repository
    arrow right: Next git repository

    return: Open the selected repository in a browser
`

type Config struct {
	wtf.WidgetConfig `yaml:",inline"`
	Username         string                 `yaml:"username"`
	ApiKey           string                 `yaml:"apiKey"`
	BaseURL          string                 `yaml:"baseURL"`
	UploadURL        string                 `yaml:"uploadURL"`
	ShowStatus       bool                   `yaml:"enableStatus"`
	Repositories     map[string]interface{} `yaml:"repositories"`
}

type Widget struct {
	*wtf.HelpfulWidget
	*wtf.TextWidget

	githubRepos []*Repo
	idx         int
	config      *Config
	fs          *wtf.FileSystem
	formatter   *wtf.Formatter
	logger      wtf.Logger
}

func (widget *Widget) Name() string {
	return "github"
}

func (widget *Widget) Init(configure wtf.UnmarshalFunc, context *wtf.AppContext) error {
	context.Logger.Debug("Github: init")

	widget.config = &Config{
		ShowStatus: false,
	}
	if err := configure(widget.config); err != nil {
		return err
	}

	widget.TextWidget = wtf.NewTextWidget(context.App, "GitHub", widget.config.WidgetConfig, true)
	widget.HelpfulWidget = wtf.NewHelpfulWidget(context.App, context.Pages, widget.TextView, HelpText)

	widget.githubRepos = widget.buildRepoCollection(widget.config.Repositories)

	widget.fs = &context.FS
	widget.formatter = &context.Formatter
	widget.logger = context.Logger

	widget.TextView.SetInputCapture(widget.keyboardIntercept)

	return nil
}

/* -------------------- Exported Functions -------------------- */

func (widget *Widget) Refresh() {
	widget.logger.Debug("Github: refresh")
	for _, repo := range widget.githubRepos {
		repo.Refresh()
	}

	widget.display()
}

func (widget *Widget) Next() {
	widget.idx = widget.idx + 1
	if widget.idx == len(widget.githubRepos) {
		widget.idx = 0
	}

	widget.display()
}

func (widget *Widget) Prev() {
	widget.idx = widget.idx - 1
	if widget.idx < 0 {
		widget.idx = len(widget.githubRepos) - 1
	}

	widget.display()
}

/* -------------------- Unexported Functions -------------------- */

func (widget *Widget) buildRepoCollection(repoData map[string]interface{}) []*Repo {
	var githubRepos []*Repo

	for name, owner := range repoData {
		repo := &Repo{
			Name:   name,
			Owner:  owner.(string),
			config: widget.config,
			logger: widget.logger,
		}
		githubRepos = append(githubRepos, repo)
	}

	return githubRepos
}

func (widget *Widget) currentGithubRepo() *Repo {
	if len(widget.githubRepos) == 0 {
		return nil
	}

	if widget.idx < 0 || widget.idx >= len(widget.githubRepos) {
		return nil
	}

	return widget.githubRepos[widget.idx]
}

func (widget *Widget) keyboardIntercept(event *tcell.EventKey) *tcell.EventKey {
	switch string(event.Rune()) {
	case "/":
		widget.ShowHelp()
		return nil
	case "h":
		widget.Prev()
		return nil
	case "l":
		widget.Next()
		return nil
	case "r":
		widget.Refresh()
		return nil
	}

	switch event.Key() {
	case tcell.KeyEnter:
		widget.openRepo()
		return nil
	case tcell.KeyLeft:
		widget.Prev()
		return nil
	case tcell.KeyRight:
		widget.Next()
		return nil
	default:
		return event
	}
}

func (widget *Widget) openRepo() {
	repo := widget.currentGithubRepo()

	if repo != nil {
		widget.fs.OpenFile(*repo.RemoteRepo.HTMLURL)
	}
}
