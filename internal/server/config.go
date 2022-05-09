package server

import (
	"os"
	"strconv"
)

const (
	DEFAULT_ADDR        = "0.0.0.0:5599"
	DEFAULT_DIFFICULTY  = 10
	DEFAULT_QUOTES_FILE = "quotes.txt"
)

type Config interface {
	GetListenAddr() string
	GetDifficulty() byte
	GetQuotesPath() string
}

type envConfig struct {
	addr       string
	difficulty byte
	quotesFile string
}

func NewEnvConfig() Config {
	addr, ok := os.LookupEnv("POW_LISTEN_ADDR")
	if !ok {
		addr = DEFAULT_ADDR
	}

	stringDifficulty, _ := os.LookupEnv("POW_DIFFICULTY")
	difficulty, err := strconv.Atoi(stringDifficulty)
	if err != nil {
		difficulty = DEFAULT_DIFFICULTY
	}

	quotesFile, ok := os.LookupEnv("POW_QUOTES_FILE")
	if !ok {
		quotesFile = DEFAULT_QUOTES_FILE
	}

	return &envConfig{
		addr:       addr,
		difficulty: byte(difficulty),
		quotesFile: quotesFile,
	}
}

func (c envConfig) GetListenAddr() string {
	return c.addr
}

func (c envConfig) GetDifficulty() byte {
	return c.difficulty
}

func (c envConfig) GetQuotesPath() string {
	return c.quotesFile
}
