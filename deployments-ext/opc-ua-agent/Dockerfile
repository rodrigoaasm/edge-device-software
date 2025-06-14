FROM docker.io/golang:1.24.2 AS builder

WORKDIR /workspace

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY cmd/main.go cmd/main.go
COPY config/ config/
COPY internal/ internal/

RUN go build -a -o opc-ua-agent cmd/main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/opc-ua-agent .
USER 65532:65532

ENTRYPOINT ["/opc-ua-agent"]