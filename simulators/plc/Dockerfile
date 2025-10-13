FROM docker.io/golang:1.24.2-alpine3.21 AS builder

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
WORKDIR /
COPY --from=builder /workspace/plc plc
COPY ./entrypoint.sh entrypoint.sh
RUN chmod +x entrypoint.sh

RUN apk add --no-cache python3 py3-pip build-base
RUN pip3 install asyncua --break-system-packages

EXPOSE 4840

ENTRYPOINT [ "./entrypoint.sh" ]