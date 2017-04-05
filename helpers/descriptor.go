package helpers

import (
	"time"

	"github.com/pborman/uuid"

	"github.com/ghmeier/gotella/gateways"
)

type Descriptor interface {
	Put(uuid.UUID, string, int) error
	Exists(uuid.UUID) (bool, error)
}

type descriptorHelper struct {
	redis gateways.Redis
}

func NewDescriptor(redis gateways.Redis) Descriptor {
	return &descriptorHelper{redis: redis}
}

func (h *descriptorHelper) Put(id uuid.UUID, ip string, port int) error {
	err := h.redis.Set(id.String(), []byte(addr(ip, port)), time.Hour*1)
	return err
}

func (h *descriptorHelper) Exists(id uuid.UUID) (bool, error) {
	return h.redis.Exists(id.String())
}

func (h *descriptorHelper) Get(id uuid.UUID) (string, error) {
	buf, err := h.redis.Get(id.String())
	if err != nil {
		return "", err
	}
	return string(buf), nil
}
