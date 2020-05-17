package restserver

import (
	"bytes"
	"crypto/subtle"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"unsafe"

	"go.uber.org/zap"

	"github.com/Naist4869/awesomeProject/tool"
	"github.com/Naist4869/log"
)

// Server 用于处理微信服务器的回调请求, 并发安全!
//  通常情况下一个 Server 实例用于处理一个公众号的消息(事件), 此时建议指定 oriId(原始ID) 和 appId(明文模式下无需指定) 用于约束消息(事件);
//  特殊情况下也可以一个 Server 实例用于处理多个公众号的消息(事件), 此时要求这些公众号的 token 是一样的, 并且 oriId 和 appId 必须设置为 "".
type Server struct {
	oriId string
	appId string

	tokenBucketPtrMutex sync.Mutex     // used only by writers
	tokenBucketPtr      unsafe.Pointer // *tokenBucket

	aesKeyBucketPtrMutex sync.Mutex     // used only by writers
	aesKeyBucketPtr      unsafe.Pointer // *aesKeyBucket
	handler              IRouter        // http句柄

	logger *log.Logger
}

func (s *Server) OriId() string {
	return s.oriId
}
func (s *Server) AppId() string {
	return s.appId
}

type tokenBucket struct {
	currentToken string
	lastToken    string
}

type aesKeyBucket struct {
	currentAESKey []byte
	lastAESKey    []byte
}

// NewServer 创建一个新的 Server.
//  oriId:        可选; 公众号的原始ID(微信公众号管理后台查看), 如果设置了值则该Server只能处理 ToUserName 为该值的公众号的消息(事件);
//  appId:        可选; 公众号的AppId, 如果设置了值则安全模式时该Server只能处理 AppId 为该值的公众号的消息(事件);
//  token:        必须; 公众号用于验证签名的token;
//  base64AESKey: 可选; aes加密解密key, 43字节长(base64编码, 去掉了尾部的'='), 安全模式必须设置;
//  handler:      必须; 处理微信服务器推送过来的消息(事件)的Handler;
//  errorHandler: 可选; 用于处理Server在处理消息(事件)过程中产生的错误, 如果没有设置则默认使用 DefaultErrorHandler.
func NewServer(oriId, appId, token, base64AESKey string, logger *log.Logger, handler IRouter) *Server {
	if token == "" {
		panic("empty token")
	}

	var (
		aesKey []byte
		err    error
	)
	if base64AESKey != "" {
		if len(base64AESKey) != 43 {
			panic("the length of base64AESKey must equal to 43")
		}
		aesKey, err = base64.StdEncoding.DecodeString(base64AESKey + "=")
		if err != nil {
			panic(fmt.Sprintf("Decode base64AESKey:%s failed", base64AESKey))
		}
	}

	return &Server{
		oriId:           oriId,
		appId:           appId,
		tokenBucketPtr:  unsafe.Pointer(&tokenBucket{currentToken: token}),
		aesKeyBucketPtr: unsafe.Pointer(&aesKeyBucket{currentAESKey: aesKey}),
		handler:         handler,
		logger:          logger,
	}
}
func (s *Server) getToken() (currentToken, lastToken string) {
	if p := (*tokenBucket)(atomic.LoadPointer(&s.tokenBucketPtr)); p != nil {
		return p.currentToken, p.lastToken
	}
	return
}

// SetToken 设置签名token.
func (s *Server) SetToken(token string) (err error) {
	if token == "" {
		return errors.New("empty token")
	}

	s.tokenBucketPtrMutex.Lock()
	defer s.tokenBucketPtrMutex.Unlock()

	currentToken, _ := s.getToken()
	if token == currentToken {
		return
	}

	bucket := tokenBucket{
		currentToken: token,
		lastToken:    currentToken,
	}
	atomic.StorePointer(&s.tokenBucketPtr, unsafe.Pointer(&bucket))
	return
}
func (s *Server) removeLastToken(lastToken string) {
	s.tokenBucketPtrMutex.Lock()
	defer s.tokenBucketPtrMutex.Unlock()

	currentToken2, lastToken2 := s.getToken()
	if lastToken != lastToken2 {
		return
	}

	bucket := tokenBucket{
		currentToken: currentToken2,
	}
	atomic.StorePointer(&s.tokenBucketPtr, unsafe.Pointer(&bucket))
	return
}

func (s *Server) getAESKey() (currentAESKey, lastAESKey []byte) {
	if p := (*aesKeyBucket)(atomic.LoadPointer(&s.aesKeyBucketPtr)); p != nil {
		return p.currentAESKey, p.lastAESKey
	}
	return
}

