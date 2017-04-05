package gateways

import (
	"time"

	"gopkg.in/redis.v3"
)

type Redis interface {
	Set(string, []byte, time.Duration) error
	Get(string) ([]byte, error)
	Exists(string) (bool, error)
	Append(string, ...string) error
	List(key string) ([]string, error)
	Len(key string) (int, error)
}

type rClient struct {
	client *redis.Client
}

func New(addr string) Redis {
	return &rClient{
		client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: "",
			DB:       0,
		}),
	}
}

func (c *rClient) Set(key string, value []byte, exp time.Duration) error {
	cmd := c.client.Set(key, value, exp)
	return cmd.Err()
}

func (c *rClient) Get(key string) ([]byte, error) {
	cmd := c.client.Get(key)
	return []byte(cmd.Val()), cmd.Err()
}

func (c *rClient) Append(key string, values ...string) error {
	cmd := c.client.LPush(key, values...)
	return cmd.Err()
}

func (c *rClient) List(key string) ([]string, error) {
	cmd := c.client.LRange(key, 0, -1)
	return cmd.Val(), cmd.Err()
}

func (c *rClient) Len(key string) (int, error) {
	cmd := c.client.LLen(key)
	return int(cmd.Val()), cmd.Err()
}

func (c *rClient) Exists(key string) (bool, error) {
	cmd := c.client.Exists(key)
	return cmd.Val(), cmd.Err()
}
