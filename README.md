# tfe-prometheus-exporter

**WORK IN PROGRESS**

A prometheus exporter for Terraform Enterprise.

## Docker

```
docker run -p 9112:9112 -e TFE_TOKEN=<some-admin-token> -e TFE_HOST=<some-tfe-address> -e HTTP_LISTENADDR=":9112" -it --rm deviavir/tfe-prometheus-exporter:latest
```
