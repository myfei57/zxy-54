package main

import (
	"flag"
	"os"
)

type Config struct {
	Port    string
	DataDir string
}

func LoadConfig() Config {
	cfg := Config{Port: "8080", DataDir: "data"}
	flag.StringVar(&cfg.Port, "port", cfg.Port, "listen port")
	flag.StringVar(&cfg.DataDir, "data", cfg.DataDir, "data directory")
	flag.Parse()
	if v := os.Getenv("MINEBELT_PORT"); v != "" {
		cfg.Port = v
	}
	if v := os.Getenv("MINEBELT_DATA"); v != "" {
		cfg.DataDir = v
	}
	return cfg
}
