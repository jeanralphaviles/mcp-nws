package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/icodealot/noaa"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var (
	httpAddr = flag.String("address", "", "Address to listen on. If not set, run MCP Server in STDIO mode.")
)

// ForecastParams are arguments for mcp-nws calls. They encode the latitude and
// longitude to obtain a forecast for.
type ForecastParams struct {
	Latitude  string `json:"latitude" jsonschema:"the latitude of the forecast location"`
	Longitude string `json:"longitude" jsonschema:"the longitude of the forecast location"`
}

// ForecastResult describes the return type of the Forecast tool.
type ForecastResult = noaa.ForecastResponse

// Forecast returns a standard weather forecast for a location covering 14 periods (day and night for 7 days).
func Forecast(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[ForecastParams]) (*mcp.CallToolResultFor[ForecastResult], error) {
	var res mcp.CallToolResultFor[ForecastResult]

	forecast, err := noaa.Forecast(params.Arguments.Latitude, params.Arguments.Longitude)
	if err != nil {
		return nil, err
	}
	res.StructuredContent = *forecast

	return &res, nil
}

func main() {
	flag.Parse()

	server := mcp.NewServer(&mcp.Implementation{Name: "mcp-nws", Title: "US National Weather Service MCP Server", Version: "v1.0.0"}, nil)
	mcp.AddTool(server, &mcp.Tool{Name: "Forecast", Description: "7d Weather Forecast"}, Forecast)

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
