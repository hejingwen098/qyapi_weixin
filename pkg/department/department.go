package department

import (
	"encoding/json"
	"net/http"
	"net/url"

	"qyapi_weixin/pkg/errorx"
)

const (
	listURL = "https://qyapi.weixin.qq.com/cgi-bin/department/list"
)

// Department 部门信息
type Department struct {
	ID         int64  `json:"id"`          // 部门 ID
	Name       string `json:"name"`        // 部门名称
	NameEn     string `json:"name_en"`     // 英文名称
	ParentID   int64  `json:"parentid"`    // 父部门 ID
	Order      uint32 `json:"order"`       // 排序
	IsInner    bool   `json:"is_inner"`    // 是否内部部门
	UpdateTime uint64 `json:"update_time"` // 更新时间
}

// Client 部门查询客户端
type Client struct {
	client *http.Client
	token  *string // 企业微信 access_token
}

// NewClient 创建部门查询客户端
func NewClient(client *http.Client, token *string) *Client {
	return &Client{
		client: client,
		token:  token,
	}
}

// ListResponse 部门列表响应
type ListResponse struct {
	errorx.QyError
	Departments []Department `json:"department"`
}

// ListAll 获取所有部门详情
func (c *Client) ListAll() ([]Department, error) {
	params := url.Values{}
	params.Add("access_token", *c.token)
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
	return result.Departments, nil
}

// GetByID 根据 ID 获取部门
func (c *Client) GetByID(deptID int64) (*Department, error) {
	depts, err := c.ListAll()
	if err != nil {
		return nil, err
	}

	for _, dept := range depts {
		if dept.ID == deptID {
			return &dept, nil
		}
	}

	return nil, errorx.ErrDepartmentNotFound
}

// GetByParentID 获取指定父部门下的子部门
func (c *Client) GetByParentID(token string, parentID int64) ([]Department, error) {
	depts, err := c.ListAll()
	if err != nil {
		return nil, err
	}

	var result []Department
	for _, dept := range depts {
		if dept.ParentID == parentID {
			result = append(result, dept)
		}
	}

	return result, nil
}
