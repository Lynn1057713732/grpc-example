package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"runtime"
	"time"

	pb "grpc-example/proto"
)

const (
	//Address监听地址
	Address string = ":8000"
	//network网络通信协议
	Network string = "tcp"
)



func main()  {
	//监听本地端口
	listener, err := net.Listen(Network, Address)
	if err != nil {
		log.Fatalf("net listen err: %v", err)
	}

	// 从输入证书文件和密钥文件为服务端构造TLS凭证
	creds, err := credentials.NewServerTLSFromFile("../tls/server.pem", "../tls/server.key")
	if err != nil {
		log.Fatalf("Failed to generate credentials %v", err)
	}

	//普通方法：一元拦截器（grpc.UnaryInterceptor）
	var interceptor grpc.UnaryServerInterceptor
	interceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		//拦截普通方法请求，验证Token
		err = Check(ctx)
		if err != nil {
			return
		}
		// 继续处理请求
		return handler(ctx, req)
	} 


	//新建gRPC服务器实例
	gRPCServer := grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(interceptor))

	//在gRPC服务器中注册我们的服务
	pb.RegisterSimpleServer(gRPCServer, &SimpleService{})
	log.Println(Address + " net.Listing whth TLS and token...")

	//用服务器Server方法以及我们的福安口信息区实现阻塞等待， 直到进程被杀死或者Stop的调用
	err = gRPCServer.Serve(listener)
	if err != nil {
		log.Fatalf("gPRC Server err: %v", err)

	}
}

type SimpleService struct{}

func (s *SimpleService) Route(ctx context.Context, req *pb.SimpleRequest) (*pb.SimpleResponse, error) {
	data := make(chan *pb.SimpleResponse, 1)
	defer close(data)

	go handle(ctx, req, data)
	select {
	case res := <-data:
		return res, nil
	case <-ctx.Done():
		return nil, status.Errorf(codes.Canceled, "Client cancelled, abandoning.")
	}
}

func handle(ctx context.Context, req *pb.SimpleRequest, data chan<- *pb.SimpleResponse) {
	select {
	case <- ctx.Done():
		log.Println(ctx.Err())
		log.Println("handle go routine exit")
		runtime.Goexit()  //超时后退出go协程
	case <- time.After(2 * time.Second): // 模拟耗时操作
		res := pb.SimpleResponse{
			Code: 200,
			Value: "hello " + req.Data,
		}
		//修改数据库前，进行超时判断
		//if ctx.Err() == context.Canceled {
		//	//如果已经超时，则退出
		//}
		data <- &res

	}
}

// Check 验证token
func Check(ctx context.Context) error {
	//从上下文中获取元数据
	md, ok := metadata.FromIncomingContext(ctx)
	log.Println(md)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "获取Token失败")
	}
	var (
		appID     string
		appSecret string
	)
	if value, ok := md["app_id"]; ok {
		appID = value[0]
	}
	if value, ok := md["app_secret"]; ok {
		appSecret = value[0]
	}
	if appID != "grpc_token" || appSecret != "123456" {
		return status.Errorf(codes.Unauthenticated, "Token无效: app_id=%s, app_secret=%s", appID, appSecret)
	}
	return nil
}

