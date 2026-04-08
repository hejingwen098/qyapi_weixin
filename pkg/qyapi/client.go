package qyapi

import (
	"net/http"

	"github.com/hejingwen098/qyapi_weixin/pkg/config"
	"github.com/hejingwen098/qyapi_weixin/pkg/department"
	"github.com/hejingwen098/qyapi_weixin/pkg/token"
	"github.com/hejingwen098/qyapi_weixin/pkg/user"
)

// QyClient 企业微信客户端
type QyClient struct {
	Config      *config.Config
	Client      *http.Client
	TokenClient *token.Client
	Token       string
	DeptClient  *department.Client
	UserClient  *user.Client
}

// NewQyClient 创建企业微信客户端
func NewQyClient(corpID, corpSecret string) (*QyClient, error) {
	cfg := config.NewConfig(corpID, corpSecret)
	client := http.DefaultClient
	tokenClient := token.NewClient(client, cfg)
	token, err := tokenClient.GetToken()
	if err != nil {
		return nil, err
	}
	deptClient := department.NewClient(client, &token)
	userClient := user.NewClient(client, &token)
	return &QyClient{
		Config:      cfg,
		Client:      client,
		TokenClient: tokenClient,
		Token:       token,
		DeptClient:  deptClient,
		UserClient:  userClient,
	}, nil
}

// GetAllDepartments 获取所有部门详情
func (c *QyClient) GetAllDepartments() ([]department.Department, error) {
	return c.DeptClient.ListAll()
}

// GetDepartmentByID 根据 ID 获取部门
func (c *QyClient) GetDepartmentByID(deptID int64) (*department.Department, error) {
	return c.DeptClient.GetByID(deptID)
}

// GetUsersByDeptID 获取指定部门的成员列表（详细信息）
func (c *QyClient) GetUsersByDeptID(deptID int64) ([]user.User, error) {
	return c.UserClient.ListByDept(deptID)
}

// GetSimpleUsersByDeptID 获取指定部门的成员列表（简单信息）
func (c *QyClient) GetSimpleUsersByDeptID(deptID int64) ([]user.SimpleUser, error) {
	return c.UserClient.SimpleListByDept(deptID)
}

// GetUserByUserID 根据 UserID 获取成员详情
func (c *QyClient) GetUserByUserID(userID string) (*user.User, error) {
	return c.UserClient.GetByUserID(userID)
}
