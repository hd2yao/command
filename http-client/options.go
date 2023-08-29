package http_client

import (
    "crypto/tls"
    "strings"
    "time"
)

type Options struct {
    Debug        bool
    BaseURI      string
    Timeout      float32
    timeout      time.Duration
    Query        interface{}
    Headers      map[string]interface{}
    Cookies      interface{}
    FormParams   map[string]interface{}
    JSON         interface{}
    XML          interface{}
    Multipart    []FormData
    Proxy        string
    Certificates []tls.Certificate
}

type FormData struct {
    Name     string
    Contents []byte
    Filename string
    Filepath string
    Headers  map[string]interface{}
}

func mergeOptions(defaultOption Options, opts ...Options) Options {
    for _, opt := range opts {
        if opt.Debug {
            defaultOption.Debug = opt.Debug
        }
        if strings.HasPrefix(opt.BaseURI, "http") {
            defaultOption.BaseURI = opt.BaseURI
        }
        if opt.Timeout > 0 {
            defaultOption.Timeout = opt.Timeout
        }
        if opt.Query != nil {
            defaultOption.Query = opt.Query
        }
        if opt.Headers != nil {
            defaultOption.Headers = opt.Headers
        }
        if opt.Cookies != nil {
            defaultOption.Cookies = opt.Cookies
        }
        if opt.FormParams != nil {
            defaultOption.FormParams = opt.FormParams
        }
        if opt.JSON != nil {
            defaultOption.JSON = opt.JSON
        }
        if opt.XML != nil {
            defaultOption.XML = opt.XML
        }
        if opt.Multipart != nil {
            defaultOption.Multipart = opt.Multipart
        }
        if opt.Proxy != "" {
            defaultOption.Proxy = opt.Proxy
        }
        if opt.Certificates != nil {
            defaultOption.Certificates = opt.Certificates
        }
    }
    return defaultOption
}
