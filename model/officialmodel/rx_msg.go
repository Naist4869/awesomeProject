package officialmodel

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
	"time"
)

type xmlRxEncryptEnvelope struct {
	ToUserName string `xml:"ToUserName"`
	Encrypt    string `xml:"Encrypt"`
}

// RxMessage 一条接收到的消息
type RxMessage struct {
	// ToUserName 开发者微信号
	ToUserName string
	// FromUserName 发送方帐号（一个OpenID）
	FromUserID string
	// CreateTime 消息创建时间 （整型）
	SendTime time.Time
	// MsgType 消息类型
	MsgType RxMessageType
	// MsgID 消息id，64位整
	MsgID int64

	extras rxMessageKind
}

// TxMessage 一条要发送的消息
type TxMessage struct {
	// ToUserName 开发者微信号
	ToUserName string
	// FromUserName 发送方帐号（一个OpenID）
	FromUserID string
	// CreateTime 消息创建时间 （整型）
	SendTime time.Time
	// MsgType 消息类型
	MsgType TxMessageType
	// MsgID 消息id，64位整
	MsgID int64

	extras txMessageKind
}

func FromEnvelope(body []byte) (rxMessage *RxMessage, err error) {
	// extract common part
	var (
		common rxMessageCommon
		extras rxMessageKind
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

func ToEnvelope(rx *RxMessage, txType TxMessageType) (tx TxMessage, err error) {
	if rx == nil {
		err = errors.New("rxmessage nil pointer")
		return
	}
	extras, ok := insertMessageExtras(txType)
	if !ok {
		err = fmt.Errorf("unknown message type '%s'", txType)
		return
	}
	tx = TxMessage{
		ToUserName: rx.FromUserID,
		FromUserID: rx.ToUserName,
		SendTime:   time.Now(),
		MsgType:    txType,
		MsgID:      rx.MsgID,
		extras:     extras,
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
func (m *TxMessage) String() string {
	var sb strings.Builder

	_, _ = fmt.Fprintf(
		&sb,
		"TxMessage { FromUserID: %#v, ToUserName: %#v, SendTime: %d, MsgType: %#v, MsgID: %d",
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
func (m *TxMessage) MessageMarshal() (reply string, err error) {
	var body, marshal, extrasMarshal []byte
	common := txMessageCommon{
		ToUserName:   cdataNode{CData: m.ToUserName},
		FromUserName: cdataNode{CData: m.FromUserID},
		CreateTime:   time.Now().Unix(),
		MsgType:      cdataNode{CData: string(m.MsgType)},
		MsgID:        m.MsgID,
	}
	marshal, err = xml.Marshal(common)
	if err != nil {
		return
	}

	marshal = bytes.TrimSuffix(marshal, []byte("</xml>"))
	extrasMarshal, err = m.extras.messageMarshal()
	if err != nil {
		return
	}
	l := len(marshal)
	el := len(extrasMarshal)
	body = make([]byte, l+el+6)
	copy(body[:], marshal)
	copy(body[l:], extrasMarshal)
	copy(body[l+el:], "</xml>")
	reply = string(body)
	return
}

// Text 如果消息为文本类型，则拿出相应的消息参数，否则返回 nil, false
func (m RxMessage) Text() (RxTextMessageExtras, bool) {
	y, ok := m.extras.(RxTextMessageExtras)
	return y, ok
}

// Image 如果消息为图片类型，则拿出相应的消息参数，否则返回 nil, false
func (m RxMessage) Image() (RxImageMessageExtras, bool) {
	y, ok := m.extras.(RxImageMessageExtras)
	return y, ok
}

// Voice 如果消息为语音类型，则拿出相应的消息参数，否则返回 nil, false
func (m RxMessage) Voice() (RxVoiceMessageExtras, bool) {
	y, ok := m.extras.(RxVoiceMessageExtras)
	return y, ok
}

// Video 如果消息为视频类型，则拿出相应的消息参数，否则返回 nil, false
func (m RxMessage) Video() (RxVideoMessageExtras, bool) {
	y, ok := m.extras.(RxVideoMessageExtras)
	return y, ok
}

// Location 如果消息为位置类型，则拿出相应的消息参数，否则返回 nil, false
func (m RxMessage) Location() (RxLocationMessageExtras, bool) {
	y, ok := m.extras.(RxLocationMessageExtras)
	return y, ok
}

// Link 如果消息为链接类型，则拿出相应的消息参数，否则返回 nil, false
func (m RxMessage) Link() (RxLinkMessageExtras, bool) {
	y, ok := m.extras.(RxLinkMessageExtras)
	return y, ok
}

// Text 如果消息为文本类型，则拿出相应的消息参数，否则返回 nil, false
func (m TxMessage) Text() (TxTextMessageExtras, bool) {
	y, ok := m.extras.(TxTextMessageExtras)
	return y, ok
}

// Image 如果消息为图片类型，则拿出相应的消息参数，否则返回 nil, false
func (m TxMessage) Image() (TxImageMessageExtras, bool) {
	y, ok := m.extras.(TxImageMessageExtras)
	return y, ok
}

// Voice 如果消息为语音类型，则拿出相应的消息参数，否则返回 nil, false
func (m TxMessage) Voice() (TxVoiceMessageExtras, bool) {
	y, ok := m.extras.(TxVoiceMessageExtras)
	return y, ok
}

// Video 如果消息为视频类型，则拿出相应的消息参数，否则返回 nil, false
func (m TxMessage) Video() (TxVideoMessageExtras, bool) {
	y, ok := m.extras.(TxVideoMessageExtras)
	return y, ok
}

// Music 如果消息为音乐类型，则拿出相应的消息参数，否则返回 nil, false
func (m TxMessage) Music() (TxMusicMessageExtras, bool) {
	y, ok := m.extras.(TxMusicMessageExtras)
	return y, ok
}

// News 如果消息为图文类型，则拿出相应的消息参数，否则返回 nil, false
func (m TxMessage) News() (TxNewsMessageExtras, bool) {
	y, ok := m.extras.(TxNewsMessageExtras)
	return y, ok
}

// CS customer_service 如果消息为客服类型，则拿出相应的消息参数，否则返回 nil, false
func (m TxMessage) CS() (TxCSMessageExtras, bool) {
	y, ok := m.extras.(TxCSMessageExtras)
	return y, ok
}
