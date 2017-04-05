package helpers

import (
	"encoding/json"
	"fmt"

	"github.com/ghmeier/gotella/gateways"
	"github.com/ghmeier/gotella/models"
)

type Peer interface {
	Put(*models.Peer) error
	Get(string, int) (*models.Peer, error)
	List() ([]string, error)
}

type peerHelper struct {
	redis   gateways.Redis
	listKey string
	listMax int
}

func NewPeer(redis gateways.Redis) Peer {
	return &peerHelper{
		redis:   redis,
		listKey: "peer_addrs",
		listMax: 7,
	}
}

func (h *peerHelper) Put(p *models.Peer) error {
	length, err := h.redis.Len(h.listKey)
	if err != nil {
		return err
	}
	if length > h.listMax {
		return nil
	}

	addr := addr(p.IP, p.Port)
	exists, err := h.redis.Exists(addr)
	if err != nil {
		return err
	}

	if !exists {
		err = h.redis.Append(h.listKey, addr)
		if err != nil {
			return err
		}
	}

	buf, _ := json.Marshal(p)
	err = h.redis.Set(addr, buf, 0)
	return err
}

func (h *peerHelper) Get(ip string, port int) (*models.Peer, error) {
	addr := addr(ip, port)
	buf, err := h.redis.Get(addr)
	if err != nil {
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

func addr(ip string, port int) string {
	return fmt.Sprintf("%s:%d", ip, port)
}
