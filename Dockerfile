ARG GO_VERSION=1.14.4

FROM golang:${GO_VERSION}-alpine AS builder

RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

RUN apk add --no-cache ca-certificates

WORKDIR /src

COPY . .

RUN go build -o /duck main.go

FROM alpine AS final

COPY --from=builder /user/group /user/passwd /etc/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /duck /duck

USER nobody:nobody

EXPOSE 8080

ENTRYPOINT ["/duck"]