// SetAESKey 设置aes加密解密key.
//  base64AESKey: aes加密解密key, 43字节长(base64编码, 去掉了尾部的'=').
func (s *Server) SetAESKey(base64AESKey string) (err error) {
	if len(base64AESKey) != 43 {
		return errors.New("the length of base64AESKey must equal to 43")
	}
	aesKey, err := base64.StdEncoding.DecodeString(base64AESKey + "=")
	if err != nil {
		return
	}

	s.aesKeyBucketPtrMutex.Lock()
	defer s.aesKeyBucketPtrMutex.Unlock()

	currentAESKey, _ := s.getAESKey()
	if bytes.Equal(aesKey, currentAESKey) {
		return
	}

	bucket := aesKeyBucket{
		currentAESKey: aesKey,
		lastAESKey:    currentAESKey,
	}
	atomic.StorePointer(&s.aesKeyBucketPtr, unsafe.Pointer(&bucket))
	return
}
func (s *Server) removeLastAESKey(lastAESKey []byte) {
	s.aesKeyBucketPtrMutex.Lock()
	defer s.aesKeyBucketPtrMutex.Unlock()

	currentAESKey2, lastAESKey2 := s.getAESKey()
	if !bytes.Equal(lastAESKey, lastAESKey2) {
		return
	}

	bucket := aesKeyBucket{
		currentAESKey: currentAESKey2,
	}
	atomic.StorePointer(&s.aesKeyBucketPtr, unsafe.Pointer(&bucket))
	return
}

// string转换[]byte 只转类型不转数据 https://www.cnblogs.com/shuiyuejiangnan/p/9707066.html
func (s *Server) Signature(c *Context) {
	query := c.Request.URL.Query()
	haveSignature := []byte(query.Get("signature"))
	if len(haveSignature) == 0 {
		c.Error(errors.New("not found signature query parameter"))
		return
	}
	timestamp := query.Get("timestamp")
	if timestamp == "" {
		c.Error(errors.New("not found timestamp query parameter"))
		return
	}
	nonce := query.Get("nonce")
	if nonce == "" {
		c.Error(errors.New("not found nonce query parameter"))
		return
	}
	echostr := query.Get("echostr")
	if echostr == "" {
		c.Error(errors.New("not found echostr query parameter"))
		return
	}
	var token string
	currentToken, lastToken := s.getToken()
	if currentToken == "" {
		c.Error(errors.New("token was not set for Server, see NewServer function or Server.SetToken method"))
		return
	}
	token = currentToken
	wantSignature := []byte(Sign(token, timestamp, nonce))

	if subtle.ConstantTimeCompare(haveSignature, wantSignature) != 1 {
		if lastToken == "" {
			c.Error(fmt.Errorf("check signature failed, have: %s, want: %s", haveSignature, wantSignature))
			return
		}
		token = lastToken
		wantSignature = []byte(Sign(token, timestamp, nonce))
		if subtle.ConstantTimeCompare(haveSignature, wantSignature) != 1 {
			c.Error(fmt.Errorf("check signature failed, have: %s, want: %s", haveSignature, wantSignature))
			return
		}
	} else if lastToken != "" {
		s.removeLastToken(lastToken)
	}
	c.String(http.StatusOK, "%s", echostr)
}

