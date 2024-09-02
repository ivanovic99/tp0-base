package common

import (
	"bufio"
	"encoding/csv"
    "net"
    "time"
    "os"
	"strconv"

    "github.com/op/go-logging"
)

const EIGHT_KB_IN_BATCHES = 10  // 8KB

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
	BatchMaxAmount	 int
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   net.Conn
    stop   chan struct{}
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
		stop:   make(chan struct{}),
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return err
	}
	c.conn = conn
	return nil
}

// ReadBetsFromFile Reads bets from a CSV file
func ReadBetsFromFile(filePath string, agencyID string) ([]Bet, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var bets []Bet
    reader := csv.NewReader(bufio.NewReader(file))
    for {
        line, err := reader.Read()
        if err != nil {
            break
        }
        agency, _ := strconv.Atoi(agencyID)
        number, _ := strconv.Atoi(line[4])
        bet := Bet{
            Agency:    agency,
            FirstName: line[0],
            LastName:  line[1],
            Document:  line[2],
            Birthdate: line[3],
            Number:    number,
        }
        bets = append(bets, bet)
    }
    return bets, nil
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
    bets, err := ReadBetsFromFile("./app/.data/dataset/agency-" + c.config.ID + ".csv", c.config.ID)
    if err != nil {
        log.Errorf("action: read_bets | result: fail | error: %v", err)
        return
    }
    log.Infof("action: read_bets | result: success | total_bets: %v", len(bets))
    if c.config.BatchMaxAmount > EIGHT_KB_IN_BATCHES {
        log.Infof("action: config | result: adjust_batch_maxAmount | original_value: %v | new_value: 10", c.config.BatchMaxAmount)
        c.config.BatchMaxAmount = EIGHT_KB_IN_BATCHES
    }

    if err := c.createClientSocket(); err != nil {
        return
    }
    protocol := NewProtocol(c.conn)
    defer c.conn.Close()

    for bet_number := 0; bet_number < len(bets); bet_number += c.config.BatchMaxAmount {
        end := bet_number + c.config.BatchMaxAmount
        if end > len(bets) {
            end = len(bets)
        }
        batch := bets[bet_number:end]

        select {
        case <-c.stop:
            log.Infof("action: stop_client_loop | result: success | client_id: %v", c.config.ID)
            return
        default:
            
            if err := protocol.SendBets(batch); err != nil {
                log.Errorf("action: send_bets | result: fail | client_id: %v | error: %v", c.config.ID, err)
                return
            }

            log.Infof("action: apuestas_enviadas | result: success | cantidad: %v", len(batch))

            time.Sleep(c.config.LoopPeriod)
        }
    }
    if err := protocol.SendOk(); err != nil {
        log.Errorf("action: send_bets | result: fail | client_id: %v | error: %v", c.config.ID, err)
        return
    }
    ok, err := protocol.ReceiveOk()
    if err != nil {
        log.Errorf("action: receive_ok | result: fail | client_id: %v | error: %v", c.config.ID, err)
        return
    }
    if ok {
        log.Infof("action: receive_ok | result: success | client_id: %v", c.config.ID)
        log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
        } else {
        log.Infof("action: receive_ok | result: success | client_id: %v", c.config.ID)
        log.Infof("action: loop_finished | result: failure | message: Not all batches could be correctly processed | client_id: %v", c.config.ID)
    }

    agencyID, _ := strconv.Atoi(c.config.ID)
    winners, err := protocol.ReceiveWinners(agencyID)
    if err != nil {
        log.Errorf("action: receive_winners | result: fail | client_id: %v | error: %v", c.config.ID, err)
        return
    }
    log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %v", len(winners))

}

// StopClientLoop Stops the client loop and closes the connection
func (c *Client) StopClientLoop() {
    close(c.stop)
    if c.conn != nil {
        c.conn.Close()
        log.Infof("action: close_connection | result: success | client_id: %v", c.config.ID)
    }
}
