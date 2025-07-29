package server

import (
	"encoding/json"
	"fmt"
	"github.com/trancecho/mundo-chat/models"
	"github.com/trancecho/mundo-chat/server/common"
	"log"
)

var (
	manager = NewManager()
)

// ProcessData 处理数据
func ProcessData(client *Client, message []byte) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("处理数据 stop", r)
		}
	}()
	var request models.Request
	if err := json.Unmarshal(message, &request); err != nil {
		fmt.Println("处理数据 json Unmarshal", err)
		client.SendMsg([]byte("数据不合法"))
		return
	}
	seq := request.Seq
	cmd := request.Cmd
	requestData, err := json.Marshal(request.Data)
	if err != nil {
		fmt.Println("处理数据 json Marshal", err)
		client.SendMsg([]byte("处理数据失败"))
		return
	}
	var (
		code uint32
		msg  string
		data any
	)

	fmt.Println("acc_request", cmd, "from", client.Addr)

	// 采用 map 注册的方式
	if value, ok := getHandlers(cmd); ok {
		code, msg, data = value(client, seq, requestData)
	} else {
		code = common.RoutingNotExist
		fmt.Println("处理数据 路由不存在", client.Addr, "cmd", cmd)
	}
	msg = common.GetErrorMessage(code, msg)
	responseHead := models.NewResponseHead(seq, cmd, code, msg, data)
	headByte, err := json.Marshal(responseHead)
	if err != nil {
		fmt.Println("处理数据 json Marshal", err)
		return
	}

	client.SendMsg(headByte)
	fmt.Println("acc_response send", client.Addr, client.AppID, client.UserID, "cmd", cmd, "code", code)
	return
}
