package graph

import (
	"context"
	"strconv"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/go-kratos/kratos/v2/log"
)

// newToolsNode component initialization function of node 'ToolsNode1' in graph 'pghelper'
func newToolsNode(ctx context.Context) (tsn *compose.ToolsNode, err error) {
	// TODO Modify component configuration here.
	config := &compose.ToolsNodeConfig{}
	toolIns11, err := newTool(ctx)
	if err != nil {
		return nil, err
	}
	config.Tools = []tool.BaseTool{toolIns11}
	tsn, err = compose.NewToolNode(ctx, config)
	if err != nil {
		return nil, err
	}
	return tsn, nil
}

type ToolImpl struct {
	config *ToolConfig
}

type ToolConfig struct {
}

func newTool(ctx context.Context) (bt tool.BaseTool, err error) {
	// TODO Modify component configuration here.
	config := &ToolConfig{}
	bt = &ToolImpl{config: config}
	return bt, nil
}

func (impl *ToolImpl) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "get_current_time",
		Desc: "获取当前时间,标准yyyy-MM-dd HH:mm:ss格式",
	}, nil
}

func (impl *ToolImpl) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	timeZone, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(timeZone).Unix()
	log.Infof("[getCurrentTimeFunc] now: %d", now)
	return `{"current_time": ` + strconv.FormatInt(now, 10) + `}`, nil
}
