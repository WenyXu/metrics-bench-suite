FROM golang:1.23 as builder

ENV LANG en_US.utf8
WORKDIR /metrics-loader

COPY . .
RUN make

FROM ubuntu:22.04 as base

WORKDIR /metrics-loader
COPY configs /configs
COPY --from=builder /metrics-loader/bin/sample_loader /usr/local/bin/

ENTRYPOINT ["sample_loader"]
