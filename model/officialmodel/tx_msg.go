package officialmodel

import (
	"encoding/xml"
)

type cdataNode struct {
	CData string `xml:",cdata"`
}

type viewTextMessage struct {
	XMLName      xml.Name  `xml:"xml"`
	ToUserName   cdataNode `xml:"ToUserName"`
	FromUserName cdataNode `xml:"FromUserName"`
	CreateTime   int64     `xml:"CreateTime"`
	MsgType      cdataNode `xml:"MsgType"`
	Content      cdataNode `xml:"Content"`
}
type xmlRxEncryptEnvelope struct {
	ToUserName string `xml:"ToUserName"`
	Encrypt    string `xml:"Encrypt"`
}
