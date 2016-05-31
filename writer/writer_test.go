package writer

import (
	"testing"
    "github.com/BlueMasters/asnl-go/message"
    "reflect"
)

func TestBasic(t *testing.T) {

    var m = message.NewMessage(32)
    var w = NewAsnlWriter(m)

    w.Init()

	w.Struct()
	w.Int(2, -3)
	w.Struct()
	w.Int(1, 'p')
	w.String("1234")
	w.EndStruct()
    w.Uint(4, 42)
	w.EndStruct()

	m.Dump();

	if !reflect.DeepEqual(m.Buffer, []byte{
		message.STRUCT, 0x15, message.INT, 0x02, 0xff, 0xfd, message.STRUCT, 0x09,
        message.INT, 0x01, 0x70, message.STRING, 0x04, 0x31, 0x32, 0x33,
        0x34, message.UINT, 0x04, 0x00, 0x00, 0x00, 0x2a,
	}) {
		t.Error("Basic Test failed")
	}
}
