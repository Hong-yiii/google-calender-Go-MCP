package calendar

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// ToolManager manages MCP tools for calendar operations
type ToolManager struct {
	service CalendarService
	config  *CalendarConfig
}

// NewToolManager creates a new tool manager
func NewToolManager(service CalendarService, config *CalendarConfig) *ToolManager {
	return &ToolManager{
		service: service,
		config:  config,
	}
}

// RegisterTools registers all calendar tools with the MCP server
func (tm *ToolManager) RegisterTools(s *server.MCPServer) {
	tm.registerCheckAvailabilityTool(s)
	tm.registerCreateEventTool(s)
	tm.registerListEventsTool(s)
	tm.registerUpdateEventTool(s)
	tm.registerDeleteEventTool(s)
	tm.registerSearchEventsTool(s)
	tm.registerGetCalendarInfoTool(s)
}

// registerCheckAvailabilityTool registers the check availability tool
func (tm *ToolManager) registerCheckAvailabilityTool(s *server.MCPServer) {
	tool := mcp.NewTool("check_google_calendar",
		mcp.WithDescription("Checks for available time slots in a Google Calendar within a specified time range."),
		mcp.WithString("start_time",
			mcp.Required(),
			mcp.Description("The start of the time window to check, in RFC3339 format (e.g., 2024-07-22T09:00:00Z)."),
		),
		mcp.WithString("end_time",
			mcp.Required(),
			mcp.Description("The end of the time window to check, in RFC3339 format (e.g., 2024-07-22T17:00:00Z)."),
		),
		mcp.WithString("calendar_id",
			mcp.Description("Specific calendar ID to check. Defaults to primary calendar if not provided."),
		),
	)

	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("Received call to 'check_google_calendar' with request: %+v", request)

		// Check if service is available
		if result := tm.checkServiceAvailability(); result != nil {
			return result, nil
		}

		startTimeStr, err := request.RequireString("start_time")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid start_time: %v", err)), nil
		}

		endTimeStr, err := request.RequireString("end_time")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid end_time: %v", err)), nil
		}

		startTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid start_time format. Please use RFC3339 format: %v", err)), nil
		}

		endTime, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid end_time format. Please use RFC3339 format: %v", err)), nil
		}

		timeSlots, err := tm.service.CheckAvailability(ctx, startTime, endTime)
		if err != nil {
			if calErr, ok := err.(CalendarError); ok {
				return mcp.NewToolResultError(calErr.Error()), nil
			}
			return mcp.NewToolResultError(fmt.Sprintf("Failed to check availability: %v", err)), nil
		}

		response, _ := json.MarshalIndent(map[string]interface{}{
			"time_range": map[string]string{
				"start": startTime.Format(time.RFC3339),
				"end":   endTime.Format(time.RFC3339),
			},
			"time_slots": timeSlots,
		}, "", "  ")

		return mcp.NewToolResultText(string(response)), nil
	})
}

