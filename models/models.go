package models

import (
	"encoding/json"

	"github.com/pborman/uuid"
)

type Descriptor struct {
	Version string            `json:"version"`
	Header  *DescriptorHeader `json:"header"`
	IP      string            `json:"ip"`
	Port    int               `json:"port"`
	Payload []byte            `json:"payload"`
}

type DescriptorHeader struct {
	ID     uuid.UUID      `json:"id"`
	Type   DescriptorType `json:"type"`
	TTL    int            `json:"ttl"`
	Hops   int            `json:"hops"`
	Length int            `json:"length"`
}

type Pong struct {
	Files int `json:"files"`
	Size  int `json:"size"`
}

func FromBuff(buf []byte) (*Descriptor, error) {
	var d Descriptor
	err := json.Unmarshal(buf, &d)

	if err != nil {
		return nil, err
	}

	return &d, nil
}

func PongDescriptor(ip string, port int, p *Pong, h *DescriptorHeader) *Descriptor {
	buf, _ := json.Marshal(p)
	return &Descriptor{
		Version: "0.4",
		IP:      ip,
		Port:    port,
		Payload: buf,
		Header: &DescriptorHeader{
			ID:     uuid.NewUUID(),
			Type:   PONG,
			TTL:    h.TTL - 1,
			Hops:   h.Hops + 1,
			Length: len(buf),
		},
	}
}

func toDescriptorType(i int) DescriptorType {
	switch i {
	case PING:
		return PING
	case PONG:
		return PONG
	case QUERY:
		return QUERY
	case QUERYHIT:
		return QUERYHIT
	case PUSH:
		return PUSH
	default:
		return INVALID
	}
}

type DescriptorType int

const (
	PING     = 1
	PONG     = 2
	QUERY    = 3
	QUERYHIT = 4
	PUSH     = 5
	INVALID  = -1
)
