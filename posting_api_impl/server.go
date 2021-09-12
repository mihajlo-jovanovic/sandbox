package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	pb "github.com/linus18/sandbox/posting_api_grpc"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

const (
	serverport = ":50051"
	host       = "postgres-postgresql"
	//host = "localhost"
	port        = 5432
	user        = "postgres"
	passwordKey = "ROOT_DB_PASSWD"
	dbname      = "bank"
)

type postingServer struct {
	pb.UnimplementedPostingServer
	db *sql.DB
}

func SavePosting(db *sql.DB, posting *pb.PostingRequest) error {
	_, err := db.Exec(`insert into postings (id, posting_date, merchant, amount, is_credit, account_id) values($1, $2, $3, $4, $5, $6)`, posting.Id, posting.PostingDate.AsTime(), posting.Merchant, posting.Amount, posting.IsCredit, posting.AccountId)
	return err
}

func (server *postingServer) CreatePosting(ctx context.Context, in *pb.PostingRequest) (*pb.PostingReply, error) {
	log.Printf("Received: %v", in)
	if err := SavePosting(server.db, in); err != nil {
		log.Printf("Error saving posting to db: %v\n", err)
		return &pb.PostingReply{ResponseCode: "Failed"}, nil
	}
	return &pb.PostingReply{ResponseCode: "Accepted"}, nil
}

func main() {
	passwd, isPresent := os.LookupEnv(passwordKey)
	if !isPresent {
		log.Fatalf("Could not find %s in environment - existing...", passwordKey)
	}
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, passwd, dbname)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected!")

	lis, err := net.Listen("tcp", serverport)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("gRPC server listening on port %s...", serverport)
	s := grpc.NewServer()
	pb.RegisterPostingServer(s, &postingServer{db: db})
	if err = s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
