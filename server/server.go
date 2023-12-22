package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"

	pb "github.com/huangjianchao95/grpc-learn/server/protogen"
)

var (
	port = flag.String("port", ":9999", "The server port")
)

type server struct {
	pb.UnimplementedLearnServiceServer
}

func (server *server) Hello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	res := &pb.HelloResponse{
		Msg: fmt.Sprintf("Hello %s", req.Name),
	}

	return res, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", *port)
	if err != nil {
		log.Fatal("rpc server, failed to listen: ", err)
	}

	defer lis.Close()
	grpcServer := grpc.NewServer()
	pb.RegisterLearnServiceServer(grpcServer, &server{})
	log.Println("rpc server, listening at:", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("rpc server, failed to serve: ", err)
	}
}
