# terraform-provider-sshtunnel

## Overview

This provider create a portforwarding tunnel to a backend target host via a jumphost using AWS instance connect.

The provider allows and will setup the tunnel on terraform state such as :
- terraform plan
- terraform refresh
- terraform apply

The project can be use together with 
- [winebarrel/mysql](https://registry.terraform.io/providers/winebarrel/mysql/latest)
- [cyrilgdn/postgresql](https://registry.terraform.io/providers/cyrilgdn/postgresql/latest)

to form terrraform self-service automation, managing users credential & schemas 

## Pre-requisites
- **ACL/Security Group**: The platform where the terraform using this provider shall be executed upon,needs to have network access to the jumphost.
- **AWS PROFILE**: The platform needs relevant AWS CREDS/ACCESS, using aws profile, that have permission to perform AWS INSTANCE CONNECT  to the jumphost.

## Other Notes

### libc
Hashicorp's Terraform Docker image are shipped with libc-musl and not glibc. Therefore the provider must be compiled with libc-musl instead of glibc. This can be ensure by setting `CGO_ENABLED=0` as shown in the [Dockefile](./Dockerfile) user for local development purposes.

### Time Sleep 
There are 2 sleep timers implemented .

- The [first](./connect/connect.go#L71), implements a 2 sec delay to allow instance connect to propogate the signature to the jumphost before establishing the ssh tunnel.

- The [second](https://github.com/gabsg/terraform-provider-sshtunnel/blob/master/sshtunnel/data_source_port.go#L32), implements a 3 sec delay, before the datasource return the listening port of the tunnel. This ensure the portforwarding is properly setup. Before a next provider that depends on this information starts to use it as an endpoint.

As the portforwarding and tunnel is executed over goroutine. The delay is a cost for each tunnel established and will be run in parallel depending how Terraform executes the setup. 



## How to use. 
See [example.](example-tf/main.tf)

