package service

import (
	"dai-engine/define"
	"dai-engine/helper"
	"dai-engine/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type Service struct {
	EtcdConfig    *EtcdConfig
	ServiceKey    string
	Port          string
	ServicePrefix string
	GrpcServer    *grpc.Server
	register      *Register
}

type FuncService func(*Service)

func SetServicePrefix(prefix string) FuncService {
	return func(service *Service) {
		service.ServicePrefix = prefix
	}
}

func NewService(es *EtcdConfig, serviceKey, port string, services ...FuncService) *Service {
	s := &Service{
		EtcdConfig: es,
		ServiceKey: serviceKey,
		Port:       port,
		GrpcServer: grpc.NewServer(),
	}
	for _, v := range services {
		v(s)
	}
	w := newDefaultWork()
	if s.ServicePrefix == "" {
		s.ServicePrefix = w.ServicePrefix
	}

	var (
		key   = s.ServicePrefix + "/" + serviceKey + "/" + helper.GetUUID()
		value = helper.GetLocalIp() + ":" + port
	)
	r, err := NewRegister(es, key, value, define.Lease)
	if err != nil {
		log.Fatalln("[NEW REGISTER ERROR] : " + err.Error())
	}
	s.register = r
	return s
}

func (s *Service) Run() {
	reflection.Register(s.GrpcServer)
	listener, err := net.Listen("tcp", ":"+s.Port)
	if err != nil {
		log.Fatalln("[NEW LISTENER ERROR] : " + err.Error())
	}
	logger.Info("[" + s.ServiceKey + "] RUN PORT [:" + s.Port + "]")
	err = s.GrpcServer.Serve(listener)
	if err != nil {
		log.Fatalln("[GRPC SERVER ERROR] : " + err.Error())
	}
}
