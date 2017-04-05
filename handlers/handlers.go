package handlers

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/ghmeier/gotella/helpers"
	"github.com/ghmeier/gotella/models"
	"github.com/ghmeier/gotella/receiver"
)

type Context struct {
	Peer       helpers.Peer
	Descriptor helpers.Descriptor
}

type handler struct {
	peer       helpers.Peer
	descriptor helpers.Descriptor
}

func newHandler(p helpers.Peer, d helpers.Descriptor) *handler {
	return &handler{peer: p, descriptor: d}
}

func HandleInvalid(ctx *Context) receiver.ReceiverFunc {
	return invalid
}
func invalid(conn *net.TCPConn, d *models.Descriptor) {}

func send(conn *net.TCPConn, d *models.Descriptor) {
	buf, _ := json.Marshal(d)

	_, err := conn.Write(buf)
	if err != nil {
		printError(err)
	}
}

func printError(err error) {
	fmt.Printf("ERROR: %s\n", err.Error())
}
