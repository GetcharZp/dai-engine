package service

import (
	"log"
	"testing"
	"time"
)

func TestDiscovery(t *testing.T) {
	d := NewDiscovery(&EtcdConfig{
		Endpoints: []string{"192.168.1.8:2379"},
		Username:  "",
		Password:  "",
	})
	defer d.Close()
	err := d.WatchService("/services")
	if err != nil {
		log.Fatalln(err)
	}
	for {
		select {
		case <-time.Tick(2 * time.Second):
			log.Println(d.serviceList)
		}
	}
}
