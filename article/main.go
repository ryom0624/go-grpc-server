package main

import (
	"context"
	"fmt"
	"go-grpc-server/article/client"
	"go-grpc-server/article/pb"
	"io"
	"log"
)

func main() {
	c, _ := client.NewClient("localhost:50051")

	create(c)
	create(c)
	create(c)
	read(c, 1)
	update(c, 1)
	del(c, 1)
	list(c)
}

func create(c *client.Client) {
	input := &pb.ArticleInput{
		Author:  "gopher",
		Title:   "golang",
		Content: "golang is so cool!",
	}
	res, err := c.Service.CreateArticle(context.Background(), &pb.CreateArticleRequest{ArticleInput: input})
	if err != nil {
		log.Fatalf("Failed to CreateArticle: %v\n", err)
	}
	fmt.Printf("CreateArticle Response: %v\n", res)
}

func read(c *client.Client, id int64) {
	res, err := c.Service.ReadArticle(context.Background(), &pb.ReadArticleRequest{Id: id})
	if err != nil {
		log.Fatalf("Failed to ReadArticle: %v\n", err)
	}
	fmt.Printf("ReadArticle Response: %v\n", res)
}

func update(c *client.Client, id int64) {
	input := &pb.ArticleInput{
		Author:  "gRPC master",
		Title:   "gRPC",
		Content: "gRPC is very smart!",
	}
	res, err := c.Service.UpdateArticle(context.Background(), &pb.UpdateArticleRequest{Id: id, Article: input})
	if err != nil {
		log.Fatalf("Failed to UpdateArticle: %v\n", err)
	}
	fmt.Printf("UpdateArticle Response: %v\n", res)
}

func del(c *client.Client, id int64) {
	res, err := c.Service.DeleteArticle(context.Background(), &pb.DeleteArticleRequest{Id: id})
	if err != nil {
		log.Fatalf("Failed to UpdateArticle: %v\n", err)
	}
	fmt.Printf("The article has been deleted (%v)\n", res)
}

func list(c *client.Client) {
	stream, err := c.Service.ListArticle(context.Background(), &pb.ListArticleRequest{})
	if err != nil {
		log.Fatalf("Failed to ListArticle: %v\n", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Failed to Server Streaming: %v\n", err)
		}
		fmt.Println(res)
	}
}