// registerCreateEventTool registers the create event tool
func (tm *ToolManager) registerCreateEventTool(s *server.MCPServer) {
	tool := mcp.NewTool("create_calendar_event",
		mcp.WithDescription("Creates a new event in a Google Calendar."),
		mcp.WithString("title",
			mcp.Required(),
			mcp.Description("The title/summary of the event."),
		),
		mcp.WithString("start_time",
			mcp.Required(),
			mcp.Description("The start time for the event, in RFC3339 format."),
		),
		mcp.WithString("end_time",
			mcp.Required(),
			mcp.Description("The end time for the event, in RFC3339 format."),
		),
		mcp.WithString("description",
			mcp.Description("A description for the event."),
		),
		mcp.WithString("location",
			mcp.Description("The location for the event."),
		),
		mcp.WithString("attendees",
			mcp.Description("Comma-separated list of attendee email addresses."),
		),
	)

	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("Received call to 'create_calendar_event' with request: %+v", request)

		// Check if service is available
		if result := tm.checkServiceAvailability(); result != nil {
			return result, nil
		}

		title, err := request.RequireString("title")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid title: %v", err)), nil
		}

		startTimeStr, err := request.RequireString("start_time")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid start_time: %v", err)), nil
		}

		endTimeStr, err := request.RequireString("end_time")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid end_time: %v", err)), nil
		}

		startTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid start_time format. Please use RFC3339 format: %v", err)), nil
		}

		endTime, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid end_time format. Please use RFC3339 format: %v", err)), nil
		}

		eventReq := &EventCreateRequest{
			Summary:   title,
			StartTime: startTime,
			EndTime:   endTime,
		}

		// Optional fields
		if description := request.GetString("description", ""); description != "" {
			eventReq.Description = description
		}

		if location := request.GetString("location", ""); location != "" {
			eventReq.Location = location
		}

		if attendeesStr := request.GetString("attendees", ""); attendeesStr != "" {
			attendees := strings.Split(attendeesStr, ",")
			for i, email := range attendees {
				attendees[i] = strings.TrimSpace(email)
			}
			eventReq.Attendees = attendees
		}

		event, err := tm.service.CreateEvent(ctx, eventReq)
		if err != nil {
			if calErr, ok := err.(CalendarError); ok {
				return mcp.NewToolResultError(calErr.Error()), nil
			}
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create event: %v", err)), nil
		}

		response, _ := json.MarshalIndent(map[string]interface{}{
			"success": true,
			"message": fmt.Sprintf("Successfully created event '%s'", event.Summary),
			"event":   event,
		}, "", "  ")

		return mcp.NewToolResultText(string(response)), nil
	})
}

// registerListEventsTool registers the list events tool
func (tm *ToolManager) registerListEventsTool(s *server.MCPServer) {
	tool := mcp.NewTool("list_calendar_events",
		mcp.WithDescription("Lists events in a Google Calendar within a specified time range."),
		mcp.WithString("start_time",
			mcp.Required(),
			mcp.Description("The start of the time window to list events, in RFC3339 format."),
		),
		mcp.WithString("end_time",
			mcp.Required(),
			mcp.Description("The end of the time window to list events, in RFC3339 format."),
		),
		mcp.WithNumber("max_results",
			mcp.Description("Maximum number of events to return (default: 50)."),
		),
	)

	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("Received call to 'list_calendar_events' with request: %+v", request)

		// Check if service is available
		if result := tm.checkServiceAvailability(); result != nil {
			return result, nil
		}

		startTimeStr, err := request.RequireString("start_time")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid start_time: %v", err)), nil
		}

		endTimeStr, err := request.RequireString("end_time")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid end_time: %v", err)), nil
		}

		startTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid start_time format. Please use RFC3339 format: %v", err)), nil
		}

		endTime, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid end_time format. Please use RFC3339 format: %v", err)), nil
		}

		events, err := tm.service.ListEvents(ctx, startTime, endTime)
		if err != nil {
			if calErr, ok := err.(CalendarError); ok {
				return mcp.NewToolResultError(calErr.Error()), nil
			}
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list events: %v", err)), nil
		}

		response, _ := json.MarshalIndent(map[string]interface{}{
			"time_range": map[string]string{
				"start": startTime.Format(time.RFC3339),
				"end":   endTime.Format(time.RFC3339),
			},
			"event_count": len(events),
			"events":      events,
		}, "", "  ")

		return mcp.NewToolResultText(string(response)), nil
	})
}

