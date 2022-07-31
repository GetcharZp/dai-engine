package service

import (
	"context"
	"dai-engine/middleware"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
)

// Register for register service
type Register struct {
	client        *clientv3.Client
	leaseID       clientv3.LeaseID
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	key           string
	val           string
}

// NewRegister register service
func NewRegister(ec *EtcdConfig, key, value string, lease int64) (*Register, error) {
	client := middleware.NewEtcdClient(ec.Endpoints, ec.Username, ec.Password)
	r := &Register{
		client: client.Client,
		key:    key,
		val:    value,
	}
	// set lease time and register service
	if err := r.register(lease); err != nil {
		return nil, err
	}
	return r, nil
}

// register set lease time and register service
func (r *Register) register(lease int64) error {
	// set lease time
	resp, err := r.client.Grant(context.Background(), lease)
	if err != nil {
		return err
	}
	// register service with lease
	_, err = r.client.Put(context.Background(), r.key, r.val, clientv3.WithLease(resp.ID))
	if err != nil {
		return err
	}
	// service keep alive
	keepAliveChan, err := r.client.KeepAlive(context.Background(), resp.ID)
	if err != nil {
		return err
	}
	r.leaseID = resp.ID
	r.keepAliveChan = keepAliveChan
	return nil
}

// ListenLeaseRespChan listen the response for lease
func (r *Register) ListenLeaseRespChan() {
	for keepAliveChan := range r.keepAliveChan {
		log.Println("lease success", keepAliveChan.String())
	}
	log.Println("close lease")
}

// Close cancel service
func (r *Register) Close() error {
	// revoke lease
	if _, err := r.client.Revoke(context.Background(), r.leaseID); err != nil {
		return err
	}
	return r.client.Close()
}
