package common

import (
    "net"
    "encoding/binary"
)

type Protocol struct {
    conn net.Conn
}

const AMMOUNT_OF_BYTES = 4
const BETS = 0x01
const OK = 0x02
const WINNERS = 0x03

func NewProtocol(conn net.Conn) *Protocol {
    return &Protocol{conn: conn}
}

func (protocol *Protocol) _sendCase(caseType byte) error {
    return protocol.writeAll([]byte{caseType})
}

func (protocol *Protocol) SendBets(bets []Bet) error {
    if err := protocol._sendCase(BETS); err != nil {
        return err
    }

    data, err := SerializeBets(bets)
    if err != nil {
        return err
    }

    length := uint32(len(data))
    lengthBuf := make([]byte, AMMOUNT_OF_BYTES)
    binary.BigEndian.PutUint32(lengthBuf, length)

    if err := protocol.writeAll(lengthBuf); err != nil {
        return err
    }

    if err := protocol.writeAll(data); err != nil {
        return err
    }

    return nil
}

func (protocol *Protocol) SendOk() error {
    return protocol._sendCase(OK)
}

func (protocol *Protocol) ReceiveOk() (bool, error) {
    response := make([]byte, 1)
    if err := protocol.readAll(response); err != nil {
        return false, err
    }
    return response[0] == OK, nil
}

func (protocol *Protocol) ReceiveWinners(agencyID int) ([]string, error) {
    if err := protocol._sendCase(WINNERS); err != nil {
        return nil, err
    }

    agencyIDBuf := make([]byte, AMMOUNT_OF_BYTES)
    binary.BigEndian.PutUint32(agencyIDBuf, uint32(agencyID))

    if err := protocol.writeAll(agencyIDBuf); err != nil {
        return nil, err
    }

    response := make([]byte, AMMOUNT_OF_BYTES)
    if err := protocol.readAll(response); err != nil {
        return nil, err
    }
    length := binary.BigEndian.Uint32(response)

    data := make([]byte, length)
    if err := protocol.readAll(data); err != nil {
        return nil, err
    }

    winners, err := DeserializeWinners(data)
    if err != nil {
        return nil, err
    }

    return winners, nil
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

func (protocol *Protocol) readAll(data []byte) error {
    totalRead := 0
    for totalRead < len(data) {
        n, err := protocol.conn.Read(data[totalRead:])
        if err != nil {
            return err
        }
        totalRead += n
    }
    return nil
}
