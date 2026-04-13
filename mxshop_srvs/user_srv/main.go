package main

import (
	"flag"
	"mxshop_srvs/user_srv/handler"
	"mxshop_srvs/user_srv/proto"

	"fmt"
	"net"

	"google.golang.org/grpc"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 50051, "端口号")
	flag.Parse()
	fmt.Println("ip: ", *IP)
	fmt.Println("port: ", *Port)

	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("failed to listen: " + err.Error())
	}

	err = server.Serve(lis)
	if err != nil {
		panic("failed to serve: " + err.Error())
	}	
}
