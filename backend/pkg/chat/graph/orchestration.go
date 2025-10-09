package graph

import (
	"context"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func Buildpghelper(ctx context.Context) (r compose.Runnable[[]*schema.Message, *schema.Message], err error) {
	const (
		ChatModel1 = "ChatModel1"
		ToolsNode1 = "ToolsNode1"
		ChatModel2 = "ChatModel2"
	)
	g := compose.NewGraph[[]*schema.Message, *schema.Message]()
	chatModel1KeyOfChatModel, err := newChatModel(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddChatModelNode(ChatModel1, chatModel1KeyOfChatModel)
	toolsNode1KeyOfToolsNode, err := newToolsNode(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddToolsNode(ToolsNode1, toolsNode1KeyOfToolsNode)
	chatModel2KeyOfChatModel, err := newChatModel1(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddChatModelNode(ChatModel2, chatModel2KeyOfChatModel)
	_ = g.AddEdge(compose.START, ChatModel1)
	_ = g.AddEdge(ChatModel2, compose.END)
	_ = g.AddEdge(ChatModel1, ToolsNode1)
	_ = g.AddEdge(ToolsNode1, ChatModel2)
	r, err = g.Compile(ctx, compose.WithGraphName("pghelper"), compose.WithNodeTriggerMode(compose.AnyPredecessor))
	if err != nil {
		return nil, err
	}
	return r, err
}
