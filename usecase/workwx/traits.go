package workwx

import (
	"net/url"

	"github.com/Naist4869/awesomeProject/model/wxmodel"
)

type urlValuer interface {
	intoURLValues() url.Values
}

type bodyer interface {
	intoBody() ([]byte, error)
}

type mediaUploader interface {
	getMedia() *wxmodel.Media
}
