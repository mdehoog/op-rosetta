FROM golang:1.18 as rosetta-builder

WORKDIR /app

ENV REPO=https://github.com/mdehoog/op-rosetta.git

RUN git clone $REPO src \
    && src \
    && make build

FROM ubuntu:20.04

RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates

WORKDIR /app

COPY --from=rosetta-builder /app/bin/op-rosetta /app/op-rosetta

EXPOSE 8080

ENTRYPOINT ["/app/op-rosetta"]
