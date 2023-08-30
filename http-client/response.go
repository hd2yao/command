package http_client

import (
    "fmt"
    "net"
    "net/http"
    "strings"

    "github.com/launchdarkly/eventsource"
    "github.com/tidwall/gjson"
)

type Response struct {
    resp   *http.Response
    req    *http.Request
    body   []byte
    stream chan []byte
    err    error
}

type ResponseBody []byte

func (r ResponseBody) String() string {
    return string(r)
}

// Read get slice of response body
func (r ResponseBody) Read(length int) []byte {
    if length > len(r) {
        length = len(r)
    }
    return r[:length]
}

func (r ResponseBody) GetContents() string {
    return string(r)
}

func (r *Response) GetRequest() *http.Request {
    return r.req
}

func (r *Response) GetBody() (ResponseBody, error) {
    return r.body, r.err
}

func (r *Response) GetParseBody() (*gjson.Result, error) {
    pb := gjson.ParseBytes(r.body)
    return &pb, nil
}

func (r *Response) GetStatusCode() int {
    return r.resp.StatusCode
}

func (r *Response) GetReasonPhrase() string {
    status := r.resp.Status
    arr := strings.Split(status, " ")
    return arr[1]
}

func (r *Response) IsTimeout() bool {
    if r.err == nil {
        return false
    }
    netErr, ok := r.err.(net.Error)
    if !ok {
        return false
    }
    if netErr.Timeout() {
        return true
    }
    return false
}

func (r *Response) GetHeaders() map[string][]string {
    return r.resp.Header
}

func (r *Response) GetHeader(name string) []string {
    headers := r.GetHeaders()
    for k, v := range headers {
        if strings.ToLower(name) == strings.ToLower(k) {
            return v
        }
    }
    return nil
}

func (r *Response) GetHeaderLine(name string) string {
    header := r.GetHeader(name)
    if len(header) > 0 {
        return header[0]
    }
    return ""
}

func (r *Response) HasHeader(name string) bool {
    headers := r.GetHeaders()
    for k := range headers {
        if strings.ToLower(name) == strings.ToLower(k) {
            return true
        }
    }
    return false
}

func (r *Response) Err() error {
    return r.err
}

func (r *Response) Stream() chan []byte {
    return r.stream
}

func (r *Response) parseSteam() {
    r.stream = make(chan []byte)
    decoder := eventsource.NewDecoder(r.resp.Body)

    go func() {
        defer r.resp.Body.Close()
        defer close(r.stream)

        for {
            event, err := decoder.Decode()
            if err != nil {
                r.err = fmt.Errorf("decode data failed: %v", err)
                return
            }

            data := event.Data()
            if data == "" || data == "[DONE]" {
                // read data finished, success return
                return
            }

            r.stream <- []byte(data)
        }
    }()
}
