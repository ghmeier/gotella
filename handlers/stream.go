package handlers

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/ghmeier/gotella/models"
	"github.com/ghmeier/gotella/receiver"
)

type streamHandler struct {
	*handler
}

func HandleStream(ctx *Context) receiver.ReceiverFunc {
	p := &streamHandler{
		handler: newHandler(ctx.Peer, ctx.Descriptor, ctx.Files),
	}
	return p.stream
}

func (h *streamHandler) stream(conn *net.TCPConn, d *models.Descriptor) {
	fmt.Printf("%s: stream request from %s:%d\n", d.Header.ID.String(), d.IP, d.Port)
	var stream models.Stream
	err := json.Unmarshal(d.Payload, &stream)
	if err != nil {
		printError(err)
		conn.Write([]byte("ERROR"))
		return
	}

	if !h.files.Exists(stream.FileName) {
		conn.Write([]byte("ERROR"))
		return
	}

	buf, err := h.files.Get(stream.FileName)
	_, err = conn.Write(buf)
	if err != nil {
		printError(err)
	}
}
