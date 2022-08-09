# ISC-DHCPD exporter

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/spagno/iscdhcpd_exporter/tree/master.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/spagno/iscdhcpd_exporter/tree/master)
[![Docker Pulls](https://img.shields.io/docker/pulls/spagno/iscdhcpd-exporter.svg?maxAge=604800)][hub]
[![Go Report Card](https://goreportcard.com/badge/github.com/spagno/iscdhcpd_exporter)][goreportcard]

## Building and running

Prerequisites:

* [Go compiler](https://golang.org/dl/)
* dhcpd-pools installed on the target server

Building:

    go get github.com/spagno/iscdhcpd_exporter
    cd ${GOPATH-$HOME/go}/src/github.com/spagno/iscdhcpd_exporter
    make
    ./iscdhcpd_exporter <flags>

To see all available configuration flags:

    ./iscdhcpd_exporter -h

[hub]: https://hub.docker.com/r/spagno/iscdhcpd-exporter/
[circleci]: https://circleci.com/gh/spagno/iscdhcpd_exporter
[goreportcard]: https://goreportcard.com/report/github.com/spagno/iscdhcpd_exporter
