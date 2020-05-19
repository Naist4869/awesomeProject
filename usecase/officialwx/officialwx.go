package officialwx

import (
	"fmt"

	"github.com/Naist4869/awesomeProject/model/officialmodel"
	"github.com/Naist4869/awesomeProject/usecase"
)

type UseCase struct {
}

func (u *UseCase) ReplyMessage(xmlMsg []byte) (reply string, err error) {
	rxMessage, err := officialmodel.FromEnvelope(xmlMsg)
	if err != nil {
		err = fmt.Errorf("解析消息失败: %w", err)
		return
	}
	text, ok := rxMessage.Text()
	if !ok {
		err = usecase.ErrMessageTypeAssertionFail
		return
	}
	content := text.GetContent()
	text.SetContent(content + "!!!")
	return rxMessage.MessageMarshal()
}
