package helpers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ghmeier/gotella/gateways"
	"github.com/ghmeier/gotella/models"
)

const MAX_PEERS = 2

type Peer interface {
	Put(*models.Peer) error
	Get(string, int) (*models.Peer, error)
	List() ([]string, error)
}

type peerHelper struct {
	redis   gateways.Redis
	listKey string
	listMax int
	write   chan *models.Peer
}

func NewPeer(redis gateways.Redis) Peer {
	p := &peerHelper{
		redis:   redis,
		listKey: "peer_addrs",
		listMax: MAX_PEERS,
		write:   make(chan *models.Peer),
	}

	go p.handle()
	return p
}

func (h *peerHelper) handle() {
	for p := range h.write {
		length, err := h.redis.Len(h.listKey)
		if err != nil {
			peerError(err)
			continue
		}
		if length >= h.listMax {
			continue
		}

		addr := addr(p.IP, p.Port)
		exists, err := h.redis.Exists(addr)
		if err != nil {
			peerError(err)
			return
		}
		if !exists {
			err = h.redis.Append(h.listKey, addr)
			if err != nil {
				peerError(err)
				return
			}
		}

		buf, _ := json.Marshal(p)
		h.redis.Set(addr, buf, time.Second*15)
	}
}

func (h *peerHelper) Put(p *models.Peer) error {
	h.write <- p
	return nil
}

func (h *peerHelper) Get(ip string, port int) (*models.Peer, error) {
	addr := addr(ip, port)
	buf, err := h.redis.Get(addr)
	if err != nil {
		if len(buf) == 0 {
			h.redis.Remove(h.listKey, addr)
			return nil, nil
		}
		return nil, err
	}

	var p models.Peer
	err = json.Unmarshal(buf, &p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (h *peerHelper) List() ([]string, error) {
	list, err := h.redis.List(h.listKey)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func peerError(err error) {
	fmt.Printf("PEER ERROR: %s\n", err.Error())
}

func addr(ip string, port int) string {
	return fmt.Sprintf("%s:%d", ip, port)
}
