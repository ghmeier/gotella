package models

import (
	"encoding/json"

	"github.com/pborman/uuid"
)

const DEFAULT_TTL = 5

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
	Files int   `json:"files"`
	Size  int64 `json:"size"`
}

type Stream struct {
	FileName string `json:"filename"`
}

type Query struct {
	Filename    string `json:"filename"`
	RequestIP   string `json:"ip"`
	RequestPort int    `json:"port"`
}

type QueryHit struct {
	Hits    int      `json:"hits"`
	IP      string   `json:"ip"`
	Port    int      `json:"port"`
	Results []Result `json:"results"`
}

type Result struct {
	Name  string `json:"name"`
	Size  int64  `json:"size"`
	Index string `json:"index"`
}

type DescriptorCache struct {
	ID   uuid.UUID      `json:"id"`
	IP   string         `json:"ip"`
	Port int            `json:"port"`
	Type DescriptorType `json:"type"`
}

type Peer struct {
	IP    string `json:"ip"`
	Port  int    `json:"port"`
	Files int    `json:"files"`
	Size  int64  `json:"size"`
}

func FromBuf(buf []byte) (*Descriptor, error) {
	var d Descriptor
	err := json.Unmarshal(buf, &d)

	if err != nil {
		return nil, err
	}

	return &d, nil
}

func PongDescriptor(ip string, port int, p *Pong, d *Descriptor) *Descriptor {
	buf, _ := json.Marshal(p)
	header := &DescriptorHeader{
		ID:     d.Header.ID,
		Type:   PONG,
		TTL:    d.Header.Hops,
		Hops:   0,
		Length: len(buf),
	}
	return NewDescriptor(ip, port, buf, header)
}

func PingDescriptor(ip string, port int) *Descriptor {
	buf := make([]byte, 0)
	header := &DescriptorHeader{
		ID:     uuid.NewUUID(),
		Type:   PING,
		TTL:    DEFAULT_TTL,
		Hops:   0,
		Length: 0,
	}
	return NewDescriptor(ip, port, buf, header)
}

func QueryDescriptor(ip string, port int, q *Query) *Descriptor {
	buf, _ := json.Marshal(q)
	header := &DescriptorHeader{
		ID:     uuid.NewUUID(),
		Type:   QUERY,
		TTL:    DEFAULT_TTL,
		Hops:   0,
		Length: len(buf),
	}
	return NewDescriptor(ip, port, buf, header)
}

func QueryHitDescriptor(ip string, port int, q *QueryHit, d *Descriptor) *Descriptor {
	buf, _ := json.Marshal(q)
	header := &DescriptorHeader{
		ID:     d.Header.ID,
		Type:   QUERYHIT,
		TTL:    d.Header.Hops,
		Hops:   0,
		Length: len(buf),
	}
	return NewDescriptor(ip, port, buf, header)
}

func StreamDescriptor(ip string, port int, s *Stream) *Descriptor {
	buf, _ := json.Marshal(s)
	header := &DescriptorHeader{
		ID:     uuid.NewUUID(),
		Type:   STREAM,
		TTL:    1,
		Hops:   0,
		Length: len(buf),
	}
	return NewDescriptor(ip, port, buf, header)
}

func (d *Descriptor) Next() *Descriptor {
	header := &DescriptorHeader{
		ID:     d.Header.ID,
		Type:   d.Header.Type,
		TTL:    d.Header.TTL - 1,
		Hops:   d.Header.Hops + 1,
		Length: d.Header.Length,
	}
	return NewDescriptor(d.IP, d.Port, d.Payload, header)
}

func NewDescriptor(ip string, port int, payload []byte, h *DescriptorHeader) *Descriptor {
	return &Descriptor{
		Version: "0.1",
		IP:      ip,
		Port:    port,
		Payload: payload,
		Header:  h,
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
	case STREAM:
		return STREAM
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
	STREAM   = 5
	INVALID  = -1
)
