FROM docker.mirrors.alibaba.ir/library/golang:alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    ENV=PRODUCTION

WORKDIR /src

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go test ./... -v

WORKDIR /src/cmd
RUN go build -o main .

WORKDIR /dist

RUN cp /src/cmd/main .
RUN cp /src/*.env .

EXPOSE 8080

CMD ["/dist/main"]