package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	"log"
	"time"

	pb "grpc-example/proto"
)

const Address string = ":8000"

var gRPCClient pb.SimpleClient


func main() {
	//连接服务器
	conn, err := grpc.Dial(Address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("met connect err: %v", err)
	}
	defer conn.Close()

	ctx := context.Background()
	//建立gRPC连接
	gRPCClient = pb.NewSimpleClient(conn)
	route(ctx, 1)
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


