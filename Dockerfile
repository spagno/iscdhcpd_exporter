ARG ARCH="amd64"
ARG OS="linux"
FROM quay.io/prometheus/busybox-${OS}-${ARCH}:latest
LABEL maintainer="Spagno <spagno@gmail.com>"

ARG ARCH="amd64"
ARG OS="linux"
COPY .build/${OS}-${ARCH}/iscdhcpd_exporter /bin/iscdhcpd_exporter

ENTRYPOINT ["/bin/iscdhcpd_exporter"]
EXPOSE     9367
