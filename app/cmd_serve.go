package app

import "github.com/gorilla/websocket"

type CmdServe struct {
	Topo          string   `short:"t" help:"topology file" type:"existingfile"`
	TemplatePaths []string `short:"p" help:"paths to search for templates" type:"existingdir"`
}

var upgrader = websocket.Upgrader{} // use default options

func (r *CmdServe) Run(ctx *Context) error {
	//websocket.
	return nil
}
