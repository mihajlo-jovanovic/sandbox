package main

import (
	"flag"
	"github.com/google/uuid"
	pb "github.com/linus18/sandbox/posting_api_grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"os"
	"strings"

	"github.com/Shopify/sarama"
)

var (
	brokers = flag.String("brokers", os.Getenv("KAFKA_PEERS"), "The Kafka brokers to connect to, as a comma-separated list")
	topic = flag.String("topic", "postings", "Name of the Kafka topic to product to")
)

func main() {
	flag.Parse()

	brokerList := strings.Split(*brokers, ",")
	log.Printf("Kafka brokers: %s", strings.Join(brokerList, ", "))

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}

	req := &pb.PostingRequest{
		Id:          uuid.New().String(),
		PostingDate: timestamppb.Now(),
		Merchant:    "Amazon",
		Amount:      100000,
		IsCredit:    false,
		AccountId:   "1"}

	reqBytes, err := proto.Marshal(req)
	if err != nil {
		log.Fatalf("Failed to marshal requst: %v", err)
	}
	msg := &sarama.ProducerMessage{Topic: *topic, Value: sarama.ByteEncoder(reqBytes)}
	p, o, err := producer.SendMessage(msg)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}
	log.Printf("Message sent successfully (partition: %d, ofset: %d).", p, o)
}
