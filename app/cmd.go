package app

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/labctl/labctl/helpers"
	"github.com/labctl/labctl/utils/tx"
	"github.com/posener/complete"
	"github.com/willabides/kongplete"

	log "github.com/sirupsen/logrus"
)

var cli struct {
	Debug    int    `short:"d" help:"Enable debug mode." type:"counter" group:"Global Flags"`
	Settings string `help:"Settings yaml file." type:"existingfile" group:"Global Flags"`

	Color  CmdColor  `cmd:"" help:"Add some color to any Linux command (i.e. ssh)." short:"c" passthrough:""`
	Config CmdConfig `cmd:"" help:"Render a template."`
	Serve  CmdServe  `cmd:"" help:"Serve the web UI."`

	InstallCompletions kongplete.InstallCompletions `cmd:"" help:"install shell completions"`
	Version            CmdVersion                   `cmd:"" help:"show the labctl version"`
}

func GetCmdLineParser(ast interface{}) *kong.Kong {
	return kong.Must(ast,
		kong.Name("labctl"),
		kong.Description("Control your network lab."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: false,
		}))
}

// App wide context
var Ctx *helpers.Context

func Main() {
	parser := GetCmdLineParser(&cli)

	// Handle completion requests
	kongplete.Complete(parser,
		kongplete.WithPredictor("topo", complete.PredictFiles("*.clab.y*ml")),
		kongplete.WithPredictor("dir", complete.PredictDirs("*")),
	)

	kctx, err := parser.Parse(os.Args[1:])
	parser.FatalIfErrorf(err)

	s := helpers.NewSettings(true)
	if err = s.AddSettings("~/.labctl.yml", true); err != nil {
		kctx.FatalIfErrorf(err)
	}
	if err = s.InitColors(); err != nil {
		kctx.FatalIfErrorf(err)
	}

	log.SetOutput(os.Stdout)
	if cli.Debug > 0 {
		log.SetLevel(log.DebugLevel)
	}
	tx.DebugCount = cli.Debug

	Ctx = &helpers.Context{
		Command:    kctx.Command(),
		DebugCount: cli.Debug,
		Settings:   s,
		Output:     &helpers.LogOutput{},
	}

	fetchLatestVersion()

	// Call the Run() method of the selected parsed command.
	err = kctx.Run(Ctx)
	kctx.FatalIfErrorf(err)

	logLatestVersion(1)
}
