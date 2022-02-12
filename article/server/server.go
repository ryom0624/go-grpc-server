package main

import (
	"go-grpc-server/article/pb"
	"go-grpc-server/article/repository"
	"go-grpc-server/article/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen %v\n", err)
	}

	defer lis.Close()

	repo, err := repository.NewsqliteRepo()
	if err != nil {
		log.Fatalf("Failed to create sqlite repository %v", err)
	}
	s := service.NewService(repo)

	server := grpc.NewServer()
	pb.RegisterArticleServiceServer(server, s)

	log.Println("Listening on port 50051...")

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to server %v", err)
	}
}
