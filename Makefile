run: ## run server
	go run ./cmd/server/main.go

request_via_connect-client: ## run this in a separate terminal while server runing
	go run ./cmd/client/main.go

request_via_curl: ## run this in a separate terminal while server runing
	curl \
    --header "Content-Type: application/json" \
    --data '{"name": "Jane"}' \
    http://localhost:8080/greet.v1.GreetService/Greet

request_via_grpcurl: ## run this in a separate terminal while server runing
	grpcurl \
    -protoset <(buf build -o -) -plaintext \
    -d '{"name": "Jane"}' \
    localhost:8080 greet.v1.GreetService/Greet

help: ## show options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

