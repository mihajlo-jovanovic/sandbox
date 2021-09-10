package main

import (
	"context"
	pb "github.com/linus18/sandbox/posting_api_grpc"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"os"
	"time"
)

const TARGET = "TARGET"

func main() {
	target, ok := os.LookupEnv(TARGET)
	if !ok {
		log.Fatalf("Missing %s environment variable", TARGET)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, target, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to dail gRPC server: %v", err)
	}
	defer conn.Close()
	client := pb.NewPostingClient(conn)
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.PostingRequest{
		Id: "b776a2d9-80f9-47de-9a5c-fff67b82cb56",
		PostingDate: timestamppb.Now(),
		Merchant:  "Amazon",
		Amount:    100000,
		IsCredit:  false,
		AccountId: "1"}
	r, err := client.CreatePosting(ctx, req)
	if err != nil {
		log.Fatalf("could not post: %v", err)
	}
	log.Printf("Got back: %s", r.ResponseCode)
}
