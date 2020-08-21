package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"

	"github.com/aewens/anda/internal/web"
	"github.com/aewens/anda/pkg/reading"
)

func cleanup() {
	fmt.Println("Closing server")
}

func handleSigterm() {
	sigterm := make(chan os.Signal)
	signal.Notify(sigterm, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigterm
		cleanup()
		os.Exit(1)
	}()
}

func parseFlags() map[string]interface{} {
	configFlag := flag.String("config", "", "Path to config file")
	flag.Parse()

	if len(*configFlag) == 0 {
		log.Println("ERROR: Config flag is missing!")
		os.Exit(2)
	}

	flags := make(map[string]interface{})
	flags["config"] = *configFlag

	return flags
}

func main() {
	handleSigterm()

	flags := parseFlags()
	configFlag, ok := flags["config"]
	if !ok {
		log.Println("ERROR: Config flag is missing!")
		os.Exit(3)
	}

	configPath, ok := configFlag.(string)
	if !ok {
		log.Fatalf("Config flag is not a string: %v", configFlag)
	}

	config, err := reading.Config(configPath)
	if err != nil {
		log.Println(err)
		os.Exit(4)
	}

	server, err := web.Create(config)
	if err != nil {
		log.Println(err)
		os.Exit(5)
	}

	server.AddRoute("GET", "/api", web.Welcome)
	server.AddRoute("GET", "/api/entries", web.GetEntries)

	server.Start()
}
