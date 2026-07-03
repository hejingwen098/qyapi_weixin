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

	// 示例: 发送文本消息
	textMessage := map[string]interface{}{
		"touser":  "UserID1|UserID2",   // 接收消息的成员ID列表，多个用|分隔
		"toparty": "PartyID1|PartyID2", // 接收消息的部门ID列表，多个用|分隔
		"totag":   "TagID1|TagID2",     // 接收消息的标签ID列表，多个用|分隔
		"agentid": 1000002,             // 企业应用的id
		"text": map[string]interface{}{
			"content": "你的奖励已经到账，请前往微信查看",
		},
	}

	response, err := client.SendMessage("text", textMessage)
	if err != nil {
		log.Printf("发送消息失败：%v", err)
	} else {
		fmt.Printf("消息发送成功，响应：%s\n", response)
	}

	err = client.TokenClient.Logout(client.Token)
	if err != nil {
		log.Fatalf("Logout failed: %v", err)
	}
}
