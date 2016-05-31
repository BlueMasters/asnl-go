package parser

import (
    "testing"
    "github.com/BlueMasters/asnl-go/message"
    "fmt"
    "bytes"
    "strconv"
)

func TestBasic(t *testing.T) {
    msg := message.NewMessage(32)
    msg.FromBytes([]byte{
        message.STRUCT, 0x15, message.INT, 0x02, 0xff, 0xfd, message.STRUCT, 0x09,
        message.INT, 0x01, 0x70, message.STRING, 0x04, 0x31, 0x32, 0x33,
        0x34, message.UINT, 0x04, 0x00, 0x00, 0x00, 0x2a,
    })

    msg.Dump()

    p := NewParser(msg)
    p.Init()

    var r bytes.Buffer

    loop:
    for {
        token, _ := p.NextToken()
        r.WriteString(string(token))
        switch token {
        case message.NIL:
            break loop
        case message.INT:
            v, _ := p.ReadInt()
            r.WriteString(strconv.FormatInt(int64(v), 10))
        case message.UINT:
            v, _ := p.ReadUint()
            r.WriteString(strconv.FormatInt(int64(v), 10))
        case message.STRING:
            v, _ := p.ReadString()
            r.WriteString(string(v))
        }
    }

    fmt.Println(r.String())
    if r.String() != "{I-3{I112\"1234}U42}0" {
        t.Error("Basic Test failed")
    }
}