// registerUpdateEventTool registers the update event tool
func (tm *ToolManager) registerUpdateEventTool(s *server.MCPServer) {
	tool := mcp.NewTool("update_calendar_event",
		mcp.WithDescription("Updates an existing event in a Google Calendar."),
		mcp.WithString("event_id",
			mcp.Required(),
			mcp.Description("The ID of the event to update."),
		),
		mcp.WithString("title",
			mcp.Description("New title for the event."),
		),
		mcp.WithString("start_time",
			mcp.Description("New start time for the event, in RFC3339 format."),
		),
		mcp.WithString("end_time",
			mcp.Description("New end time for the event, in RFC3339 format."),
		),
		mcp.WithString("description",
			mcp.Description("New description for the event."),
		),
		mcp.WithString("location",
			mcp.Description("New location for the event."),
		),
		mcp.WithString("attendees",
			mcp.Description("Comma-separated list of attendee email addresses."),
		),
	)

	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("Received call to 'update_calendar_event' with request: %+v", request)

		// Check if service is available
		if result := tm.checkServiceAvailability(); result != nil {
			return result, nil
		}

		eventID, err := request.RequireString("event_id")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid event_id: %v", err)), nil
		}

		update := &EventUpdateRequest{}

		if title := request.GetString("title", ""); title != "" {
			update.Summary = &title
		}

		if description := request.GetString("description", ""); description != "" {
			update.Description = &description
		}

		if location := request.GetString("location", ""); location != "" {
			update.Location = &location
		}

		if startTimeStr := request.GetString("start_time", ""); startTimeStr != "" {
			startTime, err := time.Parse(time.RFC3339, startTimeStr)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Invalid start_time format. Please use RFC3339 format: %v", err)), nil
			}
			update.StartTime = &startTime
		}

		if endTimeStr := request.GetString("end_time", ""); endTimeStr != "" {
			endTime, err := time.Parse(time.RFC3339, endTimeStr)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Invalid end_time format. Please use RFC3339 format: %v", err)), nil
			}
			update.EndTime = &endTime
		}

		if attendeesStr := request.GetString("attendees", ""); attendeesStr != "" {
			attendees := strings.Split(attendeesStr, ",")
			for i, email := range attendees {
				attendees[i] = strings.TrimSpace(email)
			}
			update.Attendees = attendees
		}

		event, err := tm.service.UpdateEvent(ctx, eventID, update)
		if err != nil {
			if calErr, ok := err.(CalendarError); ok {
				return mcp.NewToolResultError(calErr.Error()), nil
			}
			return mcp.NewToolResultError(fmt.Sprintf("Failed to update event: %v", err)), nil
		}

		response, _ := json.MarshalIndent(map[string]interface{}{
			"success": true,
			"message": fmt.Sprintf("Successfully updated event '%s'", event.Summary),
			"event":   event,
		}, "", "  ")

		return mcp.NewToolResultText(string(response)), nil
	})
}

// registerDeleteEventTool registers the delete event tool
func (tm *ToolManager) registerDeleteEventTool(s *server.MCPServer) {
	tool := mcp.NewTool("delete_calendar_event",
		mcp.WithDescription("Deletes an event from a Google Calendar."),
		mcp.WithString("event_id",
			mcp.Required(),
			mcp.Description("The ID of the event to delete."),
		),
	)

	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("Received call to 'delete_calendar_event' with request: %+v", request)

		// Check if service is available
		if result := tm.checkServiceAvailability(); result != nil {
			return result, nil
		}

		eventID, err := request.RequireString("event_id")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid event_id: %v", err)), nil
		}

		err = tm.service.DeleteEvent(ctx, eventID)
		if err != nil {
			if calErr, ok := err.(CalendarError); ok {
				return mcp.NewToolResultError(calErr.Error()), nil
			}
			return mcp.NewToolResultError(fmt.Sprintf("Failed to delete event: %v", err)), nil
		}

		response, _ := json.MarshalIndent(map[string]interface{}{
			"success": true,
			"message": fmt.Sprintf("Successfully deleted event with ID: %s", eventID),
		}, "", "  ")

		return mcp.NewToolResultText(string(response)), nil
	})
}

