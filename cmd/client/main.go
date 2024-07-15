package main

import (
	"context"
	"log"

	"github.com/FrostJ143/simplebank/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cc, err := grpc.Dial(
		"localhost:9090",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("can not connect to server: ", err)
	}
	defer cc.Close()

	client := pb.NewSimpleBankClient(cc)
	res, err := client.Sum(context.Background(), &pb.SumRequest{Num1: 5, Num2: 10})
	if err != nil {
		log.Fatal("call service failed")
	}

	log.Print(res.GetResult())
}
