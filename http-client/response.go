package http_client

import "net/http"

type Response struct {
    resp   *http.Response
    req    *http.Request
    body   []byte
    stream chan []byte
    err    error
}
