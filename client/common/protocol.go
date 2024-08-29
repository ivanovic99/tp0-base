package common

import (
    "bufio"
    "net"
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

    _, err = protocol.conn.Write(data)
    if err != nil {
        return err
    }

    return nil
}

func (protocol *Protocol) ReceiveResponse() (string, error) {
    msg, err := bufio.NewReader(protocol.conn).ReadString('\n')
    if err != nil {
        return "", err
    }
    return msg, nil
}
