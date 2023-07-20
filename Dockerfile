FROM golang:1.20 AS builder
COPY . /build
WORKDIR /build
RUN make test-api
RUN make build-api
RUN make build-tests

FROM alpine:latest
RUN apk --no-cache add ca-certificates make

WORKDIR /app
COPY --from=builder /build/littlejohn .

RUN mkdir /app/tests
COPY --from=builder /build/integration_tests /app/tests

RUN ls -la

EXPOSE 8080

CMD ["./littlejohn"]
