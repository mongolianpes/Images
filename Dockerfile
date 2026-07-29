FROM golang:1.25.10-alpine AS builder

WORKDIR /images

RUN apk add --no-cache build-base libwebp-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o images ./cmd/images

FROM alpine:latest

WORKDIR /images

RUN apk add --no-cache libwebp

COPY --from=builder /images/images .
COPY --from=builder /images/files ./files

EXPOSE 8080
EXPOSE 8086

CMD ["./images"]