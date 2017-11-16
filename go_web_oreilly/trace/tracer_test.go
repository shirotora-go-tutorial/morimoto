package trace

import(
	"testing"
	"bytes"
)
func TestNew(t *testing.T){
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil{
		t.Error("return value fron new is nil")
	} else {
		tracer.Trace("Hello Trace Package")
		if buf.String() != "Hello Trace Package\n" {
			t.Errorf("'%s' is miss string was inputed", buf.String())
		}
	}
}

func TestOff(t *testing.T){
	var silentTracer Tracer = Off()
	silentTracer.Trace("Data")
}

