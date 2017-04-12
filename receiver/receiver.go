package receiver

import (
	"fmt"
	"net"
	"time"

	"github.com/ghmeier/gotella/models"
)

type ReceiverFunc func(*net.TCPConn, *models.Descriptor)

type Probe interface {
	Send(net.Addr)
}

type Receiver struct {
	addr     *net.TCPAddr
	listener *net.TCPListener
	port     string
	routes   map[models.DescriptorType]ReceiverFunc
	probe    Probe
}

func New(port string) *Receiver {
	return &Receiver{
		port:   port,
		routes: make(map[models.DescriptorType]ReceiverFunc),
	}
}

func (r *Receiver) Start() error {
	ip, err := externalIP()
	if err != nil {
		fmt.Println(err)
	}

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%s", ip, r.port))
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

func (r *Receiver) Probe() {
	for {
		r.probe.Send(r.listener.Addr())
		time.Sleep(time.Second * 5)
	}
}

func (r *Receiver) Register(route models.DescriptorType, f ReceiverFunc) {
	r.routes[route] = f
}

func (r *Receiver) RegisterProbe(p Probe) {
	r.probe = p
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
	buf := make([]byte, 2048)
	n, err := conn.Read(buf)
	if err != nil {
		printErr(err)
		return
	}

	descriptor, err := models.FromBuf(buf[:n])
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

func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", fmt.Errorf("are you connected to the network?")
}
