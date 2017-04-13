package gateways

import (
	"fmt"
	"time"

	"gopkg.in/redis.v3"
)

type Redis interface {
	Set(string, []byte, time.Duration) error
	Get(string) ([]byte, error)
	Exists(string) (bool, error)
	Append(string, string) error
	List(key string) ([]string, error)
	Len(key string) (int, error)
	Remove(key, value string)
}

type rClient struct {
	client *redis.Client
	id     string
}

func New(addr string, id string) Redis {
	return &rClient{
		client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: "",
			DB:       0,
		}),
		id: id,
	}
}

func (c *rClient) Set(key string, value []byte, exp time.Duration) error {
	cmd := c.client.Set(c.key(key), value, exp)
	return cmd.Err()
}

func (c *rClient) Get(key string) ([]byte, error) {
	cmd := c.client.Get(c.key(key))
	return []byte(cmd.Val()), cmd.Err()
}

func (c *rClient) Append(key string, value string) error {
	cmd := c.client.SAdd(c.key(key), value)
	return cmd.Err()
}

func (c *rClient) Remove(key, value string) {
	c.client.SRem(c.key(key), value)
}

func (c *rClient) List(key string) ([]string, error) {
	cmd := c.client.SMembers(c.key(key))
	return cmd.Val(), cmd.Err()
}

func (c *rClient) Len(key string) (int, error) {
	cmd := c.client.SCard(c.key(key))
	return int(cmd.Val()), cmd.Err()
}

func (c *rClient) Exists(key string) (bool, error) {
	cmd := c.client.Exists(c.key(key))
	return cmd.Val(), cmd.Err()
}

func (c *rClient) key(base string) string {
	return fmt.Sprintf("%s-%s", c.id, base)
}
