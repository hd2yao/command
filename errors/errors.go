package errors

import (
    "errors"
    "fmt"
)

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

// Wrap 将上下文信息和堆栈跟踪添加到 err 中，并返回一个带有新上下文的新错误。参数处理与 fmt.Print.Wrap 相同。
// 使用 Root 恢复由一个或多个 Wrap 调用封装的原始错误
// 使用 Stack 恢复堆栈跟踪
// Wrap returns nil if err is nil.
func Wrap(err error, a ...interface{}) error {
    if err == nil {
        return nil
    }
    return wrap(err, fmt.Sprint(a...), 1)
}

// Wrapf 与 Wrap 类型，但参数处理方式与 fmt.Printf 相同。
func Wrapf(err error, format string, a ...interface{}) error {
    if err == nil {
        return nil
    }
    return wrap(err, fmt.Sprintf(format, a...), 1)
}
