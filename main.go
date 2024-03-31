package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"transaction-parser/internal/api"
	"transaction-parser/internal/config"
	"transaction-parser/internal/service"
)

func loadConfig(configFilePath string) (*config.Config, error) {
	configFile, err := os.Open(configFilePath)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()
	var cfg = &config.Config{}
	err = json.NewDecoder(configFile).Decode(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGTRAP,
	)
	terminateClientListen := make(chan bool, 1)

	configFile := flag.String("config", "config.json", "path to the config file")

	cfg, err := loadConfig(*configFile)
	if err != nil {
		log.Fatalf("failed to load the config file: %v", err)
	}
	parser, err := service.NewEthereumParser(context.Background(), cfg)
	if err != nil {
		log.Fatalf("failed to create the parser: %v", err)
	}
	//_ = parser.Subscribe("0xad03de7dbef5bf0cc4ea23285786e87dbaf69eb3")
	handler := api.NewHTTPHandler(parser)
	http.HandleFunc("/subscribe", api.Post(handler.HandleSubscribe))
	http.HandleFunc("/transactions", api.Get(handler.HandleGetTransactions))
	http.HandleFunc("/current_block", api.Get(handler.HandleGetCurrentBlock))

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil)
		if err != nil {
			log.Fatalf("failed to start the server: %v", err)
		}
		os.Exit(1)
	}()
	// Shutdown
	serveAndListenDone := make(chan bool, 1)
	go func() {
		wg.Wait()
		serveAndListenDone <- true
	}()
	log.Printf("HTTP Server started at port %d\n", cfg.Port)
	select {
	case <-serveAndListenDone:
		log.Println("server shutdown")
	case sig := <-interrupt:
		log.Printf("interrupt signal received, reason: %s\n", sig.String())
	}
	terminateClientListen <- true
}
