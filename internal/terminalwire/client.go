package terminalwire

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/gorilla/websocket"
	"github.com/vmihailenco/msgpack/v5"
)

const ProtocolVersion = "0.3.4"

// Message is the wire format for Terminalwire protocol (MessagePack-encoded maps).
type Message map[string]interface{}

// Client implements the Terminalwire client protocol over WebSocket.
type Client struct {
	conn        *websocket.Conn
	authority   string
	storagePath string
	programName string
}

// Connect establishes a WebSocket connection to a Terminalwire server.
func Connect(wsURL, programName string) (*Client, error) {
	u, err := url.Parse(wsURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	authority := u.Host
	homeDir, _ := os.UserHomeDir()
	storagePath := filepath.Join(homeDir, ".terminalwire", "authorities", authority, "storage")

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("websocket connect: %w", err)
	}

	return &Client{
		conn:        conn,
		authority:   authority,
		storagePath: storagePath,
		programName: programName,
	}, nil
}

// Run sends the initialization message and enters the main event loop.
func (c *Client) Run(args []string) error {
	defer c.conn.Close()

	if err := c.sendInit(args); err != nil {
		return fmt.Errorf("init: %w", err)
	}

	for {
		msg, err := c.readMsg()
		if err != nil {
			return fmt.Errorf("read: %w", err)
		}

		if done, err := c.dispatch(msg); done {
			return nil
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "handle error: %v\n", err)
		}
	}
}

func (c *Client) sendInit(args []string) error {
	homeDir, _ := os.UserHomeDir()
	storagePattern := filepath.Join(homeDir, ".terminalwire", "authorities", c.authority, "storage", "**/*")

	return c.writeMsg(Message{
		"event": "initialization",
		"protocol": map[string]interface{}{
			"version": ProtocolVersion,
		},
		"entitlement": map[string]interface{}{
			"authority":             c.authority,
			"schemes":              []interface{}{"http", "https"},
			"paths":                []interface{}{c.storagePath, storagePattern},
			"environment_variables": []interface{}{"TERMINALWIRE_HOME"},
		},
		"program": map[string]interface{}{
			"name":      c.programName,
			"arguments": args,
		},
	})
}

func (c *Client) dispatch(msg Message) (done bool, err error) {
	event, _ := msg["event"].(string)

	switch event {
	case "resource":
		return false, c.handleResource(msg)
	case "exit":
		exitWithStatus(msg["status"])
		return true, nil
	default:
		return false, fmt.Errorf("unknown event: %s", event)
	}
}

func (c *Client) writeMsg(msg Message) error {
	data, err := msgpack.Marshal(msg)
	if err != nil {
		return fmt.Errorf("msgpack encode: %w", err)
	}
	return c.conn.WriteMessage(websocket.BinaryMessage, data)
}

func (c *Client) readMsg() (Message, error) {
	_, data, err := c.conn.ReadMessage()
	if err != nil {
		return nil, fmt.Errorf("websocket read: %w", err)
	}

	var msg Message
	if err := msgpack.Unmarshal(data, &msg); err != nil {
		return nil, fmt.Errorf("msgpack decode: %w", err)
	}
	return msg, nil
}

func (c *Client) succeed(name string, response interface{}) error {
	return c.writeMsg(Message{
		"event":    "resource",
		"name":     name,
		"status":   "success",
		"response": response,
	})
}

func (c *Client) fail(name string, reason string) error {
	return c.writeMsg(Message{
		"event":    "resource",
		"name":     name,
		"status":   "failure",
		"response": reason,
	})
}

func exitWithStatus(status interface{}) {
	code := 0
	switch v := status.(type) {
	case int64:
		code = int(v)
	case uint64:
		code = int(v)
	case int8:
		code = int(v)
	case uint8:
		code = int(v)
	case float64:
		code = int(v)
	}
	os.Exit(code)
}
