package handlers

import (
	"fmt"
	"net"
	"strconv"

	"github.com/ghmeier/gotella/models"
	"github.com/ghmeier/gotella/receiver"
)

type probeHandler struct {
	*handler
	discovery string
}

func NewProbe(ctx *Context, addr string) receiver.Probe {
	return &probeHandler{
		handler:   newHandler(ctx.Peer, ctx.Descriptor, ctx.Files),
		discovery: addr,
	}
}

func (p *probeHandler) Send(laddr net.Addr) {
	fmt.Printf("Send from: %s\n", laddr.String())
	lIP, port, err := net.SplitHostPort(laddr.String())
	if err != nil {
		printError(err)
	}
	lPort, _ := strconv.Atoi(port)
	d := models.PingDescriptor(lIP, lPort)

	peers, err := p.peer.List()
	if err != nil {
		printError(err)
		return
	}

	err = p.descriptor.Put(d)
	if err != nil {
		// only put descriptors if we can cache them
		printError(err)
		return
	}
	for _, peer := range peers {
		fmt.Printf("Pinging %s\n", peer)
		connectAndSend(peer, d)
	}

	if len(peers) <= 0 && p.discovery != "" {
		fmt.Printf("Pinging discovery host %s\n", p.discovery)
		connectAndSend(p.discovery, d)
	}
}
