package middleware

import (
	"context"
	dai_engine "dai-engine"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	_, err = client.Status(ctx, endpoints[0])
	if err != nil {
		log.Fatalf("Etct Connect Error: %v", err)
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

// GetByPrefixKey get value by prefix key
func (c *EtcdClient) GetByPrefixKey(key string) ([]*mvccpb.KeyValue, error) {
	r, err := c.Client.Get(context.Background(), key, clientv3.WithPrefix())
	if r.Count == 0 {
		return nil, dai_engine.ErrEmptyValue
	}
	return r.Kvs, err
}
