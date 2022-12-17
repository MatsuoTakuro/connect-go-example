package main

import (
	"context"
	greetv1 "example/gen/greet/v1"
	"example/gen/greet/v1/greetv1connect"
	"fmt"
	"log"
	"net/http"

	"github.com/bufbuild/connect-go"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type GreetServer struct{}

func (s *GreetServer) Greet(
	ctx context.Context,
	req *connect.Request[greetv1.GreetRequest],
) (*connect.Response[greetv1.GreetResponse], error) {

	log.Println("Request headers: ", req.Header())
	res := connect.NewResponse(
		&greetv1.GreetResponse{
			Greeting: fmt.Sprintf("Hello, %s!", req.Msg.Name),
		})
	res.Header().Set("Greet-Version", "v1")
	return res, nil
}

func main() {
	greeter := &GreetServer{}
	api := http.NewServeMux()
	// create path and handler from implementation of grpc service
	path, handler := greetv1connect.NewGreetServiceHandler(greeter)
	// register handler with path to mux
	api.Handle(path, handler)

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "http handler")
	}))
	mux.Handle("/grpc/", http.StripPrefix("/grpc", h2c.NewHandler(api, &http2.Server{})))

	http.ListenAndServe(
		"localhost:8080",
		// Use h2c so we can serve HTTP/2 without TLS.
		mux,
	)
}
