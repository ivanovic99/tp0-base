package common

import (
	"bufio"
	"encoding/csv"
    "net"
    "time"
    "os"
	"strconv"
	"io"

    "github.com/op/go-logging"
)

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

// ProcessBetsFromFile Reads and processes bets from a CSV file in batches
func (c *Client) ProcessBetsFromFile(filePath string, protocol *Protocol) error {
    file, err := os.Open(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    reader := csv.NewReader(bufio.NewReader(file))
    batch := make([]Bet, 0, c.config.BatchMaxAmount)
    totalBets := 0

    for {
        line, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return err
        }

        agency, _ := strconv.Atoi(c.config.ID)
        number, _ := strconv.Atoi(line[4])
        bet := Bet{
            Agency:    agency,
            FirstName: line[0],
            LastName:  line[1],
            Document:  line[2],
            Birthdate: line[3],
            Number:    number,
        }

        batch = append(batch, bet)

        if len(batch) >= c.config.BatchMaxAmount {
            if err := c.sendBatchAndWait(protocol, batch); err != nil {
                return err
            }
            totalBets += len(batch)
            batch = batch[:0]
        }
    }

    // Send any remaining bets
    if len(batch) > 0 {
        if err := c.sendBatchAndWait(protocol, batch); err != nil {
            return err
        }
        totalBets += len(batch)
    }

    log.Infof("action: process_file_complete | result: success | client_id: %v | total_bets_sent: %v", c.config.ID, totalBets)
    return nil
}

func (c *Client) sendBatchAndWait(protocol *Protocol, batch []Bet) error {
    select {
    case <-c.stop:
        return nil
    default:
        if err := protocol.SendBets(batch); err != nil {
            log.Errorf("action: send_bets | result: fail | client_id: %v | error: %v", c.config.ID, err)
            return err
        }
        time.Sleep(c.config.LoopPeriod)
        return nil
    }
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
    if c.config.BatchMaxAmount > 10 {
        log.Infof("action: config | result: adjust_batch_maxAmount | original_value: %v | new_value: 10", c.config.BatchMaxAmount)
        c.config.BatchMaxAmount = 10
    }

    if err := c.createClientSocket(); err != nil {
        return
    }
    protocol := NewProtocol(c.conn)
    defer c.conn.Close()

    filePath := "./app/.data/dataset/agency-" + c.config.ID + ".csv"
    err := c.ProcessBetsFromFile(filePath, protocol)
    if err != nil {
        log.Errorf("action: process_bets | result: fail | client_id: %v | error: %v", c.config.ID, err)
        return
    }

    if err := protocol.SendOk(); err != nil {
        log.Errorf("action: send_ok | result: fail | client_id: %v | error: %v", c.config.ID, err)
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
