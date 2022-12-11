# README

This command starts a simple origin server listening on "localhost:8081" 
for testing purposes.

The request URL-encoded query string `s=<seconds>` can be used to control 
the "Expires" header in the response. The server will add an "Expires" header
with a date expiring after <seconds>. This is valid on any URL path.

Example:
```
$ curl -i http://localhost:8081/hello\?s\=10
HTTP/1.1 200 OK
Expires: Sun, 11 Dec 2022 03:33:38 GMT
```

