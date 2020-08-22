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
	r := recover()
	if r != nil {
		log.Println("[ERROR]:", r)
		os.Exit(2)
	}
	fmt.Println("Closing server")
	os.Exit(1)
}

func handleSigterm() {
	sigterm := make(chan os.Signal)
	signal.Notify(sigterm, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigterm
		cleanup()
	}()
}

func parseFlags() map[string]interface{} {
	configFlag := flag.String("config", "", "Path to config file")
	flag.Parse()

	if len(*configFlag) == 0 {
		panic("ERROR: Config flag is missing!")
	}

	flags := make(map[string]interface{})
	flags["config"] = *configFlag

	return flags
}

func main() {
	defer cleanup()
	handleSigterm()

	flags := parseFlags()
	configFlag, ok := flags["config"]
	if !ok {
		panic("ERROR: Config flag is missing!")
	}

	configPath, ok := configFlag.(string)
	if !ok {
		panic(fmt.Sprintf("Config flag is not a string: %v", configFlag))
	}

	config, err := reading.Config(configPath)
	if err != nil {
		panic(err)
	}

	server, err := web.Create(config)
	if err != nil {
		panic(err)
	}

	server.AddRoute("GET", "/api", web.Welcome)
	server.AddRoute("GET", "/api/entries", web.GetEntries)

	server.Start()
}
