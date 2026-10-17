package config

// OSS（腾讯云 COS）配置从 AppConfig 加载，避免密钥硬编码。
// 注意：本项目存储层本次未替换，仍需自行配置可用的 COS 后端；
// 未配置时 /publish/action/ 不可用（见 README）。
var (
	SecretId      = AppConfig.OSS.SecretId
	SecretKey     = AppConfig.OSS.SecretKey
	CosUrl        = AppConfig.OSS.CosURL
	ReplaceSuffix = AppConfig.OSS.ReplaceSuffix
	CosUrlAnd     = AppConfig.OSS.CosURLAnd
)
