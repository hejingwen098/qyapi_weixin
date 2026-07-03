package message

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/hejingwen098/qyapi_weixin/pkg/errorx"
)

const (
	sendURL = "https://qyapi.weixin.qq.com/cgi-bin/message/send"
)

// SendResponse 发送消息响应
type SendResponse struct {
	errorx.QyError
	InvalidUser  string `json:"invaliduser"`  // 无效用户列表
	InvalidParty string `json:"invalidparty"` // 无效部门列表
	InvalidTag   string `json:"invalidtag"`   // 无效标签列表
}

// Client 消息客户端
type Client struct {
	client *http.Client
	token  *string
}

// NewClient 创建消息客户端
func NewClient(client *http.Client, token *string) *Client {
	return &Client{
		client: client,
		token:  token,
	}
}

// SendMessage 发送消息
// msgType: 消息类型，如 text、image、voice、video、file、textcard、news、mpnews、markdown、miniprogram_notice 等
// content: 消息内容，根据消息类型传入不同的结构体
func (c *Client) SendMessage(msgType string, content interface{}) (string, error) {
	// 构建完整的消息体
	message := map[string]interface{}{
		"msgtype": msgType,
	}

	// 将 content 转换为 map 并合并到 message 中
	contentBytes, err := json.Marshal(content)
	if err != nil {
		return "", fmt.Errorf("序列化消息内容失败: %v", err)
	}

	var contentMap map[string]interface{}
	if err := json.Unmarshal(contentBytes, &contentMap); err != nil {
		return "", fmt.Errorf("解析消息内容失败: %v", err)
	}

	for k, v := range contentMap {
		message[k] = v
	}

	// 序列化消息
	jsonData, err := json.Marshal(message)
	if err != nil {
		return "", fmt.Errorf("序列化消息失败: %v", err)
	}

	// 构建请求 URL
	params := url.Values{}
	params.Add("access_token", *c.token)
	reqURL := sendURL + "?" + params.Encode()

	// 发送请求
	resp, err := c.client.Post(reqURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("发送消息请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	// 解析响应
	result := SendResponse{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	// 检查错误码
	if result.Code != 0 {
		return "", fmt.Errorf("企业微信返回错误 [%d]: %s", result.Code, result.Msg)
	}

	return string(body), nil
}
