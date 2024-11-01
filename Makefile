generate:
	go generate ./...

run:
	go run ./example

install:
	cd cmd/i18n && go install .