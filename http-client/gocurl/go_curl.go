package gocurl

import (
    "net/http"
    "net/http/cookiejar"
    "sync"
)

var curSiteCookiesJar, _ = cookiejar.New(nil)
var httpCli = sync.Pool{
    New: func() interface{} {
        return &http.Client{
            Jar: curSiteCookiesJar,
        }
    },
}

func CreateHttpClient(opts ...Options) *Request {
    var hClient = httpCli.Get().(*http.Client)
    defer httpCli.Put(hClient)

    req := &Request{
        cli: hClient,
    }
    if len(opts) > 0 {
        req.opts = mergeDefaultParams(defaultHeader(), opts[0])
    } else {
        req.opts = defaultHeader()
    }
    req.cookiesJar = curSiteCookiesJar
    return req
}
