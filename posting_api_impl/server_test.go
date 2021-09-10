package main

import (
	"database/sql"
	"fmt"
	pb "github.com/linus18/sandbox/posting_api_grpc"
	postgresql "github.com/linus18/sandbox/testhelpers/postgres"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestMyFirstTest(t *testing.T) {
	cleanup, connURL := postgresql.PrepareTestContainer(t, "11.13")
	defer cleanup()

	//let's try it out
	db, _ := sql.Open("postgres", connURL)
	defer db.Close()
	if err := db.Ping(); err != nil {
		t.Fatalf("Unable to ping db")
	}

	create := "CREATE TABLE IF NOT EXISTS postings (id VARCHAR(36), posting_date TIMESTAMP, merchant VARCHAR, amount INT, is_credit BOOL, account_id VARCHAR, PRIMARY KEY (id))"
	_, err := db.Exec(create); if err != nil {
		t.Fatalf("Counld not create table: %v", err)
	}
	posting := &pb.PostingRequest{
		Id: "test_id_123",
		PostingDate: timestamppb.Now(),
		Merchant:    "Starbucks",
		Amount:      725,
		IsCredit:    false,
		AccountId:   "123",
	}
	if err := SavePosting(db, posting); err != nil {
		t.Fatalf("Failed to save posting: %v", err)
	}

	//id := uuid.New()
	//if err != nil {
	//	t.Fatalf("Could not insert data: %v", err)
	//}
	var bal int64
	var posting_dt time.Time
	row := db.QueryRow("SELECT amount, posting_date FROM postings WHERE account_id = $1", "123")
	if err := row.Scan(&bal, &posting_dt); err != nil {
		if err == sql.ErrNoRows {
			t.Fatalf("no such account")
		}
	}
	fmt.Printf("Amount posted: %d posted on %v\n", bal, posting_dt)
	fmt.Println("Worked!")
}

