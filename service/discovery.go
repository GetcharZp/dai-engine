package service

import (
	"context"
	"dai-engine/middleware"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"sync"
)

type Discovery struct {
	client      *clientv3.Client
	lock        sync.RWMutex
	serviceList map[string]string
}

// NewDiscovery Create Service Discovery
func NewDiscovery(ec *EtcdConfig) *Discovery {
	client := middleware.NewEtcdClient(ec.Endpoints, ec.Username, ec.Password)
	return &Discovery{
		client:      client.Client,
		serviceList: make(map[string]string),
	}
}

// WatchService Init And Listen Service
func (d *Discovery) WatchService(prefix string) error {
	resp, err := d.client.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	// Init Service
	for _, v := range resp.Kvs {
		d.putServiceList(string(v.Key), string(v.Value))
	}
	go d.watch(prefix)
	return nil
}

// watch create or update serviceList
func (d *Discovery) watch(prefix string) {
	wc := d.client.Watch(context.Background(), prefix, clientv3.WithPrefix())
	for resp := range wc {
		for _, event := range resp.Events {
			switch event.Type {
			case mvccpb.PUT:
				d.putServiceList(string(event.Kv.Key), string(event.Kv.Value))
			case mvccpb.DELETE:
				d.deleteServiceList(string(event.Kv.Key))
			}
		}
	}
}

// PutServiceList add\modify one item to serviceList
func (d *Discovery) putServiceList(key, value string) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.serviceList[key] = value
	log.Println("key:", key, "value:", value)
}

// DeleteServiceList delete one item in serviceList
func (d *Discovery) deleteServiceList(key string) {
	d.lock.Lock()
	defer d.lock.Unlock()
	delete(d.serviceList, key)
}

// Close service
func (d *Discovery) Close() error {
	return d.client.Close()
}
