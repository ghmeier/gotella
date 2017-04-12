package handlers

import (
	"fmt"
	"net"

	"github.com/ghmeier/gotella/models"
	"github.com/ghmeier/gotella/receiver"
)

type pingHandler struct {
	*handler
}

func HandlePing(ctx *Context) receiver.ReceiverFunc {
	p := &pingHandler{
		handler: newHandler(ctx.Peer, ctx.Descriptor, ctx.Files),
	}
	return p.ping
}

func (h *pingHandler) ping(conn *net.TCPConn, d *models.Descriptor) {
	fmt.Printf("%s: ping from %s:%d\n", d.Header.ID.String(), d.IP, d.Port)
	e, err := h.descriptor.Exists(d.Header.ID)
	if err != nil {
		printError(err)
		return
	}
	if e {
		return
	}

	host, port := localAddr(conn)
	descriptor := models.PongDescriptor(
		host,
		port,
		&models.Pong{
			Files: h.files.Count(),
			Size:  h.files.Size(),
		},
		d)
	fmt.Printf("%s: ponging %s:%d\n", descriptor.Header.ID.String(), d.IP, d.Port)
	connectAndSend(fmt.Sprintf("%s:%d", d.IP, d.Port), descriptor)

	err = h.peer.Put(&models.Peer{
		IP:    d.IP,
		Port:  d.Port,
		Files: 0,
		Size:  0,
	})

	err = h.descriptor.Put(d)
	if err != nil {
		printError(err)
		return
	}

	h.propogate(d)
}
