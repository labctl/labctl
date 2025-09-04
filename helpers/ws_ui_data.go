package helpers

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/chigopher/pathlib"
	"github.com/gorilla/websocket"
	"gopkg.in/yaml.v3"
)

type layoutsNodes map[string]struct {
	X float32 `json:"x" yaml:"x,flow"`
	Y float32 `json:"y" yaml:"y,flow"`
}
type layouts struct {
	Nodes layoutsNodes `json:"nodes"`
}

type WsUiData struct {
	Options   map[string]any    `json:"options"`   // A generic options dictionary for persistent UI options
	Layouts   layouts           `json:"layouts"`   // layouts is a property of v-network-graph
	Templates map[string]string `json:"templates"` // Templates used by the UI etc
	Context   *ContextJson      `json:"context"`   // Context with filename, template paths etc
}

func NewWsUiData() *WsUiData {
	return &WsUiData{
		Layouts: layouts{
			Nodes: make(layoutsNodes, 0),
		},
		Options:   make(map[string]any, 0),
		Templates: make(map[string]string),
	}
}

func (u *WsUiData) Print() {
	log.Printf("%v", u)
	// for k, v := range u.Layouts.Nodes {
	// 	log.Infof("%s x:%v y:%v", k, v.X, v.Y)
	// }
}

func (u *WsUiData) WriteFile(ctx *Context) {
	labfn := ctx.LabctlFilename
	// Use the clab filemode
	s, err := os.Lstat(ctx.TopoFilename)
	fm := s.Mode()
	if err != nil {
		log.Errorf("Unable to read %s: %s", ctx.TopoFilename, err.Error())
		fm = 0o664
	}
	log.Infof("Saving as %s", labfn)

	// Migrate commands to MarkDown files
	delete(u.Options, "commands")

	data, err := yaml.Marshal(&u)
	if err != nil {
		log.Error(err)
		return
	}
	err = os.WriteFile(labfn, data, fm)
	if err != nil {
		log.Error(err)
	}
}

func (u *WsUiData) ReadFile(ctx *Context) error {
	data, err := os.ReadFile(ctx.LabctlFilename)
	if os.IsNotExist(err) {
		log.Debugf("No labctl file yet: %s", err)
		return nil
	}
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, u)
	if err != nil {
		return err
	}

	// Migrate commands to MarkDown files
	cmds, hascmd := u.Options["commands"]
	if hascmd {
		// migrate to old_command.md
		ocp := pathlib.NewPath(ctx.LabctlFilename).Parent().Join("old_command.md")
		cmdss := fmt.Sprintf("%v", cmds)
		err := ocp.WriteFile([]byte(cmdss))
		if err != nil {
			log.Warn("Could not write old_commend.md", "err", err, "value", cmdss)
		}
		delete(u.Options, "commands")
	}
	return nil
}

func WsUiUpdate(ctx *Context) *WsMessage {
	wsmsg := &WsMessage{
		Code:   WscUiData,
		UiData: NewWsUiData(),
	}
	err := wsmsg.UiData.ReadFile(ctx)
	if err != nil {
		return &WsMessage{Code: WscError, Msg: "could not read .labctl.yml"}
	}
	return wsmsg
}

func WsMakeChannel(conn *websocket.Conn) chan<- any {
	c := make(chan any)

	go func() {
		for {
			wsmsg, ok := <-c
			if !ok {
				return
			}
			if wsmsg == nil { // pong
				_ = conn.WriteMessage(websocket.TextMessage, []byte("%"))
				continue
			}
			if DebugCount > 0 {
				log.Debugf("Ws Tx: %s", wsmsg)
			}
			err := conn.WriteJSON(wsmsg)
			if err != nil {
				log.Errorf("could not write to the websocket: %s", wsmsg)
				return
			}
		}
	}()

	return c
}
