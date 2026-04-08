package token

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/hejingwen098/qyapi_weixin/pkg/config"
)

const (
	tokenURL       = "https://qyapi.weixin.qq.com/cgi-bin/gettoken"
	logoutURL      = "https://qyapi.weixin.qq.com/cgi-bin/user/logout"
	defaultTimeout = 10 * time.Second
)

// TokenInfo token 信息
type GettokenResponse struct {
	Token     string `json:"access_token"`
	ExpiresIn int64  `json:"expires_in"`
}

// Client 认证客户端
type Client struct {
	client *http.Client
	config *config.Config
}

// NewClient 创建认证客户端
func NewClient(client *http.Client, cfg *config.Config) *Client {
	return &Client{
		client: client,
		config: cfg,
	}
}

func (c *Client) GetToken() (string, error) {
	params := url.Values{}
	params.Add("corpid", c.config.CorpID)
	params.Add("corpsecret", c.config.CorpSecret)
	req, _ := http.NewRequest("GET", tokenURL+"?"+params.Encode(), nil)
	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	result := GettokenResponse{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	return result.Token, nil
}

func (c *Client) Logout(token string) error {
	params := url.Values{}
	params.Add("access_token", token)
	req, _ := http.NewRequest("GET", logoutURL+"?"+params.Encode(), nil)
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
