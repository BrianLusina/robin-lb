FROM golang:1.18 AS builder
WORKDIR /app
COPY main.go go.mod ./
RUN CGO_ENABLED=0 GOOS=linux go build -o robinlb .

FROM alpine:3.15.4

# hadolint ignore=DL3018
RUN apk --no-cache add ca-certificates

WORKDIR /root

ARG CONTAINER_USER_NAME=robinlb-user

# Create non-root user
# hadolint ignore=SC2015
RUN set -xe \
    && addgroup --system ${CONTAINER_USER_NAME} || true \
    && adduser --system --disabled-login --ingroup ${CONTAINER_USER_NAME} --home /home/${CONTAINER_USER_NAME} --gecos "${CONTAINER_USER_NAME} user" --shell /bin/false  ${CONTAINER_USER_NAME} || true

# Use non-root user and start
USER $CONTAINER_USER_NAME

COPY --from=builder /app/lb .
ENTRYPOINT [ "/root/robinlb" ]