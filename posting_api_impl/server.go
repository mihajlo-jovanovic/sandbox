package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	pb "github.com/linus18/sandbox/posting_api_grpc"
	"google.golang.org/grpc"
	"log"
	"net"

	_ "github.com/lib/pq"
)

const (
	serverport = ":50051"
	//host = "postgres-postgresql"
	host = "localhost"
	port = 5432
	user = "postgres"
	passwd = "H1aCELrAQD"
	dbname = "bank"
)


type postingServer struct {
	pb.UnimplementedPostingServer
}

func (server *postingServer) CreatePosting(ctx context.Context, in *pb.PostingRequest) (*pb.PostingReply, error) {
	log.Printf("Received: %v", in)
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
	type Balance struct {
		amount int64
		acct string
	}
	var bal int64
	row := db.QueryRow("SELECT amount FROM balance WHERE acct = $1", in.AccountId)
	if err := row.Scan(&bal); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("no such account")
		}
		return nil, fmt.Errorf("acct %s: no such account number", in.AccountId)
	}
	fmt.Printf("Current balance: %d\n", bal)
	_, err = db.Exec("UPDATE balance SET amount=$1 WHERE acct=$2", bal-in.Amount, in.AccountId)
	if err != nil {
		return nil, fmt.Errorf("acct %s: error updating balance", in.AccountId)
	}
	return &pb.PostingReply{ResponseCode: "Accepted"}, nil
}

func main() {
	lis, err := net.Listen("tcp", serverport)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Printf("gRPC server listening on port %s...", serverport)
	s := grpc.NewServer()
	pb.RegisterPostingServer(s, &postingServer{})
	if err = s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
