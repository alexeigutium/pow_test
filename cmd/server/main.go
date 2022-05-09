package main

import (
	"fmt"

	"github.com/alexeigutium/pow_test/internal/server"
)

func main() {
	cfg := server.NewEnvConfig()

	// loading quotes storage
	quotes, err := server.NewFileQuotes(cfg.GetQuotesPath())
	if err != nil {
		panic(fmt.Sprintf("can't init quotes storage: %s", err.Error()))
	}

	// lets ignore logs, we can just use stdout
	fmt.Println("server is started")
	server.StartServer(cfg, quotes)
	fmt.Println("server's gone, something was wrong")
}
