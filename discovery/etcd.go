package discovery

import (
	"context"
	dai_engine "dai-engine"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
)

type EtcdClient struct {
	Client *clientv3.Client
}

// NewEtcdClient Get EtcdClient
func NewEtcdClient(endpoints []string, username, password string) *EtcdClient {
	config := clientv3.Config{
		Endpoints: endpoints,
		Username:  username,
		Password:  password,
	}
	client, err := clientv3.New(config)
	if err != nil {
		log.Fatal(fmt.Sprintf("New Etcd Error:%v", err))
	}
	return &EtcdClient{
		Client: client,
	}
}

// Get Etcd Value By Key
func (c *EtcdClient) Get(key string) (string, error) {
	r, err := c.Client.Get(context.Background(), key)
	if err != nil {
		return "", err
	}
	if r.Count == 0 {
		return "", dai_engine.ErrEmptyValue
	}
	return string(r.Kvs[0].Value), nil
}

// Put Create Or Change Etcd Value
func (c *EtcdClient) Put(key, val string) error {
	_, err := c.Client.Put(context.Background(), key, val)
	return err
}
