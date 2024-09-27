.PHONY: start   

# Command to start Air with your configuration file
start:
	air -c cmd/.air.toml    

test:
	go test ./...   

build:
	go build -o yourapp ./cmd/main.go   


migrate:
	go run cmd/main.go migrate

rollback:
	go run cmd/main.go rollback
