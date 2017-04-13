package handlers

import (
	"fmt"
	"net"
	"strconv"
	"strings"

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
	lIP, lPort := fromAddr(laddr)
	d := models.PingDescriptor(lIP, lPort)

	p.probePeers(d)
}

func (p *probeHandler) Query(laddr net.Addr, search string) {
	lIP, lPort := fromAddr(laddr)

	q := &models.Query{
		Filename:    strings.Trim(search, "\n\t "),
		RequestIP:   lIP,
		RequestPort: lPort,
	}
	d := models.QueryDescriptor(lIP, lPort, q)

	p.probePeers(d)
}

func (p *probeHandler) probePeers(d *models.Descriptor) {
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

		// attempt to get peer to clean it up
		addr, _ := net.ResolveTCPAddr("tcp", peer)
		p.peer.Get(addr.IP.String(), addr.Port)
	}

	if len(peers) <= 0 && p.discovery != "" {
		fmt.Printf("Pinging discovery host %s\n", p.discovery)
		connectAndSend(p.discovery, d)
	}
}

func fromAddr(laddr net.Addr) (string, int) {
	lIP, port, err := net.SplitHostPort(laddr.String())
	if err != nil {
		printError(err)
	}
	lPort, _ := strconv.Atoi(port)
	return lIP, lPort
}
