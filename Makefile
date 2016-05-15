install:
	@go install

lint:
	@gometalinter

test:
	@go test

coverage:
	@rm -f cover.out cover.html
	@go test -coverprofile cover.out && go tool cover -html cover.out -o cover.html

.PHONY: coverage install test
