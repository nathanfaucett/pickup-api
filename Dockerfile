FROM golang:1.22-alpine3.18 as go-builder

RUN apk add build-base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG GOARCH=amd64
ARG GOOS=linux
ARG VERSION=0.1.0

ENV GOARCH=$GOARCH
ENV GOOS=$GOOS
ENV VERSION=$VERSION
ENV CGO_ENABLED=0
RUN go build -ldflags "-s -w -X main.Version=$VERSION -X main.Build=$(date +%s)"

FROM --platform=linux/amd64 alpine:3.18

WORKDIR /opt/app

ARG ENV=prod
ENV ENV=${ENV}

COPY --from=go-builder /app/pickup /opt/app/pickup
COPY --from=node-builder /app/html /opt/app/html

CMD [ "sh", "-c", "/opt/app/pickup" ]