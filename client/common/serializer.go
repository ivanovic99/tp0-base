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

func SerializeBets(bets []Bet) ([]byte, error) {
    return json.Marshal(bets)
}

func DeserializeWinners(data []byte) ([]string, error) {
    var winners []string
    err := json.Unmarshal(data, &winners)
    return winners, err
}
