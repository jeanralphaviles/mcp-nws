// Package forecast implements weather forecasting tools for mcp-nws.
package forecast

import (
	"context"

	"github.com/icodealot/noaa"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ForecastParams are arguments for mcp-nws calls. They encode the latitude and
// longitude to obtain a forecast for.
type ForecastParams struct {
	Latitude  string `json:"latitude" jsonschema:"The latitude of the forecast location."`
	Longitude string `json:"longitude" jsonschema:"The longitude of the forecast location."`
}

// ForecastResponse describes the return type of the Forecast tool.
type ForecastResponse = noaa.ForecastResponse

// Forecast returns a standard weather forecast for a location covering 14 periods (day and night for 7 days).
func Forecast(ctx context.Context, req *mcp.CallToolRequest, params ForecastParams) (*mcp.CallToolResult, *ForecastResponse, error) {
	forecast, err := noaa.Forecast(params.Latitude, params.Longitude)
	if err != nil {
		return nil, nil, err
	}

	return nil, forecast, nil
}

// HourlyForecastResponse describes the return type of the HourlyForecast tool.
type HourlyForecastResponse = noaa.HourlyForecastResponse

// HourlyForecast returns a standard hourly weather forecast for a location covering 7 days.
func HourlyForecast(ctx context.Context, req *mcp.CallToolRequest, params ForecastParams) (*mcp.CallToolResult, *HourlyForecastResponse, error) {
	forecast, err := noaa.HourlyForecast(params.Latitude, params.Longitude)
	if err != nil {
		return nil, nil, err
	}

	return nil, forecast, nil
}

// GridpointForecastResponse describes the return type of the GridpointForecast tool.
type GridpointForecastResponse = noaa.GridpointForecastResponse

// GridpointForecast returns a detailed 7 day weather forecast for a location with raw timeseries data.
func GridpointForecast(ctx context.Context, req *mcp.CallToolRequest, params ForecastParams) (*mcp.CallToolResult, *GridpointForecastResponse, error) {
	forecast, err := noaa.GridpointForecast(params.Latitude, params.Longitude)
	if err != nil {
		return nil, nil, err
	}

	return nil, forecast, nil
}
