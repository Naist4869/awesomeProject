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

type xmlTxEnvelope struct {
	XMLName      xml.Name `xml:"xml"`
	Encrypt      string   `xml:"Encrypt"`
	MsgSignature string   `xml:"MsgSignature"`
	Timestamp    int64    `xml:"TimeStamp"`
	Nonce        string   `xml:"Nonce"`
}
