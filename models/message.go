package models

import "github.com/trancecho/mundo-chat/server/common"

const (
	MessageTypeText  = "text"
	MessageTypeImage = "image"
	MessageTypeEnter = "enter" // 用户进入聊天室
	MessageTypeLeave = "leave" // 用户离开聊天室
)

type Message struct {
	Target  string `json:"target"`
	Content string `json:"content"`
	From    string `json:"from"`
	MsgType string `json:"msgType"`
}

// NewMsg 创建新的消息
func NewMsg(from string, Msg string) (message *Message) {
	message = &Message{
		MsgType: MessageTypeText,
		From:    from,
		Content: Msg,
	}
	return
}

func getTextMsgData(cmd, uuID, msgID, message string) string {
	textMsg := NewMsg(uuID, message)
	head := NewResponseHead(msgID, cmd, common.OK, "Ok", textMsg)
	return head.String()
}

// GetMsgData 文本消息
func GetMsgData(uuID, msgID, cmd, message string) string {
	return getTextMsgData(cmd, uuID, msgID, message)
}

// GetTextMsgData 文本消息
func GetTextMsgData(uuID, msgID, message string) string {
	return getTextMsgData("msg", uuID, msgID, message)
}

// GetTextMsgDataEnter 用户进入消息
func GetTextMsgDataEnter(uuID, msgID, message string) string {
	return getTextMsgData("enter", uuID, msgID, message)
}

// GetTextMsgDataExit 用户退出消息
func GetTextMsgDataExit(uuID, msgID, message string) string {
	return getTextMsgData("exit", uuID, msgID, message)
}
