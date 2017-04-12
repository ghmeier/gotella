package handlers

import (
	"net"

	"github.com/ghmeier/gotella/models"
	"github.com/ghmeier/gotella/receiver"
)

type queryHitHandler struct {
	*handler
}

func HandleQueryHit(ctx *Context) receiver.ReceiverFunc {
	q := &queryHitHandler{
		handler: newHandler(ctx.Peer, ctx.Descriptor, ctx.Files),
	}
	return q.queryHit
}

func (h *queryHitHandler) queryHit(conn *net.TCPConn, d *models.Descriptor) {

}
