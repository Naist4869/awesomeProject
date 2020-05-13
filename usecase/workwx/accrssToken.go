package workwx

import "time"

// getAccessToken 获取 access token
func (c *WorkwxApp) getAccessToken() (respAccessToken, error) {
	return c.execGetAccessToken(reqAccessToken{
		CorpID:     c.CorpID,
		CorpSecret: c.CorpSecret,
	})
}

// syncAccessToken 同步该客户端实例的 access token
//
// 会拿 `tokenMu` 写锁
func (c *WorkwxApp) syncAccessToken() error {
	tok, err := c.getAccessToken()
	if err != nil {
		// TODO: error_chain
		return err
	}

	c.tokenMu.Lock()
	defer c.tokenMu.Unlock()

	c.accessToken = tok.AccessToken
	c.tokenExpiresIn = time.Duration(tok.ExpiresInSecs) * time.Second
	c.lastRefresh = time.Now()

	return nil
}
