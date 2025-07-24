package main

import (
	"context"
	"fmt"
	"log"

	"google_cal_mcp_golang/calendar"

	"github.com/joho/godotenv"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using environment variables")
	}

	// Load configuration
	config, err := calendar.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Loaded configuration: %s", config.String())

	// Create calendar service (allow it to fail gracefully)
	calendarService, err := calendar.NewCalendarService(config)
	if err != nil {
		log.Printf("Warning: Failed to create calendar service: %v", err)
		log.Printf("Calendar tools will return errors until valid credentials are provided")
		// Create a nil service - tools will handle this gracefully
		calendarService = nil
	}

	// Create a new MCP server
	s := server.NewMCPServer(
		config.ServerName,
		config.ServerVersion,
		server.WithToolCapabilities(true), // Enable tool capabilities
		server.WithRecovery(),
	)

	fmt.Printf("MCP Server '%s' v%s initialized and starting...\n", config.ServerName, config.ServerVersion)

	// --- Calculator Tool (from original code) ---
	log.Println("Registering calculator tool...")
	addCalculatorTool(s)
	log.Println("Calculator tool registered successfully")

	// --- Google Calendar Tools (New) ---
	log.Println("Registering calendar tools...")
	toolManager := calendar.NewToolManager(calendarService, config)
	toolManager.RegisterTools(s)
	log.Println("Calendar tools registered successfully")

	// Log server capabilities
	log.Printf("Server capabilities enabled: tools=true, resources=false, prompts=false")
	log.Printf("Server ready to accept JSON-RPC requests on stdio")
	log.Printf("Expected request format: {\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/call\",\"params\":{\"name\":\"tool_name\",\"arguments\":{...}}}")

	// Start the server
	log.Println("Starting MCP server on stdio...")
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v\n", err)
	}
}

// addCalculatorTool adds the calculator tool and its handler to the server.
func addCalculatorTool(s *server.MCPServer) {
	calculatorTool := mcp.NewTool("calculate",
		mcp.WithDescription("Perform basic arithmetic operations"),
		mcp.WithString("operation",
			mcp.Required(),
			mcp.Description("The operation to perform (add, subtract, multiply, divide)"),
			mcp.Enum("add", "subtract", "multiply", "divide"),
		),
		mcp.WithNumber("x", mcp.Required(), mcp.Description("First number")),
		mcp.WithNumber("y", mcp.Required(), mcp.Description("Second number")),
	)

	s.AddTool(calculatorTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("Received call to 'calculate' tool with request: %+v\n", request)
		op, err := request.RequireString("operation")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		x, err := request.RequireFloat("x")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		y, err := request.RequireFloat("y")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		var result float64
		switch op {
		case "add":
			result = x + y
		case "subtract":
			result = x - y
		case "multiply":
			result = x * y
		case "divide":
			if y == 0 {
				return mcp.NewToolResultError("cannot divide by zero"), nil
			}
			result = x / y
		}

		log.Printf("Calculated result: %.2f\n", result)
		return mcp.NewToolResultText(fmt.Sprintf("%.2f", result)), nil
	})
}
