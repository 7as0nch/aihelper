package server

import (
	"github.com/example/aichat/backend/internal/service"
	"github.com/google/wire"
	"github.com/go-kratos/kratos/v2/log"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewGRPCServer, NewHTTPServer, NewWebSocketServerWrapper, NewWebSocketAppWrapper)

// NewWebSocketServerWrapper new a WebSocket server.
func NewWebSocketServerWrapper(chat *service.ChatService, logger log.Logger) *WebSocketServer {
	return NewWebSocketServer(chat, logger)
}

// NewWebSocketAppWrapper new a WebSocket app.
func NewWebSocketAppWrapper(ws *WebSocketServer, logger log.Logger) *WebSocketApp {
	return NewWebSocketApp(ws, logger)
}
