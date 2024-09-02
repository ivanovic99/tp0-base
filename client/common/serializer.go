package common

import (
    "bytes"
    "encoding/binary"
    "fmt"
)

type Bet struct {
    Agency     int
    FirstName  string
    LastName   string
    Document   string
    Birthdate  string // Format: 'YYYY-MM-DD'
    Number     int
}

func SerializeBet(bet Bet) ([]byte, error) {
    buf := new(bytes.Buffer)

    // Serialize Agency
    if err := binary.Write(buf, binary.BigEndian, int32(bet.Agency)); err != nil {
        return nil, fmt.Errorf("failed to serialize agency: %w", err)
    }

    // Serialize FirstName
    if err := serializeString(buf, bet.FirstName); err != nil {
        return nil, fmt.Errorf("failed to serialize first name: %w", err)
    }

    // Serialize LastName
    if err := serializeString(buf, bet.LastName); err != nil {
        return nil, fmt.Errorf("failed to serialize last name: %w", err)
    }

    // Serialize Document
    if err := serializeString(buf, bet.Document); err != nil {
        return nil, fmt.Errorf("failed to serialize document: %w", err)
    }

    // Serialize Birthdate
    if err := serializeString(buf, bet.Birthdate); err != nil {
        return nil, fmt.Errorf("failed to serialize birthdate: %w", err)
    }

    // Serialize Number
    if err := binary.Write(buf, binary.BigEndian, int32(bet.Number)); err != nil {
        return nil, fmt.Errorf("failed to serialize number: %w", err)
    }

    return buf.Bytes(), nil
}

func serializeString(buf *bytes.Buffer, str string) error {
    length := int32(len(str))
    if err := binary.Write(buf, binary.BigEndian, length); err != nil {
        return fmt.Errorf("failed to write string length: %w", err)
    }
    if _, err := buf.WriteString(str); err != nil {
        return fmt.Errorf("failed to write string: %w", err)
    }
    return nil
}
