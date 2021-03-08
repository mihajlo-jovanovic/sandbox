package posting_api_grpc_test

import (
	pb "github.com/linus18/sandbox/posting_api_grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io/ioutil"
	"log"
	"testing"
)

func TestMarshall(t *testing.T) {
	req := &pb.PostingRequest{
		PostingDate: timestamppb.Now(),
		Merchant:    "Starbucks",
		Amount:      725,
		IsCredit:    false,
		AccountId:   "1",
	}
	out, err := proto.Marshal(req)
	if err != nil {
		log.Fatalln("Failed to encode posting request:", err)
	}
	if err := ioutil.WriteFile("req.bin", out, 0644); err != nil {
		log.Fatalln("Failed to write posting request:", err)
	}
}

func TestUnmarshall(t *testing.T) {
	in, err := ioutil.ReadFile("req.bin")
	if err != nil {
		log.Fatalln("Failed to read file:", err)
	}
	req := &pb.PostingRequest{}
	if err := proto.Unmarshal(in, req); err != nil {
		log.Fatalln("Failed to decode request:", err)
	}
	log.Println(req)
}
