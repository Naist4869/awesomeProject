package officialmodel

import (
	"encoding/xml"
	"fmt"
	"io"
)

type rxMessageKind interface {
	messageKind
}

// NOTE: 这顺便就构成了一个封闭的 enum
type messageKind interface {
	formatInto(io.Writer)
}
type txMessageKind interface {
	messageKind
	messageMarshal() ([]byte, error)
}

func extractMessageExtras(ty RxMessageType, body []byte) (rxMessageKind, error) {
	switch ty {
	case RxMessageTypeText:
		var x rxTextMessageSpecifics
		err := xml.Unmarshal(body, &x)
		if err != nil {
			return nil, err
		}
		return &x, nil

	case RxMessageTypeImage:
		var x rxImageMessageSpecifics
		err := xml.Unmarshal(body, &x)
		if err != nil {
			return nil, err
		}
		return &x, nil

	case RxMessageTypeVoice:
		var x rxVoiceMessageSpecifics
		err := xml.Unmarshal(body, &x)
		if err != nil {
			return nil, err
		}
		return &x, nil

	case RxMessageTypeVideo:
		var x rxVideoMessageSpecifics
		err := xml.Unmarshal(body, &x)
		if err != nil {
			return nil, err
		}
		return &x, nil

	case RxMessageTypeLocation:
		var x rxLocationMessageSpecifics
		err := xml.Unmarshal(body, &x)
		if err != nil {
			return nil, err
		}
		return &x, nil

	case RxMessageTypeLink:
		var x rxLinkMessageSpecifics
		err := xml.Unmarshal(body, &x)
		if err != nil {
			return nil, err
		}
		return &x, nil
	case RxMessageTypeShortVideo:
		var x rxShortVideoMessageSpecifics
		err := xml.Unmarshal(body, &x)
		if err != nil {
			return nil, err
		}
		return &x, nil
	default:
		return nil, fmt.Errorf("unknown message type '%s'", ty)

	}

}

func insertMessageExtras(ty TxMessageType) (txMessageKind, bool) {
	switch ty {
	case TxMessageTypeText:
		return &txTextMessageSpecifics{}, true
	case TxMessageTypeImage:
		return &txImageMessageSpecifics{}, true
	case TxMessageTypeVoice:
		return &txVoiceMessageSpecifics{}, true
	case TxMessageTypeVideo:
		return &txVideoMessageSpecifics{}, true
	case TxMessageTypeMusic:
		return &txMusicMessageSpecifics{}, true
	case TxMessageTypeNews:
		return &txNewsMessageSpecifics{}, true
	case TxMessageTypeCS:
		return &txCSMessageSpecifics{}, true
	default:
		return nil, false
	}
}

// RxTextMessageExtras 文本消息的参数。
type RxTextMessageExtras interface {
	rxMessageKind

	// GetContent 返回文本消息的内容。
	GetContent() string
}

var _ RxTextMessageExtras = (*rxTextMessageSpecifics)(nil)

func (r *rxTextMessageSpecifics) formatInto(w io.Writer) {
	_, _ = fmt.Fprintf(w, "Content: %#v", r.Content)
}

func (r *rxTextMessageSpecifics) GetContent() string {
	return r.Content
}

// TxTextMessageExtras 文本消息的参数。
type TxTextMessageExtras interface {
	txMessageKind
	// GetContent 返回文本消息的内容。
	SetContent(string)
}

var _ TxTextMessageExtras = (*txTextMessageSpecifics)(nil)

func (t *txTextMessageSpecifics) formatInto(w io.Writer) {
	_, _ = fmt.Fprintf(w, "Content: %#v", t.Content.CData)
}
func (t *txTextMessageSpecifics) SetContent(content string) {
	t.Content.CData = content
}
func (t *txTextMessageSpecifics) messageMarshal() ([]byte, error) {
	marshal, err := xml.Marshal(*t)
	if err != nil {
		return nil, err
	}
	// 去除<xml></xml>
	marshal = marshal[5 : len(marshal)-6 : len(marshal)-6]
	return marshal, nil
}

// RxImageMessageExtras 图片消息的参数。
type RxImageMessageExtras interface {
	rxMessageKind

	// GetPicURL 返回图片消息的图片链接 URL。
	GetPicURL() string

	// GetMediaID 返回图片消息的图片媒体文件 ID。
	//
	// 可以调用【获取媒体文件】接口拉取，仅三天内有效。
	GetMediaID() string
}

