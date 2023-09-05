package errors

import "fmt"

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
