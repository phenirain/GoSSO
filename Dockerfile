FROM golang:1.23-alpine AS BUILDER
LABEL authors="eto_ne_ananasbi95"

WORKDIR /grpc

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN go build -o auth-sso ./cmd/GrpcSSO/main.go
RUN go build -o auth-sso-migrator ./cmd/migrator/main.go

FROM golang:1.23-alpine AS RUNNER
WORKDIR /auth-sso

COPY --from=BUILDER /grpc/auth-sso ./auth-sso
COPY --from=BUILDER /grpc/auth-sso-migrator ./auth-sso-migrator
COPY --from=BUILDER /grpc/migrations ./migrations
RUN mkdir config/

CMD ["./auth-sso --config=./config/local.yaml"]
