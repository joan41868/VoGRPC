FROM golang:alpine

ENV GO111MODULE=on

WORKDIR /app

COPY . .

CMD ["CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build"]

ENV PORT 50515

EXPOSE 50515

ENTRYPOINT ["/app/vogrpc"]