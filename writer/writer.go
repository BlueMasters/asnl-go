package writer

import (
    "github.com/BlueMasters/asnl-go/message"
)

type AsnlWriter struct {
    msg *message.AsnlMsg
    fix int
}

func NewAsnlWriter(msg *message.AsnlMsg) *AsnlWriter {
    var w = new(AsnlWriter)
    w.msg = msg
    w.fix = 0
    return w
}

func (w *AsnlWriter) Init() {
    w.msg.Buffer = w.msg.Buffer[0:0]
    w.fix = 0
}

func (w *AsnlWriter) Int(size int, value int) {
    w.msg.Buffer = append(w.msg.Buffer, message.INT, byte(size))
    x := uint32(value)
    v := make([]byte, size)
    for i := 0; i < size; i++ {
        v[i] = byte(x % 256)
        x = x / 256
    }
    for i := 0; i < size; i++ {
        w.msg.Buffer = append(w.msg.Buffer, v[size-i-1])
    }
}

func (w *AsnlWriter) Uint(size int, value uint) {
    w.msg.Buffer = append(w.msg.Buffer, message.UINT, byte(size))
    x := uint32(value)
    v := make([]byte, size)
    for i := 0; i < size; i++ {
        v[i] = byte(x % 256)
        x = x / 256
    }
    for i := 0; i < size; i++ {
        w.msg.Buffer = append(w.msg.Buffer, v[size-i-1])
    }
}

func (w *AsnlWriter) String(value string) {
    w.msg.Buffer = append(w.msg.Buffer, message.STRING, byte(len(value)))
    for i := 0; i < len(value); i++ {
        w.msg.Buffer = append(w.msg.Buffer, byte(value[i]))
    }
}

func (w *AsnlWriter) Struct() {
    w.msg.Buffer = append(w.msg.Buffer, message.STRUCT, byte(w.fix))
    w.fix = len(w.msg.Buffer) - 1
}

func (w *AsnlWriter) EndStruct() {
    if w.fix > 0 {
        i := int(w.msg.Buffer[w.fix])
        w.msg.Buffer[w.fix] = byte(len(w.msg.Buffer) - 1 - w.fix)
        w.fix = i
    }
}

func (w *AsnlWriter) Close() {
    for w.fix > 0 {
        i := int(w.msg.Buffer[w.fix])
        w.msg.Buffer[w.fix] = byte(len(w.msg.Buffer) - 1 - w.fix)
        w.fix = i
    }
}

func (w *AsnlWriter) FixOk() bool {
    return w.fix == 0;
}

func (w *AsnlWriter) Msg() *message.AsnlMsg {
    return w.msg
}