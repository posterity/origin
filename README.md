[![GoDoc reference](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/posterity/origin)

# Origin

Package origin provides simple tools and methods to compare and verify
the `Origin` header of a request on the server-side, specifically in the
context of a CORS request.

It supports simple wildcard pattern matching, and handles omitted port numbers
for the most common web protocols.

## Usage

```go
import (
  "github.com/posterity/origin"
)

// Trusted origins:
//  - example.com and its subdomains over HTTPS on port 443 (implicit);
//  - localhost on any scheme and any port.
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
