# Google Calendar MCP Server Troubleshooting Guide

This guide helps diagnose and fix common issues with the Google Calendar MCP server.

## Quick Diagnosis

### 1. Test Server Startup
```bash
go run main.go
```

**Expected output:**
```
Loaded configuration: CalendarConfig{CalendarID: primary, TimeZone: UTC, Environment: development, Debug: false}
MCP Server 'Google Calendar MCP Server' v1.0.0 initialized and starting...
Registering calculator tool...
Calculator tool registered successfully
Registering calendar tools...
Calendar tools registered successfully
Server capabilities enabled: tools=true, resources=false, prompts=false
Server ready to accept JSON-RPC requests on stdio
Expected request format: {"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"tool_name","arguments":{...}}}
Starting MCP server on stdio...
```

### 2. Test Basic Functionality
```bash
./quick_test.sh
```

## Common Issues

### Issue 1: "Method not found" Error

**Symptoms:**
```json
{"jsonrpc":"2.0","id":1,"error":{"code":-32601,"message":"Method not found"}}
```

**Cause:** Using incorrect JSON-RPC method names

**❌ Wrong:**
```bash
echo '{"jsonrpc":"2.0","id":1,"method":"calculate","params":{"operation":"add","x":5,"y":3}}' | go run main.go
```

**✅ Correct:**
```bash
echo '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"calculate","arguments":{"operation":"add","x":5,"y":3}}}' | go run main.go
```

**Solution:** Always use `"method":"tools/call"` and put tool name in `"params.name"`

### Issue 2: Server Won't Start

**Symptoms:**
```
Failed to load configuration: configuration validation failed: GOOGLE_CALENDAR_CREDENTIALS_JSON is required
```

**Cause:** Missing or invalid environment variables

**Solution:**
1. Create a `.env` file:
   ```bash
   cp env.example .env
   # Edit .env with your credentials
   ```

2. Or set environment variables directly:
   ```bash
   export GOOGLE_CALENDAR_CREDENTIALS_JSON='{"type":"service_account","project_id":"your-project"}'
   ```

3. Verify .env file:
   ```bash
   # Check if .env exists
   ls -la .env
   
   # Check .env file permissions
   chmod 600 .env
   ```

### Issue 3: Environment Variables Not Loading

**Symptoms:**
```
Warning: .env file not found, using environment variables
Failed to load configuration: configuration validation failed: GOOGLE_CALENDAR_CREDENTIALS_JSON is required
```

**Solutions:**
1. **Check .env file location:**
   - Ensure .env file is in the project root directory
   - Verify file name is exactly ".env" (not ".env.local" etc.)

2. **Check .env file format:**
   - No spaces around '=' sign
   - No quotes around values unless they contain spaces
   - Each variable on a new line

3. **Check .env file content:**
   ```bash
   # View .env content (excluding credentials)
   grep -v CREDENTIALS .env
   
   # Check file format
   cat -A .env  # Shows hidden characters
   ```

4. **Verify file permissions:**
   ```bash
   chmod 600 .env
   ```

### Issue 4: Calendar Tools Fail with Authentication Error

**Symptoms:**
```json
{"jsonrpc":"2.0","id":1,"result":{"content":[{"type":"text","text":"INVALID_CREDENTIALS: Failed to create authenticated client"}]}}
```

**Causes & Solutions:**

1. **Invalid credentials file:**
   - Verify the JSON file is valid service account credentials
   - Check file permissions: `chmod 600 credentials.json`

2. **Calendar API not enabled:**
   ```bash
   gcloud services enable calendar-json.googleapis.com
   ```

3. **Calendar not shared with service account:**
   - Share your Google Calendar with the service account email
   - Grant "Make changes to events" permission

### Issue 5: "Invalid time format" Error

**Symptoms:**
```json
{"error": "Invalid start_time format. Please use RFC3339 format"}
```

**Cause:** Incorrect time format

**❌ Wrong formats:**
- `"2024-01-15 14:00:00"`
- `"2024-01-15T14:00:00"`
- `"01/15/2024 2:00 PM"`

**✅ Correct formats:**
- `"2024-01-15T14:00:00Z"` (UTC)
- `"2024-01-15T14:00:00-05:00"` (with timezone)
- `"2024-01-15T14:00:00.000Z"` (with milliseconds)

### Issue 6: Build Errors

