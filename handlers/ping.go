package handlers

import (
	"net"

	"github.com/ghmeier/gotella/models"
	"github.com/ghmeier/gotella/receiver"
)

type pingHandler struct {
	*handler
}

func HandlePing(ctx *Context) receiver.ReceiverFunc {
	return &pingHandler{
		handler: newHandler(ctx.Peer, ctx.Descriptor),
	}.ping
}

func (h *pingHandler) ping(conn *net.TCPConn, d *models.Descriptor) {
	addr := conn.LocalAddr()
	tcpAddr, _ := net.ResolveTCPAddr(addr.Network(), addr.String())

	descriptor := models.PongDescriptor(
		tcpAddr.IP.String(),
		tcpAddr.Port,
		&models.Pong{},
		d.Header)
	send(conn, descriptor)

	e, err := h.descriptor.Exists(d.Header.ID)
	if err != nil {
		printError(err)
		return
	}
	if e {
		return
	}

	err = h.descriptor.Put(d.Header.ID, d.IP, d.Port)
	if err != nil {
		printError(err)
		return
	}

	list, err := h.peer.List()
	if err != nil {
		printError(err)
		return
	}

}
