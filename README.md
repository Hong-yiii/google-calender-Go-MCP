# Google Calendar MCP Server

A Model Context Protocol (MCP) server that provides comprehensive Google Calendar integration for AI assistants and applications. This server enables AI agents to interact with Google Calendar through a clean, well-structured API.

## Features

- **Calendar Availability Checking**: Check for free/busy time slots
- **Event Management**: Create, read, update, and delete calendar events
- **Event Search**: Search events by query with optional time filtering
- **Calendar Information**: Retrieve basic calendar metadata
- **Comprehensive Error Handling**: Detailed error messages and proper error types
- **Secure Authentication**: Support for Google service account credentials
- **Time Zone Support**: Proper handling of time zones and RFC3339 formatting

## Architecture

This implementation follows clean architecture principles with proper separation of concerns:

```
google_cal_mcp_golang/
├── main.go                     # Application entry point
├── calendar/                   # Calendar package
│   ├── service.go             # Calendar service interface & implementation
│   ├── models.go              # Data structures
│   ├── config.go              # Configuration management
│   ├── tools.go               # MCP tool definitions
│   ├── errors.go              # Custom error types
│   └── auth.go                # Authentication handling
├── tests/                      # Test files
└── docs/                       # Documentation
```

## Prerequisites

- Go 1.19 or higher
- Google Cloud Project with Calendar API enabled
- Service Account credentials (JSON file)

## Setup

### 1. Google Cloud Setup

1. Create a Google Cloud Project or use an existing one
2. Enable the Google Calendar API:
   ```bash
   gcloud services enable calendar-json.googleapis.com
   ```
3. Create a Service Account:
   ```bash
   gcloud iam service-accounts create calendar-mcp-server \
       --display-name="Calendar MCP Server"
   ```
4. Create and download credentials:
   ```bash
   gcloud iam service-accounts keys create credentials.json \
       --iam-account=calendar-mcp-server@YOUR_PROJECT_ID.iam.gserviceaccount.com
   ```
5. Share your Google Calendar with the service account email

### 2. Environment Configuration

1. Copy the example environment file:
   ```bash
   cp env.example .env
   ```

2. Edit `.env` with your configuration:
   ```bash
   # Google Calendar API Configuration
   GOOGLE_CALENDAR_CREDENTIALS_JSON=./credentials.json
   GOOGLE_CALENDAR_ID=primary
   GOOGLE_CALENDAR_TIMEZONE=America/New_York

   # MCP Server Configuration
   MCP_SERVER_NAME=Google Calendar MCP Server
   MCP_SERVER_VERSION=1.0.0
   LOG_LEVEL=info

   # Development Settings
   ENVIRONMENT=development
   DEBUG=false
   ```

### 3. Installation

1. Clone or download this repository
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Build the application:
   ```bash
   go build -o calendar-mcp-server
   ```

## Usage

### Running the Server

```bash
# Run the compiled binary
./calendar-mcp-server
```

Or directly with Go:
```bash
go run main.go
```

The server will automatically load environment variables from the `.env` file if it exists. If no `.env` file is found, it will use the system environment variables.

### Testing the Server

The server follows the MCP (Model Context Protocol) JSON-RPC specification. Use the provided test scripts:

**Quick Test (without real calendar credentials):**
```bash
./quick_test.sh
```

**Full Test Suite (requires real Google Calendar credentials):**
```bash
./test_mcp_server.sh
```

**Manual Testing Examples:**

1. **Initialize the connection:**
```bash
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{"tools":{}},"clientInfo":{"name":"test","version":"1.0.0"}}}' | go run main.go
```

2. **List available tools:**
```bash
echo '{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}' | go run main.go
```

3. **Call the calculator tool:**
```bash
echo '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"calculate","arguments":{"operation":"add","x":5,"y":3}}}' | go run main.go
```

4. **Check calendar availability:**
```bash
echo '{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"check_google_calendar","arguments":{"start_time":"2024-01-15T09:00:00Z","end_time":"2024-01-15T17:00:00Z"}}}' | go run main.go
```

**⚠️ Important Notes:**
- Always use `"method":"tools/call"` for tool invocations, NOT the tool name directly
- The MCP protocol requires initialization before tools can be called
- Parameters must be in the `"arguments"` object within `"params"`
- Each command line invocation handles one JSON-RPC request and exits (this is normal MCP behavior)
- For interactive testing, use the provided test scripts or an MCP client

### Available Tools

The server provides the following MCP tools:

#### 1. `check_google_calendar`
Check for available time slots in a specified time range.

**Parameters:**
- `start_time` (required): Start time in RFC3339 format
- `end_time` (required): End time in RFC3339 format
- `calendar_id` (optional): Specific calendar ID

**Example:**
```json
{
  "name": "check_google_calendar",
  "arguments": {
    "start_time": "2024-01-15T09:00:00Z",
    "end_time": "2024-01-15T17:00:00Z"
  }
}
```

#### 2. `create_calendar_event`
Create a new calendar event.

**Parameters:**
- `title` (required): Event title
- `start_time` (required): Start time in RFC3339 format
- `end_time` (required): End time in RFC3339 format
- `description` (optional): Event description
- `location` (optional): Event location
- `attendees` (optional): Comma-separated email addresses

