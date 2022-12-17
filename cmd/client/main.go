package main

import (
	"context"
	"errors"
	greetv1 "example/gen/greet/v1"
	"example/gen/greet/v1/greetv1connect"
	"fmt"
	"log"
	"net/http"

	"github.com/bufbuild/connect-go"
)

func main() {
	client := greetv1connect.NewGreetServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
		// connect.WithGRPC(),
	)
	req := connect.NewRequest(&greetv1.GreetRequest{Name: "Jane"})
	req.Header().Set("Acme-Tenant-Id", "1234")

	res, err := client.Greet(
		context.Background(),
		req,
	)
	if err != nil {
		log.Println(err)
		fmt.Println(connect.CodeOf(err))
		if connectErr := new(connect.Error); errors.As(err, &connectErr) {
			fmt.Println("connectErr.Message: ", connectErr.Message())
			fmt.Println("connectErr.Details: ", connectErr.Details())
		}
		return
	}
	log.Println(res.Msg.Greeting)
}
