package main

import (
	"context"
	"log"
	"net/http"

	greetv1 "connect-go/gen/greet/v1"
	"connect-go/gen/greet/v1/greetv1connect"

	"github.com/bufbuild/connect-go"
)

func main() {
	ctx := context.Background()

	client := greetv1connect.NewGreetServiceClient(http.DefaultClient, "http://localhost:8080", connect.WithGRPC())
	res, err := client.Greet(ctx, connect.NewRequest(&greetv1.GreetRequest{Name: "Jane"}))
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(res.Msg.Greeting)
}