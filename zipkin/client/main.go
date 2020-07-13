package main

import (
	"context"
	greet "demo/zipkin/proto"
	"google.golang.org/grpc"
	"log"
)

func main() {

	conn, err := grpc.Dial(":1234", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc dial error: %v", err)
	}
	defer conn.Close()
	client := greet.NewServiceClient(conn)
	req, err := client.Say(context.Background(), &greet.SayRequest{
		Name: "liuning",
		Msg:  "hello world!",
	})
	if err != nil {
		log.Printf("call say error: %v", err)
		return
	}
	log.Printf("call say response: (code: %d msg: %s)", req.Code, req.Msg)
}
