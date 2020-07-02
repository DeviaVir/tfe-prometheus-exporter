FROM golang:alpine AS builder

ENV GO111MODULE="on"
ENV CGO_ENABLED="0"

RUN apk add --update git

RUN mkdir -p /go/src/github.com/DeviaVir/tfe-prometheus-exporter

COPY . /go/src/github.com/DeviaVir/tfe-prometheus-exporter

RUN cd /go/src/github.com/DeviaVir/tfe-prometheus-exporter \
 && go mod vendor \
 && go build \
      -mod vendor \
      -o /go/bin/tfe-prometheus-exporter

FROM alpine
COPY --from=builder /go/bin/tfe-prometheus-exporter /usr/local/bin/tfe-prometheus-exporter
CMD ["/usr/local/bin/tfe-prometheus-exporter"]
