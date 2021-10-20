FROM golang:alpine

ENV GO111MODULE=on

WORKDIR /app

COPY . .

RUN go mod download

#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

RUN go build

ENV PORT 50515

EXPOSE 50515

ENTRYPOINT ["/app/vogrpc"]