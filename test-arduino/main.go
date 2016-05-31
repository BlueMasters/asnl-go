package main

import (
    "fmt"
    "github.com/tarm/serial"
    "github.com/BlueMasters/asnl-go/message"
    "github.com/BlueMasters/asnl-go/writer"
    "github.com/BlueMasters/asnl-go/parser"
    "log"
    "bytes"
    "strconv"
    _ "time"
)

func dumpVars(w *writer.AsnlWriter, s *serial.Port) {
    w.Init()
    w.Struct()
    w.Int(1, 'd')
    w.EndStruct()

    m := w.Msg();
    m.WriteToSerial(s)

    err := m.ReadFromSerial(s)
    if (err != nil) {
        panic(err)
    }

    printMsg(m)
}

func setDi(w *writer.AsnlWriter, s *serial.Port, di int) {
    w.Init()
    w.Struct()
    w.Int(1, 'i')
    w.Int(2, di)
    w.EndStruct()

    m := w.Msg();
    m.WriteToSerial(s)

    err := m.ReadFromSerial(s)
    if (err != nil) {
        panic(err)
    }

    printMsg(m)
}

func setDf(w *writer.AsnlWriter, s *serial.Port, df int) {
    w.Init()
    w.Struct()
    w.Int(1, 'f')
    w.Int(2, df)
    w.EndStruct()

    m := w.Msg();
    m.WriteToSerial(s)

    err := m.ReadFromSerial(s)
    if (err != nil) {
        panic(err)
    }

    printMsg(m)
}

func setPin(w *writer.AsnlWriter, s *serial.Port, pin string) {
    w.Init()
    w.Struct()
    w.Int(1, 'p')
    w.String(pin)
    w.EndStruct()

    m := w.Msg();
    m.WriteToSerial(s)

    err := m.ReadFromSerial(s)
    if (err != nil) {
        panic(err)
    }

    printMsg(m)
}

func exit(w *writer.AsnlWriter, s *serial.Port) {
    w.Init()
    w.Struct()
    w.Int(1, 'x')
    w.EndStruct()

    m := w.Msg();
    m.WriteToSerial(s)

    err := m.ReadFromSerial(s)
    if (err != nil) {
        panic(err)
    }

    printMsg(m)
}


func printMsg(m *message.AsnlMsg) {
    p := parser.NewParser(m)
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
}

func main() {
    c := &serial.Config{Name: "/dev/cu.usbserial-A704BJR9", Baud: 9600}
    buf := make([]byte, 128)
    s, err := serial.OpenPort(c)

    s.Flush()

    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Ready")
    n, err := s.Read(buf)

    if n == 1 && buf[0] == '>' {
        fmt.Println("Configuring...")
        n, _ = s.Write([]byte{'!'})
        var m = message.NewMessage(32)
        var w = writer.NewAsnlWriter(m)

        fmt.Println("Dump:")
        dumpVars(w, s)

        fmt.Println("Set Pin to 123:")
        setPin(w, s, "123")

        fmt.Println("Set DF to 2000:")
        setDf(w, s, 2000)

        fmt.Println("Set DI to 1000:")
        setDi(w, s, 1000)

        fmt.Println("Exit:")
        exit(w, s)
    }

}
