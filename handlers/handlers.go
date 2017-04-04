package handlers

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/ghmeier/gotella/models"
	"github.com/ghmeier/gotella/receiver"
)

func HandlePing( /*TODO: context*/ ) receiver.ReceiverFunc { return ping }
func ping(conn *net.TCPConn, d *models.Descriptor) {
	addr := conn.LocalAddr()
	tcpAddr, _ := net.ResolveTCPAddr(addr.Network(), addr.String())

	descriptor := models.PongDescriptor(
		tcpAddr.IP.String(),
		tcpAddr.Port,
		&models.Pong{},
		d.Header)
	send(conn, descriptor)
}

func HandlePong( /*TODO: context*/ ) receiver.ReceiverFunc { return pong }
func pong(conn *net.TCPConn, d *models.Descriptor) {

}

func HandleQuery( /*TODO: context*/ ) receiver.ReceiverFunc { return query }
func query(conn *net.TCPConn, d *models.Descriptor) {

}

func HandleQueryHit( /*TODO: context*/ ) receiver.ReceiverFunc { return queryHit }
func queryHit(conn *net.TCPConn, d *models.Descriptor) {

}

func HandlePush( /*TODO: context*/ ) receiver.ReceiverFunc { return push }
func push(conn *net.TCPConn, d *models.Descriptor) {

}

func HandleInvalid( /*TODO: context*/ ) receiver.ReceiverFunc { return invalid }
func invalid(conn *net.TCPConn, d *models.Descriptor)         {}

func send(conn *net.TCPConn, d *models.Descriptor) {
	buf, _ := json.Marshal(d)

	_, err := conn.Write(buf)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
	}
}
