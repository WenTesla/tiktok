package config

// Secret 是 JWT 签名密钥，从配置加载（外部化，避免硬编码）。
var Secret = AppConfig.JWT.Secret
