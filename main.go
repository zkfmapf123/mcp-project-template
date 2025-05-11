package main

import (
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"golang.org/x/net/context"
)

func main() {
	s := server.NewMCPServer(
		"calculator Demo",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)

	calculatorTool := mcp.NewTool("calculator",
		mcp.WithDescription("Perform basic arthmetic operations"),
		mcp.WithString("operation",
			mcp.Required(),
			mcp.Description("The operation to perform (add, subtract, multiply, divide)"),
			mcp.Enum("add", "subtract", "multiply", "divide"),
		),
		mcp.WithNumber("x", mcp.Required(), mcp.Description("First number")),
		mcp.WithNumber("y", mcp.Required(), mcp.Description("Second number")),
	)

	s.AddTool(calculatorTool, calculatorHandler)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error : %v\n", err)
	}
}

func calculatorHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	op := req.Params.Arguments["operation"].(string)
	x := req.Params.Arguments["x"].(float64)
	y := req.Params.Arguments["y"].(float64)

	var res float64

	switch op {
	case "add":
		res = x + y
	case "subtract":
		res = x - y
	case "multiply":
		res = x * y
	case "divide":
		res = x / y
	}

	return mcp.NewToolResultText(fmt.Sprintf("%.2f", res)), nil
}
