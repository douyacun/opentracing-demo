package main

import (
	"context"
	toiletProto "demo/proto"
	"flag"
	"google.golang.org/grpc"
	"log"
)

func main() {
	id := flag.Int("id", 2, "toilet slot id")
	flag.Parse()
	conn, err := grpc.Dial(":1234", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc dial error: %v", err)
	}
	defer conn.Close()
	client := toiletProto.NewToiletClient(conn)
	resp, err := client.Find(context.Background(), &toiletProto.FindRequest{
		Id: int32(*id),
	})
	if err != nil {
		log.Printf("call say error: %v", err)
		return
	}
	log.Printf("%+v", resp)
}
