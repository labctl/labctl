package app

import (
	"fmt"
	"strings"

	"github.com/google/shlex"
	"github.com/labctl/labctl/helpers"
	"github.com/labctl/labctl/utils"
	log "github.com/sirupsen/logrus"
	"github.com/srl-labs/containerlab/clab/config"
	"golang.org/x/exp/slices"
)

type CmdConfig struct {
	Topo          string   `short:"t" help:"topology file" type:"existingfile" predictor:"topo"`
	TemplatePaths []string `short:"p" help:"paths to search for templates" type:"existingdir" predictor:"dir"`
	TemplateList  []string `short:"l" help:"template names to render"`
	Filter        []string `short:"f" help:"nodes to include" optional:""`

	Commit   struct{} `cmd:"" help:"Commit configuration on the node."`
	Compare  struct{} `cmd:"" help:"Compare rendered templates to the node's config."`
	Send     struct{} `cmd:"" help:"Send commands to the node (i.e. show)"`
	Template struct{} `cmd:"" help:"Show rendered templates."`
	Vars     struct{} `cmd:"" help:"Show variable inputs."`
}

func (r *CmdConfig) Run(ctx *helpers.Context) error {
	var err error

	r.TemplateList = utils.Unique(r.TemplateList)

	if strings.HasPrefix(ctx.Command, "config ") {
		// Executed from the commandline
		ctx.Command = strings.TrimPrefix(ctx.Command, "config ")
		var err error
		r.Topo, err = ctx.InitPaths(r.Topo, r.TemplatePaths)
		if err != nil {
			return err
		}
	} else {
		// Executed form a websocket
		log.Debugf("Websocket %s %v", ctx.Command, ctx)
		if r.Topo != "" {
			return fmt.Errorf("--topo/-t not allowed. Fixed at %s", Ctx.TopoFilename)
		}
		if len(r.TemplatePaths) > 0 {
			return fmt.Errorf(
				"--template-paths/-p not allowed. Fixed at %s",
				strings.Join(Ctx.TemplatePathsSlice(), ", "),
			)
		}
	}
	if len(r.TemplateList) == 0 {
		return fmt.Errorf("you need to specify at least one template (-l/--template-list)")
	}

	ctx.NodeFilter = r.Filter

	// Setup containerlab's config engine
	config.TemplateNames = r.TemplateList
	config.TemplatePaths = ctx.TemplatePathsSlice()
	log.Debugf("Search path: %v", config.TemplatePaths)

	configs, err := LoadAndPrep(&ctx.NodeFilter, ctx.TopoFilename, ctx.Command != "vars")
	if err != nil {
		return err
	}

	switch ctx.Command {
	case "send", "commit", "compare":
		return ConfigTx(configs, ctx)
	case "vars":
		for nn, cfg := range configs {
			if !slices.Contains(ctx.NodeFilter, nn) {
				continue
			}
			if ctx.Output.PreferStdout() {
				cfg.Print(true, false)
			} else {
				ctx.Output.LogResponses(nil, cfg)
			}
		}
	case "template":
		for nn, cfg := range configs {
			if !slices.Contains(ctx.NodeFilter, nn) {
				continue
			}
			if ctx.Output.PreferStdout() {
				cfg.Print(false, true)
			} else {
				ctx.Output.LogResponses(nil, cfg)
			}
		}
	default:
		return fmt.Errorf("unknown command %s", ctx.Command)
	}
	return nil
}

func WebConfigDone(ws chan<- interface{}, retry int) {
	ws <- &helpers.WsMessage{
		Code:   helpers.WscConfig,
		Config: &helpers.WsConfig{Cmd: "done"},
	}
}

func RunWebConfig(ws chan<- interface{}, cmd string) {
	defer WebConfigDone(ws, 3)

	args, err := shlex.Split(cmd)
	if err != nil {
		helpers.WsErrorf(ws, "%s\n\nwhile trying to split:\n\t%s", err, cmd)
		return
	}

	// Parse using config
	var cliConfig CmdConfig
	p := GetCmdLineParser(&cliConfig)
	kctx, err := p.Parse(args)
	if err != nil {
		helpers.WsErrorf(ws, "%s\n\nwhile trying to parse:\n\t%s", err, cmd)
		return
	}

	// Copy the global context, update the output & command
	nctx := Ctx
	nctx.Output = &helpers.WebSocketOutput{
		Ws: ws,
	}
	nctx.Async = true
	nctx.Command = kctx.Command()

	err = kctx.Run(nctx)
	if err != nil {
		helpers.WsErrorf(ws, err.Error())
		return
	}
}
