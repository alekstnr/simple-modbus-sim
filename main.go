package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/alekstnr/simple-modbus-sim/handler"
	"github.com/alekstnr/simple-modbus-sim/types"
	"github.com/simonvetter/modbus"
	"gopkg.in/yaml.v3"
)

func main() {
	configPath := "config.yaml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	var cfg types.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}

	if len(cfg.Devices) == 0 {
		log.Fatal("No devices configured")
	}

	url := cfg.Server.URL
	if url == "" {
		url = "tcp://localhost:5502"
	}

	timeout := 30 * time.Second
	if cfg.Server.Timeout > 0 {
		timeout = time.Duration(cfg.Server.Timeout) * time.Second
	}

	maxClients := uint(5)
	if cfg.Server.MaxClients > 0 {
		maxClients = uint(cfg.Server.MaxClients)
	}

	h := handler.NewModbusHandler(cfg.Devices, cfg.Server.AllowUndefinedRegisters)

	server, err := modbus.NewServer(&modbus.ServerConfiguration{
		URL:        url,
		Timeout:    timeout,
		MaxClients: maxClients,
	}, h)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	fmt.Printf("Modbus TCP server started on %s\n", url)
	fmt.Printf("Allow undefined registers: %v\n", cfg.Server.AllowUndefinedRegisters)
	fmt.Printf("Configured devices:\n")
	for _, d := range cfg.Devices {
		fmt.Printf("  - %s (ID: %d)\n", d.Name, d.ModbusId)
	}
	fmt.Println("Press Ctrl+C to stop")

	select {}
}
