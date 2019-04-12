package pcsconfig

type (
	// WebConfig web 配置
	WebConfig struct {
		Addr string `json:"addr"`
		Pwd  string `json:"pwd"`
	}
)
