FROM alpine:3.12.0 as build_pools

# dhcpd-pools build is based on https://sourceforge.net/projects/dhcpd-pools/files/
ADD https://netix.dl.sourceforge.net/project/dhcpd-pools/dhcpd-pools-3.0.tar.xz /tmp/dhcpd-pools.tar.xz
ADD https://github.com/troydhanson/uthash/archive/v2.1.0.zip /tmp/uthash.zip

    # Install build dependencies
RUN apk add --no-cache build-base=0.5-r2 && \
    # Create target directory and extract dhcpd pools
    mkdir /dhcpd-pools && \
    tar xvf /tmp/dhcpd-pools.tar.xz -C /dhcpd-pools --strip 1 && \
    # Extract UTHash and move it version agnostic location
    unzip /tmp/uthash.zip -d / && \
    mv /uthash-* /uthash

WORKDIR /dhcpd-pools

    # Build dhcpd pools
RUN ./configure --with-uthash=/uthash/include && \
    make && \
    make check && \
    make install

FROM golang:1.14.4 as build

COPY . /go/src/github.com/homecentr/docker-dhcp-exporter

WORKDIR /go/src/github.com/homecentr/docker-dhcp-exporter

RUN make build

FROM homecentr/base:2.4.3-alpine

ENV DHCP_EXPORTER_ARGS="--dhcpd.config-file /config/dhcpd.conf --dhcpd.lease-file /leases/dhcpd.leases"

COPY --from=build_pools /dhcpd-pools/dhcpd-pools /usr/bin/dhcpd-pools
COPY --from=build /go/src/github.com/homecentr/docker-dhcp-exporter/docker-dhcp-exporter /usr/bin/dhcp-exporter

RUN chmod a+x /usr/bin/dhcp-exporter

COPY ./fs/ /

VOLUME "/config"
VOLUME "/leases"

EXPOSE 9367