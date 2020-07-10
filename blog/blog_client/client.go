package main

import (
	"context"
	"fmt"
	"go-grpc-course/blog/blogpb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Blog Client Started")

	opts := grpc.WithInsecure()
	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	fmt.Println("Creating the blog")
	blog := &blogpb.Blog{
		AuthorId: "Milo",
		Title:    "My First Blog",
		Content:  "Content of the first blog",
	}
	createBlogRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected error: %v\n", err)
	}
	fmt.Printf("Blog has been created: %v\n", createBlogRes)
}
