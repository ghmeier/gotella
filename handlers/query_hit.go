package handlers

import (
	"net"

	"github.com/ghmeier/gotella/models"
	"github.com/ghmeier/gotella/receiver"
)

func HandleQueryHit(ctx *Context) receiver.ReceiverFunc {
	return newHandler(ctx.Peer, ctx.Descriptor).queryHit
}

func (h *handler) queryHit(conn *net.TCPConn, d *models.Descriptor) {

}
