[![Project status](https://badgen.net/badge/project%20status/stable%20%26%20actively%20maintaned?color=green)](https://github.com/homecentr/docker-dhcp-exporter/graphs/commit-activity) [![](https://badgen.net/github/label-issues/homecentr/docker-dhcp-exporter/bug?label=open%20bugs&color=green)](https://github.com/homecentr/docker-dhcp-exporter/labels/bug) [![](https://badgen.net/github/release/homecentr/docker-dhcp-exporter)](https://hub.docker.com/repository/docker/homecentr/dhcp-exporter)
[![](https://badgen.net/docker/pulls/homecentr/dhcp-exporter)](https://hub.docker.com/repository/docker/homecentr/dhcp-exporter) 
[![](https://badgen.net/docker/size/homecentr/dhcp-exporter)](https://hub.docker.com/repository/docker/homecentr/dhcp-exporter)

![CI/CD on master](https://github.com/homecentr/docker-dhcp-exporter/workflows/CI/CD%20on%20master/badge.svg)
![Regular Docker image vulnerability scan](https://github.com/homecentr/docker-dhcp-exporter/workflows/Regular%20Docker%20image%20vulnerability%20scan/badge.svg)


# HomeCentr - ISC DHCP Prometheus Exporter

## Usage

```yml
version: "3.7"
services:
  dhcp_exporter:
    image: homecentr/dhcp
    volumes:
        - ./example/config:/config # Make sure both containers share the same configuration
        - dhcp_leases:/leases:ro   # And leases directory

  dhcp_exporter:
    image: homecentr/dhcp-exporter
    volumes:
        - ./example/config:/config
        - dhcp_leases:/leases
    
volumes:
  dhcp_leases:
```

## Environment variables

| Name | Default value | Description |
|------|---------------|-------------|
| PUID | 7077 | UID of the user dhcp-exporter should be running as. |
| PGID | 7077 | GID of the user dhcp-exporter should be running as. |
| DHCP_EXPORTER_ARGS | `--dhcpd.config-file /config/dhcpd.conf --dhcpd.lease-file /leases/dhcpd.leases` | If you mount the configuration file or the lease file to different location, you need to adjust the arguments accordingly. |

## Exposed ports

| Port | Protocol | Description |
|------|------|-------------|
| 80 | TCP | Some useful details |

## Volumes

| Container path | Description |
|------------|---------------|
| /config | Some useful details |

## Security
The container is regularly scanned for vulnerabilities and updated. Further info can be found in the [Security tab](https://github.com/homecentr/docker-dhcp-exporter/security).

### Container user
The container supports privilege drop. Even though the container starts as root, it will use the permissions only to perform the initial set up. The dhcp-exporter process runs as UID/GID provided in the PUID and PGID environment variables.

:warning: Do not change the container user directly using the `user` Docker compose property or using the `--user` argument. This would break the privilege drop logic.