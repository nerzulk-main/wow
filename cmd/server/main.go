package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"wisdomserver/pkg/server"
)

func main() {
	addr := os.Getenv("SERVER_ADDR")

	srv := server.NewServer(func(cfg *server.Cfg) {
		cfg.Addr = addr
	})
	err := srv.Start()
	if err != nil {
		log.Panicf("[PANIC] server start panic: %v", err)
	}

	log.Printf("[INFO] Server started at %v \n", addr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	err = srv.Stop()
	if err != nil {
		log.Panicf("[PANIC] server stop panic: %v", err)
	}
	log.Println("[INFO] Server switched off")
}
