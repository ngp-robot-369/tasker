build:
	go build ./cli

tests:
	go test ./...

static_check:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	GOOS=windows "${GOPATH}/bin/staticcheck" ./...
	GOOS=linux   "${GOPATH}/bin/staticcheck" ./...