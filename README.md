# QyAPI Weixin

企业微信 Go SDK - 简洁高效的企业微信 API 客户端库

## 📖 简介

`qyapi_weixin` 是一个用 Go 语言开发的企业微信 API 客户端库。它提供了简洁的接口来访问企业微信的核心功能，包括成员管理、部门管理等。通过封装复杂的 HTTP 请求和认证流程，让开发者能够专注于业务逻辑的实现。

## ✨ 特性

- **简洁的 API 设计**：提供直观的链式调用接口
- **自动 Token 管理**：内置 access_token 获取和管理机制
- **类型安全**：完整的类型定义，编译期错误检查
- **模块化设计**：按功能模块划分，易于扩展和维护
- **错误处理完善**：统一的企业微信错误处理机制
- **零外部依赖**：仅依赖标准库

## 🚀 快速开始

### 安装

```bash
go get qyapi_weixin
```

### 基本使用

```go
package main

import (
    "flag"
    "fmt"
    "log"

    "qyapi_weixin/pkg/qyapi"
)

func main() {
    // 配置企业微信凭证
    corpID := flag.String("corpid", "your_corp_id", "企业 ID")
    corpSecret := flag.String("corpsecret", "your_corp_secret", "应用凭证密钥")
    flag.Parse()

    // 创建客户端（自动完成认证）
    client, err := qyapi.NewQyClient(*corpID, *corpSecret)
    if err != nil {
        log.Fatalf("创建客户端失败：%v", err)
    }

    // 获取所有部门
    depts, err := client.GetAllDepartments()
    if err != nil {
        log.Fatalf("获取部门失败：%v", err)
    }
    fmt.Printf("部门数量：%d\n", len(depts))
    for _, dept := range depts {
        fmt.Printf("  - ID:%d, Name:%s, ParentID:%d\n", dept.ID, dept.Name, dept.ParentID)
        
        // 获取部门成员列表
        users, _ := client.GetUsersByDeptID(dept.ID)
        fmt.Printf("  - 用户数量：%d\n", len(users))
        for _, user := range users {
            fmt.Printf("    - 姓名：%s, UserID: %s, 部门：%v\n", user.Name, user.UserID, user.Departments[0])
        }
    }

    // 退出时清理 token
    defer client.TokenClient.Logout(client.Token)
}
```

### 运行示例

```bash
go run examples/main.go -corpid "your_corp_id" -corpsecret "your_corp_secret"
```

## 📦 功能模块

### 1. 认证模块 (`pkg/token`)

负责企业微信 access_token 的获取和管理。

```go
// 获取 access_token
token, err := client.TokenClient.GetToken()

// 注销 token（可选）
err = client.TokenClient.Logout(token)
```

### 2. 部门管理模块 (`pkg/department`)

提供部门相关的查询功能。

```go
// 获取所有部门
depts, err := client.GetAllDepartments()

// 根据 ID 获取部门
dept, err := client.GetDepartmentByID(1)

// 部门信息结构
Department {
    ID         int64  // 部门 ID
    Name       string // 部门名称
    NameEn     string // 英文名称
    ParentID   int64  // 父部门 ID
    Order      uint32 // 排序
    IsInner    bool   // 是否内部部门
    UpdateTime uint64 // 更新时间
}
```

### 3. 成员管理模块 (`pkg/user`)

提供成员的查询功能，支持简单信息和详细信息两种模式。

```go
// 获取指定部门的成员列表（详细信息）
users, err := client.GetUsersByDeptID(deptID)

// 根据 UserID 获取成员详情
user, err := client.GetUserByUserID("zhangsan")

// 成员信息结构（详细）
User {
    UserID         string  // 用户 ID
    Name           string  // 姓名
    Position       string  // 职位
    Departments    []int64 // 部门 ID 列表
    MainDepartment int64   // 主部门
    Email          string  // 邮箱
    IsLeader       bool    // 是否上级
    Mobile         string  // 手机号
    Gender         string  // 性别
    Avatar         string  // 头像
    WeChatAccount  string  // 微信账号
    Status         int     // 状态：1=已关注，2=已禁用，4=未关注
    ExtAttr        ExtAttr // 扩展属性
}
```

