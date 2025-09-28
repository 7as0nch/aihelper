package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/spf13/viper"
	"log"
)

// MCPServer 表示一个MCP服务器配置
type MCPServer struct {
	Name    string        `json:"name"`
	URL     string        `json:"url"`
	Timeout time.Duration `json:"timeout"`	
}

// MCPClient MCP服务客户端
type MCPClient struct {
	servers map[string]*MCPServer
}

// NewMCPClient 创建新的MCP客户端
func NewMCPClient() *MCPClient {
	client := &MCPClient{
		servers: make(map[string]*MCPServer),
	}

	// 从配置中加载MCP服务器
	enabled := viper.GetBool("mcp.enabled")
	if !enabled {
		log.Println("MCP service is disabled")
		return client
	}

	serversConfig := viper.Get("mcp.servers")
	if serversConfig == nil {
		log.Println("No MCP servers configured")
		return client
	}

	// 解析配置
	servers, ok := serversConfig.([]interface{})
	if !ok {
		log.Println("Invalid MCP servers configuration")
		return client
	}

	for _, server := range servers {
		serverMap, ok := server.(map[string]interface{})
		if !ok {
			log.Println("Invalid MCP server configuration")
			continue
		}

		name, _ := serverMap["name"].(string)
		url, _ := serverMap["url"].(string)
		timeoutStr, _ := serverMap["timeout"].(string)

		if name == "" || url == "" {
			log.Println("MCP server name or URL is empty")
			continue
		}

		timeout := 30 * time.Second
		if timeoutStr != "" {
			parsedTimeout, err := time.ParseDuration(timeoutStr)
			if err == nil {
				timeout = parsedTimeout
			}
		}

		client.servers[name] = &MCPServer{
			Name:    name,
			URL:     url,
			Timeout: timeout,
		}
	}

	log.Printf("Loaded %d MCP servers", len(client.servers))
	return client
}

// Call 调用MCP服务的API
func (c *MCPClient) Call(serverName, toolName string, params map[string]interface{}) ([]byte, error) {
	// 检查服务器是否存在
	server, exists := c.servers[serverName]
	if !exists {
		return nil, fmt.Errorf("MCP server %s not found", serverName)
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), server.Timeout)
	defer cancel()

	// 构建请求体
	body := map[string]interface{}{
		"tool_name": toolName,
		"args":      params,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "POST", server.URL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// 检查响应状态
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("MCP server returned non-200 status code: %d, body: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// ListServers 列出所有可用的MCP服务器
func (c *MCPClient) ListServers() []*MCPServer {
	servers := make([]*MCPServer, 0, len(c.servers))
	for _, server := range c.servers {
		servers = append(servers, server)
	}
	return servers
}