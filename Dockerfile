FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN go build -o /usermanagement

EXPOSE 8585
ENTRYPOINT ["/usermanagement"]