package main

import (
	"bufio"
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

	go r.Probe()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("ENTER A FILENAME TO SEARCH AND PRESS ENTER:")
		text, _ := reader.ReadString('\n')
		r.Query(text)
	}
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

	if len(args) > 0 && args[0] != "" {
		return args[0]
	}

	fmt.Println("ERROR: Enter server port as first argument")
	return ""
}

func redis() string {
	args := os.Args[1:]

	if len(args) > 1 && args[1] != "" {
		return "127.0.0.1:" + args[1]
	}

	fmt.Println("ERROR: Enter redis host:port as second argument")
	return ""
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
	r.Register(models.STREAM, handlers.HandleStream(ctx))
	r.Register(models.INVALID, handlers.HandleInvalid(ctx))
	r.RegisterProbe(handlers.NewProbe(ctx, discovery))
}
