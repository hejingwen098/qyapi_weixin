package config

// Config 企业微信配置
type Config struct {
	CorpID     string // 企业 ID
	CorpSecret string // 应用凭证密钥
	Proxy      string // 代理地址
}
