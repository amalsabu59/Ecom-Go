.PHONY: start   

# Command to start Air with your configuration file
start:
	air -c cmd/.air.toml    

test:
	go test ./...   

build:
	go build -o yourapp ./cmd/main.go   
