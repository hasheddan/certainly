# certainly

`certainly` generates a root certificate, and a specified number of leaf
certificates and private keys. The root certficate can be used to verify that
leaf certificates were signed by the root's private key.

## Status

`certainly` is currently intended for simple testing scenarios and should not be
utilized in any production environments.

## Installation

`certainly` is not currently packaged for distribution, but you may install it
with `go`:

```
go install github.com/hasheddan/certainly
```

Alternatively, you can clone this repository and build from source.

## Usage

Currently, `certainly` requires that both an organization and a number of leaf
certificates to generate are supplied.

```
certainly <org-name> <leaf-cert-count>
```

The following example demonstrates expected output.

```
$ certainly my-cool-org 5
$ ls
dev-0.crt.pem  dev-0.key.pem  dev-1.crt.pem  dev-1.key.pem  dev-2.crt.pem  dev-2.key.pem  dev-3.crt.pem  dev-3.key.pem  dev-4.crt.pem  dev-4.key.pem  root.crt.pem
```
