package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"wisdomserver/pkg/client"
)

func main() {
	cli, err := client.NewClient(context.Background(), os.Getenv("SERVER_ADDR"))
	if err != nil {
		log.Panicf("[PANIC] client start panic: %v", err)
	}

	log.Println("[INFO] Client connected")

	quote, err := cli.WaitQuote()
	if err != nil {
		log.Panicf("[PANIC] client wait qoute: %v", err)
	}

	fmt.Println(quote)
}
