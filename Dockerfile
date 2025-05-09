FROM golang:1.23 as builder

ENV LANG en_US.utf8

WORKDIR /metrics-loader

# RUN go env -w  GOPROXY=https://goproxy.cn,direct

COPY . .
RUN make

FROM ubuntu:22.04 as base

WORKDIR /metrics-loader
COPY configs /configs
COPY --from=builder /metrics-loader/bin/* /usr/local/bin/

ENTRYPOINT ["sample_loader"]
