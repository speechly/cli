FROM golang:1.16 AS builder
WORKDIR /src/
COPY . .
RUN make -e CGO_ENABLED=0 build

FROM gcr.io/distroless/static
COPY --from=builder /src/bin/speechly /usr/local/bin/speechly
ENTRYPOINT ["/usr/local/bin/speechly"]
