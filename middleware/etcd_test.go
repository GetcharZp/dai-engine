package middleware

import (
	"fmt"
	"testing"
)

var etcdEndpoints = []string{"119.27.186.148:2379"}

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
	resp, err := etcdClient.Get("/services")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp)
}

func TestGetByPrefixKey(t *testing.T) {
	etcdClient := NewEtcdClient(etcdEndpoints, "", "")
	resp, err := etcdClient.GetByPrefixKey("/services")
	if err != nil {
		t.Fatal(err)
	}
	for _, value := range resp {
		fmt.Println(string(value.Value))
	}
}
