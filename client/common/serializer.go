package common

import (
    "encoding/json"
)

type Bet struct {
    Nombre     string `json:"nombre"`
    Apellido   string `json:"apellido"`
    DNI        string `json:"dni"`
    Nacimiento string `json:"nacimiento"`
    Numero     string `json:"numero"`
}

func SerializeBet(bet Bet) ([]byte, error) {
    return json.Marshal(bet)
}

func DeserializeBet(data []byte) (Bet, error) {
    var bet Bet
    err := json.Unmarshal(data, &bet)
    return bet, err
}
