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
	Latitude  string `json:"latitude" jsonschema:"The latitude of the forecast location."`
	Longitude string `json:"longitude" jsonschema:"The longitude of the forecast location."`
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

// HourlyForecastResult describes the return type of the HourlyForecast tool.
type HourlyForecastResult = noaa.HourlyForecastResponse

// HourlyForecast returns a standard hourly weather forecast for a location covering 7 days.
func HourlyForecast(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[ForecastParams]) (*mcp.CallToolResultFor[HourlyForecastResult], error) {
	var res mcp.CallToolResultFor[HourlyForecastResult]

	forecast, err := noaa.HourlyForecast(params.Arguments.Latitude, params.Arguments.Longitude)
	if err != nil {
		return nil, err
	}
	res.StructuredContent = *forecast

	return &res, nil
}

// GridpointForecastResult describes the return type of the GridpointForecast tool.
type GridpointForecastResult = noaa.GridpointForecastResponse

// GridpointForecast returns a detailed 7 day weather forecast for a location with raw timeseries data.
func GridpointForecast(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[ForecastParams]) (*mcp.CallToolResultFor[GridpointForecastResult], error) {
	var res mcp.CallToolResultFor[GridpointForecastResult]

	forecast, err := noaa.GridpointForecast(params.Arguments.Latitude, params.Arguments.Longitude)
	if err != nil {
		return nil, err
	}
	res.StructuredContent = *forecast

	return &res, nil
}

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
	mcp.AddTool(server, &mcp.Tool{Name: "Forecast", Description: "Basic 7 Day Weather Forecast"}, Forecast)
	mcp.AddTool(server, &mcp.Tool{Name: "HourlyForecast", Description: "Basic Hourly 7 Day Weather Forecast"}, GridpointForecast)
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "GridpointForecast",
			Description: "Detailed 7 Day Weather Forecast with Raw Timeseries Data",
		},
		GridpointForecast,
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
