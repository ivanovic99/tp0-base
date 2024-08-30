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

    _, err = protocol.conn.Write(lengthBuf)
    if err != nil {
        return err
    }

    _, err = protocol.conn.Write(data)
    if err != nil {
        return err
    }

    return nil
}
