package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"

	pb "grpc-example/proto"
)

const (
	//Address监听地址
	Address string = ":8000"
	//network网络通信协议
	Network string = "tcp"
)

type SimpleService struct{}

func (s *SimpleService) Route(ctx context.Context, req *pb.SimpleRequest) (*pb.SimpleResponse, error) {
	res := pb.SimpleResponse{
		Code: 200,
		Value: "hello " + req.Data,
	}

	return &res, nil

}

func main()  {
	//监听本地端口
	listener, err := net.Listen(Network, Address)
	if err != nil {
		log.Fatalf("net listen err: %v", err)
	}

	log.Println(Address + " net listening...")

	//新建gRPC服务器实例
	gRPCServer := grpc.NewServer()
	//在gRPC服务器中注册我们的服务
	pb.RegisterSimpleServer(gRPCServer, &SimpleService{})

	//用服务器Server方法以及我们的福安口信息区实现阻塞等待， 直到进程被杀死或者Stop的调用
	err = gRPCServer.Serve(listener)
	if err != nil {
		log.Fatalf("gPRC Server err: %v", err)

	}



}


