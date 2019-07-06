package tracer

import (
	"fmt"
	"io"
)

// Tracer is an interface that describe an object that capable to trace an event
type Tracer interface {
	Trace(a ...interface{})
}

// tracer implements the Tracer interface
type tracer struct {
	out io.Writer
}

func (t *tracer) Trace(a ...interface{}) {
	fmt.Fprint(t.out, a...) // print a without newline
	fmt.Fprintln(t.out)     // print the newline
}

// New creates new Tracer intance
func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

type nilTracer struct{}

func (t *nilTracer) Trace(a ...interface{}) {}

// Off return Tracer that ignore calls to Trace
func Off() Tracer {
	return &nilTracer{}
}
