package discovery

import (
	"fmt"
	"testing"
)

var etcdEndpoints = []string{"192.168.1.8:2379"}

func TestPut(t *testing.T) {
	etcdClient := NewEtcdClient(etcdEndpoints, "", "")
	err := etcdClient.Put("my-key", "my-value")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Success")
}

func TestGet(t *testing.T) {
	etcdClient := NewEtcdClient(etcdEndpoints, "", "")
	resp, err := etcdClient.Get("my-key")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp)
}
