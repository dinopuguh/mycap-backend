FROM golang:alpine as builder
LABEL maintainer="Dino Puguh <dinopuguh@gmail.com>"
RUN apk update && apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
RUN go build -o main .
RUN go build -o migrate ./cmd

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/migrate .
ENV MYCAP_PORT=3000
EXPOSE 3000
CMD [ "./main", "-migrate" ]