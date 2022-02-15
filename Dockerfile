FROM golang:1.17.6-alpine3.12 as builder
COPY go.mod go.sum /go/src/github.com/OscarClemente/go-noob/
WORKDIR /go/src/github.com/OscarClemente/go-noob/
RUN go mod download
COPY . /go/src/github.com/OscarClemente/go-noob/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/go-noob github.com/OscarClemente/go-noob

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/github.com/OscarClemente/go-noob/build/bucketeer /usr/bin/go-noob
EXPOSE 8080 8080
ENTRYPOINT ["/usr/bin/go-noob"]