.PHONY: tidy fmt test vet build run-dev verify stats

tidy:
	go mod tidy

fmt:
	gofmt -w .

test:
	go test ./...

vet:
	go vet ./...

build:
	go build ./...

run-dev:
	./scripts/run-dev.sh

verify: fmt tidy test vet build
	./scripts/stats.sh

stats:
	./scripts/stats.sh

