# README

A simple HTTP reverse proxy written in Go.

## Usage

This single-host reverse proxy listens on `localhost:8080` and forwards requests to `localhost:8081`. Responses are stored in an in-memory cache if they include an "Expires" header. The cached responses are indexed by the request URL path.

An example origin server is included in `cmd/origin`. It allows you to control the response's `Expires` header based on the request's URL-encoded query string (see [README.md](cmd/origin/README.md)).

## Disclaimer

This is part of a coding challenge to write an HTTP reverse proxy with a caching feature. The only restriction is to not use `net/http/httputil`.

For a full-featured open-source reverse proxy and load balancer you should check [Traefik Proxy](https://traefik.io/traefik/) ;)