## 🏗️ 项目结构

```
qyapi_weixin/
├── pkg/                      # 核心代码包
│   ├── qyapi/               # 主客户端入口
│   │   └── client.go        # QyClient 实现
│   ├── config/              # 配置管理
│   │   └── config.go        # Config 结构定义
│   ├── token/               # 认证模块
│   │   └── token.go         # Token 获取和管理
│   ├── department/          # 部门管理模块
│   │   └── department.go    # 部门 CRUD 操作
│   ├── user/                # 成员管理模块
│   │   └── user.go          # 成员 CRUD 操作
│   └── errorx/              # 错误处理
│       └── error.go         # 自定义错误类型
├── examples/                 # 示例代码
│   └── main.go              # 完整使用示例
├── go.mod                    # Go 模块定义
└── README.md                 # 项目文档
```

## 🔧 核心类型

### QyClient

企业微信客户端，提供所有 API 的访问入口。

```go
type QyClient struct {
    config      *config.Config
    client      *http.Client
    TokenClient *token.Client
    Token       string
    deptClient  *department.Client
    userClient  *user.Client
}
```

#### 主要方法

| 方法 | 描述 | 参数 | 返回值 |
|------|------|------|--------|
| `NewQyClient(corpID, corpSecret)` | 创建客户端 | 企业 ID 和凭证 | `*QyClient, error` |
| `GetAllDepartments()` | 获取所有部门 | 无 | `[]Department, error` |
| `GetDepartmentByID(id)` | 根据 ID 获取部门 | 部门 ID | `*Department, error` |
| `GetUsersByDeptID(deptID)` | 获取部门成员列表 | 部门 ID | `[]User, error` |
| `GetUserByUserID(userID)` | 根据 UserID 获取成员 | 用户 ID | `*User, error` |

## ❌ 错误处理

项目定义了统一的企业微信错误类型。

```go
type QyError struct {
    Code int    `json:"errcode"`
    Msg  string `json:"errmsg"`
}
```

### 预定义错误

- `ErrInvalidToken` (40014): access_token 无效
- `ErrTokenExpired` (42001): access_token 过期
- `ErrDepartmentNotFound` (60001): 部门不存在

### 错误处理示例

```go
depts, err := client.GetAllDepartments()
if err != nil {
    if qyErr, ok := err.(*errorx.QyError); ok {
        fmt.Printf("企业微信错误码：%d, 错误信息：%s\n", qyErr.Code, qyErr.Msg)
    } else {
        fmt.Printf("其他错误：%v\n", err)
    }
}
```

## 🛠️ 开发指南

### 添加新的 API 支持

1. 在对应的模块目录下创建文件（如 `pkg/user/user.go`）
2. 定义请求和响应结构体
3. 实现客户端方法
4. 在 `QyClient` 中添加暴露的方法

### API 测试

可以使用 `API/` 目录下的 `.http` 文件进行接口测试（需要安装 REST Client 插件或类似工具）。

## 📝 注意事项

1. **Token 有效期**：access_token 的有效期通常为 7200 秒（2 小时），建议实现自动刷新机制
2. **调用频率限制**：企业微信对 API 调用频率有限制，请注意避免超限
3. **权限管理**：确保使用的 corpSecret 具有相应的权限
4. **并发安全**：当前实现未考虑并发安全，多线程环境下需自行加锁

## 🔗 相关资源

- [企业微信官方文档](https://developer.work.weixin.qq.com/document)
- [Go 语言官网](https://golang.org/)

## 📄 License

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。

## 🤝 贡献

欢迎提交 Issue 和 Pull Request 来帮助改进这个项目！