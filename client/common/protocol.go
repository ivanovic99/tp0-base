package common

import (
    "net"
    "encoding/binary"

)

type Protocol struct {
    conn net.Conn
}

const AMMOUNT_OF_BYTES = 4

func NewProtocol(conn net.Conn) *Protocol {
    return &Protocol{conn: conn}
}

func (protocol *Protocol) SendBets(bets []Bet) error {
    data, err := SerializeBets(bets)
    if err != nil {
        return err
    }

    length := uint32(len(data))
    lengthBuf := make([]byte, AMMOUNT_OF_BYTES)
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

func (protocol *Protocol) AmountOfBets(totalBets uint32) error {
    // Convert the total number of bets to a byte slice
    totalBetsBytes := make([]byte, AMMOUNT_OF_BYTES)
    binary.BigEndian.PutUint32(totalBetsBytes, totalBets)
    
    // Send the byte slice over the connection
    _, err := protocol.conn.Write(totalBetsBytes)
    return err
}

func (protocol *Protocol) ReceiveOk() (bool, error) {
    // Read the response from the server, only one byte
    response := make([]byte, 1)
    _, err := protocol.conn.Read(response)
    if err != nil {
        return false, err
    }
    if response[0] == 0x01 {
        return true, nil
    } else {
        return false, nil
    }
}