var _ RxImageMessageExtras = (*rxImageMessageSpecifics)(nil)

func (r *rxImageMessageSpecifics) formatInto(w io.Writer) {
	_, _ = fmt.Fprintf(w, "PicURL: %#v, MediaID: %#v", r.PicURL, r.MediaID)
}

func (r *rxImageMessageSpecifics) GetPicURL() string {
	return r.PicURL
}

func (r *rxImageMessageSpecifics) GetMediaID() string {
	return r.MediaID
}

// TxImageMessageExtras 图片消息的参数。
type TxImageMessageExtras interface {
	txMessageKind
	// 通过素材管理中的接口上传多媒体文件，得到的id。
	SetMediaID(string)
}

var _ TxImageMessageExtras = (*txImageMessageSpecifics)(nil)

func (t *txImageMessageSpecifics) formatInto(w io.Writer) {
	_, _ = fmt.Fprintf(w, "MediaID: %#v", t.MediaID.CData)
}

func (t *txImageMessageSpecifics) messageMarshal() ([]byte, error) {
	return xml.Marshal(*t)

}
func (t *txImageMessageSpecifics) SetMediaID(mediaID string) {
	t.MediaID.CData = mediaID
}

// RxVoiceMessageExtras 语音消息的参数。
type RxVoiceMessageExtras interface {
	rxMessageKind

	// GetMediaID 返回语音消息的语音媒体文件 ID。
	//
	// 可以调用【获取媒体文件】接口拉取，仅三天内有效。
	GetMediaID() string

	// GetFormat 返回语音消息的语音格式，如 "amr"、"speex" 等。
	GetFormat() string

	// GetRecognition 语音识别结果，UTF8编码  默认不开启
	GetRecognition() string
}

var _ RxVoiceMessageExtras = (*rxVoiceMessageSpecifics)(nil)

func (r *rxVoiceMessageSpecifics) formatInto(w io.Writer) {
	_, _ = fmt.Fprintf(w, "MediaID: %#v, Format: %#v", r.MediaID, r.Format)
}

func (r *rxVoiceMessageSpecifics) GetMediaID() string {
	return r.MediaID
}

func (r *rxVoiceMessageSpecifics) GetFormat() string {
	return r.Format
}

func (r *rxVoiceMessageSpecifics) GetRecognition() string {
	return r.Recognition
}

// TxVoiceMessageExtras 语音消息的参数
type TxVoiceMessageExtras interface {
	txMessageKind
	SetMediaID(string)
}

func (t *txVoiceMessageSpecifics) formatInto(w io.Writer) {
	_, _ = fmt.Fprintf(w, "MediaID: %#v", t.MediaID.CData)
}

func (t *txVoiceMessageSpecifics) messageMarshal() ([]byte, error) {
	return xml.Marshal(*t)
}
func (t *txVoiceMessageSpecifics) SetMediaID(mediaID string) {
	t.MediaID.CData = mediaID
}

// RxVideoMessageExtras 视频消息的参数。
type RxVideoMessageExtras interface {
	rxMessageKind

	// GetMediaID 返回视频消息的视频媒体文件 ID。
	//
	// 可以调用【获取媒体文件】接口拉取，仅三天内有效。
	GetMediaID() string

	// GetThumbMediaID 返回视频消息缩略图的媒体 ID。
	//
	// 可以调用【获取媒体文件】接口拉取，仅三天内有效。
	GetThumbMediaID() string
}

var _ RxVideoMessageExtras = (*rxVideoMessageSpecifics)(nil)

func (r *rxVideoMessageSpecifics) formatInto(w io.Writer) {
	_, _ = fmt.Fprintf(w, "MediaID: %#v, ThumbMediaID: %#v", r.MediaID, r.ThumbMediaID)
}

func (r *rxVideoMessageSpecifics) GetMediaID() string {
	return r.MediaID
}

func (r *rxVideoMessageSpecifics) GetThumbMediaID() string {
	return r.ThumbMediaID
}

type TxVideoMessageExtras interface {
	txMessageKind
	SetTitle(string)
	SetMediaID(string)
	SetDescription(string)
}

var _ TxVideoMessageExtras = (*txVideoMessageSpecifics)(nil)

