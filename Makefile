lint:
	@staticcheck ./...

build:
	@cd cmd/bot && go build -o zneixbot

check: lint
