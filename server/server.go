package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"io"
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

func (server *server) Add(stream pb.LearnService_AddServer) error {
	res := &pb.AddResponse{
		Sum: 0,
	}
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(res)
		}
		if err != nil {
			log.Println("rpc server, Add error: ", err)
			return err
		}
		for _, num := range req.Nums {
			res.Sum += num
		}
	}
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
