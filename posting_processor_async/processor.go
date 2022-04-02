package main

import (
	"flag"
	"github.com/Shopify/sarama"
	pb "github.com/linus18/sandbox/posting_api_grpc"
	"google.golang.org/protobuf/proto"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	brokers = flag.String("brokers", os.Getenv("KAFKA_PEERS"), "The Kafka brokers to connect to, as a comma-separated list")
	topic   = flag.String("topic", "postings", "Name of the Kafka topic to consume from")
	partition = flag.Int("partition", 0, "Kafka partition to consume from")
)

func main() {
	flag.Parse()

	brokerList := strings.Split(*brokers, ",")
	log.Printf("Kafka brokers: %s", strings.Join(brokerList, ", "))

	config := sarama.NewConfig()
	consumer, err := sarama.NewConsumer(brokerList, config)
	if err != nil {
		log.Fatalf("Could not create Kafka consumer: %v", err)
	}
	partitionConsumer, err := consumer.ConsumePartition(*topic, int32(*partition), sarama.OffsetOldest)
	if err != nil {
		log.Fatalf("Could not create Kafka partition consumer: %v", err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			receivedPosting := &pb.PostingRequest{}
			err := proto.Unmarshal(msg.Value, receivedPosting)
			if err != nil {
				log.Fatalf("Failed to unmarshall posting request: %v", err)
			}
			log.Printf("Posting: %v", receivedPosting)
		case <-signals:
			log.Println("Received termination signal - exiting...")
			return
		}
	}
}
