# DaiEngine

> 基于 Grpc、Etcd的微服务架子

## 安装

```shell
go get gitee.com/up-zero/dai-engine
```

## 快速开始

1. 启动 worker

```go
w := service.NewWork([]string{"192.168.1.8:2379"}, "", "")
w.Run()
```

2. 编写并启动服务

```go
s := service.NewService(&service.EtcdConfig{Endpoints: []string{"192.168.1.8:2379"}}, "my", "13110")
// 注册grpc服务
proto.RegisterMyServer(s.GrpcServer, my.ServiceMy)
s.Run()
```