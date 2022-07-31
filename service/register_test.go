package service

import (
	"log"
	"testing"
	"time"
)

func TestRegisterService(t *testing.T) {
	r, err := NewRegister(&EtcdConfig{
		Endpoints: []string{"192.168.1.8:2379:2379"},
		Username:  "",
		Password:  "",
	}, "/services", "127.0.0.1:8080", 5)
	if err != nil {
		log.Fatalln(err)
	}
	go r.ListenLeaseRespChan()
	select {
	case <-time.After(100 * time.Second):
		r.Close()
	}
}
