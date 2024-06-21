.PHONY: lint test traefik_bearer_token_plugin clean

export GO111MODULE=on

default: lint test

lint:
	golangci-lint run

test:
	go test -v -cover ./...

yaegi_test:
	yaegi test -v .

vendor:
	go mod traefik_bearer_token_plugin

clean:
	rm -rf ./traefik_bearer_token_plugin