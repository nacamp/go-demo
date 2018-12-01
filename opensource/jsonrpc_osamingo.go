package opensource

import (
	"context"
	"log"
	"net/http"

	"github.com/intel-go/fastjson"

	"github.com/osamingo/jsonrpc"
)

type JsonHandler interface {
	jsonrpc.Handler
	Name() string
	Params() interface{}
	Result() interface{}
}

type Handler1 struct {
}

func (h *Handler1) Name() string {
	return "handler1"
}
func (h *Handler1) ServeJSONRPC(c context.Context, params *fastjson.RawMessage) (interface{}, *jsonrpc.Error) {
	return []string{"string slice"}, nil
}
func (h *Handler1) Params() interface{} {
	return []string{}
}
func (h *Handler1) Result() interface{} {
	return []string{}
}

type Handler2 struct {
}

func (h *Handler2) Name() string {
	return "handler2"
}
func (h *Handler2) ServeJSONRPC(c context.Context, params *fastjson.RawMessage) (interface{}, *jsonrpc.Error) {
	p := []string{}
	if err := jsonrpc.Unmarshal(params, &p); err != nil {
		return nil, err
	}
	return p[0], nil
}
func (h *Handler2) Params() interface{} {
	return []string{}
}
func (h *Handler2) Result() interface{} {
	return ""
}

type Handler3 struct {
}

type (
	EchoHandler struct{}
	EchoParams  struct {
		Name string `json:"name"`
	}
	EchoResult struct {
		Message string `json:"message"`
	}
)

func (h *Handler3) Name() string {
	return "handler3"
}
func (h *Handler3) ServeJSONRPC(c context.Context, params *fastjson.RawMessage) (interface{}, *jsonrpc.Error) {
	// var p struct {
	// 	Name string `json:"name"`
	// }
	var p EchoParams
	if err := jsonrpc.Unmarshal(params, &p); err != nil {
		return nil, err
	}

	return struct {
		Message string `json:"message"`
	}{Message: p.Name}, nil
}
func (h *Handler3) Params() interface{} {
	// return EchoParams{}
	return struct {
		Name string `json:"name"`
	}{}
}
func (h *Handler3) Result() interface{} {
	return EchoResult{}
}

type RpcService struct {
	server *RpcServer
}

func (rs *RpcService) Setup(server *RpcServer) {
	rs.server = server
	rs.server.RegisterHandler(&Handler1{})
	rs.server.RegisterHandler(&Handler2{})
	rs.server.RegisterHandler(&Handler3{})
}

type RpcServer struct {
	mr *jsonrpc.MethodRepository
}

func NewRpcServer() *RpcServer {
	return &RpcServer{
		mr: jsonrpc.NewMethodRepository(),
	}
}

func (js *RpcServer) RegisterHandler(handler JsonHandler) {
	if err := js.mr.RegisterMethod(handler.Name(), handler, handler.Params, handler.Result); err != nil {
		log.Fatalln(err)
	}

}

func (js *RpcServer) Start() {

	http.Handle("/jrpc", js.mr)
	// http.HandleFunc("/jrpc/debug", js.mr.ServeDebug)
	go func() {
		if err := http.ListenAndServe(":8080", http.DefaultServeMux); err != nil {
			log.Fatalln(err)
		}
	}()
}

func Jsonrpc_osamingo_main() {
	server := NewRpcServer()
	service := RpcService{}
	service.Setup(server)
	server.Start()
	select {}
}

/*
curl -X POST -H "Content-Type: application/json" -d '{"jsonrpc": "2.0",   "method": "handler1", "params":[]}' http://localhost:8080/jrpc
{"jsonrpc":"2.0","result":["string slice"]}

curl -X POST -H "Content-Type: application/json" -d '{"jsonrpc": "2.0",   "method": "handler2", "params":["param"]}' http://localhost:8080/jrpc
{"jsonrpc":"2.0","result":"param"}

curl -X POST -H "Content-Type: application/json" -d '{"jsonrpc": "2.0",   "method": "handler3", "params": {"name": "John Doe"}}' http://localhost:8080/jrpc
{"jsonrpc":"2.0","result":{"message":"John Doe"}}
*/
