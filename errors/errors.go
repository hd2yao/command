package errors

import "errors"

func New(text string) error {
    return errors.New(text)
}

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

// Root 返回被一次或多次调用 Wrap 所包裹的原始错误
// 如果 e 没有封装其他错误，则将按原样返回
func Root(e error) error {
    if wErr, ok := e.(wrapperError); ok {
        return wErr.root
    }
    return e
}

// wrap 将上下文信息和堆栈跟踪添加到 err 中，并返回一个包含新上下文的新错误。该函数用于与其他导出函数（如 Wrap 和 WithDetail）一起使用。
// 参数 stackSkip 是生成堆栈跟踪时要上升的堆栈帧数，其中 0 表示 wrap 的调用者。
func wrap(err error, msg string, stackSkip int) error {
    if err == nil {
        return nil
    }

    wErr, ok := err.(wrapperError)
    if !ok {
        wErr.root = err
        wErr.msg = err.Error()
        wErr.stack = getStack(stackSkip+2, stackTraceSize)
    }
    if msg != "" {
        wErr.msg = msg + ": " + wErr.msg
    }

    return wErr
}
