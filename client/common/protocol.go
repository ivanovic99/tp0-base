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
    _, err := protocol.conn.Write([]byte{caseType})
    return err
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

func (protocol *Protocol) SendOk() error {
    if err := protocol._sendCase(OK); err != nil {
        return err
    }
    return nil
}

func (protocol *Protocol) ReceiveOk() (bool, error) {   
    response := make([]byte, 1)
    _, err := protocol.conn.Read(response)
    if err != nil {
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

    _, err := protocol.conn.Write(agencyIDBuf)
    if err != nil {
        return nil, err
    }

    response := make([]byte, AMMOUNT_OF_BYTES)
    _, err = protocol.conn.Read(response)
    if err != nil {
        return nil, err
    }
    length := binary.BigEndian.Uint32(response)

    data := make([]byte, length)
    _, err = protocol.conn.Read(data)
    if err != nil {
        return nil, err
    }
    
    winners, err := DeserializeWinners(data)
    if err != nil {
        return nil, err
    }

    return winners, nil
}
