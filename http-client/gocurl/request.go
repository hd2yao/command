package gocurl

import (
    "bytes"
    "crypto/tls"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "net/http"
    "net/http/cookiejar"
    "net/url"
    "strings"
    "time"
)

type Request struct {
    opts                 Options
    cli                  *http.Client
    req                  *http.Request
    body                 io.Reader
    subGetFormDataParams string
    cookiesJar           *cookiejar.Jar
}

// Request send request
func (r *Request) Request(method, uri string, opts ...Options) (*Response, error) {
    if len(opts) > 0 {
        r.opts = mergeDefaultParams(defaultHeader(), opts[0], r.opts)
    }
    switch method {
    case http.MethodGet, http.MethodDelete:
        uri = r.opts.BaseURI + uri + r.parseGetFormData()
        req, err := http.NewRequest(method, uri, nil)
        if err != nil {
            return nil, err
        }
        r.req = req
    case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodOptions:
        // parse body
        r.parseBody()
        uri = r.opts.BaseURI + uri + r.parseGetFormData()
        req, err := http.NewRequest(method, uri, r.body)
        if err != nil {
            return nil, err
        }
        r.req = req
    default:
        return nil, errors.New("invalid request method")
    }
    r.opts.Headers["Host"] = fmt.Sprintf("%v", r.req.Host)

    // parseTimeout
    r.parseTimeout()

    // parseClient
    r.parseClient()

    // parse headers
    r.parseHeaders()

    // parse cookies
    r.parseCookies()

    _resp, err := r.cli.Do(r.req)
    resp := &Response{
        resp:          _resp,
        req:           r.req,
        cookiesJar:    r.cookiesJar,
        err:           err,
        setResCharset: r.opts.SetResCharset,
    }
    if err != nil {
        return resp, err
    }
    return resp, nil
}

func (r *Request) parseTimeout() {
    if r.opts.Timeout > 0 {
        r.opts.timeout = time.Duration(r.opts.Timeout*1000) * time.Millisecond
    } else {
        r.opts.Timeout = 0
    }
}

func (r *Request) parseClient() {
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }

    if r.opts.Proxy != "" {
        proxy, err := url.Parse(r.opts.Proxy)
        if err == nil {
            tr.Proxy = http.ProxyURL(proxy)
        } else {
            fmt.Println(r.opts.Proxy+proxyError, err.Error())
        }
    }

    r.cli = &http.Client{
        Timeout:   r.opts.timeout,
        Transport: tr,
        Jar:       r.cookiesJar,
    }
}

func (r *Request) parseCookies() {
    switch r.opts.Cookies.(type) {
    case string:
        cookies := r.opts.Cookies.(string)
        r.req.Header.Add("Cookies", cookies)
    case map[string]string:
        cookies := r.opts.Cookies.(map[string]string)
        for k, v := range cookies {
            if strings.ReplaceAll(v, " ", "") != "" {
                r.req.AddCookie(&http.Cookie{
                    Name:  k,
                    Value: v,
                })
            }
        }
    case []*http.Cookie:
        cookies := r.opts.Cookies.([]*http.Cookie)
        for _, cookie := range cookies {
            if cookie != nil {
                r.req.AddCookie(cookie)
            }
        }
    }
}

func (r *Request) parseHeaders() {
    if r.opts.Headers != nil {
        for k, v := range r.opts.Headers {
            if vv, ok := v.([]string); ok {
                for _, vvv := range vv {
                    if strings.ReplaceAll(vvv, " ", "") != "" {
                        r.req.Header.Add(k, vvv)
                    }
                }
                continue
            }
            vv := fmt.Sprintf("%v", v)
            r.req.Header.Set(k, vv)
        }
    }
}

func (r *Request) parseBody() {
    // application/x-www-form-urlencoded
    if r.opts.FormParams != nil {
        values := url.Values{}
        for k, v := range r.opts.FormParams {
            if vv, ok := v.([]string); ok {
                for _, vvv := range vv {
                    if strings.ReplaceAll(vvv, " ", "") != "" {
                        values.Add(k, vvv)
                    }
                }
                continue
            }
            vv := fmt.Sprintf("%v", v)
            values.Set(k, vv)
        }
        r.body = strings.NewReader(values.Encode())
        return
    }

    // application/json
    if r.opts.JSON != nil {
        b, err := json.Marshal(r.opts.JSON)
        if err == nil {
            r.body = bytes.NewReader(b)
            return
        }
    }

    // text/xml
    if r.opts.XML != "" {
        r.body = strings.NewReader(r.opts.XML)
        return
    }
    return
}

// 解析 get 方式传递的 formData(application/x-www-form-urlencoded)
func (r *Request) parseGetFormData() string {
    if r.opts.FormParams != nil {
        values := url.Values{}
        for k, v := range r.opts.FormParams {
            if vv, ok := v.([]string); ok {
                for _, vvv := range vv {
                    if strings.ReplaceAll(vvv, " ", "") != "" {
                        values.Add(k, vvv)
                    }
                }
                continue
            }
            vv := fmt.Sprintf("%v", v)
            values.Set(k, vv)
        }
        r.subGetFormDataParams = values.Encode()
        return "?" + r.subGetFormDataParams
    } else {
        return ""
    }
}
