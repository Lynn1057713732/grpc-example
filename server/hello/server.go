package main

import (
	"context"
	//"log"
	//"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	//"google.golang.org/grpc"

	pb "grpc-example/proto/hello"
)

// GreeterServer is the server API for Greeter service.
// 定义结构体，在调用注册api的时候作为入参，
// 该结构体会带上SayHello方法，里面是业务代码
// 这样远程调用时就执行了业务代码了
type GreeterServer struct {
	// pb.go中自动生成的，是个空结构体
	pb.UnimplementedGreeterServer
}

// SayHello implement to say hello
// 业务代码在此写，客户端远程调用SayHello时，
// 会执行这里的代码
func (h *GreeterServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	// 实例化结构体HelloReply，作为返回值
	return &pb.HelloReply{
		Message: "hello " + req.Name,
	}, nil
}


func main() {
	ctx := context.TODO()
	mux := runtime.NewServeMux()
	// Register generated routes to mux
	err := pb.RegisterGreeterHandlerServer(ctx, mux, &GreeterServer{})
	if err != nil {
		panic(err)
	}
	// Register custom route for  GET /hello/{name}
	//err = mux.HandlePath("GET", "/hello/{name}", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	//	w.Write([]byte("hello " + pathParams["name"]))
	//})
	//if err != nil {
	//	panic(err)
	//}
	http.ListenAndServe(":8090", mux)
}



//type server struct{
//	helloworldpb.UnimplementedGreeterServer
//}
//
//func NewServer() *server {
//	return &server{}
//}
//
//func (s *server) SayHello(ctx context.Context, in *helloworldpb.HelloRequest) (*helloworldpb.HelloReply, error) {
//	return &helloworldpb.HelloReply{Message: in.Name + " world"}, nil
//}
//
//func main() {
//	// Create a listener on TCP port
//	lis, err := net.Listen("tcp", ":8080")
//	if err != nil {
//		log.Fatalln("Failed to listen:", err)
//	}
//
//	// Create a gRPC server object
//	s := grpc.NewServer()
//
//	// Attach the Greeter service to the server
//	helloworldpb.RegisterGreeterServer(s, &server{})
//	// Serve gRPC server
//	log.Println("Serving gRPC on 0.0.0.0:8080")
//	go func() {
//		log.Fatalln(s.Serve(lis))
//	}()
//
//	// Create a client connection to the gRPC server we just started
//	// This is where the gRPC-Gateway proxies the requests
//	conn, err := grpc.DialContext(
//		context.Background(),
//		"0.0.0.0:8080",
//		grpc.WithBlock(),
//		grpc.WithInsecure(),
//	)
//	if err != nil {
//		log.Fatalln("Failed to dial server:", err)
//	}
//
//	gwmux := runtime.NewServeMux()
//	// Register Greeter
//	err = helloworldpb.RegisterGreeterHandler(context.Background(), gwmux, conn)
//	if err != nil {
//		log.Fatalln("Failed to register gateway:", err)
//	}
//
//	gwServer := &http.Server{
//		Addr:    ":8090",
//		Handler: gwmux,
//	}
//
//	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090")
//	log.Fatalln(gwServer.ListenAndServe())
//}

