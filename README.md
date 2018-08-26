# ISC-DHCPD exporter

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
