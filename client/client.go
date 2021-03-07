package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	status "google.golang.org/grpc/status"
	"log"
	"time"

	"grpc-example/client/auth"
	pb "grpc-example/proto"
)

const Address string = ":8000"

var gRPCClient pb.SimpleClient


func main() {
	//从输入的证书文件中为客户端构造TLS凭证
	creds, err := credentials.NewClientTLSFromFile("../tls/server.pem", "go-grpc-example")
	if err != nil {
		log.Fatalf("Failed to create TLS credentials %v", err)
	}

	//构建Token
	token := auth.Token{
		Value: "bearer grpc.auth.token",
	}

	//连接服务器
	conn, err := grpc.Dial(Address, grpc.WithTransportCredentials(creds), grpc.WithPerRPCCredentials(&token))
	if err != nil {
		log.Fatalf("met connect err: %v", err)
	}
	defer conn.Close()


	ctx := context.Background()
	//灵活添加gRPC的metadata
	md := metadata.Pairs(
			"traceId", "11111", "meId", "22222",
			)
	//NewOutgoingContext底层调用了的withValue方法
	ctx = metadata.NewOutgoingContext(ctx, md)

	//建立gRPC连接
	gRPCClient = pb.NewSimpleClient(conn)
	route(ctx, 3)
}

//route调用服务端你的Route方法
func route(ctx context.Context, deadlines time.Duration) {
	//设置3秒超时时间
	clientDeadline := time.Now().Add(time.Duration(deadlines * time.Second))
	ctx, cancel := context.WithDeadline(ctx, clientDeadline)
	defer cancel()


	//创建发送结构体
	req := pb.SimpleRequest{
		Data: "gRPC",
	}

	//调用我们的服务Route方法
	//同时传入了一个超时截止的context.Context，在需要时可以让我们改变RPC的行为，比如超时/取消一个正在运行的RPC
	res, err := gRPCClient.Route(ctx, &req)
	if err != nil {
		//获取错误状态
		errStatus, ok := status.FromError(err)
		if ok {
			//判断是否为调用超时
			if errStatus.Code() == codes.DeadlineExceeded {
				log.Fatalln("Route Timeout")
			}
		}
		log.Fatalf("Call Route err: %v", err)
	}

	log.Println(res)
}


