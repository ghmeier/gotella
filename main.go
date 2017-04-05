package main

import (
	"fmt"
	"os"

	"github.com/ghmeier/gotella/gateways"
	"github.com/ghmeier/gotella/handlers"
	"github.com/ghmeier/gotella/helpers"
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

func redis() string {
	args := os.Args[1:]

	redis := "8080"
	if len(args) > 1 && args[1] != "" {
		redis = args[1]
	}

	return redis
}

func register(r *receiver.Receiver) {
	redis := gateways.New(redis())
	ctx := &handlers.Context{
		Peer:       helpers.NewPeer(redis),
		Descriptor: helpers.NewDescriptor(redis),
	}

	r.Register(models.PING, handlers.HandlePing(ctx))
	r.Register(models.PONG, handlers.HandlePong(ctx))
	r.Register(models.QUERY, handlers.HandleQuery(ctx))
	r.Register(models.QUERYHIT, handlers.HandleQueryHit(ctx))
	r.Register(models.PUSH, handlers.HandlePush(ctx))
	r.Register(models.INVALID, handlers.HandleInvalid(ctx))
}
