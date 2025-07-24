#!/bin/bash

# MCP Server Test Script
# This script tests the Google Calendar MCP server using proper JSON-RPC protocol

set -e

SERVER_CMD="go run main.go"
TEMP_DIR=$(mktemp -d)
LOG_FILE="$TEMP_DIR/test.log"

echo "🧪 Starting MCP Server Test Suite"
echo "📁 Temp directory: $TEMP_DIR"
echo "📝 Log file: $LOG_FILE"
echo ""

# Function to send JSON-RPC request and show response
send_request() {
    local description="$1"
    local request="$2"
    local expected_success="$3"  # true/false
    
    echo "🔄 Testing: $description"
    echo "📤 Request: $request"
    
    # Send request and capture response
    response=$(echo "$request" | $SERVER_CMD 2>>"$LOG_FILE")
    echo "📥 Response: $response"
    
    # Check if response contains error
    if echo "$response" | grep -q '"error"'; then
        if [ "$expected_success" = "true" ]; then
            echo "❌ FAILED: Expected success but got error"
            return 1
        else
            echo "✅ Expected error received"
        fi
    else
        if [ "$expected_success" = "false" ]; then
            echo "⚠️  WARNING: Expected error but got success"
        else
            echo "✅ SUCCESS"
        fi
    fi
    echo "---"
    return 0
}

echo "🚀 Step 1: Initialize MCP Connection"
send_request "Initialize connection" \
'{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{"tools":{}},"clientInfo":{"name":"test-client","version":"1.0.0"}}}' \
true

echo ""
echo "📋 Step 2: List Available Tools"
send_request "List tools" \
'{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}' \
true

echo ""
echo "🧮 Step 3: Test Calculator Tool"
send_request "Calculator - Addition" \
'{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"calculate","arguments":{"operation":"add","x":5,"y":3}}}' \
true

send_request "Calculator - Division" \
'{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"calculate","arguments":{"operation":"divide","x":10,"y":2}}}' \
true

send_request "Calculator - Division by zero (should fail)" \
'{"jsonrpc":"2.0","id":5,"method":"tools/call","params":{"name":"calculate","arguments":{"operation":"divide","x":10,"y":0}}}' \
false

echo ""
echo "📅 Step 4: Test Calendar Tools"

echo "📅 Step 4a: Get Calendar Info"
send_request "Get calendar info" \
'{"jsonrpc":"2.0","id":6,"method":"tools/call","params":{"name":"get_calendar_info","arguments":{}}}' \
true

echo "📅 Step 4b: Check Calendar Availability"
send_request "Check availability" \
'{"jsonrpc":"2.0","id":7,"method":"tools/call","params":{"name":"check_google_calendar","arguments":{"start_time":"2024-01-15T09:00:00Z","end_time":"2024-01-15T17:00:00Z"}}}' \
true

echo "📅 Step 4c: List Calendar Events"
send_request "List events" \
'{"jsonrpc":"2.0","id":8,"method":"tools/call","params":{"name":"list_calendar_events","arguments":{"start_time":"2024-01-15T00:00:00Z","end_time":"2024-01-15T23:59:59Z"}}}' \
true

echo "📅 Step 4d: Create Calendar Event"
send_request "Create event" \
'{"jsonrpc":"2.0","id":9,"method":"tools/call","params":{"name":"create_calendar_event","arguments":{"title":"Test Event","start_time":"2024-01-15T14:00:00Z","end_time":"2024-01-15T15:00:00Z","description":"Test event created by MCP server","location":"Test Location"}}}' \
true

echo "📅 Step 4e: Search Calendar Events"
send_request "Search events" \
'{"jsonrpc":"2.0","id":10,"method":"tools/call","params":{"name":"search_calendar_events","arguments":{"query":"test"}}}' \
true

echo ""
echo "🧪 Step 5: Test Error Conditions"

send_request "Invalid tool name (should fail)" \
'{"jsonrpc":"2.0","id":11,"method":"tools/call","params":{"name":"nonexistent_tool","arguments":{}}}' \
false

send_request "Missing required parameter (should fail)" \
'{"jsonrpc":"2.0","id":12,"method":"tools/call","params":{"name":"calculate","arguments":{"operation":"add","x":5}}}' \
false

send_request "Invalid time format (should fail)" \
'{"jsonrpc":"2.0","id":13,"method":"tools/call","params":{"name":"check_google_calendar","arguments":{"start_time":"invalid-time","end_time":"2024-01-15T17:00:00Z"}}}' \
false

echo ""
echo "📊 Test Summary"
echo "✅ Test completed successfully!"
echo "📝 Check log file for detailed server output: $LOG_FILE"
echo ""

# Interactive mode
echo "🔧 Interactive Testing Mode"
echo "You can now test individual commands. Examples:"
echo ""
echo "1. List tools:"
echo 'echo '\''{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}'\'' | GOOGLE_CALENDAR_CREDENTIALS_JSON=./credentials.json go run main.go'
echo ""
echo "2. Call calculator:"
echo 'echo '\''{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"calculate","arguments":{"operation":"multiply","x":6,"y":7}}}'\'' | GOOGLE_CALENDAR_CREDENTIALS_JSON=./credentials.json go run main.go'
echo ""
echo "3. Check calendar availability:"
echo 'echo '\''{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"check_google_calendar","arguments":{"start_time":"2024-01-15T09:00:00Z","end_time":"2024-01-15T17:00:00Z"}}}'\'' | GOOGLE_CALENDAR_CREDENTIALS_JSON=./credentials.json go run main.go'
echo ""

# Cleanup
echo "🧹 Cleaning up temporary files..."
rm -rf "$TEMP_DIR"
echo "✨ Done!" 