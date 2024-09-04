package common

import (
    "net"
    "encoding/binary"
)

type Protocol struct {
    conn net.Conn
}

func NewProtocol(conn net.Conn) *Protocol {
    return &Protocol{conn: conn}
}

func (protocol *Protocol) SendBet(bet Bet) error {
    data, err := SerializeBet(bet)
    if err != nil {
        return err
    }

    length := uint32(len(data))
    lengthBuf := make([]byte, 4)
    binary.BigEndian.PutUint32(lengthBuf, length)

    if err := protocol.writeAll(lengthBuf); err != nil {
        return err
    }

    if err := protocol.writeAll(data); err != nil {
        return err
    }

    return nil
}

func (protocol *Protocol) writeAll(data []byte) error {
    totalWritten := 0
    for totalWritten < len(data) {
        n, err := protocol.conn.Write(data[totalWritten:])
        if err != nil {
            return err
        }
        totalWritten += n
    }
    return nil
}
