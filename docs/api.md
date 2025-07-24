# Google Calendar MCP Server API Documentation

This document provides detailed API documentation for all available MCP tools in the Google Calendar server.

## Tool Responses

All tools return JSON responses with consistent structure:

### Success Response Format
```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": { /* tool-specific data */ }
}
```

### Error Response Format
```json
{
  "code": "ERROR_CODE",
  "message": "Human-readable error message",
  "type": "ERROR_TYPE",
  "details": "Additional error details (optional)"
}
```

## Available Tools

### 1. check_google_calendar

**Description**: Checks for available time slots in a Google Calendar within a specified time range.

**Parameters**:
- `start_time` (string, required): Start time in RFC3339 format
- `end_time` (string, required): End time in RFC3339 format  
- `calendar_id` (string, optional): Specific calendar ID

**Example Request**:
```json
{
  "name": "check_google_calendar",
  "arguments": {
    "start_time": "2024-01-15T09:00:00Z",
    "end_time": "2024-01-15T17:00:00Z"
  }
}
```

**Example Response**:
```json
{
  "time_range": {
    "start": "2024-01-15T09:00:00Z",
    "end": "2024-01-15T17:00:00Z"
  },
  "time_slots": [
    {
      "start": "2024-01-15T09:00:00Z",
      "end": "2024-01-15T10:00:00Z",
      "free": true
    },
    {
      "start": "2024-01-15T10:00:00Z",
      "end": "2024-01-15T11:00:00Z",
      "free": false
    },
    {
      "start": "2024-01-15T11:00:00Z",
      "end": "2024-01-15T17:00:00Z",
      "free": true
    }
  ]
}
```

---

### 2. create_calendar_event

**Description**: Creates a new event in a Google Calendar.

**Parameters**:
- `title` (string, required): Event title/summary
- `start_time` (string, required): Start time in RFC3339 format
- `end_time` (string, required): End time in RFC3339 format
- `description` (string, optional): Event description
- `location` (string, optional): Event location
- `attendees` (string, optional): Comma-separated email addresses

**Example Request**:
```json
{
  "name": "create_calendar_event",
  "arguments": {
    "title": "Team Standup",
    "start_time": "2024-01-15T14:00:00Z",
    "end_time": "2024-01-15T14:30:00Z",
    "description": "Daily team standup meeting",
    "location": "Conference Room A",
    "attendees": "john@example.com,jane@example.com"
  }
}
```

**Example Response**:
```json
{
  "success": true,
  "message": "Successfully created event 'Team Standup'",
  "event": {
    "id": "abc123def456",
    "summary": "Team Standup",
    "description": "Daily team standup meeting",
    "start_time": "2024-01-15T14:00:00Z",
    "end_time": "2024-01-15T14:30:00Z",
    "location": "Conference Room A",
    "attendees": ["john@example.com", "jane@example.com"],
    "status": "confirmed",
    "created_at": "2024-01-15T12:00:00Z",
    "updated_at": "2024-01-15T12:00:00Z"
  }
}
```

---

### 3. list_calendar_events

**Description**: Lists events in a Google Calendar within a specified time range.

**Parameters**:
- `start_time` (string, required): Start time in RFC3339 format
- `end_time` (string, required): End time in RFC3339 format
- `max_results` (number, optional): Maximum number of events (default: 50)

**Example Request**:
```json
{
  "name": "list_calendar_events",
  "arguments": {
    "start_time": "2024-01-15T00:00:00Z",
    "end_time": "2024-01-15T23:59:59Z",
    "max_results": 10
  }
}
```

**Example Response**:
```json
{
  "time_range": {
    "start": "2024-01-15T00:00:00Z",
    "end": "2024-01-15T23:59:59Z"
  },
  "event_count": 3,
  "events": [
    {
      "id": "event1",
      "summary": "Morning Meeting",
      "description": "Weekly planning meeting",
      "start_time": "2024-01-15T09:00:00Z",
      "end_time": "2024-01-15T10:00:00Z",
      "location": "Room 101",
      "attendees": ["team@example.com"],
      "status": "confirmed",
      "created_at": "2024-01-14T15:00:00Z",
      "updated_at": "2024-01-14T15:00:00Z"
    }
  ]
}
```

---

### 4. update_calendar_event

**Description**: Updates an existing event in a Google Calendar.

**Parameters**:
- `event_id` (string, required): ID of the event to update
- `title` (string, optional): New event title
- `start_time` (string, optional): New start time in RFC3339 format
- `end_time` (string, optional): New end time in RFC3339 format
- `description` (string, optional): New event description
- `location` (string, optional): New event location
- `attendees` (string, optional): New comma-separated email addresses

**Example Request**:
```json
{
  "name": "update_calendar_event",
  "arguments": {
    "event_id": "abc123def456",
    "title": "Updated Team Meeting",
    "start_time": "2024-01-15T15:00:00Z",
    "end_time": "2024-01-15T16:00:00Z",
    "location": "Conference Room B"
  }
}
```

**Example Response**:
```json
{
  "success": true,
  "message": "Successfully updated event 'Updated Team Meeting'",
  "event": {
    "id": "abc123def456",
    "summary": "Updated Team Meeting",
    "description": "Daily team standup meeting",
    "start_time": "2024-01-15T15:00:00Z",
    "end_time": "2024-01-15T16:00:00Z",
    "location": "Conference Room B",
    "attendees": ["john@example.com", "jane@example.com"],
    "status": "confirmed",
    "created_at": "2024-01-15T12:00:00Z",
    "updated_at": "2024-01-15T13:00:00Z"
  }
}
```

