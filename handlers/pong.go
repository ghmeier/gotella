package handlers

import (
	"net"

	"github.com/ghmeier/gotella/models"
	"github.com/ghmeier/gotella/receiver"
)

func HandlePong(ctx *Context) receiver.ReceiverFunc {
	return newHandler(ctx.Peer, ctx.Descriptor).pong
}

func (h *handler) pong(conn *net.TCPConn, d *models.Descriptor) {

}
