package message

import (
    "fmt"
    "github.com/tarm/serial"
)

const (
    INT    = 'I'
    UINT   = 'U'
    STRING = '"'
    STRUCT = '{'
)

const (
    END_STRUCT = '}'
    NIL        = '0'
)

type AsnlMsg struct {
    Buffer []byte
}

func NewMessage(bufferSize int) *AsnlMsg {
    msg := new(AsnlMsg)
    msg.Buffer = make([]byte, 0, bufferSize)
    return msg
}

func (msg *AsnlMsg) Length() int {
    return len(msg.Buffer)
}

func (msg *AsnlMsg) FromBytes(src []byte) {
    if len(msg.Buffer) < len(src) {
        msg.Buffer = make([]byte, len(src))
    }
    copy(msg.Buffer, src)
}

func (msg *AsnlMsg) WriteToSerial(c *serial.Port) {
    c.Write(msg.Buffer)
}

func (msg *AsnlMsg) ReadFromSerial(c *serial.Port) error {
    header := make([]byte, 0, 2)
    headerLen := 0;
    for headerLen < 2 {
        buffer := make([]byte, 2 - headerLen)
        n, err := c.Read(buffer)
        if err != nil {
            return err
        }
        headerLen += n
        header = append(header, buffer[0:n]...)
    }
    bodyExpectedLen := int(header[1])
    body := make([]byte, 0, bodyExpectedLen)
    bodyLen := 0
    for bodyLen < bodyExpectedLen {
        buffer := make([]byte, bodyExpectedLen - bodyLen)
        n, err := c.Read(buffer)
        if err != nil {
            return err
        }
        bodyLen += n
        body = append(body, buffer[0:n]...)
    }
    msg.Buffer = append(header, body...)
    return nil
}


func (msg *AsnlMsg) Dump() {
    for i := 0; i < len(msg.Buffer); i++ {
        fmt.Printf("%02x ", msg.Buffer[i])
        if i % 8 == 7 {
            fmt.Println();
        }
    }
    fmt.Println()
}