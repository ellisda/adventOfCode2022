.PHONY: dep test check build run profile clean

BINARY_NAME="day24"

dep:
	go mod download

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

check:
	go fmt ./...
	go vet ./...

build: dep check test main.go
	go build -o ${BINARY_NAME} ./...

run: build
	./${BINARY_NAME}

profile: build
	./${BINARY_NAME} --cpuprofile=day24.prof

pprof-web: day24.prof
	go tool pprof -http=":" day24.prof

clean:
	go clean
	rm -f ${BINARY_NAME}