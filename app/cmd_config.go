package app

import (
	"fmt"
	"strings"

	"github.com/google/shlex"
	"github.com/gorilla/websocket"
	"github.com/labctl/labctl/helpers"
	log "github.com/sirupsen/logrus"
	"github.com/srl-labs/containerlab/clab/config"
)

type CmdConfig struct {
	Topo          string   `short:"t" help:"topology file" type:"existingfile"`
	TemplatePaths []string `short:"p" help:"paths to search for templates" type:"existingdir"`
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

	if strings.HasPrefix(ctx.Command, "config ") { // Executed from the commandline
		ctx.Command = strings.TrimPrefix(ctx.Command, "config ")
		ctx.TopoFile = r.Topo
		// Ensure paths are valid valid
		ctx.TemplatePaths, err = helpers.InitTemplatePaths(r.TemplatePaths)
		if err != nil {
			return err
		}

	} else { // Executed form a websocket
		log.Debugf("Websocket %s %v", ctx.Command, ctx)
		if r.Topo != "" {
			return fmt.Errorf("--topo/-t not allowed. Fixed at %s", Ctx.TopoFile)
		}
		if len(r.TemplatePaths) > 0 {
			return fmt.Errorf(
				"--template-paths/-p not allowed. Fixed at %s",
				strings.Join(Ctx.TemplatePathsSlice(), ", "),
			)
		}
	}
	ctx.Output.Info("TST", "aa")

	ctx.NodeFilter = r.Filter

	// Setup containerlab's config engine
	config.TemplateNames = r.TemplateList
	config.TemplatePaths = ctx.TemplatePathsSlice()

	switch ctx.Command {
	case "send", "commit", "compare":
		return ConfigRun(ctx.Command, ctx)

	case "template", "vars":
		return ConfigView(ctx.Command, ctx)

	default:
		return fmt.Errorf("unknown command %s", ctx.Command)
	}
}

func ParseWebString(wsconn *websocket.Conn, cmd string) error {
	args, err := shlex.Split(cmd)
	if err != nil {
		return err
	}

	// Parse using only the
	var cli CmdConfig
	p := GetCmdLineParser(&cli)
	kctx, err := p.Parse(args)
	if err != nil {
		return err
	}

	// Copy the global context, update the output & command
	nctx := Ctx
	nctx.Output = &helpers.WebSocketOutput{
		Conn: wsconn,
	}
	nctx.Command = kctx.Command()

	return kctx.Run(nctx)
}
