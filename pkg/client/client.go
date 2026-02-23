package client

import (
	"bufio"
	"fmt"
	"RayhanDB/pkg/core"
	"RayhanDB/pkg/protocol"
	"net"
	"sync"
	"sync/atomic"

	"github.com/google/uuid"
)

// Client represents a RayhanDB client connection
type Client struct {
	conn          net.Conn
	reader        *bufio.Reader
	writer        *bufio.Writer
	username      string
	authenticated bool
	mu            sync.Mutex
	requestID     atomic.Uint64
}

// Connect establishes a connection to RayhanDB server
func Connect(address, username, password string) (*Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	client := &Client{
		conn:     conn,
		reader:   bufio.NewReader(conn),
		writer:   bufio.NewWriter(conn),
		username: username,
	}

	// Authenticate
	if err := client.authenticate(username, password); err != nil {
		conn.Close()
		return nil, err
	}

	return client, nil
}

// authenticate sends authentication credentials
func (c *Client) authenticate(username, password string) error {
	// Send a dummy query with auth credentials to authenticate
	req := &protocol.Request{
		ID:    c.nextRequestID(),
		Query: "FRISSHOW FRISKABLES;", // Simple query to test auth
		Auth: protocol.Auth{
			Username: username,
			Password: password,
		},
	}

	resp, err := c.sendRequest(req)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("authentication failed: %s", resp.Error)
	}

	c.authenticated = true
	return nil
}

// Query executes a SELECT query and returns rows
func (c *Client) Query(query string) ([]core.Row, error) {
	if !c.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	req := &protocol.Request{
		ID:    c.nextRequestID(),
		Query: query,
		Auth: protocol.Auth{
			Username: c.username,
		},
	}

	resp, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, fmt.Errorf("query failed: %s", resp.Error)
	}

	// Convert response data to rows
	if resp.Data == nil {
		return []core.Row{}, nil
	}

	// Data comes as []interface{} of map[string]interface{}
	dataSlice, ok := resp.Data.([]interface{})
	if !ok {
		return []core.Row{}, nil
	}

	rows := make([]core.Row, len(dataSlice))
	for i, item := range dataSlice {
		if rowMap, ok := item.(map[string]interface{}); ok {
			rows[i] = core.Row(rowMap)
		}
	}

	return rows, nil
}

// Exec executes a non-SELECT query (INSERT, UPDATE, DELETE, CREATE, DROP)
func (c *Client) Exec(query string) (string, error) {
	if !c.authenticated {
		return "", fmt.Errorf("not authenticated")
	}

	req := &protocol.Request{
		ID:    c.nextRequestID(),
		Query: query,
		Auth: protocol.Auth{
			Username: c.username,
		},
	}

	resp, err := c.sendRequest(req)
	if err != nil {
		return "", err
	}

	if !resp.Success {
		return "", fmt.Errorf("query failed: %s", resp.Error)
	}

	return resp.Message, nil
}

// sendRequest sends a request and waits for response
func (c *Client) sendRequest(req *protocol.Request) (*protocol.Response, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Encode and send request
	data, err := protocol.EncodeRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	c.writer.Write(data)
	c.writer.WriteByte('\n')
	if err := c.writer.Flush(); err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// Read response
	respData, err := c.reader.ReadBytes('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	resp, err := protocol.DecodeResponse(respData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return resp, nil
}

// nextRequestID generates a unique request ID
func (c *Client) nextRequestID() string {
	return uuid.New().String()
}

// Close closes the connection
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
