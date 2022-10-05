package webpty

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/hairyhenderson/gomplate/v3/base64"
	"github.com/labctl/labctl/utils"
	"github.com/pkg/errors"
	"github.com/sorenisanerd/gotty/backend/localcommand"
	"github.com/sorenisanerd/gotty/webtty"

	log "github.com/sirupsen/logrus"
)

type InitMessage struct {
	Cmd string `json:"cmd"`
}

var wsupgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func Websock(w http.ResponseWriter, r *http.Request) {
	wsconn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer wsconn.Close()

	err = handle(context.TODO(), wsconn)
	if err != nil {
		log.Error(err.Error())
	}
}

func handle(ctx context.Context, wsconn *websocket.Conn) error {
	typ, initMsg, err := wsconn.ReadMessage()
	if err != nil {
		return errors.Wrapf(err, "failed to establish websocket pty connection")
	}
	if typ != websocket.TextMessage {
		return errors.New("failed to establish websocket pty connection: invalid message type")
	}
	log.Infof("rx: %s", initMsg)

	var init InitMessage
	err = json.Unmarshal(initMsg, &init)
	if err != nil {
		return errors.Wrapf(err, "failed to establish websocket pty connection")
	}

	cmd, args, err := allow(init.Cmd)
	if err != nil {
		// wsconn.WriteMessage(websocket.TextMessage, []byte(string(webtty.UnknownOutput)+"Command not allowed"))
		b64 := utils.Must(base64.Encode([]byte(err.Error())))
		_ = wsconn.WriteMessage(websocket.TextMessage, []byte(string(webtty.Output)+b64))
		return err
	}

	opts := []webtty.Option{
		// webtty.WithWindowTitle(titleBuf.Bytes()),
		// webtty.WithReconnect(server.options.ReconnectTime)
		// webtty.WithFixedColumns(80),
		// webtty.WithFixedRows(24),
		webtty.WithPermitWrite(),
	}

	var slave webtty.Slave
	slave, err = localcommand.New(cmd, args)

	if err != nil {
		return errors.Wrapf(err, "failed to create backend")
	}
	// defer slave.Close()

	tty, err := webtty.New(&wsWrapper{wsconn}, slave, opts...)
	if err != nil {
		return errors.Wrapf(err, "failed to create webtty")
	}

	return tty.Run(ctx)
}
