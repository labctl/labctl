package app

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/labctl/labctl/helpers"
	"github.com/posener/complete"
	"github.com/willabides/kongplete"

	log "github.com/sirupsen/logrus"
)

var cli struct {
	Debug    int    `help:"Enable debug mode." short:"d" type:"counter" group:"Global Flags"`
	Settings string `help:"Settings yaml file." type:"existingfile" group:"Global Flags"`

	Color  CmdColor  `cmd:"" help:"Add some color to any Linux command (i.e. ssh)." short:"c" passthrough:""`
	Config CmdConfig `cmd:"" help:"Render a template."`
	Serve  CmdServe  `cmd:"" help:"Serve the web UI."`

	InstallCompletions kongplete.InstallCompletions `cmd:"" help:"install shell completions"`
}

func Main() {
	parser := kong.Must(&cli,
		kong.Name("labctl"),
		kong.Description("Control your network lab."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: false,
		}))

	// Handle completion requests
	kongplete.Complete(parser,
		kongplete.WithPredictor("file", complete.PredictFiles("*")),
	)

	ctx, err := parser.Parse(os.Args[1:])
	parser.FatalIfErrorf(err)

	s := helpers.NewSettings(true)
	// if err = s.AddSettings("~/.labctl.yml", true); err != nil {
	// 	ctx.FatalIfErrorf(err)
	// }

	log.SetOutput(os.Stdout)
	if cli.Debug > 0 {
		log.SetLevel(log.DebugLevel)
	}

	// Call the Run() method of the selected parsed command.
	err = ctx.Run(&helpers.Context{
		Command:    ctx.Command(),
		DebugCount: cli.Debug,
		Settings:   s,
	})
	ctx.FatalIfErrorf(err)
}
