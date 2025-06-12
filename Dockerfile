FROM golang:1.23 as builder

ENV LANG en_US.utf8

WORKDIR /metrics-loader

# RUN go env -w  GOPROXY=https://goproxy.cn,direct

COPY . .
RUN make

FROM ubuntu:22.04 as base

RUN apt-get update && apt-get install -y mysql-client && rm -rf /var/lib/apt/lists/*

WORKDIR /metrics-loader
COPY configs /configs
COPY --from=builder /metrics-loader/bin/* /usr/local/bin/

ENTRYPOINT ["sample_loader"]