func (t *txVideoMessageSpecifics) formatInto(w io.Writer) {
	_, _ = fmt.Fprintf(w, "MediaID: %#v, Title: %#v, Description: %#v", t.MediaID.CData, t.Title.CData, t.Description.CData)
}

func (t *txVideoMessageSpecifics) messageMarshal() ([]byte, error) {
	return xml.Marshal(*t)
}
func (t *txVideoMessageSpecifics) SetTitle(s string) {
	t.Title.CData = s
}

func (t *txVideoMessageSpecifics) SetMediaID(s string) {
	t.MediaID.CData = s
}

func (t *txVideoMessageSpecifics) SetDescription(s string) {
	t.Description.CData = s
}

type RxShortVideoMessageExtras RxVideoMessageExtras

var _ RxShortVideoMessageExtras = (*rxShortVideoMessageSpecifics)(nil)

func (r *rxShortVideoMessageSpecifics) formatInto(w io.Writer) {
	_, _ = fmt.Fprintf(w, "MediaID: %#v, ThumbMediaID: %#v", r.MediaID, r.ThumbMediaID)
}

func (r *rxShortVideoMessageSpecifics) GetMediaID() string {
	return r.MediaID
}

func (r *rxShortVideoMessageSpecifics) GetThumbMediaID() string {
	return r.ThumbMediaID
}

// RxLocationMessageExtras 位置消息的参数。
type RxLocationMessageExtras interface {
	rxMessageKind

	// GetLatitude 返回位置消息的纬度（角度值；北纬为正）。
	GetLatitude() float64

	// GetLongitude 返回位置消息的经度（角度值；东经为正）。
	GetLongitude() float64

	// GetScale 返回位置消息的地图缩放大小。
	GetScale() int

	// GetLabel 返回位置消息的地理位置信息。
	GetLabel() string

	// 不知道这个有啥用，先不暴露
	// GetAppType() string
}

var _ RxLocationMessageExtras = (*rxLocationMessageSpecifics)(nil)

func (r *rxLocationMessageSpecifics) formatInto(w io.Writer) {
	_, _ = fmt.Fprintf(
		w,
		"Latitude: %#v, Longitude: %#v, Scale: %d, Label: %#v",
		r.Lat,
		r.Lon,
		r.Scale,
		r.Label,
	)
}

func (r *rxLocationMessageSpecifics) GetLatitude() float64 {
	return r.Lat
}

func (r *rxLocationMessageSpecifics) GetLongitude() float64 {
	return r.Lon
}

func (r *rxLocationMessageSpecifics) GetScale() int {
	return r.Scale
}

func (r *rxLocationMessageSpecifics) GetLabel() string {
	return r.Label
}

// RxLinkMessageExtras 链接消息的参数。
type RxLinkMessageExtras interface {
	rxMessageKind

	// GetTitle 返回链接消息的标题。
	GetTitle() string

	// GetDescription 返回链接消息的描述。
	GetDescription() string

	// GetURL 返回链接消息的跳转 URL。
	GetURL() string
}

var _ RxLinkMessageExtras = (*rxLinkMessageSpecifics)(nil)

func (r *rxLinkMessageSpecifics) formatInto(w io.Writer) {
	_, _ = fmt.Fprintf(
		w,
		"Title: %#v, Description: %#v, URL: %#v",
		r.Title,
		r.Description,
		r.URL,
	)
}

func (r *rxLinkMessageSpecifics) GetTitle() string {
	return r.Title
}

func (r *rxLinkMessageSpecifics) GetDescription() string {
	return r.Description
}

func (r *rxLinkMessageSpecifics) GetURL() string {
	return r.URL
}

// txNewsMessageSpecifics 发送的音乐消息，特有字段 MsgType 图文为news
type TxNewsMessageExtras interface {
	txMessageKind
	// ArticleCount 图文消息个数；当用户发送文本、图片、视频、图文、地理位置这五种消息时，开发者只能回复1条图文消息；其余场景最多可回复8条图文消息 如果图文数超过限制，则将只发限制内的条数
	SetArticleCount(articleCount int)
	// Title 图文消息标题
	SetTitle(index int, title string)
	// Description 图文消息描述

	SetDescription(index int, description string)
	// PicURL 图片链接，支持JPG、PNG格式，较好的效果为大图360*200，小图200*200

	SetPicURL(index int, picURL string)
	// URL 点击图文消息跳转链接

	SetURL(index int, url string)
}

