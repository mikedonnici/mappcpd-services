## webd/

Contains executable for the web server api.

REST(ish) web services and GraphQL run on the same deployment.

GraphQL is simply answering on the `/graphql` endpoint.

**flags**

`-p` override default port (5000) - env var `PORT` takes precedence over all

```bash
# start on dev machine with default port (5000)
$  go run cmd/webd/main.go

# start on dev machine with custom port
$ go run cmd/webd/main.go - p 8081
```