package main

import (
	"context"
	"google.golang.org/grpc"
	"log"

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

	//建立gRPC连接
	gRPCClient = pb.NewSimpleClient(conn)
	route()
}

func route() {
	//创建发送结构体
	req := pb.SimpleRequest{
		Data: "gRPC",
	}

	//调用我们的服务Route方法
	//同时传入了一个context.Context，在需要时可以让我们改变RPC的行为，比如超时/取消一个正在运行的RPC
	res, err := gRPCClient.Route(context.Background(), &req)
	if err != nil {
		log.Fatalf("Call Route err: %v", err)
	}

	log.Println(res)
}


