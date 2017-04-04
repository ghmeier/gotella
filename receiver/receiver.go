package receiver

import (
	"fmt"
	"net"

	"github.com/ghmeier/gotella/models"
)

type ReceiverFunc func(*net.TCPConn, *models.Descriptor)

type Receiver struct {
	addr     *net.TCPAddr
	listener *net.TCPListener
	port     string
	routes   map[models.DescriptorType]ReceiverFunc
}

func New(port string) *Receiver {
	return &Receiver{
		port:   port,
		routes: make(map[models.DescriptorType]ReceiverFunc),
	}
}

func (r *Receiver) Start() error {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%s", r.port))
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}

	r.addr = addr
	r.listener = listener
	go r.listen()

	return nil
}

func (r *Receiver) Register(route models.DescriptorType, f ReceiverFunc) {
	r.routes[route] = f
}

func (r *Receiver) listen() {
	defer r.listener.Close()
	for {
		conn, err := r.listener.AcceptTCP()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		go r.handle(conn)
	}
}

func (r *Receiver) handle(conn *net.TCPConn) {
	defer conn.Close()
	buf := make([]byte, 1024)

	_, err := conn.Read(buf)
	if err != nil {
		printErr(err)
		return
	}

	descriptor, err := models.NewDescriptor(buf)
	if err != nil {
		printErr(err)
		return
	}

	if r.routes[descriptor.Header.Type] == nil {
		printErr(fmt.Errorf("No router for %d", descriptor.Header.Type))
		return
	}

	r.routes[descriptor.Header.Type](conn, descriptor)
}

func printErr(err error) {
	fmt.Printf("ERROR: %s\n", err.Error())
}
