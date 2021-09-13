FROM gcr.io/distroless/static
COPY speechly /usr/local/bin/speechly
ENTRYPOINT ["/usr/local/bin/speechly"]
