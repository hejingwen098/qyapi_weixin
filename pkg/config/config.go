package config

// Config 企业微信配置
type Config struct {
	CorpID     string // 企业 ID
	CorpSecret string // 应用凭证密钥
}

// NewConfig 创建配置实例
func NewConfig(corpID, corpSecret string) *Config {
	cfg := &Config{
		CorpID:     corpID,
		CorpSecret: corpSecret,
	}
	return cfg
}
