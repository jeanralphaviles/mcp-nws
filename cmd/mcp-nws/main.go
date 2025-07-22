package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/jeanralphaviles/mcp-nws/internal/forecast"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var (
	httpAddr = flag.String("address", "", "Address to listen on. If not set, run MCP Server in STDIO mode.")
)

func main() {
	flag.Parse()

	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "mcp-nws",
			Title:   "US National Weather Service MCP Server",
			Version: "v1.0.0",
		},
		nil,
	)
	mcp.AddTool(server, &mcp.Tool{Name: "Forecast", Description: "Basic 7 Day Weather Forecast"}, forecast.Forecast)
	mcp.AddTool(server, &mcp.Tool{Name: "HourlyForecast", Description: "Basic Hourly 7 Day Weather Forecast"}, forecast.HourlyForecast)
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "GridpointForecast",
			Description: "Detailed 7 Day Weather Forecast with Raw Timeseries Data",
		},
		forecast.GridpointForecast,
	)

	if *httpAddr != "" {
		handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
			return server
		}, nil)
		log.Printf("MCP handler listening at %s", *httpAddr)
		if err := http.ListenAndServe(*httpAddr, handler); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
	} else {
		t := mcp.NewLoggingTransport(mcp.NewStdioTransport(), os.Stderr)
		if err := server.Run(context.Background(), t); err != nil {
			log.Fatalf("Server failed: %s", err)
		}
	}
}
