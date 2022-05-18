FROM golang:1.18.2-alpine3.15 AS builder

# hadolint ignore=DL3018
RUN apk update && \
    apk --update --no-cache add git make

WORKDIR /app

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -mod=mod -o robinlb app/cmd/server/main.go

FROM alpine:3.15.4

# hadolint ignore=DL3018
RUN apk --no-cache add ca-certificates

WORKDIR /app

EXPOSE 3030

COPY --from=builder /app/robinlb .

ENTRYPOINT [ "/app/robinlb" ]