package common

import (
	"bufio"
	"fmt"
	"net"
	"time"

	"github.com/op/go-logging"
	"github.com/ivanovic99/tp0-base/client/common"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
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

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	bet := common.Bet{
        Nombre:     os.Getenv("NOMBRE"),
        Apellido:   os.Getenv("APELLIDO"),
        DNI:        os.Getenv("DOCUMENTO"),
        Nacimiento: os.Getenv("NACIMIENTO"),
        Numero:     os.Getenv("NUMERO"),
    }
	// There is an autoincremental msgID to identify every message sent
	// Messages if the message amount threshold has not been surpassed
	for msgID := 1; msgID <= c.config.LoopAmount; msgID++ {
		select {
        case <-c.stop:
            log.Infof("action: stop_client_loop | result: success | client_id: %v", c.config.ID)
            return
        default:
			// Create the connection the server in every loop iteration. Send an
			if err := c.createClientSocket(); err != nil {
                return
            }
            protocol := common.NewProtocol(c.conn)
            if err := protocol.SendBet(bet); err != nil {
                log.Errorf("action: send_bet | result: fail | client_id: %v | error: %v", c.config.ID, err)
                return
            }

            msg, err := protocol.ReceiveResponse()
            c.conn.Close()

			if err != nil {
				log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
					c.config.ID,
					err,
				)
				return
			}

			log.Infof("action: receive_message | result: success | client_id: %v | msg: %v",
				c.config.ID,
				msg,
			)

			// Wait a time between sending one message and the next one
			time.Sleep(c.config.LoopPeriod)

		}
	}
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}

// StopClientLoop Stops the client loop and closes the connection
func (c *Client) StopClientLoop() {
    close(c.stop)
    if c.conn != nil {
        c.conn.Close()
        log.Infof("action: close_connection | result: success | client_id: %v", c.config.ID)
    }
}
