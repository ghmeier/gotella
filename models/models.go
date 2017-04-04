package models

import (
	"encoding/json"

	"github.com/pborman/uuid"
)

type Descriptor struct {
	Version string           `json:"version"`
	Header  DescriptorHeader `json:"header"`
	IP      string           `json:"ip"`
	Port    string           `json:"port"`
	Payload []byte           `json:"payload"`
}

type DescriptorHeader struct {
	ID     uuid.UUID      `json:"id"`
	Type   DescriptorType `json:"type"`
	TTL    int            `json:"ttl"`
	Hops   int            `json:"hops"`
	Length int            `json:"length"`
}

func NewDescriptor(buf []byte) (*Descriptor, error) {
	var d Descriptor
	err := json.Unmarshal(buf, &d)

	if err != nil {
		return nil, err
	}

	return &d, nil
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
