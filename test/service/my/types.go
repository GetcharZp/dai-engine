package my

import (
	"context"
	my "dai-engine/test/proto"
)

var (
	ServiceMy = &serviceMy{}
)

type serviceMy struct{}

// SayHello method
func (s *serviceMy) SayHello(ctx context.Context, in *my.HelloRequest) (*my.HelloReply, error) {
	return &my.HelloReply{
		Message: in.GetName() + " : Hello World",
	}, nil
}