func (t *txNewsMessageSpecifics) formatInto(w io.Writer) {
	_, _ = fmt.Fprintf(
		w,
		"ArticleCount: %#v",
		t.ArticleCount,
	)
	for i := 0; i < t.ArticleCount; i++ {
		_, _ = fmt.Fprintf(
			w,
			"item { Title: %#v, Description: %#v, PicURL: %#v, URL: %#v",
			t.Articles[i].Item.Title.CData,
			t.Articles[i].Item.Description.CData,
			t.Articles[i].Item.PicURL.CData,
			t.Articles[i].Item.URL.CData,
		)
		_, _ = fmt.Fprint(
			w,
			" }")
	}

}

func (t *txNewsMessageSpecifics) messageMarshal() ([]byte, error) {
	marshal, err := xml.Marshal(*t)
	if err != nil {
		return nil, err
	}
	// 去除<xml></xml>
	marshal = marshal[5 : len(marshal)-6 : len(marshal)-6]
	return marshal, nil
}

func (t *txNewsMessageSpecifics) SetArticleCount(articleCount int) {
	t.ArticleCount = articleCount
	t.Articles = make([]Items, articleCount)
}

func (t *txNewsMessageSpecifics) SetTitle(index int, title string) {
	t.Articles[index].Item.Title.CData = title

}

func (t *txNewsMessageSpecifics) SetDescription(index int, description string) {
	t.Articles[index].Item.Description.CData = description
}

func (t *txNewsMessageSpecifics) SetPicURL(index int, picURL string) {
	t.Articles[index].Item.PicURL.CData = picURL
}

func (t *txNewsMessageSpecifics) SetURL(index int, url string) {
	t.Articles[index].Item.URL.CData = url
}

// TxMusicMessageExtras 发送的音乐消息，特有字段 MsgType 音乐为music
type TxMusicMessageExtras interface {
	txMessageKind

	// Title 音乐标题
	SetTitle(title string)
	// Description 音乐描述
	SetDescription(description string)
	// MusicURL 音乐链接
	SetMusicURL(musicURL string)
	// HQMusicURL 高质量音乐链接，WIFI环境优先使用该链接播放音乐
	SetHQMusicURL(HQMusicURL string)
	// ThumbMediaID 缩略图的媒体id，通过素材管理中的接口上传多媒体文件，得到的id,必须字段
	SetThumbMediaID(thumbMediaID string)
}

func (t *txMusicMessageSpecifics) formatInto(w io.Writer) {
	_, _ = fmt.Fprintf(
		w,
		"Title: %#v, Description: %#v, MusicURL: %#v, HQMusicURL: %#v, ThumbMediaID: %#v",
		t.Title.CData,
		t.Description.CData,
		t.MusicURL.CData,
		t.HQMusicURL.CData,
		t.ThumbMediaID.CData,
	)
}

func (t *txMusicMessageSpecifics) messageMarshal() ([]byte, error) {
	return xml.Marshal(*t)
}

func (t *txMusicMessageSpecifics) SetTitle(title string) {
	t.Title.CData = title
}

func (t *txMusicMessageSpecifics) SetDescription(description string) {
	t.Description.CData = description
}

func (t *txMusicMessageSpecifics) SetMusicURL(musicURL string) {
	t.MusicURL.CData = musicURL
}

func (t *txMusicMessageSpecifics) SetHQMusicURL(HQMusicURL string) {
	t.HQMusicURL.CData = HQMusicURL
}

func (t *txMusicMessageSpecifics) SetThumbMediaID(thumbMediaID string) {
	t.ThumbMediaID.CData = thumbMediaID
}

type TxCSMessageExtras interface {
	txMessageKind
	// 客服账号是个人微信号@公众号微信号
	SetKfAccount(account string)
}

func (t *txCSMessageSpecifics) formatInto(w io.Writer) {
	_, _ = fmt.Fprintf(
		w,
		"KfAccount: %#v",
		t.KfAccount.CData,
	)
}

func (t *txCSMessageSpecifics) messageMarshal() ([]byte, error) {
	return xml.Marshal(*t)
}

func (t *txCSMessageSpecifics) SetKfAccount(account string) {
	t.KfAccount.CData = account
}
