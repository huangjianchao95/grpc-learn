package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	pb "github.com/huangjianchao95/grpc-learn/client/protogen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:9999", "The addr to connect to")
)

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("client, connect error: ", err)
	}
	defer conn.Close()
	client := pb.NewLearnServiceClient(conn)
	req := &pb.HelloRequest{
		Name: "jack",
	}
	res, err := client.Hello(context.Background(), req)
	if err != nil {
		log.Println("client request error: ", err)
	} else {
		fmt.Println(res.Msg)
	}
}
