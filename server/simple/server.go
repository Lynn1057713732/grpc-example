package main

import (
	"context"
	"crypto/tls"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"grpc-example/server/gateway"
	"log"
	"net"
	//"runtime"
	//"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"

	pb "grpc-example/proto/simple"
	"grpc-example/server/middleware/auth"
	"grpc-example/server/middleware/cred"
	"grpc-example/server/middleware/recovery"
	"grpc-example/server/middleware/zap"
)

const (
	// Address 监听地址
	Address string = "127.0.0.1:8000"

	// Network 网络通信协议
	Network string = "tcp"
)



func main()  {
	//监听本地端口
	listener, err := net.Listen(Network, Address)
	if err != nil {
		log.Fatalf("net listen err: %v", err)
	}

	//新建gRPC服务器实例
	gRPCServer := grpc.NewServer(
		cred.TLSInterceptor(),
		grpc.ChainStreamInterceptor(
			grpc_middleware.ChainStreamServer(
				// grpc_ctxtags.StreamServerInterceptor(),
				// grpc_opentracing.StreamServerInterceptor(),
				// grpc_prometheus.StreamServerInterceptor,
				grpc_zap.StreamServerInterceptor(zap.ZapInterceptor()),
				grpc_auth.StreamServerInterceptor(auth.AuthInterceptor),
				grpc_recovery.StreamServerInterceptor(recovery.RecoveryInterceptor()),
				),
			),
		grpc.ChainUnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				// grpc_ctxtags.StreamServerInterceptor(),
				// grpc_opentracing.StreamServerInterceptor(),
				// grpc_prometheus.StreamServerInterceptor,
				grpc_zap.UnaryServerInterceptor(zap.ZapInterceptor()),
				grpc_auth.UnaryServerInterceptor(auth.AuthInterceptor),
				grpc_recovery.UnaryServerInterceptor(recovery.RecoveryInterceptor()),
				),
			),
		)

	//在gRPC服务器中注册我们的服务
	pb.RegisterSimpleServer(gRPCServer, &SimpleService{})
	log.Println(Address + " gRPC net.Listing with TLS and Token...")

	//使用gateway把grpcServer转成httpServer
	httpServer := gateway.ProvideHTTP(Address, gRPCServer)
	//用服务器 Serve() 方法以及我们的端口信息区实现阻塞等待，直到进程被杀死或者 Stop() 被调用
	if err = httpServer.Serve(tls.NewListener(listener, httpServer.TLSConfig)); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	//用服务器Server方法以及我们的端口信息区实现阻塞等待， 直到进程被杀死或者Stop的调用
	//err = gRPCServer.Serve(listener)
	//if err != nil {
	//	log.Fatalf("gPRC Server err: %v", err)
	//
	//}
}

type SimpleService struct{
	pb.UnimplementedSimpleServer
}

// Route 实现Route方法
func (s *SimpleService) Route(ctx context.Context, req *pb.InnerMessage) (*pb.OuterMessage, error) {
	res := pb.OuterMessage{
		ImportantString: "hello grpc validator",
		Inner:           req,
	}
	return &res, nil
}

//func (s *SimpleService) Route(ctx context.Context, req *pb.SimpleRequest) (*pb.SimpleResponse, error) {
//	data := make(chan *pb.SimpleResponse, 1)
//	defer close(data)
//
//	//从上下文中获取特殊的元数据,做一些特殊的处理
//	md, _ := metadata.FromIncomingContext(ctx)
//	log.Println(md)
//
//
//	go handle(ctx, req, data)
//	select {
//	case res := <-data:
//		return res, nil
//	case <-ctx.Done():
//		return nil, status.Errorf(codes.Canceled, "Client cancelled, abandoning.")
//	}
//}

//func handle(ctx context.Context, req *pb.SimpleRequest, data chan<- *pb.SimpleResponse) {
//	select {
//	case <- ctx.Done():
//		log.Println(ctx.Err())
//		log.Println("handle go routine exit")
//		runtime.Goexit()  //超时后退出go协程
//	case <- time.After(2 * time.Second): // 模拟耗时操作
//		res := pb.SimpleResponse{
//			Code: 200,
//			Value: "hello " + req.Data,
//		}
//		//修改数据库前，进行超时判断
//		//if ctx.Err() == context.Canceled {
//		//	//如果已经超时，则退出
//		//}
//		data <- &res
//
//	}
//}

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

