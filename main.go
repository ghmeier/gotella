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
	register(r, peer())

	err := r.Start()
	if err != nil {
		fmt.Println("ERROR creating reciever: %s\n", err.Error())
	}

	r.Probe()
}

func peer() string {
	args := os.Args[1:]

	addr := ""
	if len(args) > 2 {
		addr = args[2]
	}
	return addr
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

	redis := "127.0.0.1:8080"
	if len(args) > 1 && args[1] != "" {
		redis = "127.0.0.1:" + args[1]
	}

	return redis
}

func register(r *receiver.Receiver, discovery string) {
	redis := gateways.New(redis(), port())
	ctx := &handlers.Context{
		Peer:       helpers.NewPeer(redis),
		Descriptor: helpers.NewDescriptor(redis),
		Files:      helpers.NewFiles("public"),
	}

	r.Register(models.PING, handlers.HandlePing(ctx))
	r.Register(models.PONG, handlers.HandlePong(ctx))
	r.Register(models.QUERY, handlers.HandleQuery(ctx))
	r.Register(models.QUERYHIT, handlers.HandleQueryHit(ctx))
	r.Register(models.PUSH, handlers.HandlePush(ctx))
	r.Register(models.INVALID, handlers.HandleInvalid(ctx))
	r.RegisterProbe(handlers.NewProbe(ctx, discovery))
}
