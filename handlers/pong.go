package handlers

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/ghmeier/gotella/models"
	"github.com/ghmeier/gotella/receiver"
)

type pongHandler struct {
	*handler
}

func HandlePong(ctx *Context) receiver.ReceiverFunc {
	p := &pongHandler{
		handler: newHandler(ctx.Peer, ctx.Descriptor, ctx.Files),
	}
	return p.pong
}

func (h *pongHandler) pong(conn *net.TCPConn, d *models.Descriptor) {
	fmt.Printf("%s: pong from %s:%d\n", d.Header.ID.String(), d.IP, d.Port)
	cache, err := h.descriptor.Get(d.Header.ID)
	if err != nil {
		printError(err)
		return
	}
	if cache == nil || cache.Type == d.Header.Type {
		// pongs must have corresponding ping cached
		return
	}

	var pong models.Pong
	err = json.Unmarshal(d.Payload, &pong)
	if err != nil {
		printError(err)
		return
	}
	peer := &models.Peer{
		IP:    d.IP,
		Port:  d.Port,
		Files: pong.Files,
		Size:  pong.Size,
	}
	err = h.peer.Put(peer)
	if err != nil {
		printError(err)
		return
	}

	next := d.Next()
	if next.Header.TTL < 0 {
		return
	}
	fmt.Printf("%s: passing pong to %s:%d\n", d.Header.ID.String(), cache.IP, cache.Port)
	connectAndSend(fmt.Sprintf("%s:%d", cache.IP, cache.Port), next)
}
