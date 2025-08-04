test:
	go test ./... -v

run:
	go run main.go
	
bench:
	go test -bench='.' ./...

cover:
	go test ./... -coverprofile coverage.out
	go tool cover -html coverage.out -o coverage.html