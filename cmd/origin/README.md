# README

This command starts a simple "echo" server on `localhost:8081` for testing
purposes.

The query string of the request can be used to control some specific behavior:

- `expires=<N>`: will add the "Expires" header to the response, expiring
  after `N` seconds.
- `chunked=<N>`: will send `N` chunks of data (chunked transfer encoding) before
  actual echo response.

Examples:
```console
$ curl -i http://localhost:8081/hello
HTTP/1.1 200 OK
...

/hello

$ curl -i http://localhost:8081/hello\?expires\=3
HTTP/1.1 200 OK
...
Expires: Sun, 11 Dec 2022 03:33:38 GMT

/hello

$ curl -i http://localhost:8081/hello\?chunked\=3
HTTP/1.1 200 OK
...
Transfer-Encoding: chunked

Chunk 1
Chunk 2
Chunk 3
/hello
```

