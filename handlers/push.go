package handlers

import (
	"net"

	"github.com/ghmeier/gotella/models"
	"github.com/ghmeier/gotella/receiver"
)

func HandlePush(ctx *Context) receiver.ReceiverFunc {
	return newHandler(ctx.Peer, ctx.Descriptor).push
}

func (h *handler) push(conn *net.TCPConn, d *models.Descriptor) {

}
