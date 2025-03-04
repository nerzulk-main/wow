package server

import (
	"context"
	"crypto/rand"
	_ "embed"
	"fmt"
	"log"
	rand2 "math/rand"
	"net"
	"strings"
	"time"
	"wisdomserver/pkg/pow"

	"wisdomserver/pkg/messages"
)

//go:embed quotes.txt
var quotesData string

type Cfg struct {
	Addr          string
	Logger        *log.Logger
	PowTimeout    time.Duration
	PowComplexity int

	WisdomQuotes []string
}

type Server struct {
	cfg *Cfg

	active bool
}

type CfgApplier func(cfg *Cfg)

func NewServer(cfgs ...CfgApplier) *Server {
	cfg := defaultServerCfg()

	for _, c := range cfgs {
		c(cfg)
	}

	return &Server{cfg: cfg}
}

func (s *Server) Start() error {
	tcp, err := net.Listen("tcp", s.cfg.Addr)
	if err != nil {
		return err
	}

	s.active = true

	go func() {
		for {
			if !s.active {
				return
			}

			conn, err := tcp.Accept()
			if err != nil {
				log.Printf("[ERROR] failed accepting TCP connection: %v \n", err)
			}

			deadline := time.Now().Add(s.cfg.PowTimeout)

			err = conn.SetDeadline(deadline)
			if err != nil {
				log.Printf("[ERROR] failed conn set deadline: %v \n", err)
			}

			connCtx, cancel := context.WithDeadline(context.Background(), deadline)

			go func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("[ERROR] panic during handling TCP connection: %v \n", r)
					}
				}()

				err := s.handleConnection(connCtx, conn)
				if err != nil {
					log.Printf("[ERROR] handle connection failure: %v \n", err)
				}

				cancel()
			}()
		}
	}()

	return nil
}

func (s *Server) Stop() error {
	s.active = false
	// should be canceled global context for all working goroutines
	// but as we have poor functionality of only receiving fast response sleep will be enough (lazy)
	time.Sleep(s.cfg.PowTimeout)
	return nil
}

func (s *Server) handleConnection(_ context.Context, conn net.Conn) error {
	rBytes := make([]byte, 32)
	_, err := rand.Read(rBytes)
	if err != nil {
		return fmt.Errorf("rand read: %v", err)
	}

	challenge, err := messages.NewChallengeRequest(s.cfg.PowComplexity, rBytes)
	if err != nil {
		return fmt.Errorf("new challenge message: %v", err)
	}

	// reuse external buffer to save memory allocations
	buffer := make([]byte, 65535)

	err = challenge.Encode(conn, buffer)
	if err != nil {
		return fmt.Errorf("challenge encode: %v", err)
	}

	res, err := messages.Decode(conn, buffer)
	if err != nil {
		return fmt.Errorf("message decode: %v", err)
	}

	if res.Type != messages.ChallengeResponseType {
		errMsg := messages.NewError(fmt.Errorf("wrong challenge response"))
		err := errMsg.Encode(conn, buffer)
		if err != nil {
			return fmt.Errorf("error message encode: %v", err)
		}

		return nil
	}

	msg, err := messages.GetBody[messages.Text](res)
	if err != nil {
		errMsg := messages.NewError(fmt.Errorf("wrong message format %v", err))
		err := errMsg.Encode(conn, buffer)
		if err != nil {
			return fmt.Errorf("error message encode: %v", err)
		}

		return nil
	}

	if !pow.ValidateChallenge(s.cfg.PowComplexity, rBytes, msg.Text) {
		errMsg := messages.NewError(fmt.Errorf("wrong challenge response"))
		err := errMsg.Encode(conn, buffer)
		if err != nil {
			return fmt.Errorf("error message encode: %v", err)
		}

		return nil
	}

	wisdomQuote := messages.NewText(fmt.Sprintf("QOUTE: %v", s.getRandomQuote()))
	err = wisdomQuote.Encode(conn, buffer)
	if err != nil {
		return fmt.Errorf("error message encode: %v", err)
	}

	return nil
}

func (s *Server) getRandomQuote() string {
	return s.cfg.WisdomQuotes[rand2.Intn(len(s.cfg.WisdomQuotes))]
}

func defaultServerCfg() *Cfg {
	return &Cfg{
		Addr:          "localhost:8080",
		Logger:        log.Default(),
		PowTimeout:    time.Second * 10,
		PowComplexity: 6,
		WisdomQuotes:  strings.Split(strings.TrimSpace(quotesData), "\n"),
	}
}
