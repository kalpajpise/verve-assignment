FROM golang:1.22-bullseye AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

FROM build AS build-production

RUN useradd -u 1001 runner

COPY . .

RUN go build \
    -ldflags="-linkmode external -extldflags -static" \
    -tags netgo \
    -o server cmd/main.go

FROM scratch

WORKDIR /

COPY --from=build-production /etc/passwd /etc/passwd
COPY --from=build-production /app/server server
COPY --from=build-production /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER runner

EXPOSE 8080

CMD ["/server"]
