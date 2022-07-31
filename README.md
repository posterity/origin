[![GoDoc reference](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/posterity/origin)

# Origin

Package origin provides simple tools and methods to compare and verify
the `Origin` header of a request on the server-side, specifically in the
context of Cross-Origin Resource Sharing (CORS).

It supports simple wildcard pattern-matching, and handles omitted port numbers
for the most common web protocols.

## Patterns

The patterns to be checked must be formatted as following:

```text
scheme://hostname:port
```

A wildcard `*` is valid in any position, `scheme`, `hostname` or `port`
(e.g. `*://example.com:*`).

`port` can be omitted if `scheme` is a common web protocol. The value
will default to the standard port associated with it (e.g. `443` for `HTTPS`).

`hostname` can contain multiple wildcards to target subdomains. For example,
`*.*.example.com` will match any sub-subdomain of `example.com`.

`*` is a valid pattern value, and is the equivalent of `*://*:*`.

## Usage

### Single pattern

```go
import (
  "fmt"

  "github.com/posterity/origin"
)

func Main() {
  ok, err := origin.Match("https://subdomain.example.com:443", "https://*.example.com")
  if err != nil {
    panic(err) // Either the origin or the pattern is mis-formatted.
  }
  fmt.Println("is is a match? %v", ok)
}
```

### List of patterns

```go
import (
  "io"

  "github.com/posterity/origin"
)

var patterns = origin.Patterns{
  "https://example.com",
  "https://*.example.com",
  "*://localhost:*",
}

func handler(w http.ResponseWriter, r *http.Request) {
  ok, err := patterns.Match(origin.Get(r))
  if err != nil {
    panic(err) // Either the origin or the pattern is mis-formatted.
  }
  if !ok {
    w.WriteHeader(401)
    io.WriteString(w, "This request is not from a trusted origin")
    return
  }

  io.WriteString(w, "Hello, World!")
}
```

## Contributions

Contributions are welcome via Pull Requests.

## About us

What if you're hit by a bus tomorrow? [Posterity](https://posterity.life) helps
you make a plan in the event something happens to you.
