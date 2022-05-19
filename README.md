# Robin LB

[![License](https://img.shields.io/github/license/brianlusina/robin-lb)](https://github.com/brianlusina/robin-lb/blob/main/LICENSE)
[![Version](https://img.shields.io/github/v/release/brianlusina/robin-lb?color=%235351FB&label=version)](https://github.com/brianlusina/robin-lb/releases)
[![Tests](https://github.com/BrianLusina/robin-lb/actions/workflows/tests.yml/badge.svg)](https://github.com/BrianLusina/robin-lb/actions/workflows/tests.yml)
[![Lint](https://github.com/BrianLusina/robin-lb/actions/workflows/lint.yml/badge.svg)](https://github.com/BrianLusina/robin-lb/actions/workflows/lint.yml)
[![Build](https://github.com/BrianLusina/robin-lb/actions/workflows/build_app.yml/badge.svg)](https://github.com/BrianLusina/robin-lb/actions/workflows/build_app.yml)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/b78a2ad36e184c0eb39d5b5bcc721b8b)](https://www.codacy.com/gh/BrianLusina/robin-lb/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=BrianLusina/robin-lb&amp;utm_campaign=Badge_Grade)
[![Go](https://img.shields.io/badge/Go-1.18-blue.svg)](https://go.dev/)

A simple load balancer in Go. It uses RoundRobin algorithm to send requests into set of backends and support
retries too.

It also performs active cleaning and passive recovery for unhealthy backends.

Since its simple it assume if / is reachable for any host its available

## How to use

Run the application with:

```bash
go run app/cmd/server/main.go --backends=<BACKEND_URL>
```

> Ensure that you have a backend running on the specified url BACKEND_URL

Options available are beloe:

```bash
Usage:
  -backends string
        Load balanced backends, use commas to separate
  -port int
        Port to serve (default 3030)
```

Example:

To add followings as load balanced backends

- http://localhost:3031
- http://localhost:3032
- http://localhost:3033
- http://localhost:3034

```bash
robinlb --backends=http://localhost:3031,http://localhost:3032,http://localhost:3033,http://localhost:3034
```

You can optionally run with the provided [docker-compose.yml](./docker-compose.yml).

```bash
docker-compose up
```

> Will build the load balancer and run it on port 3030 and run the other services in the docker-compose file as well.

Afterwards, you can access the load balancer at http://localhost:3030/ and requests will be sent to the appropriate backend.
