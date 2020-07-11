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
	"net/url"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/Naist4869/awesomeProject/api"

	"github.com/Naist4869/awesomeProject/usecase"

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

	secretBucketPtrMutex sync.Mutex     // used only by writers
	secretBucketPtr      unsafe.Pointer // *secretBucket

	handler IRouter // http句柄

	wx usecase.IOfficialWx

	logger *log.Logger
}

func (s *Server) OriId() string {
	return s.oriId
}
func (s *Server) AppId() string {
	return s.appId
}

type secretBucket struct {
	currentSecret string
	lastSecret    string
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
func NewServer(oriId, appId, token, base64AESKey, secret string, logger *log.Logger, useCaseImpl usecase.IOfficialWx, handler IRouter) *Server {
	s := &Server{
		oriId:   oriId,
		appId:   appId,
		handler: handler,
		wx:      useCaseImpl,
		logger:  logger,
	}
	if err := s.SetAESKey(base64AESKey); err != nil {
		panic(err)
	}
	if err := s.SetToken(token); err != nil {
		panic(err)
	}
	if err := s.SetSecret(secret); err != nil {
		panic(err)
	}
	s.Http()
	return s
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
func (s *Server) getSecret() (currentSecret, lastSecret string) {
	if p := (*secretBucket)(atomic.LoadPointer(&s.secretBucketPtr)); p != nil {
		return p.currentSecret, p.lastSecret
	}
	return
}

func (s *Server) SetSecret(secret string) (err error) {
	if secret == "" {
		return errors.New("empty secret")
	}

	s.secretBucketPtrMutex.Lock()
	defer s.secretBucketPtrMutex.Unlock()

	currentSecret, _ := s.getSecret()
	if secret == currentSecret {
		return
	}

	bucket := secretBucket{
		currentSecret: secret,
		lastSecret:    currentSecret,
	}
	atomic.StorePointer(&s.secretBucketPtr, unsafe.Pointer(&bucket))
	return
}
func (s *Server) removeLastSecret(lastSecret string) {
	s.secretBucketPtrMutex.Lock()
	defer s.secretBucketPtrMutex.Unlock()

	currentSecret2, lastSecret2 := s.getSecret()
	if lastSecret != lastSecret2 {
		return
	}

	bucket := secretBucket{
		currentSecret: currentSecret2,
	}
	atomic.StorePointer(&s.secretBucketPtr, unsafe.Pointer(&bucket))
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

func (s *Server) Router() IRouter {
	return s.handler
}

func (s *Server) Verify(c *Context) {
	query := c.Request.URL.Query()
	encryptType := query.Get("encrypt_type")
	switch encryptType {
	case "aes":
		if !s.aesVerity(c, query) {
			return
		}
	case "", "raw":
		if !s.rawVerity(c, query) {
			return
		}

	default:
		c.Error(errors.New("unknown encrypt_type: " + encryptType))
		return
	}
}

func (s *Server) rawVerity(c *Context, query url.Values) bool {
	haveSignature := []byte(query.Get("signature"))
	if len(haveSignature) == 0 {
		c.Error(errors.New("not found signature query parameter"))
		return true
	}
	timestamp := query.Get("timestamp")
	if timestamp == "" {
		c.Error(errors.New("not found timestamp query parameter"))
		return true
	}
	nonce := query.Get("nonce")
	if nonce == "" {
		c.Error(errors.New("not found nonce query parameter"))
		return true
	}
	var token string
	currentToken, lastToken := s.getToken()
	if currentToken == "" {
		c.Error(errors.New("token was not set for Server, see NewServer function or Server.SetToken method"))
		return true
	}
	token = currentToken
	wantSignature := []byte(Sign(token, timestamp, nonce))
	if subtle.ConstantTimeCompare(haveSignature, wantSignature) != 1 {
		if lastToken == "" {
			c.Error(fmt.Errorf("check signature failed, have: %s, want: %s", haveSignature, wantSignature))
			return true
		}
		token = lastToken
		wantSignature = []byte(Sign(token, timestamp, nonce))
		if subtle.ConstantTimeCompare(haveSignature, wantSignature) != 1 {
			c.Error(fmt.Errorf("check signature failed, have: %s, want: %s", haveSignature, wantSignature))
			return true
		}
	} else if lastToken != "" {
		s.removeLastToken(lastToken)
	}
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.Error(errors.New("read request body fail"))
		return true
	}
	defer c.Request.Body.Close()
	c.Set("xmlMsg", body)
	c.Set("token", token)
	return false
}

func (s *Server) aesVerity(c *Context, query url.Values) (ok bool) {
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
	data := &xmlRxEncryptEnvelope{}
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
	c.Set("aesKey", aesKey)
	c.Set("token", token)
	ok = true
	return
}
func (s *Server) UseCase(c *Context) {
	var (
		err          error
		replyMessage string
		signature    string
		ok           bool
	)
	encrypt, xmlMsg, token, random, aesKey, e := s.prepare(c)
	if e != nil {
		//todo view
		return
	}
	s.logger.Info("回复消息", zap.ByteString("收到消息", xmlMsg), zap.Bool("是否需要加密回复", encrypt))
	replyMessage, err = s.wx.ReplyMessage(xmlMsg)
	if err != nil {
		s.logger.Error("回复消息", zap.Error(err))
		return
	}
	s.logger.Info("回复消息", zap.String("回复消息", replyMessage))
	nonce := makeNonce()
	unix := time.Now().Unix()
	timestamp := strconv.FormatInt(unix, 10)
	if !encrypt {
		signature = Sign(token, timestamp, nonce)
		c.Bytes(http.StatusOK, "application/xml; charset=utf-8", []byte(replyMessage))
		return
	}
	replyMessage, signature, ok = s.encryptMessage(random, replyMessage, aesKey, signature, token, timestamp, nonce)
	if !ok {
		//todo view
		return
	}
	envelope := &xmlTxEncryptEnvelope{
		Encrypt:      replyMessage,
		MsgSignature: signature,
		Timestamp:    unix,
		Nonce:        nonce,
	}
	c.XML(http.StatusOK, envelope)
}

func (s *Server) encryptMessage(random []byte, replyMessage string, aesKey []byte, signature string, token string, timestamp string, nonce string) (string, string, bool) {
	base64EncryptedMsg, err := AESEncryptMsg(random, []byte(replyMessage), s.appId, aesKey)
	if err != nil {
		s.logger.Error("加密消息错误", zap.Error(err))
		return "", "", false
	}
	replyMessage = base64EncryptedMsg
	signature = Sign(token, timestamp, nonce, replyMessage)
	return replyMessage, signature, true
}

func (s *Server) prepare(c *Context) (encrypt bool, xmlMsg []byte, token string, random []byte, aesKey []byte, err error) {
	encrypt = true
	if xm, exists := c.Get("xmlMsg"); !exists {
		err = errors.New("缺少必要参数")
		return
	} else {
		xmlMsg = xm.([]byte)
	}

	if t, exists := c.Get("token"); !exists {
		err = errors.New("缺少必要参数")
		return
	} else {
		token = t.(string)

	}

	if r, exists := c.Get("random"); !exists {
		encrypt = false
	} else {
		random = r.([]byte)

	}

	if ak, exists := c.Get("aesKey"); !exists {
		encrypt = false
	} else {
		aesKey = ak.([]byte)
	}
	return
}
func (s *Server) HTTPApis() []api.Api {
	return []api.Api{
		api.NewApi("获取Secret", api.APIGetSecret, func() interface{} {
			return &SecretReq{}
		}, func(argument interface{}, merchantID int64) (interface{}, error) {
			secret, _ := s.getSecret()
			return SecretResp{
				AppID:  s.appId,
				Secret: secret,
			}, nil
		}),
	}
}

type xmlTxEncryptEnvelope struct {
	XMLName      xml.Name `xml:"xml"`
	Encrypt      string   `xml:"Encrypt"`
	MsgSignature string   `xml:"MsgSignature"`
	Timestamp    int64    `xml:"TimeStamp"`
	Nonce        string   `xml:"Nonce"`
}
type xmlRxEncryptEnvelope struct {
	ToUserName string `xml:"ToUserName"`
	Encrypt    string `xml:"Encrypt"`
}
type SecretResp struct {
	AppID  string `json:"appid"`
	Secret string `json:"secret"`
}
type SecretReq struct {
}
