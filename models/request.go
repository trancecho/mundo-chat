package models

import "encoding/json"

// Request 通用请求数据格式
type Request struct {
	Seq  string `json:"seq"`            // 消息的唯一ID
	Cmd  string `json:"cmd"`            // 请求命令字
	Data any    `json:"data,omitempty"` // 数据 json
}

// Login 登录请求数据
type Login struct {
	ServiceToken string `json:"serviceToken"` // 验证用户是否登录
	AppID        uint32 `json:"appID,omitempty"`
}

// HeartBeat 心跳请求数据
type HeartBeat struct {
	UserID string `json:"userID,omitempty"`
}

type Head struct {
	Seq      string    `json:"seq"`      // 消息的ID
	Cmd      string    `json:"cmd"`      // 消息的cmd 动作
	Response *Response `json:"response"` // 消息体
}

// Response 响应数据体
type Response struct {
	Code    uint32      `json:"code"`
	CodeMsg string      `json:"codeMsg"`
	Data    interface{} `json:"data"` // 数据 json
}

// PushMsg 数据结构体
type PushMsg struct {
	Seq  string `json:"seq"`
	Uuid uint64 `json:"uuid"`
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

// NewResponseHead 设置返回消息
func NewResponseHead(seq string, cmd string, code uint32, codeMsg string, data interface{}) *Head {
	response := NewResponse(code, codeMsg, data)

	return &Head{Seq: seq, Cmd: cmd, Response: response}
}

// String to string
func (h *Head) String() (headStr string) {
	headBytes, _ := json.Marshal(h)
	headStr = string(headBytes)

	return
}

// NewResponse 创建新的响应
func NewResponse(code uint32, codeMsg string, data interface{}) *Response {
	return &Response{Code: code, CodeMsg: codeMsg, Data: data}
}
