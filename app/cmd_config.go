package app

import (
	"fmt"

	"github.com/labctl/labctl/helpers"
	"github.com/srl-labs/containerlab/clab/config"
)

const unexpected = "unexpected command: %v"

type CmdConfig struct {
	Topo          string   `short:"t" help:"topology file" type:"existingfile"`
	TemplatePaths []string `short:"p" help:"paths to search for templates" type:"existingdir"`
	TemplateList  []string `short:"l" help:"template names to render"`
	Filter        []string `short:"f" help:"nodes to include" optional:""`

	Commit struct {
	} `cmd:"" help:"Commit configuration on the node."`
	Compare struct {
	} `cmd:"" help:"Compare rendered templates to the node's config."`
	Send struct {
	} `cmd:"" help:"Send commands to the node (i.e. show)"`
	Template struct {
	} `cmd:"" help:"Show rendered templates."`
	Vars struct {
	} `cmd:"" help:"Show variable inputs."`
}

func (r *CmdConfig) Run(ctx *helpers.Context) error {
	ctx.Command = ctx.Command[7:]
	ctx.NodeFilter = r.Filter
	ctx.TopoFile = r.Topo

	config.TemplateNames = r.TemplateList
	config.TemplatePaths = r.TemplatePaths

	switch ctx.Command {
	case "send", "commit", "compare":
		return ConfigRun(ctx.Command, ctx)

	case "template", "vars":
		return ConfigView(ctx.Command, ctx)

	default:
		return fmt.Errorf("unknown command %s", ctx.Command)
	}
}
