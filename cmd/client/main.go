package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/jacob-alt-del/explore-service/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var addr = flag.String("addr", "localhost:50051", "the address to connect to")

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewExploreServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	callListLikedYou(ctx, c)
	callListNewLikedYou(ctx, c)
	callCountLikedYou(ctx, c)
	callPutDecision(ctx, c)
}

func callListLikedYou(ctx context.Context, c pb.ExploreServiceClient) {
	req := pb.ListLikedYouRequest{
		RecipientUserId: "2b13bf3c-b7e3-11f0-add8-627f4e32ceb4",
	}
	resp, err := c.ListLikedYou(ctx, &req)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("ListLikedYou response: %s\n\n", resp)
}

func callListNewLikedYou(ctx context.Context, c pb.ExploreServiceClient) {
	req := pb.ListLikedYouRequest{
		RecipientUserId: "2b13bf3c-b7e3-11f0-add8-627f4e32ceb4",
	}
	resp, err := c.ListNewLikedYou(ctx, &req)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("ListNewLikedYou response: %s\n\n", resp)
}

func callCountLikedYou(ctx context.Context, c pb.ExploreServiceClient) {
	req := pb.CountLikedYouRequest{
		RecipientUserId: "2b13bf3c-b7e3-11f0-add8-627f4e32ceb4",
	}
	resp, err := c.CountLikedYou(ctx, &req)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("CountLikedYou response: %s\n\n", resp)
}

func callPutDecision(ctx context.Context, c pb.ExploreServiceClient) {
	req := pb.PutDecisionRequest{
		ActorUserId:     "2b13bf3c-b7e3-11f0-add8-627f4e32ceb4",
		RecipientUserId: "2b14608b-b7e3-11f0-add8-627f4e32ceb4",
		LikedRecipient:  true,
	}
	resp, err := c.PutDecision(ctx, &req)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("PutDecision response: %s\n\n", resp)
}