---

### 5. delete_calendar_event

**Description**: Deletes an event from a Google Calendar.

**Parameters**:
- `event_id` (string, required): ID of the event to delete

**Example Request**:
```json
{
  "name": "delete_calendar_event",
  "arguments": {
    "event_id": "abc123def456"
  }
}
```

**Example Response**:
```json
{
  "success": true,
  "message": "Successfully deleted event with ID: abc123def456"
}
```

---

### 6. search_calendar_events

**Description**: Searches for events in a Google Calendar matching a query.

**Parameters**:
- `query` (string, required): Search query to match against event titles and descriptions
- `start_time` (string, optional): Search start time in RFC3339 format
- `end_time` (string, optional): Search end time in RFC3339 format
- `max_results` (number, optional): Maximum number of events (default: 50)

**Example Request**:
```json
{
  "name": "search_calendar_events",
  "arguments": {
    "query": "team meeting",
    "start_time": "2024-01-01T00:00:00Z",
    "end_time": "2024-01-31T23:59:59Z",
    "max_results": 20
  }
}
```

**Example Response**:
```json
{
  "query": "team meeting",
  "time_range": {
    "start": "2024-01-01T00:00:00Z",
    "end": "2024-01-31T23:59:59Z"
  },
  "event_count": 2,
  "events": [
    {
      "id": "search1",
      "summary": "Team Meeting - Q1 Planning",
      "description": "Quarterly planning session",
      "start_time": "2024-01-10T14:00:00Z",
      "end_time": "2024-01-10T16:00:00Z",
      "location": "Main Conference Room",
      "attendees": ["team@example.com"],
      "status": "confirmed",
      "created_at": "2024-01-05T10:00:00Z",
      "updated_at": "2024-01-05T10:00:00Z"
    }
  ]
}
```

---

### 7. get_calendar_info

**Description**: Gets basic information about the configured Google Calendar.

**Parameters**: None

**Example Request**:
```json
{
  "name": "get_calendar_info",
  "arguments": {}
}
```

**Example Response**:
```json
{
  "id": "primary",
  "summary": "John Doe",
  "description": "Primary calendar for John Doe",
  "timezone": "America/New_York",
  "location": "New York, NY"
}
```

## Error Codes

### Authentication Errors
- `MISSING_CREDENTIALS`: No credentials provided
- `INVALID_CREDENTIALS`: Invalid or malformed credentials
- `AUTHENTICATION_ERROR`: Failed to authenticate with Google Calendar API

### Permission Errors
- `PERMISSION_DENIED`: Access denied to calendar or event
- `CALENDAR_NOT_FOUND`: Specified calendar not found or not accessible

### Input Validation Errors
- `INVALID_TIME_FORMAT`: Invalid time format (must be RFC3339)
- `INVALID_TIME_RANGE`: Start time must be before end time
- `INVALID_EVENT_DATA`: Missing or invalid event data
- `EVENT_NOT_FOUND`: Specified event not found

### API Errors
- `QUOTA_EXCEEDED`: Google Calendar API quota exceeded
- `SERVICE_UNAVAILABLE`: Google Calendar API temporarily unavailable
- `NETWORK_TIMEOUT`: Network request timed out
- `INTERNAL_ERROR`: Internal server error

### Configuration Errors
- `CONFIGURATION_ERROR`: Invalid server configuration

## Time Format Requirements

All time parameters must be in RFC3339 format:
- `2024-01-15T14:30:00Z` (UTC)
- `2024-01-15T14:30:00-05:00` (with timezone offset)
- `2024-01-15T14:30:00.000Z` (with milliseconds)

## Rate Limits

The Google Calendar API has the following rate limits:
- 1,000 requests per 100 seconds per user
- 10,000 requests per 100 seconds (total)

The server implements automatic retry with exponential backoff for rate limit errors.

## Best Practices

1. **Time Zones**: Always specify time zones explicitly or use UTC
2. **Event IDs**: Store event IDs returned from create operations for future updates/deletes
3. **Batch Operations**: For multiple operations, consider the API rate limits
4. **Error Handling**: Always check for and handle error responses appropriately
5. **Attendees**: Use valid email addresses for attendees
6. **Search Queries**: Use specific search terms for better results

## Examples

### Creating a Recurring Meeting
```json
{
  "name": "create_calendar_event",
  "arguments": {
    "title": "Weekly Team Standup",
    "start_time": "2024-01-15T09:00:00Z",
    "end_time": "2024-01-15T09:30:00Z",
    "description": "Weekly team standup meeting to discuss progress and blockers",
    "location": "Conference Room A or Zoom",
    "attendees": "team-lead@company.com,dev1@company.com,dev2@company.com"
  }
}
```

### Checking Availability for Multiple Time Slots
```json
{
  "name": "check_google_calendar",
  "arguments": {
    "start_time": "2024-01-15T09:00:00Z",
    "end_time": "2024-01-15T17:00:00Z"
  }
}
```

### Finding All Events for a Day
```json
{
  "name": "list_calendar_events",
  "arguments": {
    "start_time": "2024-01-15T00:00:00Z",
    "end_time": "2024-01-15T23:59:59Z"
  }
}
``` 