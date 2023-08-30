package gocurl

import (
    "io"
    "net/http"
    "net/http/cookiejar"
)

type Request struct {
    opts                 Options
    cli                  *http.Client
    req                  *http.Request
    body                 io.Reader
    subGetFormDataParams string
    cookiesJar           *cookiejar.Jar
}