**Symptoms:**
```
go: module google_cal_mcp_golang: cannot find module providing package google_cal_mcp_golang/calendar
```

**Solution:**
```bash
go mod tidy
go build
```

### Issue 7: Permission Denied Errors

**Symptoms:**
```json
{"error": "PERMISSION_DENIED: Access denied to calendar"}
```

**Solutions:**
1. **Check calendar sharing:**
   - Go to Google Calendar settings
   - Share calendar with service account email
   - Grant appropriate permissions

2. **Verify service account permissions:**
   - Ensure service account has Calendar API access
   - Check OAuth scopes in credentials

### Issue 8: Network/Timeout Issues

**Symptoms:**
```json
{"error": "NETWORK_TIMEOUT: Network request timed out"}
```

**Solutions:**
1. Check internet connectivity
2. Verify firewall settings
3. Check Google API status: https://status.cloud.google.com/

## Debugging Steps

### 1. Enable Debug Mode
```bash
export DEBUG=true
export LOG_LEVEL=debug
go run main.go
```

### 2. Test Individual Components

**Test configuration loading:**
```bash
go test ./tests/ -v
```

**Test server startup:**
```bash
timeout 5s go run main.go
```

**Test tool registration:**
```bash
echo '{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}' | go run main.go
```

### 3. Use MCP Inspector
```bash
# Install MCP Inspector (if available)
npm install -g @modelcontextprotocol/inspector

# Run with inspector
mcp-inspector go run main.go
```

### 4. Manual JSON-RPC Testing

**Step-by-step protocol test:**
```bash
# 1. Initialize
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{"tools":{}},"clientInfo":{"name":"test","version":"1.0.0"}}}' | go run main.go

# 2. List tools
echo '{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}' | go run main.go

# 3. Call tool
echo '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"calculate","arguments":{"operation":"add","x":5,"y":3}}}' | go run main.go
```

## Performance Issues

### High Memory Usage
- Check for memory leaks in long-running processes
- Monitor goroutine count
- Use `go tool pprof` for profiling

### Slow Response Times
- Check Google Calendar API quotas
- Monitor network latency
- Consider implementing caching

## Error Code Reference

| Error Code | Description | Solution |
|------------|-------------|----------|
| `MISSING_CREDENTIALS` | No credentials provided | Set GOOGLE_CALENDAR_CREDENTIALS_JSON |
| `INVALID_CREDENTIALS` | Invalid credential format | Check JSON format and service account |
| `CALENDAR_NOT_FOUND` | Calendar not accessible | Share calendar with service account |
| `PERMISSION_DENIED` | Access denied | Check calendar permissions |
| `INVALID_TIME_FORMAT` | Wrong time format | Use RFC3339 format |
| `QUOTA_EXCEEDED` | API quota exceeded | Wait or increase quota |
| `SERVICE_UNAVAILABLE` | Google API unavailable | Check API status, retry later |

## Getting Help

### 1. Check Logs
Server logs contain detailed error information:
```bash
go run main.go 2>&1 | tee server.log
```

### 2. Verify Environment
```bash
# Check Go version
go version

# Check environment variables
env | grep GOOGLE_CALENDAR

# Check credentials file
cat credentials.json | jq .
```

### 3. Test with Minimal Configuration
```bash
# Test with minimal fake credentials
GOOGLE_CALENDAR_CREDENTIALS_JSON='{"type":"service_account","project_id":"test"}' ./quick_test.sh
```

### 4. Compare with Working Examples
Review the implementation against the MCP-Go documentation:
- https://github.com/mark3labs/mcp-go
- Check example servers in the repository

## Advanced Debugging

### 1. Network Debugging
```bash
# Monitor network requests
sudo tcpdump -i any -A 'host googleapis.com'

# Check DNS resolution
nslookup googleapis.com
```

### 2. API Debugging
```bash
# Test Google Calendar API directly
curl -H "Authorization: Bearer $(gcloud auth print-access-token)" \
     "https://www.googleapis.com/calendar/v3/calendars/primary"
```

### 3. Protocol Debugging
```bash
# Capture JSON-RPC messages
echo '{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}' | \
  strace -e write go run main.go 2>&1 | grep -A5 -B5 "jsonrpc"
```

Remember: Most issues are related to incorrect JSON-RPC format or missing/invalid credentials. Always verify these first! 