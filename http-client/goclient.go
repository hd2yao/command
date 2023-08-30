package http_client

func NewClient(opts ...Options) *Request {
    req := &Request{}
    defaultOpt := Options{}
    if len(opts) > 0 {
        defaultOpt = opts[0]
    }
    req.SetOptions(defaultOpt)
    return req
}

func Get(uri string, opts ...Options) (*Response, error) {
    r := NewClient()
    return r.Request("GET", uri, opts...)
}

func Post(uri string, opts ...Options) (*Response, error) {
    r := NewClient()
    return r.Request("POST", uri, opts...)
}

func Put(uri string, opts ...Options) (*Response, error) {
    r := NewClient()
    return r.Request("PUT", uri, opts...)
}

func Patch(uri string, opts ...Options) (*Response, error) {
    r := NewClient()
    return r.Request("PATCH", uri, opts...)
}

func Delete(uri string, opts ...Options) (*Response, error) {
    r := NewClient()
    return r.Request("DELETE", uri, opts...)
}
