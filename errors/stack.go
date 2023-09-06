package errors

import (
    "fmt"
    "runtime"
)

const stackTraceSize = 10

// StackFrame 表示堆栈跟踪中的单个条目
type StackFrame struct {
    Func string
    File string
    Line int
}

// String 满足 fmt.Stringer 接口
func (f StackFrame) String() string {
    return fmt.Sprintf("%s:%d - %s", f.File, f.Line, f.Func)
}

// Stack 返回错误的堆栈跟踪
// 错误必须包含堆栈跟踪，或封装一个有堆栈跟踪的错误
func Stack(err error) []StackFrame {
    if wErr, ok := err.(wrapperError); ok {
        return wErr.stack
    }
    return nil
}

// getStack 是 runtime.Callers 的格式化包装器
// 以 []StackFrame 的形式返回堆栈跟踪
func getStack(skip int, size int) []StackFrame {
    var (
        pc    = make([]uintptr, size)
        calls = runtime.Callers(skip+1, pc)
        trace []StackFrame
    )

    for i := 0; i < calls; i++ {
        f := runtime.FuncForPC(pc[i])
        file, line := f.FileLine(pc[i] - 1)
        trace = append(trace, StackFrame{
            Func: f.Name(),
            File: file,
            Line: line,
        })
    }

    return trace
}
