package officialwx

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/Naist4869/awesomeProject/dataservice"

	"github.com/Naist4869/awesomeProject/model/officialmodel"
)

type HandlerFunc func(message *officialmodel.RxMessage)

type UseCase struct {
	apiService dataservice.IOfficialWxRpcService
}

func NewUseCase(apiService dataservice.IOfficialWxRpcService) *UseCase {
	u := &UseCase{
		apiService: apiService,
	}
	return u
}

func (u *UseCase) ReplyMessage(xmlMsg []byte) (reply string, err error) {
	var (
		rxMessage *officialmodel.RxMessage
	)

	rxMessage, err = officialmodel.FromEnvelope(xmlMsg)
	if err != nil {
		err = fmt.Errorf("解析消息失败: %w", err)
		return
	}
	return u.handler(rxMessage)
}

func (u *UseCase) handler(rxMessage *officialmodel.RxMessage) (reply string, err error) {
	var (
		ToEnvelope officialmodel.TxMessage
	)
	if text, ok := rxMessage.Text(); ok {
		switch msg := text.GetContent(); msg {
		case "【收到不支持的消息类型，暂无法显示】":
			ToEnvelope, err = officialmodel.ToEnvelope(rxMessage, officialmodel.TxMessageTypeImage)
			if err != nil {
				return
			}
			image, ok := ToEnvelope.Image()
			if !ok {
				err = errors.New("不能断言成图片")
				return
			}
			ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
			if mediaIDResp, err := u.apiService.MediaIDGet(ctx, officialmodel.MediaIDReq{
				FakeID:    rxMessage.FromUserID,
				Timestamp: rxMessage.SendTime.Unix(),
			}); err != nil {
				return "", err
			} else {
				image.SetMediaID(mediaIDResp.MediaID)
			}

		case "我要表情包":
			ToEnvelope, err = officialmodel.ToEnvelope(rxMessage, officialmodel.TxMessageTypeText)
			if err != nil {
				return
			}
			text, ok := ToEnvelope.Text()
			if !ok {
				err = errors.New("不能断言成文本")
				return
			}
			text.SetContent("不给")
		case "我要消费券":

		default:
			ToEnvelope, err = officialmodel.ToEnvelope(rxMessage, officialmodel.TxMessageTypeText)
			if err != nil {
				return
			}
			text, ok := ToEnvelope.Text()
			if !ok {
				err = errors.New("不能断言成文本")
				return
			}
			ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
			re := regexp.MustCompile(`([$€₤₳¢¤฿฿₵₡₫ƒ₲₭£₥₦₱〒₮₩₴₪៛﷼₢M₰₯₠₣₧ƒ][a-zA-Z0-9]{9,11}[$€₤₳¢¤฿฿₵₡₫ƒ₲₭£₥₦₱〒₮₩₴₪៛﷼₢M₰₯₠₣₧ƒ])`)
			submatch := re.FindSubmatch([]byte(msg))
			if len(submatch) <= 1 {
				text.SetContent("未匹配到淘口令")
				reply, err = ToEnvelope.MessageMarshal()
				return
			}
			fromKey := string(submatch[1])
			var keyConvertKeyResp officialmodel.KeyConvertKeyResp
			keyConvertKeyResp, err = u.apiService.KeyConvertKey(ctx, officialmodel.KeyConvertKeyReq{
				FromKey: fromKey,
				UserID:  ToEnvelope.ToUserName,
			})

			if err == nil && keyConvertKeyResp.ToKey != "" {
				ToEnvelope, err = officialmodel.ToEnvelope(rxMessage, officialmodel.TxMessageTypeNews)
				if err != nil {
					return
				}
				news, ok := ToEnvelope.News()
				if !ok {
					err = errors.New("不能断言成图文")
					return
				}

				priceIntegerPart, priceDecimalPart := separate(keyConvertKeyResp.Price)
				rebateIntegerPart, rebateDecimalPart := separate(keyConvertKeyResp.Rebate)
				var ArticleCount = 1
				news.SetArticleCount(ArticleCount)
				for i := 0; i < ArticleCount; i++ {
					news.SetTitle(i, fmt.Sprintf("约返:%s.%s  付费价:%s.%s", rebateIntegerPart, rebateDecimalPart, priceIntegerPart, priceDecimalPart))
					news.SetPicURL(i, keyConvertKeyResp.PicURL)
					news.SetDescription(i, keyConvertKeyResp.Title)
					news.SetURL(i, "http://suo.im/6qr3r4")
				}
				reply, err = ToEnvelope.MessageMarshal()
				return
			}
			text.SetContent(`<a href="weixin://bizmsgmenu?msgmenucontent=我要表情包&msgmenuid=1">我要表情包</a>
<a href="weixin://bizmsgmenu?msgmenucontent=我要消费券&msgmenuid=1">我要消费券</a>`)
		}
		reply, err = ToEnvelope.MessageMarshal()
		return
	}
	if voice, ok := rxMessage.Voice(); ok {
		mediaID := voice.GetMediaID()
		ToEnvelope, err = officialmodel.ToEnvelope(rxMessage, officialmodel.TxMessageTypeVoice)
		if err != nil {
			return
		}
		voice, ok := ToEnvelope.Voice()
		if !ok {
			err = errors.New("不能断言成语音")
		}
		voice.SetMediaID(mediaID)
		reply, err = ToEnvelope.MessageMarshal()
		return
	}

	err = errors.New("不会处理")
	return
}

func separate(number string) (integerPart string, decimalPart string) {
	switch len(number) {
	case 0:
		decimalPart = "00"
		integerPart = "0"
	case 1:
		decimalPart = "0" + number
		integerPart = "0"
	case 2:
		decimalPart = number
		integerPart = "0"
	default:
		integerPart = number[:len(number)-2]
		decimalPart = number[len(number)-2:]
	}

	return
}
