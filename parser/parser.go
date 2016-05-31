package parser

import (
    "github.com/BlueMasters/asnl-go/message"
    "errors"
    "bytes"
)

type AsnlParser struct {
    msg *message.AsnlMsg
    fix int
    pos int
    valPtr int
}

func NewParser(m *message.AsnlMsg) *AsnlParser {
    var p = new(AsnlParser)
    p.msg = m
    return p
}

func (p *AsnlParser) Init() {
    p.fix = 0
    p.pos = 0
    p.valPtr = 0
}

func (p *AsnlParser) NextToken() (token int, err error) {
    err = nil
    if p.fix > 0 && p.pos > (p.fix + int(p.msg.Buffer[p.fix])) {
        // We reached the end of a structure. So we "pop" the address
        // of the previous "struct".
        i := int(p.msg.Buffer[p.fix-1])
        p.msg.Buffer[p.fix-1] = message.STRUCT
        p.fix = i;
        token = message.END_STRUCT;
        return
    }
    if p.pos >= p.msg.Length() {
        token = message.NIL;
        return
    }
    token = int(p.msg.Buffer[p.pos])
    len := int(p.msg.Buffer[p.pos+1])
    p.pos += 2
    p.valPtr = p.pos;
    switch token {
    case message.INT, message.UINT, message.STRING:
        p.pos += len
    case message.STRUCT:
        p.msg.Buffer[p.pos-2] = byte(p.fix)
        p.fix = p.pos - 1
    }
    return
}

func (p *AsnlParser) ReadInt() (result int, err error) {
    if p.valPtr <= 0 || p.valPtr >= p.msg.Length() {
        return 0, errors.New("Invalid valPtr")
    }
    valLen := int(p.msg.Buffer[p.valPtr-1])
    if p.valPtr + valLen > p.msg.Length() {
        return 0, errors.New("Invalid length")
    }
    var v uint32 = 0
    for i := 0; i < valLen; i++ {
        v = v * 256 + uint32(p.msg.Buffer[p.valPtr + i])
    }
    switch valLen {
    case 1:
        return int(int8(v)), nil
    case 2:
        return int(int16(v)), nil
    case 4:
        return int(int32(v)), nil
    }
    return int(v), nil
}

func (p *AsnlParser) ReadUint() (result uint, err error) {
    if p.valPtr <= 0 || p.valPtr >= p.msg.Length() {
        return 0, errors.New("Invalid valPtr")
    }
    valLen := int(p.msg.Buffer[p.valPtr-1])
    if p.valPtr + valLen > p.msg.Length() {
        return 0, errors.New("Invalid length")
    }
    var v uint32 = 0
    for i := 0; i < valLen; i++ {
        v = v * 256 + uint32(p.msg.Buffer[p.valPtr + i])
    }
    return uint(v), nil
}

func (p *AsnlParser) ReadString() (result string, err error) {
    if p.valPtr <= 0 || p.valPtr >= p.msg.Length() {
        return "", errors.New("Invalid valPtr")
    }
    valLen := int(p.msg.Buffer[p.valPtr-1])
    if p.valPtr + valLen > p.msg.Length() {
        return "", errors.New("Invalid length")
    }
    var buffer bytes.Buffer
    for i := 0; i < valLen; i++ {
        buffer.WriteByte(p.msg.Buffer[p.valPtr + i])
    }
    return buffer.String(), nil
}

func (p *AsnlParser) Close() {
    for p.fix > 0 {
        i := int(p.msg.Buffer[p.fix-1])
        p.msg.Buffer[p.fix-1] = message.STRUCT
        p.fix = i;
    }
}

func (p *AsnlParser) FixOk() bool {
    return p.fix == 0
}