func (s *Server) Http() {
	s.handler.GET("/callback", s.Signature)
	s.handler.POST("/callback", s.Verify, s.UseCase)
}
func (s *Server) Verify(c *Context) {
	query := c.Request.URL.Query()
	encryptType := query.Get("encrypt_type")
	switch encryptType {
	case "aes":
		haveSignature := []byte(query.Get("signature"))
		if len(haveSignature) == 0 {
			c.Error(errors.New("not found signature query parameter"))
			return
		}
		haveMsgSignature := []byte(query.Get("msg_signature"))
		if len(haveMsgSignature) == 0 {
			c.Error(errors.New("not found msg_signature query parameter"))
			return
		}
		timestamp := query.Get("timestamp")
		if timestamp == "" {
			c.Error(errors.New("not found timestamp query parameter"))
			return
		}
		nonce := query.Get("nonce")
		if nonce == "" {
			c.Error(errors.New("not found nonce query parameter"))
			return
		}
		var token string
		currentToken, lastToken := s.getToken()
		if currentToken == "" {
			c.Error(errors.New("token was not set for Server, see NewServer function or Server.SetToken method"))
			return
		}
		token = currentToken
		wantSignature := []byte(Sign(token, timestamp, nonce))
		if subtle.ConstantTimeCompare(haveSignature, wantSignature) != 1 {
			if lastToken == "" {
				c.Error(fmt.Errorf("check signature failed, have: %s, want: %s", haveSignature, wantSignature))
				return
			}
			token = lastToken
			wantSignature = []byte(Sign(token, timestamp, nonce))
			if subtle.ConstantTimeCompare(haveSignature, wantSignature) != 1 {
				c.Error(fmt.Errorf("check signature failed, have: %s, want: %s", haveSignature, wantSignature))
				return
			}
		} else if lastToken != "" {
			s.removeLastToken(lastToken)
		}
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.Error(errors.New("read request body fail"))
			return
		}
		defer c.Request.Body.Close()
		data := &xmlRxEnvelope{}
		if err := xml.Unmarshal(body, data); err != nil {
			c.Error(errors.New("xmlRxEnvelope unmarshal fail"))
			return
		}
		haveToUserName := data.ToUserName
		wantToUserName := s.oriId
		if strings.Compare(haveToUserName, wantToUserName) != 0 {
			c.Error(fmt.Errorf("the message ToUserName mismatch, have: %s, want: %s",
				haveToUserName, wantToUserName))
			return
		}
		wantMsgSignature := []byte(Sign(token, timestamp, nonce, data.Encrypt))
		if subtle.ConstantTimeCompare(haveMsgSignature, wantMsgSignature) != 1 {
			c.Error(fmt.Errorf("check msg_signature failed, have: %s, want: %s", haveMsgSignature, wantMsgSignature))
			return
		}
		var aesKey []byte
		currentAESKey, lastAESKey := s.getAESKey()
		if currentAESKey == nil {
			c.Error(errors.New("aes key was not set for Server, see NewServer function or Server.SetAESKey method"))
			return
		}
		aesKey = currentAESKey
		random, xmlMsg, haveAppIdBytes, err := AESDecryptMsg(tool.StrToBytes(data.Encrypt), aesKey)
		if err != nil {
			if lastAESKey == nil {
				c.Error(err)
				return
			}
			aesKey = lastAESKey
			random, xmlMsg, haveAppIdBytes, err = AESDecryptMsg(tool.StrToBytes(data.Encrypt), aesKey)
			if err != nil {
				c.Error(err)
				return
			}
		} else {
			if lastAESKey != nil {
				s.removeLastAESKey(lastAESKey)
			}
		}
		wantAppId := s.appId
		haveAppId := string(haveAppIdBytes)
		if len(wantAppId) != 0 && strings.Compare(haveAppId, wantAppId) != 0 {
			c.Error(fmt.Errorf("the message AppId mismatch, have: %s, want: %s", haveAppId, wantAppId))
			return
		}
		c.Set("random", random)
		c.Set("xmlMsg", xmlMsg)
	case "", "raw":
		haveSignature := []byte(query.Get("signature"))
		if len(haveSignature) == 0 {
			c.Error(errors.New("not found signature query parameter"))
			return
		}
		timestamp := query.Get("timestamp")
		if timestamp == "" {
			c.Error(errors.New("not found timestamp query parameter"))
			return
		}
		nonce := query.Get("nonce")
		if nonce == "" {
			c.Error(errors.New("not found nonce query parameter"))
			return
		}
		var token string
		currentToken, lastToken := s.getToken()
		if currentToken == "" {
			c.Error(errors.New("token was not set for Server, see NewServer function or Server.SetToken method"))
			return
		}
		token = currentToken
		wantSignature := []byte(Sign(token, timestamp, nonce))
		if subtle.ConstantTimeCompare(haveSignature, wantSignature) != 1 {
			if lastToken == "" {
				c.Error(fmt.Errorf("check signature failed, have: %s, want: %s", haveSignature, wantSignature))
				return
			}
			token = lastToken
			wantSignature = []byte(Sign(token, timestamp, nonce))
			if subtle.ConstantTimeCompare(haveSignature, wantSignature) != 1 {
				c.Error(fmt.Errorf("check signature failed, have: %s, want: %s", haveSignature, wantSignature))
				return
			}
		} else if lastToken != "" {
			s.removeLastToken(lastToken)
		}
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.Error(errors.New("read request body fail"))
			return
		}
		defer c.Request.Body.Close()
		c.Set("xmlMsg", body)
	default:
		c.Error(errors.New("unknown encrypt_type: " + encryptType))
		return
	}
}
func (s *Server) UseCase(c *Context) {
	xmlMsg, exists := c.Get("xmlMsg")
	if !exists {
		// view
		return
	}
	c.Get("random")
	s.logger.Info("收到消息", zap.ByteString("XML消息", xmlMsg.([]byte)))
}

type xmlRxEnvelope struct {
	ToUserName string `xml:"ToUserName"`
	Encrypt    string `xml:"Encrypt"`
}

type EnvelopeHandler interface {
	OnIncomingEnvelope(xmlMsg []byte) error
}
