package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/TestaVivaDK/e-conomic-connector/internal/economic"
	"github.com/TestaVivaDK/e-conomic-connector/internal/logger"
	"github.com/TestaVivaDK/e-conomic-connector/internal/tools"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	verbose := flag.Bool("verbose", false, "Enable verbose logging to stderr")
	verify := flag.Bool("verify", false, "Verify connection and exit")
	flag.Parse()

	logger.Init(*verbose)

	appSecret := os.Getenv("ECONOMIC_APP_SECRET_TOKEN")
	agreementGrant := os.Getenv("ECONOMIC_AGREEMENT_GRANT_TOKEN")
	if appSecret == "" || agreementGrant == "" {
		fmt.Fprintln(os.Stderr, "Missing ECONOMIC_APP_SECRET_TOKEN or ECONOMIC_AGREEMENT_GRANT_TOKEN")
		os.Exit(1)
	}

	ec := economic.NewClient(appSecret, agreementGrant)

	if *verify {
		raw, err := ec.TestConnection()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Connection failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(raw))
		return
	}

	s := server.NewMCPServer("e-conomic", "1.0.0", server.WithToolCapabilities(false), server.WithRecovery())
	tools.RegisterAll(s, ec)

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
