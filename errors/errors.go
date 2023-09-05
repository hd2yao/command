package errors

// wrapperError 满足 error 接口
type wrapperError struct {
    msg    string
    detail []string
    data   map[string]interface{}
    stack  []StackFrame
    root   error
}

func (e wrapperError) Error() string {
    return e.msg
}