**Example:**
```json
{
  "name": "create_calendar_event",
  "arguments": {
    "title": "Team Meeting",
    "start_time": "2024-01-15T14:00:00Z",
    "end_time": "2024-01-15T15:00:00Z",
    "description": "Weekly team sync",
    "location": "Conference Room A",
    "attendees": "john@example.com,jane@example.com"
  }
}
```

#### 3. `list_calendar_events`
List events in a specified time range.

**Parameters:**
- `start_time` (required): Start time in RFC3339 format
- `end_time` (required): End time in RFC3339 format
- `max_results` (optional): Maximum number of events (default: 50)

#### 4. `update_calendar_event`
Update an existing calendar event.

**Parameters:**
- `event_id` (required): Event ID to update
- `title` (optional): New title
- `start_time` (optional): New start time
- `end_time` (optional): New end time
- `description` (optional): New description
- `location` (optional): New location
- `attendees` (optional): New attendees list

#### 5. `delete_calendar_event`
Delete a calendar event.

**Parameters:**
- `event_id` (required): Event ID to delete

#### 6. `search_calendar_events`
Search for events matching a query.

**Parameters:**
- `query` (required): Search query
- `start_time` (optional): Search start time
- `end_time` (optional): Search end time
- `max_results` (optional): Maximum results

#### 7. `get_calendar_info`
Get basic information about the configured calendar.

**Parameters:** None

## Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `GOOGLE_CALENDAR_CREDENTIALS_JSON` | Path to service account JSON file | - | Yes |
| `GOOGLE_CALENDAR_ID` | Calendar ID to use | `primary` | No |
| `GOOGLE_CALENDAR_TIMEZONE` | Default timezone | `UTC` | No |
| `MCP_SERVER_NAME` | Server name | `Google Calendar MCP Server` | No |
| `MCP_SERVER_VERSION` | Server version | `1.0.0` | No |
| `LOG_LEVEL` | Log level (debug, info, warn, error, fatal) | `info` | No |
| `ENVIRONMENT` | Environment (development, staging, production, test) | `development` | No |
| `DEBUG` | Enable debug mode | `false` | No |

### Calendar ID Options

- `primary`: Use the primary calendar of the service account
- Specific calendar ID: Use a specific shared calendar
- Email address: Use a calendar shared with the service account

## Error Handling

The server provides comprehensive error handling with specific error types:

- `NOT_FOUND`: Resource not found (calendar, event)
- `PERMISSION_DENIED`: Access denied
- `INVALID_INPUT`: Invalid parameters or data
- `QUOTA_EXCEEDED`: API quota exceeded
- `INTERNAL_ERROR`: Internal server errors
- `AUTHENTICATION_ERROR`: Authentication failures
- `NETWORK_ERROR`: Network-related errors
- `TIMEOUT_ERROR`: Request timeouts
- `CONFLICT_ERROR`: Resource conflicts

## Development

### Project Structure

- **`main.go`**: Application entry point and server setup
- **`calendar/service.go`**: Core calendar service implementation
- **`calendar/models.go`**: Data structures and models
- **`calendar/config.go`**: Configuration management
- **`calendar/tools.go`**: MCP tool definitions and handlers
- **`calendar/errors.go`**: Custom error types and handling
- **`calendar/auth.go`**: Google Calendar API authentication

### Adding New Features

1. Define new models in `calendar/models.go`
2. Add service methods in `calendar/service.go`
3. Create MCP tools in `calendar/tools.go`
4. Update error handling in `calendar/errors.go`
5. Add tests in the `tests/` directory

### Testing

Run tests with:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

## Security Considerations

- **Credentials**: Never commit service account credentials to version control
- **Environment Variables**: Use secure methods to manage environment variables
- **Access Control**: Limit service account permissions to minimum required
- **Network Security**: Use HTTPS for all API communications
- **Input Validation**: All inputs are validated before processing

## Troubleshooting

### Common Issues

1. **Authentication Errors**
   - Verify service account credentials are correct
   - Ensure Calendar API is enabled
   - Check that calendar is shared with service account

2. **Permission Denied**
   - Verify calendar sharing permissions
   - Check service account has necessary scopes

3. **Time Format Errors**
   - Ensure all times are in RFC3339 format
   - Check timezone configurations

4. **Configuration Errors**
   - Verify all required environment variables are set
   - Check file paths for credentials

### Debug Mode

Enable debug mode for detailed logging:
```bash
export DEBUG=true
export LOG_LEVEL=debug
```

## Performance Considerations

- **Rate Limiting**: Google Calendar API has rate limits
- **Caching**: Consider implementing caching for frequently accessed data
- **Batch Operations**: Use batch requests for multiple operations
- **Connection Pooling**: HTTP client uses connection pooling

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make changes following the existing architecture
4. Add tests for new functionality
5. Update documentation
6. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For issues and questions:
1. Check the troubleshooting section
2. Review the implementation plan documentation
3. Create an issue with detailed information

## Acknowledgments

- Built with [mcp-go](https://github.com/mark3labs/mcp-go)
- Uses Google Calendar API
- Follows Model Context Protocol specification 
