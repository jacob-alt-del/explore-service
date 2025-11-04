package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/jacob-alt-del/explore-service/internal/dataaccess"
	pb "github.com/jacob-alt-del/explore-service/internal/proto"
	"github.com/jacob-alt-del/explore-service/internal/service"
	"google.golang.org/grpc"
)

func main() {
	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "50051"
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "root"
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "secret"
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost:3306"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "explore"
	}

	repo, err := dataaccess.SetupRepository(dbUser, dbPassword, dbHost, dbName)
	if err != nil {
		log.Fatalf("failed to setup database: %v", err)
	}
	defer repo.Close()

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	exploreService := service.NewExploreServiceServer(repo)
	pb.RegisterExploreServiceServer(grpcServer, exploreService)

	log.Printf("server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
