package helpers

import (
	"encoding/json"
	"time"

	"github.com/pborman/uuid"

	"github.com/ghmeier/gotella/gateways"
	"github.com/ghmeier/gotella/models"
)

type Descriptor interface {
	Put(*models.Descriptor) error
	Exists(uuid.UUID) (bool, error)
	Get(uuid.UUID) (*models.DescriptorCache, error)
}

type descriptorHelper struct {
	redis gateways.Redis
}

func NewDescriptor(redis gateways.Redis) Descriptor {
	return &descriptorHelper{redis: redis}
}

func (h *descriptorHelper) Put(d *models.Descriptor) error {
	cache := &models.DescriptorCache{
		IP:   d.IP,
		Port: d.Port,
		Type: d.Header.Type,
	}
	buf, err := json.Marshal(cache)
	if err != nil {
		return err
	}
	err = h.redis.Set(d.Header.ID.String(), buf, time.Second*30)
	return err
}

func (h *descriptorHelper) Exists(id uuid.UUID) (bool, error) {
	return h.redis.Exists(id.String())
}

func (h *descriptorHelper) Get(id uuid.UUID) (*models.DescriptorCache, error) {
	buf, err := h.redis.Get(id.String())
	if err != nil {
		return nil, err
	}
	var cache models.DescriptorCache
	err = json.Unmarshal(buf, &cache)
	if err != nil {
		return nil, err
	}
	return &cache, nil
}
