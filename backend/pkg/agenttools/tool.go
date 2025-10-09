package agenttools

import (
	"context"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

/* *
 * @Author: chengjiang
 * @Date: 2025-10-05 10:11:52
 * @Description:
**/

// mcpTools：查询pg数据库内容。
func GetTools() []tool.BaseTool {
	queryUserListFunc, err := utils.InferTool(

		"query_user_list", 
		"根据时间范围（有可能是前三天：表示今天现在这个时间到前面三天，这个不用掉工具，直接计算），查询或生成用户报表，时间：Asia/Shanghai格式", queryUserListFunc)
	if err != nil {
		panic(err)
	}
	getLocalTimeFunc, err := utils.InferTool(
		"get_current_time", 
		"获取当前时间", getCurrentTimeFunc)
	if err != nil {
		panic(err)
	}
	// GetChannelRoi
	getChannelRoiFunc, err := utils.InferTool(
		"get_channel_roi", 
		"根据渠道ID和时间范围，查询渠道ROI", GetChannelRoi)
	if err != nil {
		panic(err)
	}
	return []tool.BaseTool{
		utils.WrapInvokableToolWithErrorHandler(queryUserListFunc, func(ctx context.Context, err error) string {
			return err.Error()
		}),
		utils.WrapInvokableToolWithErrorHandler(getLocalTimeFunc, func(ctx context.Context, err error) string {
			return err.Error()
		}),
		utils.WrapInvokableToolWithErrorHandler(getChannelRoiFunc, func(ctx context.Context, err error) string {
			return err.Error()
		}),
	}
}

func GetStreamTools() []tool.StreamableTool {
	return []tool.StreamableTool{
		// {},
	}
}