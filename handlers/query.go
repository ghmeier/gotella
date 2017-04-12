package handlers

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/ghmeier/gotella/models"
	"github.com/ghmeier/gotella/receiver"
)

type queryHandler struct {
	*handler
}

func HandleQuery(ctx *Context) receiver.ReceiverFunc {
	q := &queryHandler{
		handler: newHandler(ctx.Peer, ctx.Descriptor, ctx.Files),
	}
	return q.query
}

func (h *queryHandler) query(conn *net.TCPConn, d *models.Descriptor) {
	fmt.Println("%s: query from %s:%d\n", d.Header.ID.String(), d.IP, d.Port)
	e, err := h.descriptor.Exists(d.Header.ID)
	if err != nil {
		printError(err)
	}
	if e {
		return
	}

	var query models.Query
	err = json.Unmarshal(d.Payload, &query)
	if err != nil {
		printError(err)
		return
	}

	h.descriptor.Put(d)
	if err != nil {
		printError(err)
	}

	host, port := localAddr(conn)
	if h.files.Exists(query.Filename) {
		h.hit(host, port, query, d)
		return
	}

	next := d.Next()
	next.IP = host
	next.Port = port
	h.propogate(next)
}

func (h *queryHandler) hit(host string, port int, query *models.Query, d *models.Descriptor) {
	q := &models.QueryHit{
		IP:      host,
		Port:    port,
		Hits:    1,
		Results: make([]Result, 0),
	}
	result := models.Result{
		Name:  query.Filename,
		Size:  h.files.FileSize(query.Filename),
		Index: query.Filename,
	}
	q.Results = append(q.Results, result)

	descriptor := models.QueryHitDescriptor(host, port, q, d)
	fmt.Printf("%s: query hit to %s:%d\n", d.Header.ID.String(), d.IP, d.Port)
	connectAndSend(fmt.Sprintf("%s:%d", d.IP, d.Port), descriptor)
}
