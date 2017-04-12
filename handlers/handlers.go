package handlers

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	"github.com/ghmeier/gotella/helpers"
	"github.com/ghmeier/gotella/models"
	"github.com/ghmeier/gotella/receiver"
)

type Context struct {
	Peer       helpers.Peer
	Descriptor helpers.Descriptor
	Files      helpers.FileSearch
}

type handler struct {
	peer       helpers.Peer
	descriptor helpers.Descriptor
	files      helpers.FileSearch
}

func newHandler(p helpers.Peer, d helpers.Descriptor, f helpers.FileSearch) *handler {
	return &handler{peer: p, descriptor: d, files: f}
}

func (h *handler) propogate(d *models.Descriptor) {
	if d.Header.TTL == 0 {
		// do not propogate if no TTL
		return
	}
	next := d.Next()
	peers, err := h.peer.List()
	if err != nil {
		fmt.Println("Error: unable to get peers")
	}

	for _, peer := range peers {
		connectAndSend(peer, next)
	}
}

func HandleInvalid(ctx *Context) receiver.ReceiverFunc {
	return invalid
}
func invalid(conn *net.TCPConn, d *models.Descriptor) {}

func connectAndSend(destAddr string, d *models.Descriptor) {
	addr, _ := net.ResolveTCPAddr("tcp", destAddr)
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		printError(fmt.Errorf("unable to dial %s\n%s", destAddr, err.Error()))
		return
	}
	defer conn.Close()
	send(conn, d)
}

func send(conn *net.TCPConn, d *models.Descriptor) {
	buf, err := json.Marshal(d)
	if err != nil {
		printError(err)
		return
	}

	_, err = conn.Write(buf)
	if err != nil {
		printError(err)
		return
	}
}

func printError(err error) {
	fmt.Printf("ERROR: %s\n", err.Error())
}

func localAddr(conn *net.TCPConn) (string, int) {
	host, lPort, _ := net.SplitHostPort(conn.LocalAddr().String())
	port, _ := strconv.Atoi(lPort)
	return host, port
}
