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

	if err := ctx.Err(); err != nil {
		return nil, err // automatically coded correctly
	}
	if err := validateGreetRequest(req.Msg); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	greeting, err := doGreetWork(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	log.Println("Request headers: ", req.Header())
	fmt.Println(req.Header().Get("Acme-Tenant-Id"))

	res := connect.NewResponse(
		&greetv1.GreetResponse{
			Greeting: greeting,
		})
	res.Header().Set("Greet-Version", "v1")
	return res, nil
}

func validateGreetRequest(msg *greetv1.GreetRequest) error {
	return nil

	// for causing an error intentionally
	// return errors.New("invalid arguments")
}

func doGreetWork(ctx context.Context, msg *greetv1.GreetRequest) (string, error) {
	return fmt.Sprintf("Hello, %s!", msg.Name), nil

	// for causing an error intentionally
	// return "", errors.New("internal error")
}

func main() {
	greeter := &GreetServer{}
	mux := http.NewServeMux()
	// create path and handler from implementation of grpc service
	path, handler := greetv1connect.NewGreetServiceHandler(greeter)
	// register handler with path to mux
	mux.Handle(path, handler)
	http.ListenAndServe(
		"localhost:8080",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
