package officialmodel

type MediaIDReq struct {
	FakeID    string
	Timestamp int64
}

type MediaIDResp struct {
	MediaID string
}
type TitleConvertTBKeyReq struct {
	Title string
}
type TitleConvertTBKeyResp struct {
	TBKey string
}

type KeyConvertKeyReq struct {
	FromKey string
	UserID  string
}
type KeyConvertKeyResp struct {
	ToKey  string
	Price  string
	Rebate string
	Title  string
	PicURL string
}
