.PHONY: format format-check vet test check build

format:
	gofmt -w $$(find . -name '*.go' -not -path './.git/*')

format-check:
	@test -z "$$(gofmt -l $$(find . -name '*.go' -not -path './.git/*'))" || (gofmt -l $$(find . -name '*.go' -not -path './.git/*'); exit 1)

vet:
	go vet ./...

test:
	go test -race -count=1 ./...
	go test -race -count=1 ./internal/tenancy -rapid.checks=200

build:
	go build ./cmd/api ./cmd/migrate

check: format-check vet test build
