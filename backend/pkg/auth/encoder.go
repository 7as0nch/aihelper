// Package auth @author <chengjiang@buffalo-robot.com>
// @date 2023/1/10
// @note
package auth

import (
	"encoding/json"
	"reflect"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/transport/http"
)

func DefaultResponseEncoder(w http.ResponseWriter, r *http.Request, v interface{}) error {
	var response struct {
		Code     int         `json:"code"`
		Data     interface{} `json:"data"`
		Msg      string      `json:"msg"`
		NewToken string      `json:"newToken"`
		DateTime time.Time   `json:"datetime"`
	}

	response.Code = 200
	response.DateTime = time.Now()
	response.Msg = "成功啦，宝子，你可真棒！"
	if v == nil {
		return nil
	}
	var data []byte
	if v != nil && !reflect.ValueOf(v).IsNil() {
		codec, _ := http.CodecForRequest(r, "Accept")
		data, err := codec.Marshal(v)
		response.Data = json.RawMessage(data)
		if err != nil {
			return err
		}
	}
	codec, _ := http.CodecForRequest(r, "Accept")
	data, err := codec.Marshal(&response)

	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", ContentType(codec.Name()))
	_, err = w.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func DefaultErrorEncoder(w http.ResponseWriter, r *http.Request, err error) {
	var response struct {
		Msg      string `json:"msg"`
		Code     int    `json:"code"`
		DateTime int64  `json:"datetime"`
	}

	response.Msg = err.Error()
	response.DateTime = time.Now().UnixMilli()
	response.Code = 500

	codec, _ := http.CodecForRequest(r, "Accept")
	body, err := json.Marshal(&response)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", ContentType(codec.Name()))
	w.WriteHeader(200)
	w.Write(body)
}

const (
	baseContentType = "application"
)

func ContentType(subtype string) string {
	return strings.Join([]string{baseContentType, subtype}, "/")
}
