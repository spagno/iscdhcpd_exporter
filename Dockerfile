FROM quay.io/prometheus/busybox:latest

COPY iscdhcpd_exporter /bin/iscdhcpd_exporter

ENTRYPOINT ["/bin/iscdhcpd_exporter"]
EXPOSE     9367
