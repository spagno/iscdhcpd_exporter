FROM golang as build

ADD . /go/src/github.com/homecentr/iscdhcpd_exporter

RUN apt-get update \
    && apt-get install dhcpd-pools \ 
    && cd /go/src/github.com/homecentr/iscdhcpd_exporter \
    && make build

FROM alpine

COPY --from=build /go/src/github.com/homecentr/iscdhcpd_exporter/iscdhcpd_exporter /usr/bin/iscdhcpd_exporter

ENTRYPOINT ["iscdhcpd_exporter"]
EXPOSE     9367