package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"qyapi_weixin/pkg/errorx"
)

const (
	simpleListURL = "https://qyapi.weixin.qq.com/cgi-bin/user/simplelist"
	listURL       = "https://qyapi.weixin.qq.com/cgi-bin/user/list"
	getURL        = "https://qyapi.weixin.qq.com/cgi-bin/user/get"
)

// SimpleUser 简单成员信息
type SimpleUser struct {
	UserID   string `json:"userid"`     // 用户 ID
	Name     string `json:"name"`       // 姓名
	Position string `json:"position"`   // 职位
	DeptIDs  []int  `json:"department"` // 部门 ID 列表
}

// User 详细成员信息
type User struct {
	UserID         string  `json:"userid"`          // 用户 ID
	Name           string  `json:"name"`            // 姓名
	Position       string  `json:"position"`        // 职位
	Departments    []int64 `json:"department"`      // 部门 ID 列表
	MainDepartment int64   `json:"main_department"` // 主部门
	Email          string  `json:"email"`           // 邮箱
	IsLeader       bool    `json:"is_leader"`       // 是否上级
	Mobile         string  `json:"mobile"`          // 手机号
	Gender         string  `json:"gender"`          // 性别
	Avatar         string  `json:"avatar"`          // 头像
	WeChatAccount  string  `json:"wechat_account"`  // 微信账号
	Status         int     `json:"status"`          // 状态：1=已关注，2=已禁用，4=未关注
	ExtAttr        ExtAttr `json:"extattr"`         // 扩展属性
}

// ExtAttr 扩展属性
type ExtAttr struct {
	Attrs []ExtAttrItem `json:"attrs"`
}

// ExtAttrItem 扩展属性项
type ExtAttrItem struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// Client 成员查询客户端
type Client struct {
	client *http.Client
	token  *string
}

// NewClient 创建成员查询客户端
func NewClient(client *http.Client, token *string) *Client {
	return &Client{
		client: client,
		token:  token,
	}
}

// SimpleListResponse 简单成员列表响应
type SimpleListResponse struct {
	errorx.QyError
	Users []SimpleUser `json:"userlist"`
}

// ListResponse 详细成员列表响应
type ListResponse struct {
	errorx.QyError
	Users []User `json:"userlist"`
}

// SimpleListByDept 获取指定部门的成员列表（简单信息）
func (c *Client) SimpleListByDept(deptID int64) ([]SimpleUser, error) {
	params := url.Values{}
	params.Add("access_token", *c.token)
	params.Add("department_id", fmt.Sprintf("%d", deptID))
	req, _ := http.NewRequest("GET", simpleListURL+"?"+params.Encode(), nil)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result := SimpleListResponse{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result.Users, nil
}

// ListByDept 获取指定部门的成员列表（详细信息）
func (c *Client) ListByDept(deptID int64) ([]User, error) {
	params := url.Values{}
	params.Add("access_token", *c.token)
	params.Add("department_id", fmt.Sprintf("%d", deptID))
	req, _ := http.NewRequest("GET", listURL+"?"+params.Encode(), nil)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result := ListResponse{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.Users, nil
}

// GetByUserID 根据 UserID 获取成员详情
func (c *Client) GetByUserID(userID string) (*User, error) {
	params := url.Values{}
	params.Add("access_token", *c.token)
	params.Add("userid", userID)
	req, _ := http.NewRequest("GET", getURL+"?"+params.Encode(), nil)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	result := User{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
