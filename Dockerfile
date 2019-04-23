# Processing service Go builder
FROM golang:1.12 as build
WORKDIR /go/src/github.com/SharperShape/vanadia
COPY . .
RUN bin/setup.sh

# Server container
FROM alpine:latest as server
WORKDIR /data
COPY --from=build \
    /go/src/github.com/SharperShape/vanadia/vanadia \
    /opt/bin/vanadia
RUN addgroup -S app && \
    adduser -S -G app app && \
    chmod +x /opt/bin/vanadia && \
    chown app:app /data
USER app:app
ENTRYPOINT [ "/opt/bin/vanadia" ]
