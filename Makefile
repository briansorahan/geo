IMG=bsorahan/geo

install:
	go install

image: .image

.image: Dockerfile
	@docker build -t $(IMG) .
	@touch $@

test:
	docker run -w /go/src/github.com/briansorahan/geo $(IMG) gometalinter --deadline=30s
	docker run -w /go/src/github.com/briansorahan/geo $(IMG) go test

coverage:
	rm -f cover.out cover.html
	go test -coverprofile cover.out && go tool cover -html cover.out

clean:
	rm -f .image cover.out cover.html

lint:
	gometalinter --deadline=60s

.PHONY: coverage install test clean
