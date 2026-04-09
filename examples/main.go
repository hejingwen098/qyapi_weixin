package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/hejingwen098/qyapi_weixin/pkg/config"
	"github.com/hejingwen098/qyapi_weixin/pkg/qyapi"
)

func main() {
	// 配置信息（请替换为实际的企业微信配置）

	corpID := flag.String("corpid", "your_corp_id", "Corp ID")
	corpSecret := flag.String("corpsecret", "your_corp_secret", "Corp Secret")
	proxy := flag.String("proxy", "", "Proxy URL")
	flag.Parse()
	// 创建客户端,并认证
	cfg := config.Config{
		CorpID:     *corpID,
		CorpSecret: *corpSecret,
		Proxy:      *proxy,
	}
	client, err := qyapi.NewQyClient(&cfg)
	if err != nil {
		log.Fatalf("创建客户端失败：%v", err)
	}

	// 示例: 获取所有部门
	depts, err := client.GetAllDepartments()
	if err != nil {
		log.Fatalf("获取部门失败：%v", err)
	}
	fmt.Printf("部门数量：%d\n", len(depts))
	for _, dept := range depts {
		fmt.Printf("  - ID:%d, Name:%s, ParentID:%d\n", dept.ID, dept.Name, dept.ParentID)
		// 获取部门用户
		users, _ := client.GetUsersByDeptID(dept.ID)
		fmt.Printf("  - 用户数量：%d\n", len(users))
		for _, user := range users {
			fmt.Printf("    - 用户名：%s, userid: %s， Department: %v\n", user.Name, user.UserID, user.Departments[0])
		}
	}
	fmt.Println()
	err = client.TokenClient.Logout(client.Token)
	if err != nil {
		log.Fatalf("Logout failed: %v", err)
	}
}
