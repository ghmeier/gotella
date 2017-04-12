package handlers

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/ghmeier/gotella/models"
	"github.com/ghmeier/gotella/receiver"
)

type queryHitHandler struct {
	*handler
}

func HandleQueryHit(ctx *Context) receiver.ReceiverFunc {
	q := &queryHitHandler{
		handler: newHandler(ctx.Peer, ctx.Descriptor, ctx.Files),
	}
	return q.queryHit
}

func (h *queryHitHandler) queryHit(conn *net.TCPConn, d *models.Descriptor) {
	fmt.Println("%s: query hit from %s:%d\n", d.Header.ID.String(), d.IP, d.Port)
	cache, err := h.descriptor.Get(d.Header.ID)
	if err != nil {
		printError(err)
	}
	if cache == nil || cache.Type == d.Header.Type {
		// queryhits must have corresponding query
		return
	}

	if cache.Port == 0 && cache.IP == "" {
		// query hit was meant for us, attempt to download
		conn.Close()
		h.streamFile(d)
		return
	}

	next := d.Next()
	if next.Header.TTL < 0 {
		return
	}
	fmt.Printf("%s: passing query hit to %s:%d\n", d.Header.ID.String(), cache.IP, cache.Port)
	connectAndSend(fmt.Sprintf("%s:%d", cache.IP, cache.Port), next)
}

func (h *queryHitHandler) streamFile(d *models.Descriptor) {
	var hit models.QueryHit
	err := json.Unmarshal(d.Payload, &hit)
	if err != nil {
		printError(err)
		return
	}

	if len(hit.Results) <= 0 {
		return
	}
	result := hit.Results[0]

	addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", hit.IP, hit.Port))
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		printError(err)
		return
	}
	defer conn.Close()

	host, port := localAddr(conn)
	descriptor := models.StreamDescriptor(host, port, &models.Stream{FileName: result.Name})
	send(conn, descriptor)

	buf := make([]byte, result.Size)
	_, err = conn.Read(buf)
	if err != nil {
		printError(err)
		return
	}
	h.files.Save(result.Name, buf)
}
