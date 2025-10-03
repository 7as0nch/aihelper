package server

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// WebSocketApp wraps the WebSocket server to implement the kratos.Server interface
type WebSocketApp struct {
	*WebSocketServer
	log *log.Helper
}

// NewWebSocketApp creates a new WebSocketApp
func NewWebSocketApp(ws *WebSocketServer, logger log.Logger) *WebSocketApp {
	return &WebSocketApp{
		WebSocketServer: ws,
		log:             log.NewHelper(logger),
	}
}

// Start starts the WebSocket server
func (w *WebSocketApp) Start(ctx context.Context) error {
	w.log.Info("Starting WebSocket application")
	return w.WebSocketServer.Start(ctx)
}

// Stop stops the WebSocket server
func (w *WebSocketApp) Stop(ctx context.Context) error {
	w.log.Info("Stopping WebSocket application")
	return w.WebSocketServer.Stop(ctx)
}