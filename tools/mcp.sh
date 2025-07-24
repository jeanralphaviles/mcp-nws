#!/bin/bash
# https://modelcontextprotocol.io/specification/2025-03-26/basic/lifecycle
#
# Dependencies
#   * https://github.com/jqlang/jq

SERVER="localhost:3000/mcp"

function initialize() {
	local headers_file
	headers_file="$(mktemp)"
	trap 'rm -f "${headers_file}"' EXIT
	curl \
		"${SERVER}" \
		--silent \
		--output /dev/null \
		--request POST \
		--header 'Accept: application/json, text/event-stream' \
		--header 'Content-Type: application/json' \
		--dump-header "${headers_file}" \
		--data @- <<EOF
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "protocolVersion": "2025-03-26",
    "capabilities": {},
    "clientInfo": {
      "name": "curl",
      "version": "0"
    }
  }
}
EOF
	grep -i 'mcp-session-id' "${headers_file}" | cut -d' ' -f2 | tr -d '\r'
}

function initialized() {
	local mcp_session_id="$1"
	curl \
		"${SERVER}" \
		--silent \
		--output /dev/null \
		--request POST \
		--header 'Accept: application/json, text/event-stream' \
		--header 'Content-Type: application/json' \
		--header "Mcp-Session-Id: ${mcp_session_id}" \
		--data @- <<EOF
{
  "jsonrpc": "2.0",
  "method": "notifications/initialized"
}
EOF
}

function tools_list() {
	local mcp_session_id="$1"
	curl \
		"${SERVER}" \
		--silent \
		--request POST \
		--header 'Accept: application/json, text/event-stream' \
		--header 'Content-Type: application/json' \
		--header "Mcp-Session-Id: ${mcp_session_id}" \
		--data @- <<EOF | sed -n 's/^data: \(.*\)/\1/p' | jq -r .
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/list"
}
EOF
}

function forecast() {
	local mcp_session_id="$1"
	local latitude="$2"
	local longitude="$3"
	curl \
		"${SERVER}" \
		--silent \
		--request POST \
		--header 'Accept: application/json, text/event-stream' \
		--header 'Content-Type: application/json' \
		--header "Mcp-Session-Id: ${mcp_session_id}" \
		--data @- <<EOF | sed -n 's/^data: \(.*\)/\1/p' | jq -r .
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "Forecast",
    "arguments": {
      "latitude": "${latitude}",
      "longitude": "${longitude}"
    }
  }
}
EOF
}

function hourly_forecast() {
	local mcp_session_id="$1"
	local latitude="$2"
	local longitude="$3"
	curl \
		"${SERVER}" \
		--silent \
		--request POST \
		--header 'Accept: application/json, text/event-stream' \
		--header 'Content-Type: application/json' \
		--header "Mcp-Session-Id: ${mcp_session_id}" \
		--data @- <<EOF | sed -n 's/^data: \(.*\)/\1/p' | jq -r .
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "HourlyForecast",
    "arguments": {
      "latitude": "${latitude}",
      "longitude": "${longitude}"
    }
  }
}
EOF
}

function gridpoint_forecast() {
	local mcp_session_id="$1"
	local latitude="$2"
	local longitude="$3"
	curl \
		"${SERVER}" \
		--silent \
		--request POST \
		--header 'Accept: application/json, text/event-stream' \
		--header 'Content-Type: application/json' \
		--header "Mcp-Session-Id: ${mcp_session_id}" \
		--data @- <<EOF | sed -n 's/^data: \(.*\)/\1/p' | jq -r .
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "GridpointForecast",
    "arguments": {
      "latitude": "${latitude}",
      "longitude": "${longitude}"
    }
  }
}
EOF
}

MCP_SESSION_ID="$(initialize)"
initialized "${MCP_SESSION_ID}"
tools_list "${MCP_SESSION_ID}"
# forecast "${MCP_SESSION_ID}" "37.3918" "-122.0601"
# hourly_forecast "${MCP_SESSION_ID}" "37.3918" "-122.0601"
# gridpoint_forecast "${MCP_SESSION_ID}" "37.3918" "-122.0601"
