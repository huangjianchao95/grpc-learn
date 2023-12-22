package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"time"

	pb "github.com/huangjianchao95/grpc-learn/server/protogen"

	"google.golang.org/grpc"
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

func (server *server) StockPrice(req *pb.StockRequest, stream pb.LearnService_StockPriceServer) error {
	log.Println("StockPrice, stockId: ", req.StockId)
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	for i := 1; i <= 10; i++ {
		res := &pb.StockResponse{
			Price: r.Int31n(10000),
		}
		if err := stream.Send(res); err != nil {
			log.Println("StockPrice send error: ", err)
			return err
		}
		time.Sleep(1 * time.Second)
	}
	return nil
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
