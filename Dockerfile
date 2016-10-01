FROM    golang:1.7
RUN	apt-get update
RUN	apt-get install -y libgeos-dev
RUN	go get github.com/alecthomas/gometalinter && gometalinter --install
RUN     mkdir /go/src/github.com/briansorahan
ADD     . /go/src/github.com/briansorahan/geo
WORKDIR /go/src/github.com/briansorahan/geo
RUN     go get -t ./...