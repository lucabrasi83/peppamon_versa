FROM golang:1.13.4-alpine as builder
RUN apk add --no-cache build-base git ca-certificates && update-ca-certificates 2>/dev/null || true
COPY . /go/src/github.com/lucabrasi83/peppamon_versa
WORKDIR /go/src/github.com/lucabrasi83/peppamon_versa
ENV GO111MODULE on
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -ldflags="-X github.com/lucabrasi83/peppamon_versa/initializer.Commit=$(git rev-parse --short HEAD) \
    -X github.com/lucabrasi83/peppamon_versa/initializer.Version=$(git describe --tags) \
    -X github.com/lucabrasi83/peppamon_versa/initializer.BuiltAt=$(date +%FT%T%z) \
    -X github.com/lucabrasi83/peppamon_versa/initializer.BuiltOn=$(hostname)" -o peppamon-versa-collector

FROM scratch
LABEL maintainer="sebastien.pouplin@tatacommunications.com"
USER 1001
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/lucabrasi83/peppamon_versa/peppamon-versa-collector /
CMD ["./peppamon-versa-collector"]
