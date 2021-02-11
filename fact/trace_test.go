package fact

import (
	"github.com/golang/protobuf/proto"
	"testing"
)

func TestTrace(t *testing.T) {
	trace := NewTrace()
	_, err := proto.Marshal(&trace)
	if err != nil{
		t.Fatal(err)
	}
}
