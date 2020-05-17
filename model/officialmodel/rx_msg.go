package officialmodel

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

// RxMessage 一条接收到的消息
type RxMessage struct {
	// ToUserName 开发者微信号
	ToUserName string
	// FromUserName 发送方帐号（一个OpenID）
	FromUserID string
	// CreateTime 消息创建时间 （整型）
	SendTime time.Time
	// MsgType 消息类型
	MsgType MessageType
	// MsgID 消息id，64位整
	MsgID int64

	extras messageKind
}

func fromEnvelope(body []byte) (rxMessage *RxMessage, err error) {
	// extract common part
	var (
		common rxMessageCommon
		extras messageKind
	)

	if err = xml.Unmarshal(body, &common); err != nil {
		return
	}
	// deal with polymorphic message types
	extras, err = extractMessageExtras(common.MsgType, body)
	if err != nil {
		return
	}
	// 等go 1.15 有timezone就不用这么恶心了
	sendTime := time.Unix(common.CreateTime, 0) // in time.Local

	rxMessage = &RxMessage{
		FromUserID: common.FromUserName,
		ToUserName: common.ToUserName,
		SendTime:   sendTime,
		MsgType:    common.MsgType,
		MsgID:      common.MsgID,

		extras: extras,
	}

	return
}
func (m *RxMessage) String() string {
	var sb strings.Builder

	_, _ = fmt.Fprintf(
		&sb,
		"RxMessage { FromUserID: %#v, ToUserName: %#v, SendTime: %d, MsgType: %#v, MsgID: %d",
		m.FromUserID,
		m.ToUserName,
		m.SendTime.UnixNano(),
		m.MsgType,
		m.MsgID,
	)

	m.extras.formatInto(&sb)

	sb.WriteString(" }")

	return sb.String()
}

// Text 如果消息为文本类型，则拿出相应的消息参数，否则返回 nil, false
func (m *RxMessage) Text() (TextMessageExtras, bool) {
	y, ok := m.extras.(TextMessageExtras)
	return y, ok
}

// Image 如果消息为图片类型，则拿出相应的消息参数，否则返回 nil, false
func (m *RxMessage) Image() (ImageMessageExtras, bool) {
	y, ok := m.extras.(ImageMessageExtras)
	return y, ok
}

// Voice 如果消息为语音类型，则拿出相应的消息参数，否则返回 nil, false
func (m *RxMessage) Voice() (VoiceMessageExtras, bool) {
	y, ok := m.extras.(VoiceMessageExtras)
	return y, ok
}

// Video 如果消息为视频类型，则拿出相应的消息参数，否则返回 nil, false
func (m *RxMessage) Video() (VideoMessageExtras, bool) {
	y, ok := m.extras.(VideoMessageExtras)
	return y, ok
}

// Location 如果消息为位置类型，则拿出相应的消息参数，否则返回 nil, false
func (m *RxMessage) Location() (LocationMessageExtras, bool) {
	y, ok := m.extras.(LocationMessageExtras)
	return y, ok
}

// Link 如果消息为链接类型，则拿出相应的消息参数，否则返回 nil, false
func (m *RxMessage) Link() (LinkMessageExtras, bool) {
	y, ok := m.extras.(LinkMessageExtras)
	return y, ok
}
