package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {

	hooks := &server.Hooks{}

	hooks.OnBeforeCallTool = append(hooks.OnBeforeCallTool, func(ctx context.Context, id any, message *mcp.CallToolRequest) {
		log.Printf("[tool call] ToolName: %s, Params: %v", message.Params.Name, message.Params.Arguments)
	})

	hooks.OnAfterCallTool = append(hooks.OnAfterCallTool, func(ctx context.Context, id any, message *mcp.CallToolRequest, result *mcp.CallToolResult) {
		if result != nil && result.IsError {
			log.Printf("[tool call error] ToolName: %s, Params: %v, error msg: %v", message.Params.Name, message.Params.Arguments, result.Content)
		}
	})
	// Create a new MCP server
	s := server.NewMCPServer(
		"Hello World Server",
		mcp.LATEST_PROTOCOL_VERSION,
		server.WithToolCapabilities(true),
		server.WithLogging(),
		server.WithHooks(hooks),
	)

	// Define a simple tool
	tool := mcp.NewTool("hello_world",
		mcp.WithDescription("Say hello to someone"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the person to greet"),
		),
	)
	// addr := "localhost:8080"
	// Add tool handler
	s.AddTool(tool, helloHandler)
	// Start the stdio server
	// httpServer := server.NewStreamableHTTPServer(s,
	// 	server.WithStateLess(true),
	// 	server.WithHTTPContextFunc(func(ctx context.Context, r *http.Request) context.Context {
	// 		auth := r.Header.Get("Authorization")
	// 		if len(auth) > 7 && auth[:7] == "Bearer " {
	// 			token := auth[7:]
	// 			ctx = context.WithValue(ctx, "access_token", token)
	// 		}
	// 		return ctx
	// 	}),
	// )
	// log.Printf("HTTP server listening on %s", addr)
	// if err := httpServer.Start(addr); err != nil {
	// 	if err == context.Canceled {
	// 		log.Println(err)
	// 	}
	// 	log.Fatalf("server error: %v", err)
	// }

	err := server.NewSSEServer(s,
		server.WithBaseURL("http://localhost:8080"),
		server.WithSSEContextFunc(func(ctx context.Context, r *http.Request) context.Context {
			auth := r.Header.Get("Authorization")
			if len(auth) > 7 && auth[:7] == "Bearer " {
				token := auth[7:]
				ctx = context.WithValue(ctx, "access_token", token)
			}
			return ctx
		}),
	).Start("localhost:8080")

	if err != nil {
		log.Fatal(err)
	}
}

func helloHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	arguments := request.GetArguments()
	name, ok := arguments["name"].(string)
	if !ok {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Error: name parameter is required and must be a string",
				},
			},
			IsError: true,
		}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Hello, %s! 👋", name),
			},
		},
	}, nil
}
