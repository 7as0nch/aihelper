package test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
	"github.com/example/aichat/backend/pkg/agenttools"
	"github.com/example/aichat/backend/pkg/chat"
	"github.com/example/aichat/backend/pkg/chat/graph"
	"github.com/example/aichat/backend/pkg/prints"
)

/* *
 * @Author: chengjiang
 * @Date: 2025-10-06 19:57:33
 * @Description:
**/

func TestGraph_pghelper(t *testing.T) {
	ctx := context.Background()
	agent, err := graph.Buildpghelper(ctx)
	if err != nil {
		t.Fatalf("Buildpghelper failed: %v", err)
	}
	res, err := agent.Stream(ctx, []*schema.Message{
		{
			Role:    schema.User,
			Content: "现在几点了？",
		},
	})
	if err != nil {
		t.Fatalf("Stream failed: %v", err)
	}
	for {
		msg, err := res.Recv()
		if err != nil {
			if err == schema.ErrNoValue {
				break
			}
			t.Fatalf("Read failed: %v", err)
		}
		t.Logf("Message: %+v", msg)
	}
}

func TestPgAgent(t *testing.T) {
	ctx := context.Background()
	cm := chat.NewAiAgent()
	agent := cm.GetReActAgent()
	res, err := agent.Stream(ctx, []*schema.Message{
		{
			Role:    schema.User,
			Content: "现在几点了",
		},
	})
	if err != nil {
		t.Fatalf("Stream failed: %v", err)
	}
	// t.Logf("Message: %+v", res)
	for {
		msg, err := res.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			t.Errorf("Read failed: %v", err)
			// return
		}
		// if msg.ResponseMeta.FinishReason == "stop" {
		// 	break
		// }
		// t.Logf("Message: %+v, thinking: %v", msg, msg.ReasoningContent)
		if msg.ReasoningContent != "" {
			fmt.Print(msg.ReasoningContent)
		}else{
			fmt.Print(msg.Content)
		}
	}
}


func TestDB(t *testing.T) {
	ctx := context.Background()
	data, _, err := agenttools.NewData()
	if err != nil {
		t.Fatalf("NewData failed: %v", err)
	}
	t.Logf("Data: %+v", data)
	channelRois := []*agenttools.ChannelRoi{}
	err = data.GetDB().WithContext(ctx).Raw(`
		SELECT
			aho.sub_channel_id AS 'channel_id',
			SUM(price) AS 'half_profit',
			SUM(income) AS 'half_net_profit'
		FROM
			api_halfpro_order AS aho
			LEFT JOIN api_halfpro_order_line AS ahol ON aho.id = ahol.order_id
		WHERE
			aho.channel_id = ?
			AND aho.created_at BETWEEN ?
			AND ?
			AND ahol.check_status = 1
			AND ahol.push_status = 1 
		GROUP BY
			aho.sub_channel_id;
	`, 2729, "2025-11-01 00:00:00", "2025-11-04 23:59:59").Scan(&channelRois).Error
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}
	t.Logf("ChannelRoi: %+v", channelRois)
}

func TestGraph_adk(t *testing.T) { 
	event := chat.NewAdkAgent().Run(context.Background(), []adk.Message{
		chat.PgHelperPrompt(),
		{
			Role:    schema.User,
			Content: "现在几点了",
		},
	})
	for {
		ev, ok := event.Next()
		if !ok {
			break
		}
		// fmt.Println("is streaming?", ev.Output.MessageOutput.IsStreaming)
		// if ev.Output.MessageOutput.MessageStream != nil {
		// 	for {
		// 		msg, err := ev.Output.MessageOutput.MessageStream.Recv()
		// 		if err != nil {
		// 			if err == io.EOF {
		// 				break
		// 			}
		// 			t.Errorf("Recv failed: %v", err)
		// 		}
		// 		if msg.ReasoningContent != "" {
		// 			fmt.Print(msg.ReasoningContent)
		// 		}
		// 		if msg.Content != "" {
		// 			fmt.Print(msg.Content)
		// 		}
		// 	}
		// }
		prints.Event(ev)
	}
}