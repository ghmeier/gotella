package main

import (
	"fmt"
	"os"

	"github.com/ghmeier/gotella/handlers"
	"github.com/ghmeier/gotella/models"
	"github.com/ghmeier/gotella/receiver"
)

func main() {
	r := receiver.New(port())
	register(r)

	err := r.Start()
	if err != nil {
		fmt.Println("ERROR creating reciever: %s\n", err.Error())
	}
}

func port() string {
	args := os.Args[1:]

	port := "8080"
	if len(args) > 0 && args[0] != "" {
		port = args[0]
	}

	return port
}

func register(r *receiver.Receiver) {
	r.Register(models.PING, handlers.HandlePing())
	r.Register(models.PONG, handlers.HandlePong())
	r.Register(models.QUERY, handlers.HandleQuery())
	r.Register(models.QUERYHIT, handlers.HandleQueryHit())
	r.Register(models.PUSH, handlers.HandlePush())
	r.Register(models.INVALID, handlers.HandleInvalid())
}
