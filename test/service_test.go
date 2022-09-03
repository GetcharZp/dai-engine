package test

import (
	"dai-engine/service"
	proto "dai-engine/test/proto"
	"dai-engine/test/service/my"
	"testing"
)

// TestNewService Create new service
func TestNewService(t *testing.T) {
	s := service.NewService(&service.EtcdConfig{Endpoints: []string{"192.168.1.8:2379"}}, "my", "13110")
	proto.RegisterMyServer(s.GrpcServer, my.ServiceMy)
	s.Run()
}
