package common

import (
    "encoding/json"
)

type Bet struct {
    Agency     int    `json:"agency"`
    FirstName  string `json:"first_name"`
    LastName   string `json:"last_name"`
    Document   string `json:"document"`
    Birthdate  string `json:"birthdate"` // Format: 'YYYY-MM-DD'
    Number     int    `json:"number"`
}

func SerializeBet(bet Bet) ([]byte, error) {
    return json.Marshal(bet)
}

func DeserializeBet(data []byte) (Bet, error) {
    var bet Bet
    err := json.Unmarshal(data, &bet)
    return bet, err
}
