build: clean
	go build -o bin/home-assignment cmd/home-assignment/home-assignment.go

mocks: clear-mocks
	go generate ./...

test: mocks
	go test -v ./...

clear-mocks:
	find . -type d -name "mocks" -exec rm -rf {} +

clean:
	rm -rf bin/*
	rm -rf download/*
