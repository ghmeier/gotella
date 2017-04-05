package handlers

import (
	"net"

	"github.com/ghmeier/gotella/models"
	"github.com/ghmeier/gotella/receiver"
)

func HandleQuery(ctx *Context) receiver.ReceiverFunc {
	return newHandler(ctx.Peer, ctx.Descriptor).query
}

func (h *handler) query(conn *net.TCPConn, d *models.Descriptor) {

}
