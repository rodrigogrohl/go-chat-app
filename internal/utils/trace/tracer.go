package trace

import (
	"fmt"
	"io"
)

type Tracer interface {
	Trace(...interface{})
}

type tracer struct {
	out io.Writer
}

type nilTracer struct {}
func (t *nilTracer) Trace(i ...interface{}) {}

func (t *tracer) Trace(i ...interface{}) {
	_, _ = fmt.Fprint(t.out, i...)
	_, _ = fmt.Fprintln(t.out)
}

func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

func Off() Tracer {
	return &nilTracer{}
}