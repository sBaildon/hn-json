FROM golang:1.10

WORKDIR /go/src/hn-json
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

ENTRYPOINT ["hn-json"]
