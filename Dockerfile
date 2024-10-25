FROM golang:1.23-alpine

WORKDIR /nms
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ./bin/api ./cmd/api \
    && go build -o ./bin/migrate ./cmd/migrate

EXPOSE 8080
CMD ["/nms/bin/api"]
