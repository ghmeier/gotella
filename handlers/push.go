package handlers

import (
	"net"

	"github.com/ghmeier/gotella/models"
	"github.com/ghmeier/gotella/receiver"
)

type pushHandler struct {
	*handler
}

func HandlePush(ctx *Context) receiver.ReceiverFunc {
	p := &pushHandler{
		handler: newHandler(ctx.Peer, ctx.Descriptor, ctx.Files),
	}
	return p.push
}

func (h *pushHandler) push(conn *net.TCPConn, d *models.Descriptor) {

}
