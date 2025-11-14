package server

import (
	"github.com/example/aichat/backend/internal/service"
	"github.com/google/wire"
	"go.uber.org/zap"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(
	NewGRPCServer, NewHTTPServer, NewWebSocketServerWrapper, NewWebSocketAppWrapper)

// NewWebSocketServerWrapper new a WebSocket server.
func NewWebSocketServerWrapper(chat *service.ChatService, logger *zap.Logger) *WebSocketServer {
	return NewWebSocketServer(chat, logger)
}

// NewWebSocketAppWrapper new a WebSocket app.
func NewWebSocketAppWrapper(ws *WebSocketServer, logger *zap.Logger) *WebSocketApp {
	return NewWebSocketApp(ws, logger)
}
