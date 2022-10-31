FROM golang:1.14-alpine AS builder

WORKDIR /app
COPY . /app

RUN go build -o csp_exporter main.go

FROM alpine:3.16

ENV COLLECTOR_BIND_ADDR=0.0.0.0:8000
ENV PROM_BIND_ADDR=0.0.0.0:9477

COPY --from=builder /app/csp_exporter /usr/bin/
RUN apk add --no-cache tzdata \
 && chmod +x /usr/bin/csp_exporter

USER nobody:nobody

CMD ["csp_exporter"]
