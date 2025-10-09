/* *
 * @Author: chengjiang
 * @Date: 2025-10-07 10:05:59
 * @Description:
**/
package chat

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/schema"
	"github.com/go-kratos/kratos/v2/log"
)

func GetCallback() callbacks.Handler {
	return callbacks.NewHandlerBuilder().OnStartFn(func(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
		log.Infof("OnStartFn: %+v, input: %+v", info, input)
		return ctx
	}).OnEndFn(func(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
		log.Infof("OnEndFn: %+v, output: %+v", info, output)
		return ctx
	}).OnErrorFn(func(ctx context.Context, info *callbacks.RunInfo, err error) context.Context {
		log.Errorf("OnErrorFn: %+v, err: %v", info, err)
		return ctx
	}).OnEndWithStreamOutputFn(func(ctx context.Context, info *callbacks.RunInfo, output *schema.StreamReader[callbacks.CallbackOutput]) context.Context {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Errorf("[OnEndStream] panic err:", err)
				}
			}()

			defer output.Close() // remember to close the stream in defer

			log.Infof("=========[OnEndStream]=========")
			for {
				frame, err := output.Recv()
				if errors.Is(err, io.EOF) {
					// finish
					break
				}
				if err != nil {
					log.Errorf("internal error: %s\n", err)
					return
				}

				s, err := json.Marshal(frame)
				if err != nil {
					log.Errorf("internal error: %s\n", err)
					return
				}
				log.Infof("%s: %s\n", info.Name, string(s))
			}

		}()
		return ctx
	}).Build()
}
