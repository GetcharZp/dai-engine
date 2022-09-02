package service

import (
	"context"
	"dai-engine/define"
	"dai-engine/helper"
	"dai-engine/logger"
	"dai-engine/middleware"
	"github.com/golang/protobuf/jsonpb"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	rpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type Work struct {
	Port          string
	Method        string // load method
	ServicePrefix string
	EtcdConfig    *EtcdConfig
	EtcdClient    *middleware.EtcdClient
	Discovery     *Discovery
}

type Worker func(*Work)

const WorkMethodRandom = "random"

// newDefaultWork get default work
func newDefaultWork() *Work {
	return &Work{
		Port:          "13100",
		Method:        WorkMethodRandom,
		ServicePrefix: "/DaiServices",
	}
}

// SetWorkPort change work default port
func SetWorkPort(port string) Worker {
	return func(work *Work) {
		work.Port = port
	}
}

// SetWorkMethod change work default method
func SetWorkMethod(method string) Worker {
	return func(work *Work) {
		work.Method = method
	}
}

// SetWorkServicePrefix change work default service prefix
func SetWorkServicePrefix(prefix string) Worker {
	return func(work *Work) {
		work.ServicePrefix = prefix
	}
}

// SetEtcdConfig change work default service prefix
func SetEtcdConfig(etcd *EtcdConfig) Worker {
	return func(work *Work) {
		work.EtcdConfig = etcd
	}
}

// NewWork get worker struct
func NewWork(endpoints []string, username, password string, workers ...Worker) *Work {
	work := newDefaultWork()
	for _, v := range workers {
		v(work)
	}
	work.EtcdConfig = &EtcdConfig{
		Endpoints: endpoints,
		Username:  username,
		Password:  password,
	}
	return work
}

// Run start gateway
func (w *Work) Run() {
	w.EtcdClient = middleware.NewEtcdClient(w.EtcdConfig.Endpoints, w.EtcdConfig.Username, w.EtcdConfig.Password)
	logger.Info("Address : http://localhost:" + w.Port)
	w.httpServer()
}

// httpServer http proxy
func (w *Work) httpServer() {
	switch w.Method {
	case WorkMethodRandom:
		http.HandleFunc("/", w.HttpHandleRandomServer)
	default:
		panic("[SYSTEM ERROR] : Error Work Method")
	}
	err := http.ListenAndServe(":"+w.Port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}

func (w *Work) HttpHandleRandomServer(writer http.ResponseWriter, req *http.Request) {
	// CORS
	writer.Header().Set("content-type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
	writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, X-Extra-Header, Content-Type, Accept, Authorization,Request-Id")
	// service
	if req.RequestURI == "/" {
		writer.Write(helper.MapToByte(define.M{"code": 200, "msg": "Pong"}))
		return
	}
	uri := strings.Split(strings.Trim(req.RequestURI, "/"), "/")
	// judge service exits
	serviceList, err := w.EtcdClient.GetByPrefixKey(w.ServicePrefix + "/" + uri[0])
	if err != nil {
		logger.Error("[ETCD ERROR] : " + err.Error())
		writer.Write(helper.MapToByte(define.M{"code": -1, "msg": "Service Empty"}))
		return
	}
	// proxy
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(serviceList))
	resp := httpProxy(uri, string(serviceList[index].Value), req.Body)
	defer req.Body.Close()
	writer.Write(resp)
}

// httpProxy system http proxy
// uri include systemKey,service,method
func httpProxy(uri []string, serverAddr string, body io.Reader) []byte {
	if len(uri) < 3 {
		return helper.MapToByte(define.M{
			"code": -1,
			"msg":  "Error URI",
		})
	}

	// dial grpc, can optimized use rpc pool
	grpcClient, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		logger.Error("[GRPC ERROR] : " + err.Error())
		return helper.MapToByte(define.M{
			"code": -1,
			"msg":  "ERROR URI",
		})
	}
	defer grpcClient.Close()

	// reflect client
	client := grpcreflect.NewClient(context.Background(), rpb.NewServerReflectionClient(grpcClient))

	serviceDes, err := client.ResolveService(uri[0] + "." + uri[1])
	if err != nil {
		logger.Error("[RESOLVE ERROR] : " + err.Error())
		return helper.MapToByte(define.M{
			"code": -1,
			"msg":  "ERROR SERVICE",
		})
	}

	methodDes := serviceDes.FindMethodByName(uri[2])
	if methodDes == nil {
		logger.Error("[FIND METHOD ERROR] Method : " + uri[2])
		return helper.MapToByte(define.M{
			"code": -1,
			"msg":  "ERROR METHOD",
		})
	}

	b, err := ioutil.ReadAll(body)
	pp := dynamic.NewMessage(methodDes.GetInputType())
	err = pp.UnmarshalJSONPB(&jsonpb.Unmarshaler{AllowUnknownFields: true}, helper.If(len(b) == 0, []byte("{}"), b).([]byte))
	if err != nil {
		logger.Error("[UnmarshalJSONPB ERROR] : " + err.Error())
		return helper.MapToByte(define.M{
			"code": -1,
			"msg":  "UnmarshalJSONPB ERROR",
		})
	}

	stub := grpcdynamic.NewStub(grpcClient)
	respInvoke, err := stub.InvokeRpc(context.Background(), methodDes, pp)
	if err != nil {
		logger.Error("[INVOKE ERROR] : " + err.Error())
		return helper.MapToByte(define.M{
			"code": -1,
			"msg":  "INVOKE ERROR",
		})
	}

	ri, err := dynamic.AsDynamicMessage(respInvoke)
	if err != nil {
		logger.Error("[DYNAMIC MESSAGE ERROR] : " + err.Error())
		return helper.MapToByte(define.M{
			"code": -1,
			"msg":  "DYNAMIC MESSAGE ERROR",
		})
	}

	riBytes, err := ri.MarshalJSONPB(&jsonpb.Marshaler{EmitDefaults: true})
	if err != nil {
		logger.Error("[MarshalJSONPB ERROR] : " + err.Error())
		return helper.MapToByte(define.M{
			"code": -1,
			"msg":  "MarshalJSONPB ERROR",
		})
	}

	return riBytes
}
