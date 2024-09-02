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

func SerializeBets(bets []Bet) ([]byte, error) {
    buf := new(bytes.Buffer)
    for _, bet := range bets {
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
    }
    return buf.Bytes(), nil
}

func DeserializeWinners(data []byte) ([]string, error) {
    buf := bytes.NewReader(data)
    var winners []string

    for buf.Len() > 0 {
        winner, err := deserializeString(buf)
        if err != nil {
            return nil, fmt.Errorf("failed to deserialize winner: %w", err)
        }
        winners = append(winners, winner)
    }

    return winners, nil
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

func deserializeString(buf *bytes.Reader) (string, error) {
    var length int32
    if err := binary.Read(buf, binary.BigEndian, &length); err != nil {
        return "", fmt.Errorf("failed to read string length: %w", err)
    }
    str := make([]byte, length)
    if _, err := buf.Read(str); err != nil {
        return "", fmt.Errorf("failed to read string: %w", err)
    }
    return string(str), nil
}
