package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/huangjianchao95/grpc-learn/client/protogen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:9999", "The addr to connect to")
)

func unaryRpc(conn *grpc.ClientConn) {
	client := pb.NewLearnServiceClient(conn)
	req := &pb.HelloRequest{
		Name: "jack",
	}
	res, err := client.Hello(context.Background(), req)
	if err != nil {
		log.Println("client request error: ", err)
	} else {
		log.Println(res.Msg)
	}
}

func clientStreamRpc(conn *grpc.ClientConn) {
	client := pb.NewLearnServiceClient(conn)
	stream, err := client.Add(context.TODO())
	if err != nil {
		log.Fatalln("error while calling Add: ", err)
	}

	requests := []*pb.AddRequest{
		&pb.AddRequest{
			Nums: []int32{1, 2, 3, 4, 5},
		},
		&pb.AddRequest{
			Nums: []int32{6, 7, 8, 9, 10},
		},
	}
	for _, req := range requests {
		stream.Send(req)
		time.Sleep(1 * time.Second)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalln("error while receiving response from Add: ", err)
	}
	log.Println("Add response, sum: ", res.Sum)

}

func serverStreamRpc(conn *grpc.ClientConn) {
	client := pb.NewLearnServiceClient(conn)
	ctx := context.Background()
	req := &pb.StockRequest{
		StockId: 1,
	}
	stream, err := client.StockPrice(ctx, req)
	if err != nil {
		log.Fatalln("StockPrice stream err: ", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln("failed to receive stock price: ", err)
		}
		log.Printf("stock price, id: %d, price: %d\n", req.StockId, res.Price)
	}
}

func streamRpc(conn *grpc.ClientConn) {
	client := pb.NewLearnServiceClient(conn)
	waitCh := make(chan struct{})

	stream, err := client.Chat(context.Background())
	if err != nil {
		log.Fatalln("chat open stream err: ", err)
	}
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				close(waitCh)
				break
			}
			if err != nil {
				log.Fatalln("receive err: ", err)
			}
			log.Println("chat got message: ", res.Msg)
		}
	}()
	for i := 1; i <= 10; i++ {
		req := &pb.ChatRequest{
			Msg: fmt.Sprintf("chat num %d", i),
		}
		if err := stream.Send(req); err != nil {
			log.Fatalln("chat send err: ", err)
		}
		log.Println("send msg: ", req.Msg)
		time.Sleep(100 * time.Millisecond)
	}
	stream.CloseSend()
	<-waitCh
}

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("client, connect error: ", err)
	}
	defer conn.Close()
	unaryRpc(conn)
	clientStreamRpc(conn)
	serverStreamRpc(conn)
	streamRpc(conn)
}
