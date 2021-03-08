package main

import (
	"context"
	pb "github.com/linus18/sandbox/posting_api_grpc"
	"google.golang.org/grpc"
	"log"
	"net"
)

const port = ":50051"

type postingServer struct {
	pb.UnimplementedPostingServer
}

func (server *postingServer) CreatePosting(ctx context.Context, in *pb.PostingRequest) (*pb.PostingReply, error) {
	log.Printf("Received: %v", in)
	return &pb.PostingReply{ResponseCode: "Accepted"}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Printf("gRPC serv er listening on port %s...", port)
	s := grpc.NewServer()
	pb.RegisterPostingServer(s, &postingServer{})
	if err = s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
