package http_client

import (
    "io"
    "net/http"
)

type Request struct {
    opts Options
    cli  *http.Client
    req  *http.Request
    body io.Reader
}
