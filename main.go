package main

import (
	"github.com/senorprogrammer/wtf/cfg"
	wtfFlags "github.com/senorprogrammer/wtf/flags"
	"github.com/senorprogrammer/wtf/widget/clocks"
	"github.com/senorprogrammer/wtf/widget/github"
	"github.com/senorprogrammer/wtf/widget/logger"
	"github.com/senorprogrammer/wtf/widget/security"
	"github.com/senorprogrammer/wtf/widget/status"
	"github.com/senorprogrammer/wtf/widget/system"
	"github.com/senorprogrammer/wtf/widget/textfile"
	"github.com/senorprogrammer/wtf/wtf"
	"log"
	"path"
)

// config parses the config.yml file and makes available the settings within
var (
	commit  = "dev"
	date    = "dev"
	version = "dev"
)

const configDir = "~/.config/wtf"

func main() {
	cfg.MigrateOldConfig()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flags := wtfFlags.NewFlags()
	flags.Parse()
	flags.Display(version)

	//configFile := flags.Config
	configFile := path.Join(configDir, "example.yml")
	watchConfigChanges := true

	app, err := wtf.NewApp(configDir, configFile, watchConfigChanges)
	checkErr(err)

	app.RegisterWidget("system", system.CreateConstructor(date, version))
	app.RegisterWidget("github", github.New)
	app.RegisterWidget("logger", logger.New)
	app.RegisterWidget("clocks", clocks.New)
	app.RegisterWidget("security", security.New)
	app.RegisterWidget("status", status.New)
	app.RegisterWidget("textfile", textfile.New)

	err = app.Run()
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		log.Fatalf("%+v\n", err)
	}
}

/*
func addWidget(app *tview.Application, pages *tview.Pages, widgetName string) {
	// Always in alphabetical order
	switch widgetName {
	case "bamboohr":
		widgets = append(widgets, bamboohr.NewWidget(app))
	case "bargraph":
		widgets = append(widgets, bargraph.NewWidget())
	case "bittrex":
		widgets = append(widgets, bittrex.NewWidget(app))
	case "blockfolio":
		widgets = append(widgets, blockfolio.NewWidget(app))
	case "circleci":
		widgets = append(widgets, circleci.NewWidget(app))
	case "clocks":
		widgets = append(widgets, clocks.NewWidget(app))
	case "cmdrunner":
		widgets = append(widgets, cmdrunner.NewWidget(app))
	case "cryptolive":
		widgets = append(widgets, cryptolive.NewWidget(app))
	case "datadog":
		widgets = append(widgets, datadog.NewWidget(app))
	case "gcal":
		widgets = append(widgets, gcal.NewWidget(app))
	case "gerrit":
		widgets = append(widgets, gerrit.NewWidget(app, pages))
	case "git":
		widgets = append(widgets, git.NewWidget(app, pages))
	case "github":
		widgets = append(widgets, github.NewWidget(app, pages))
	case "gitlab":
		widgets = append(widgets, gitlab.NewWidget(app, pages))
	case "gitter":
		widgets = append(widgets, gitter.NewWidget(app, pages))
	case "gspreadsheets":
		widgets = append(widgets, gspreadsheets.NewWidget(app))
	case "hackernews":
		widgets = append(widgets, hackernews.NewWidget(app, pages))
	case "ipapi":
		widgets = append(widgets, ipapi.NewWidget(app))
	case "ipinfo":
		widgets = append(widgets, ipinfo.NewWidget(app))
	case "jenkins":
		widgets = append(widgets, jenkins.NewWidget(app, pages))
	case "jira":
		widgets = append(widgets, jira.NewWidget(app, pages))
	case "logger":
		widgets = append(widgets, logger.NewWidget(app))
	case "newrelic":
		widgets = append(widgets, newrelic.NewWidget(app))
	case "opsgenie":
		widgets = append(widgets, opsgenie.NewWidget(app))
	case "power":
		widgets = append(widgets, power.NewWidget(app))
	case "prettyweather":
		widgets = append(widgets, prettyweather.NewWidget(app))
	case "security":
		widgets = append(widgets, security.NewWidget(app))
	case "status":
		widgets = append(widgets, status.NewWidget(app))
	case "system":
		widgets = append(widgets, system.NewWidget(app, date, version))
	case "spotify":
		widgets = append(widgets, spotify.NewWidget(app, pages))
	case "textfile":
		widgets = append(widgets, textfile.NewWidget(app, pages))
	case "todo":
		widgets = append(widgets, todo.NewWidget(app, pages))
	case "todoist":
		widgets = append(widgets, todoist.NewWidget(app, pages))
	case "travisci":
		widgets = append(widgets, travisci.NewWidget(app, pages))
	case "trello":
		widgets = append(widgets, trello.NewWidget(app))
	case "twitter":
		widgets = append(widgets, twitter.NewWidget(app, pages))
	case "weather":
		widgets = append(widgets, weather.NewWidget(app, pages))
	case "zendesk":
		widgets = append(widgets, zendesk.NewWidget(app))
	default:
	}
}
*/
