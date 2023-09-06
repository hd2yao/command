package errors

import "io"

func NewWriter(w io.Writer) *Writer {
    return &Writer{w: w}
}

// Writer 每次调用 Write 时都会对底层 writer 进行一次调用，直到返回错误。
// 从此时，不再调用底层编写器，并且每次都返回相同的错误值。
type Writer struct {
    w   io.Writer
    n   int64
    err error
}

// Write 对底层编写器进行一次调用
func (w *Writer) Write(buf []byte) (n int, err error) {
    if w.err != nil {
        return 0, w.err
    }
    n, w.err = w.w.Write(buf)
    w.n += int64(n)
    return n, w.err
}

func (w *Writer) Err() error {
    return w.err
}

// Written 返回写入底层写入器的字节数
func (w *Writer) Written() int64 {
    return w.n
}
