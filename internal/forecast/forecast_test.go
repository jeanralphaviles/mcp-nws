package forecast

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/icodealot/noaa"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func setupTestServer(t *testing.T, statusCode int, payload string) {
	t.Helper()
	var url string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/points/") {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(
				w,
				`{
					"forecast": "%[1]s",
					"forecastHourly": "%[1]s",
					"forecastGridData": "%[1]s"
				}`, url)
		} else {
			w.WriteHeader(statusCode)
			fmt.Fprintln(w, payload)
		}
	}))

	url = server.URL
	noaa.SetBaseURL(server.URL)

	t.Cleanup(func() {
		server.Close()
	})
}

func TestForecast(t *testing.T) {
	for _, c := range []struct {
		err      bool
		status   int
		expected float64
	}{
		{false, http.StatusOK, 75},
		{true, http.StatusNotFound, -1},
	} {
		setupTestServer(t, c.status, `{"periods":[{"temperature":{"value":75}}]}`)

		params := &mcp.CallToolParamsFor[ForecastParams]{
			Arguments: ForecastParams{Latitude: "37.3918", Longitude: "-122.0601"},
		}

		result, err := Forecast(context.Background(), nil, params)
		got := -1.0
		if result != nil && len(result.StructuredContent.Periods) > 0 {
			got = result.StructuredContent.Periods[0].Temperature
		}
		if (err != nil) != c.err || got != c.expected {
			t.Errorf("Expected Forecast() = (%v, %v), got (%v, %v)", c.expected, "error", result, err)
		}
		if result != nil {
			content, ok := result.Content[0].(*mcp.TextContent)
			if !ok {
				t.Errorf("Error: result.Content[0] is not of type *mcp.TextContent")
			}
			structuredContent, err := json.Marshal(result.StructuredContent)
			if err != nil {
				t.Errorf("Error JSON marshalling StructuredContent: %v", err)
			}
			if content.Text != string(structuredContent) {
				t.Errorf("result.Content should match result.StructuredContent. %v != %v", content.Text, string(structuredContent))
			}
		}
	}
}

func TestHourlyForecast(t *testing.T) {
	for _, c := range []struct {
		err      bool
		status   int
		expected float64
	}{
		{false, http.StatusOK, 75},
		{true, http.StatusNotFound, -1},
	} {
		setupTestServer(t, c.status, `{"periods":[{"temperature":{"value":75}}]}`)

		params := &mcp.CallToolParamsFor[ForecastParams]{
			Arguments: ForecastParams{Latitude: "37.3918", Longitude: "-122.0601"},
		}

		result, err := HourlyForecast(context.Background(), nil, params)
		got := -1.0
		if result != nil && len(result.StructuredContent.Periods) > 0 {
			got = result.StructuredContent.Periods[0].Temperature
		}
		if (err != nil) != c.err || got != c.expected {
			t.Errorf("Expected HourlyForecast() = (%v, %v), got (%v, %v)", c.expected, "error", result, err)
		}
		if result != nil {
			content, ok := result.Content[0].(*mcp.TextContent)
			if !ok {
				t.Errorf("Error: result.Content[0] is not of type *mcp.TextContent")
			}
			structuredContent, err := json.Marshal(result.StructuredContent)
			if err != nil {
				t.Errorf("Error JSON marshalling StructuredContent: %v", err)
			}
			if content.Text != string(structuredContent) {
				t.Errorf("result.Content should match result.StructuredContent. %v != %v", content.Text, string(structuredContent))
			}
		}
	}
}

func TestGridpointForecast(t *testing.T) {
	for _, c := range []struct {
		err      bool
		status   int
		expected float64
	}{
		{false, http.StatusOK, 23.88},
		{true, http.StatusNotFound, -1},
	} {
		setupTestServer(t, c.status, `{"temperature":{"values":[{"value":23.88}]}}`)

		params := &mcp.CallToolParamsFor[ForecastParams]{
			Arguments: ForecastParams{Latitude: "37.3918", Longitude: "-122.0601"},
		}

		result, err := GridpointForecast(context.Background(), nil, params)
		got := -1.0
		if result != nil && len(result.StructuredContent.Temperature.Values) > 0 {
			got = result.StructuredContent.Temperature.Values[0].Value
		}
		if (err != nil) != c.err || got != c.expected {
			t.Errorf("Expected GridpointForecast() = (%v, %v), got (%v, %v)", c.expected, "error", result, err)
		}
		if result != nil {
			content, ok := result.Content[0].(*mcp.TextContent)
			if !ok {
				t.Errorf("Error: result.Content[0] is not of type *mcp.TextContent")
			}
			structuredContent, err := json.Marshal(result.StructuredContent)
			if err != nil {
				t.Errorf("Error JSON marshalling StructuredContent: %v", err)
			}
			if content.Text != string(structuredContent) {
				t.Errorf("result.Content should match result.StructuredContent. %v != %v", content.Text, string(structuredContent))
			}
		}
	}
}
