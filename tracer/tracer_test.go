package tracer

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	var buff bytes.Buffer
	tracer := New(&buff)
	if tracer == nil {
		t.Error("New shouldn't return nil")
	} else {
		tracer.Trace("Hello tracer")
		if buff.String() != "Hello tracer\n" {
			t.Errorf("Trace should not write '%s'", buff.String())
		}
	}
}

func TestOff(t *testing.T) {
	var silentTracer Tracer = Off()
	silentTracer.Trace("hello")
}