// registerSearchEventsTool registers the search events tool
func (tm *ToolManager) registerSearchEventsTool(s *server.MCPServer) {
	tool := mcp.NewTool("search_calendar_events",
		mcp.WithDescription("Searches for events in a Google Calendar matching a query."),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("Search query to match against event titles and descriptions."),
		),
		mcp.WithString("start_time",
			mcp.Description("Optional start time to limit search, in RFC3339 format."),
		),
		mcp.WithString("end_time",
			mcp.Description("Optional end time to limit search, in RFC3339 format."),
		),
		mcp.WithNumber("max_results",
			mcp.Description("Maximum number of events to return (default: 50)."),
		),
	)

	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("Received call to 'search_calendar_events' with request: %+v", request)

		// Check if service is available
		if result := tm.checkServiceAvailability(); result != nil {
			return result, nil
		}

		query, err := request.RequireString("query")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid query: %v", err)), nil
		}

		var startTime, endTime time.Time

		if startTimeStr := request.GetString("start_time", ""); startTimeStr != "" {
			parsedTime, err := time.Parse(time.RFC3339, startTimeStr)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Invalid start_time format. Please use RFC3339 format: %v", err)), nil
			}
			startTime = parsedTime
		}

		if endTimeStr := request.GetString("end_time", ""); endTimeStr != "" {
			parsedTime, err := time.Parse(time.RFC3339, endTimeStr)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Invalid end_time format. Please use RFC3339 format: %v", err)), nil
			}
			endTime = parsedTime
		}

		events, err := tm.service.SearchEvents(ctx, query, startTime, endTime)
		if err != nil {
			if calErr, ok := err.(CalendarError); ok {
				return mcp.NewToolResultError(calErr.Error()), nil
			}
			return mcp.NewToolResultError(fmt.Sprintf("Failed to search events: %v", err)), nil
		}

		responseData := map[string]interface{}{
			"query":       query,
			"event_count": len(events),
			"events":      events,
		}

		if !startTime.IsZero() || !endTime.IsZero() {
			timeRange := make(map[string]string)
			if !startTime.IsZero() {
				timeRange["start"] = startTime.Format(time.RFC3339)
			}
			if !endTime.IsZero() {
				timeRange["end"] = endTime.Format(time.RFC3339)
			}
			responseData["time_range"] = timeRange
		}

		response, _ := json.MarshalIndent(responseData, "", "  ")
		return mcp.NewToolResultText(string(response)), nil
	})
}

// registerGetCalendarInfoTool registers the get calendar info tool
func (tm *ToolManager) registerGetCalendarInfoTool(s *server.MCPServer) {
	tool := mcp.NewTool("get_calendar_info",
		mcp.WithDescription("Gets basic information about the configured Google Calendar."),
	)

	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("Received call to 'get_calendar_info' with request: %+v", request)

		// Check if service is available
		if result := tm.checkServiceAvailability(); result != nil {
			return result, nil
		}

		calendarInfo, err := tm.service.GetCalendarInfo(ctx)
		if err != nil {
			if calErr, ok := err.(CalendarError); ok {
				return mcp.NewToolResultError(calErr.Error()), nil
			}
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get calendar info: %v", err)), nil
		}

		response, _ := json.MarshalIndent(calendarInfo, "", "  ")
		return mcp.NewToolResultText(string(response)), nil
	})
}

// Helper functions

// checkServiceAvailability checks if the calendar service is available
func (tm *ToolManager) checkServiceAvailability() *mcp.CallToolResult {
	if tm.service == nil {
		return mcp.NewToolResultError("Calendar service unavailable. Please check your Google Calendar credentials configuration.")
	}
	return nil
}

// parseMaxResults parses the max_results parameter with a default value
func parseMaxResults(request mcp.CallToolRequest) int {
	maxResultsFloat := request.GetFloat("max_results", 0)
	if maxResults := int(maxResultsFloat); maxResults > 0 && maxResults <= 250 {
		return maxResults
	}
	return DefaultMaxResults
}

// formatErrorResponse formats an error response for MCP tools
func formatErrorResponse(err error) string {
	if calErr, ok := err.(CalendarError); ok {
		errorResp := ToErrorResponse(calErr)
		response, _ := json.MarshalIndent(errorResp, "", "  ")
		return string(response)
	}

	errorResp := &ErrorResponse{
		Code:    "INTERNAL_ERROR",
		Message: err.Error(),
		Type:    string(ErrorTypeInternal),
	}
	response, _ := json.MarshalIndent(errorResp, "", "  ")
	return string(response)
}
