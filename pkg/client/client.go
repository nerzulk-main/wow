package client

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"
	"wisdomserver/pkg/messages"
	"wisdomserver/pkg/pow"
)

type Client struct {
	conn net.Conn
	buff []byte
}

func NewClient(ctx context.Context, addr string) (*Client, error) {
	dial, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("net dial: %v", err)
	}

	buffer := make([]byte, 65535)

	msg, err := messages.Decode(dial, buffer)
	if err != nil {
		return nil, fmt.Errorf("challenge decode failed: %v", err)
	}

	if msg.Type != messages.ChallengeRequestType {
		return nil, fmt.Errorf("not expected server response: %v", err)
	}

	challenge, err := messages.GetBody[messages.ChallengeRequest](msg)
	if err != nil {
		return nil, fmt.Errorf("challenge format error: %v", err)
	}

	t := time.Now()
	solution, err := pow.Solve(ctx, challenge.Data, challenge.Complexity)
	if err != nil {
		return nil, fmt.Errorf("solve failed: %v", err)
	}
	log.Printf("[INFO] POW took %v ms \n", time.Since(t).Milliseconds())

	res := messages.NewChallengeResponse(solution)
	err = res.Encode(dial, buffer)
	if err != nil {
		return nil, fmt.Errorf("challenge encode failed: %v", err)
	}

	return &Client{conn: dial, buff: buffer}, nil
}

func (c *Client) WaitQuote() (string, error) {
	msg, err := messages.Decode(c.conn, c.buff)
	if err != nil {
		return "", fmt.Errorf("quote decode failed: %v", err)
	}

	if msg.Type != messages.TextType {
		return "", fmt.Errorf("not expected server response after handshake: %v", err)
	}

	quote, err := messages.GetBody[messages.Text](msg)
	if err != nil {
		return "", fmt.Errorf("challenge format error: %v", err)
	}

	return quote.Text, nil
}
