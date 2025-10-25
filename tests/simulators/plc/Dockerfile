FROM docker.io/golang:1.24.2-alpine3.21 AS builder

RUN apk add --no-cache openssl

WORKDIR /workspace
# RUN apk add --no-cache gcc musl-dev sqlite-dev
COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY cmd/main.go cmd/main.go
COPY config/ config/
COPY internal/ internal/

RUN go build -a -o plc cmd/main.go

FROM alpine:3.21
WORKDIR /workspace
COPY --from=builder plc plc
COPY ./entrypoint.sh entrypoint.sh
RUN chmod +x entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]