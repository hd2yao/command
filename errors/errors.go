package errors

import (
    "errors"
    "fmt"
    "strings"
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

// WithDetail 会返回一个新错误，该错误会将 err 包装为包含文本的链式错误信息作为附加上下文
// Detail 在调用新错误值时，Detail 将返回给定文本
func WithDetail(err error, text string) error {
    if err == nil {
        return nil
    }
    if text == "" {
        return err
    }
    e1 := wrap(err, text, 1).(wrapperError)
    e1.detail = append(e1.detail, text)
    return e1
}

// WithDetailf 与WithDetail 类似，但会像 fmt.Printf 一样格式化详细信息
// Detail 在调用新错误值时，Detail 将返回格式化后的文本
func WithDetailf(err error, format string, v ...interface{}) error {
    if err == nil {
        return nil
    }
    text := fmt.Sprintf(format, v...)
    e1 := wrap(err, text, 1).(wrapperError)
    e1.detail = append(e1.detail, text)
    return e1
}

// Detail 返回 err 中包含的详细信息
// 如果错误是由 WithDetail 或 WithDetailf 生成的，则该错误有详细消息
func Detail(err error) string {
    wrapper, ok := err.(wrapperError)
    if !ok {
        return err.Error()
    }
    return strings.Join(wrapper.detail, ";")
}

// withData 会返回一个新错误，该错误会将 err 包装成包含 v 作为额外数据项的链式错误消息。
// 对返回的错误调用 Data 会产生 v.
// 请注意，如果 err 已经有一个数据项，则无法通过返回的错误值访问该数据项。
func withData(err error, v map[string]interface{}) error {
    if err == nil {
        return nil
    }
    e1 := wrap(err, "", 1).(wrapperError)
    e1.data = v
    return e1
}

// WithData 返回一个新错误，它将 err 包装为链式错误消息，其中包含类型为 map[string]interface{} 的值作为额外数据项
// 该 map 包含 err 中 map 的值，以及 keyval 中的值
// 对返回的错误调用 Data 会生成 map
// 请注意，如果 err 已经具有任何其他类型的数据项，则将无法通过返回的错误值访问它。
func WithData(err error, keyval ...interface{}) error {
    if err == nil {
        return nil
    }
    newKV := make(map[string]interface{})
    for k, v := range Data(err) {
        newKV[k] = v
    }
    for i := 0; i < len(keyval); i += 2 {
        newKV[keyval[i].(string)] = keyval[i+1]
    }
    return withData(err, newKV)
}

func Data(err error) map[string]interface{} {
    wrapper, _ := err.(wrapperError)
    return wrapper.data
}

// Sub 将 root 错误作为其主要错误，并从另一个错误 err 中获取所有其他元数据（例如堆栈跟踪、详细信息、消息和数据项）。
// root 或 err 中有任何一个是 nil 时， Sub 返回 nil
// 使用此函数的情况是，当你需要用一个新的根错误替代可能已经包含堆栈跟踪或其他元数据的现有错误时
func Sub(root, err error) error {
    if wrapper, ok := err.(wrapperError); ok && root != nil {
        wrapper.root = Root(root)
        wrapper.msg = root.Error()
        root = wrapper
    }
    if err == nil {
        return nil
    }
    return Wrap(root, err.Error())
}
