package gocurl

import (
    "net/http"
    "net/http/cookiejar"
)

type Response struct {
    resp          *http.Response
    req           *http.Request
    cookiesJar    *cookiejar.Jar
    err           error
    setResCharset string
}
