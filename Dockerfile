FROM golang:1.12 as build-env
WORKDIR /go/src/github.com/ftob/golang-test-task
COPY . .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/main.go

FROM scratch
WORKDIR /usr/bin

COPY --from=build-env /go/src/github.com/ftob/golang-test-task/app .

EXPOSE 8081

ENTRYPOINT ["./app"]