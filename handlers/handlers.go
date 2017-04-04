package handlers

import (
	"net"

	"github.com/ghmeier/gotella/models"
	"github.com/ghmeier/gotella/receiver"
)

func HandlePing( /*TODO: context*/ ) receiver.ReceiverFunc { return ping }
func ping(conn *net.TCPConn, d *models.Descriptor) {

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
