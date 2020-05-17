# 接收消息格式

## Models

### `rxMessageCommon` 接收消息的公共部分

  Name           | XML            | Type     | Doc
  :------------- | :------------- | :------- | :-----------------------
  `ToUserName`   | `ToUserName`   | `string` | 开发者微信号
  `FromUserName` | `FromUserName` | `string` | 发送方帐号（一个OpenID）
  `CreateTime`   | `CreateTime`   | `int64`  | 消息创建时间 （整型）
  `MsgType`      | `MsgType`      | `MessageType` | 消息类型，文本为text
  `MsgID`        | `MsgId`        | `int64`  | 消息id，64位整

 ```go
// MessageType 消息类型
type MessageType string

// MessageTypeText 文本消息
const MessageTypeText MessageType = "text"

// MessageTypeImage 图片消息
const MessageTypeImage MessageType = "image"

// MessageTypeVoice 语音消息
const MessageTypeVoice MessageType = "voice"

// MessageTypeVideo 视频消息
const MessageTypeVideo MessageType = "video"

// MessageTypeShortVideo 小视频消息
const MessageTypeShortVideo MessageType = "shortvideo"

// MessageTypeLocation 位置消息
const MessageTypeLocation MessageType = "location"

// MessageTypeLink 链接消息
const MessageTypeLink MessageType = "link"
```

### `rxTextMessageSpecifics` 接收的文本消息，特有字段

Name|XML|Type|Doc
:---|:--|:---|:--
`Content`|`Content`|`string`|文本消息内容

### `rxImageMessageSpecifics` 接收的图片消息，特有字段

Name|XML|Type|Doc
:---|:--|:---|:--
`PicURL`|`PicUrl`|`string`|图片链接（由系统生成）
`MediaID`|`MediaId`|`string`|图片媒体文件id，可以调用获取媒体文件接口拉取，仅三天内有效

### `rxVoiceMessageSpecifics` 接收的语音消息，特有字段

Name|XML|Type|Doc
:---|:--|:---|:--
`MediaID`|`MediaId`|`string`|语音媒体文件id，可以调用获取媒体文件接口拉取数据，仅三天内有效
`Format`|`Format`|`string`|语音格式，如amr，speex等

### `rxVideoMessageSpecifics` 接收的视频消息，特有字段

Name|XML|Type|Doc
:---|:--|:---|:--
`MediaID`|`MediaId`|`string`|视频媒体文件id，可以调用获取媒体文件接口拉取数据，仅三天内有效
`ThumbMediaID`|`ThumbMediaId`|`string`|视频消息缩略图的媒体id，可以调用获取媒体文件接口拉取数据，仅三天内有效

### `rxShortVideoMessageSpecifics` 接收的小视频消息，特有字段

Name|XML|Type|Doc
:---|:--|:---|:--
`MediaID`|`MediaId`|`string`|视频媒体文件id，可以调用获取媒体文件接口拉取数据，仅三天内有效
`ThumbMediaID`|`ThumbMediaId`|`string`|视频消息缩略图的媒体id，可以调用获取媒体文件接口拉取数据，仅三天内有效

### `rxLocationMessageSpecifics` 接收的位置消息，特有字段

Name|XML|Type|Doc
:---|:--|:---|:--
`Lat`|`Location_X`|`float64`|地理位置纬度
`Lon`|`Location_Y`|`float64`|地理位置经度
`Scale`|`Scale`|`int`|地图缩放大小
`Label`|`Label`|`string`|地理位置信息

### `rxLinkMessageSpecifics` 接收的链接消息，特有字段

Name|XML|Type|Doc
:---|:--|:---|:--
`Title`|`Title`|`string`|标题
`Description`|`Description`|`string`|描述
`URL`|`Url`|`string`|链接跳转的url
